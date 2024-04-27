package liveprogress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/hekmon/liveterm/v2"
	"github.com/muesli/termenv"
)

var (
	// Config values (used by Start())
	RefreshInterval = 100 * time.Millisecond // RefreshInterval is the time between each refresh of the terminal. Recommended value, setting it lower might flicker the terminal and increase CPU usage.
	Output          = os.Stdout              // Output is the writer the live progress will write to.
)

var (
	items       []fmt.Stringer
	mainItem    fmt.Stringer
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

// BaseStyle returns a base termenv style with its terminal profile correctly set.
// You can use it to create your own styles by modifying the returned style and use it in decorators.
// You should call this function after Start() if you have changed default Output value.
func BaseStyle() termenv.Style {
	return liveterm.GetTermProfile().String()
}

// GetTermProfile returns the termenv profile used by liveprogress (actually by liveterm).
// It can be used to create styles and colors that will be compatible with the terminal. See BaseStyle() for a more high level helper.
// You should call this function after Start() if you have changed default Output value.
func GetTermProfile() termenv.Profile {
	return liveterm.GetTermProfile()
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
	liveterm.SetMultiLinesUpdateFx(updater)
	return liveterm.Start()
}

// Stop stops the live progress and remove all registered bars and custom lines from its internal state.
// Set clear to true to clear the liveprogress output. After this call, Output can be used directly again (no need to use ByPass() anymore).
func Stop(clear bool) (err error) {
	// if clear is false, liveterm will call updater one last time
	err = liveterm.Stop(clear)
	RemoveAll()
	return
}

func updater() (lines []string) {
	itemsAccess.Lock()
	nbLines := len(items)
	if mainItem != nil {
		nbLines++
	}
	lines = make([]string, nbLines)
	for index, item := range items {
		lines[index] = item.String()
	}
	if mainItem != nil {
		lines[nbLines-1] = mainItem.String()
	}
	itemsAccess.Unlock()
	return
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
