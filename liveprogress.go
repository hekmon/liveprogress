package liveprogress

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/hekmon/liveterm/v2"
)

var (
	// Config values (used by Start())
	RefreshInterval = 100 * time.Millisecond // RefreshInterval is the time between each refresh of the terminal. Recommended value, setting it lower might flicker the terminal and increase CPU usage.
	Output          = os.Stdout              // Output is the writer the live progress will write to.
)

var (
	items       []fmt.Stringer
	mainItem    fmt.Stringer
	output      bytes.Buffer
	itemsAccess sync.Mutex
)

// AddBar adds a new progress bar to the live progress. Only call it after Start() has been called.
func AddBar(opts ...BarOption) (pb *Bar) {
	if pb = newBar(opts...); pb == nil {
		return
	}
	// Register the bar
	itemsAccess.Lock()
	items = append(items, pb)
	itemsAccess.Unlock()
	return
}

// RemoveAll removes all bars and custom lines from the live progress but does not stop the liveprogress itself.
func RemoveAll() {
	itemsAccess.Lock()
	items = make([]fmt.Stringer, 0, 1)
	mainItem = nil
	itemsAccess.Unlock()
}

// RemoveBar removes a bar from the live progress.
// This is needed only if you want to remove a bar while leaving liveprogress running, otherwise use Stop(true).
func RemoveBar(pb *Bar) {
	if pb == nil {
		return
	}
	defer itemsAccess.Unlock()
	itemsAccess.Lock()
	// Is it the main item?
	if mainItemBar, ok := mainItem.(*Bar); ok && mainItemBar == pb {
		mainItem = nil
		return
	}
	// Search for the bar
	for index, item := range items {
		if item, ok := item.(*Bar); ok && item == pb {
			items = append(items[:index], items[index+1:]...)
			break
		}
	}
}

// SetMainLineAsBar sets the main line as a bar. MainLine will always be the last line.
// Only call it after Start() has been called.
func SetMainLineAsBar(opts ...BarOption) (pb *Bar) {
	if pb = newBar(opts...); pb == nil {
		return
	}
	// Register the bar
	itemsAccess.Lock()
	mainItem = pb
	itemsAccess.Unlock()
	return
}

// Start starts the live progress. It will render every bars and custom lines added after.
// It is important to note that Output (default to os.Stdout) should not be used directly (for example with fmt.Print*()) after Start() is called and until Stop() is called.
// See ByPass() to get a writer that will bypass the live progress and write definitive lines directly to the output without disrupting live progress.
func Start() (err error) {
	liveterm.RefreshInterval = RefreshInterval
	liveterm.Output = Output
	liveterm.SetRawUpdateFx(updater)
	liveterm.HideCursor = true
	return liveterm.Start()
}

// Stop stops the live progress and remove all registered bars and custom lines from its internal state.
// Set clear to true to clear the liveprogress output. After this call, Output can be used directly again (no need to use ByPass() anymore).
func Stop(clear bool) (err error) {
	// if clear is false, liveterm will call updater one last time
	err = liveterm.Stop(clear)
	// Add a newline to separate the live progress output if needed
	if !clear {
		if output.Len() > 0 && output.Bytes()[output.Len()-1] != '\n' {
			fmt.Fprint(Output, "\n")
		}
	}
	RemoveAll()
	return
}

func updater() []byte {
	output.Reset()
	defer itemsAccess.Unlock()
	itemsAccess.Lock()
	// Choose mode
	var autoSizeSameSize int
	for _, item := range items {
		if bar, ok := item.(*Bar); ok {
			if bar.barWidth == 0 && bar.sameAutoSize {
				autoSizeSameSize++
			}
		}
	}
	if mainItem != nil {
		if mainBar, ok := mainItem.(*Bar); ok && mainBar.barWidth == 0 && mainBar.sameAutoSize {
			autoSizeSameSize++
		}
	}
	// Compute
	if autoSizeSameSize < 2 {
		// Regular 1 pass mode
		for index, item := range items {
			output.WriteString(item.String())
			if index < len(items)-1 {
				output.WriteRune('\n')
			}
		}
		if mainItem != nil {
			if len(items) > 0 {
				output.WriteRune('\n')
			}
			output.WriteString(mainItem.String())
		}
		return output.Bytes()
	}
	// 2 pass mode for bar autosize
	//// 1st pass as full auto bar size
	pfxWidths := make([]int, len(items)+1)
	barWidths := make([]int, len(items)+1)
	afxWidths := make([]int, len(items)+1)
	for index, item := range items {
		if bar, ok := item.(*Bar); ok {
			if bar.barWidth == 0 && bar.sameAutoSize {
				_, pfxWidths[index], barWidths[index], afxWidths[index] = bar.render(0, 0, 0)
			}
		}
	}
	if mainItem != nil {
		if mainBar, ok := mainItem.(*Bar); ok && mainBar.barWidth == 0 && mainBar.sameAutoSize {
			_, pfxWidths[len(barWidths)-1], barWidths[len(barWidths)-1], afxWidths[len(barWidths)-1] = mainBar.render(0, 0, 0)
		}
	}
	var (
		smallestBarPfxWidth int
		smallestBarWidth    int
		smallestBarAfxWidth int
	)
	for i := range barWidths {
		if barWidths[i] > 0 {
			if smallestBarWidth == 0 || barWidths[i] < smallestBarWidth {
				smallestBarPfxWidth = pfxWidths[i]
				smallestBarWidth = barWidths[i]
				smallestBarAfxWidth = afxWidths[i]
			}
		}
	}
	// 2nd pass as fixed bar size
	for index, item := range items {
		if bar, ok := item.(*Bar); ok {
			if bar.barWidth == 0 && bar.sameAutoSize {
				pfxPadding := smallestBarPfxWidth - pfxWidths[index]
				apfPadding := smallestBarAfxWidth - afxWidths[index]
				barStr, _, _, _ := bar.render(pfxPadding, smallestBarWidth, apfPadding)
				output.WriteString(barStr)
			} else {
				// progress bar but without same auto size
				output.WriteString(item.String())
			}
		} else {
			// custom line
			output.WriteString(item.String())
		}
		if index < len(items)-1 {
			output.WriteRune('\n')
		}
	}
	if mainItem != nil {
		if len(items) > 0 {
			output.WriteRune('\n')
		}
		if mainBar, ok := mainItem.(*Bar); ok && mainBar.barWidth == 0 && mainBar.sameAutoSize {
			pfxPadding := smallestBarPfxWidth - pfxWidths[len(barWidths)-1]
			apfPadding := smallestBarAfxWidth - afxWidths[len(barWidths)-1]
			barStr, _, _, _ := mainBar.render(pfxPadding, smallestBarWidth, apfPadding)
			output.WriteString(barStr)
		} else {
			output.WriteString(mainItem.String())
		}
	}
	return output.Bytes()
}

/*
	Specials
*/

// Bypass returns a writer that will bypass the live progress and write directly to the output without being wiped by the next refresh.
func Bypass() io.Writer {
	return liveterm.Bypass()
}

// CustomLine is a custom line to add to the live progress.
// Do not instantiate it directly, use AddCustomLine() instead.
type CustomLine struct {
	generator func() string
}

// Implements fmt.Stringer needed as a liveprogress item.
func (cl *CustomLine) String() string {
	return cl.generator()
}

// AddCustomLine adds a custom line to the live progress. Only call it after Start() has been called.
func AddCustomLine(generator func() string) (cl *CustomLine) {
	if generator == nil {
		return
	}
	itemsAccess.Lock()
	cl = &CustomLine{
		generator: generator,
	}
	items = append(items, cl)
	itemsAccess.Unlock()
	return
}

// RemoveCustomLine removes a custom line from the live progress.
func RemoveCustomLine(cl *CustomLine) {
	if cl == nil {
		return
	}
	defer itemsAccess.Unlock()
	itemsAccess.Lock()
	// Is it main item?
	if mainItemCustomLine, ok := mainItem.(*CustomLine); ok && mainItemCustomLine == cl {
		mainItem = nil
		return
	}
	// Search in other lines
	for index, item := range items {
		if item, ok := item.(*CustomLine); ok && item == cl {
			items = append(items[:index], items[index+1:]...)
			break
		}
	}
}

// SetMainLineAsCustomLine sets the main line as a custom line. MainLine will always be the last line.
// Only call it after Start() has been called.
func SetMainLineAsCustomLine(generator func() string) (cl *CustomLine) {
	if generator == nil {
		return
	}
	itemsAccess.Lock()
	cl = &CustomLine{
		generator: generator,
	}
	mainItem = cl
	itemsAccess.Unlock()
	return
}
