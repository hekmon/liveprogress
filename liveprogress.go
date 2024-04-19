package liveprogress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/hekmon/liveterm"
)

var (
	// Config
	Output          = os.Stdout
	RefreshInterval = 100 * time.Millisecond
)

var (
	items       []fmt.Stringer
	mainItem    fmt.Stringer
	itemsAccess sync.Mutex
)

// AddBar adds a new progress bar to the live progress. This does not start the live progress itself, see Start().
func AddBar(total uint64, config BarConfig, decorators ...DecoratorAddition) (pb *Bar) {
	if pb = newBar(total, config, decorators...); pb == nil {
		return
	}
	// Register the bar
	itemsAccess.Lock()
	items = append(items, pb)
	itemsAccess.Unlock()
	return
}

// RemoveAll removes all bars and custom lines from the live progress but do not stops the liveprogress itself.
func RemoveAll() {
	itemsAccess.Lock()
	items = make([]fmt.Stringer, 0, 1)
	mainItem = nil
	itemsAccess.Unlock()
}

// RemoveBar removes a bar from the live progress.
// This is needed only if you want to remove a bar while leaving the live progress running, otherwise use Stop(true).
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

// SetMainLineAsBar sets the main line as a bar. This does not start the live progress itself, see Start().
// Main item will always be the last line.
func SetMainLineAsBar(total uint64, config BarConfig, decorators ...DecoratorAddition) (pb *Bar) {
	if pb = newBar(total, config, decorators...); pb == nil {
		return
	}
	// Register the bar
	itemsAccess.Lock()
	mainItem = pb
	itemsAccess.Unlock()
	return
}

// Start starts the live progress. It will render every bar and custom line added.
// It is imported to note that output (default to os.Stdout) should not be used after Start() is called and until Stop() is called.
// Se ByPass() to get a writer that will bypass the live progress and write directly to the output without being wiped between Start() and Stop().
func Start() (err error) {
	liveterm.RefreshInterval = RefreshInterval
	liveterm.Output = Output
	liveterm.SetMultiLinesUpdateFx(updater)
	return liveterm.Start()
}

// Stop stops the live progress and remove. Set clear to true to clear the liveprogress output.
// After this output can be used directly again.
func Stop(clear bool) {
	// if clear is false, liveterm will call updater one last time (and thus locking the mutex)
	liveterm.Stop(clear)
	RemoveAll()
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

// Bypass returns a writer that will bypass the live progress and write directly to the output without being wiped.
func Bypass() io.Writer {
	return liveterm.Bypass()
}

// CustomLine is a custom line to add to the live progress.
// Do not instantiate it directly, use AddCustomLine instead.
type CustomLine struct {
	generator func() string
}

// Implements fmt.Stringer needed as a liveprogress item.
func (cl *CustomLine) String() string {
	return cl.generator()
}

// AddCustomLine adds a custom line to the live progress.
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
