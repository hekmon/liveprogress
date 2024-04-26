package liveprogress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/hekmon/liveterm"
	"github.com/muesli/termenv"
)

const (
	// DefaultRefreshInterval is the recommended value for Start().
	// Setting it lower might flicker the terminal and increase CPU usage.
	DefaultRefreshInterval = 100 * time.Millisecond
)

var (
	output      io.Writer
	termOutput  *termenv.Output
	items       []fmt.Stringer
	mainItem    fmt.Stringer
	itemsAccess sync.Mutex
)

func init() {
	// Even if liveprogress is not started, bar options will need termOutput to be set
	// so we set it to os.Stdout by default, user can change it with SetOutput()
	SetOutput(os.Stdout)
}

// SetOutput sets the output writer for the live progress.
// Default is os.Stdout and if you wish to change it (for example by os.Stderr),
// you should do it before creating any bar, custom line or starting the live progress.
func SetOutput(out io.Writer) {
	output = out
	termOutput = termenv.NewOutput(output)
}

// AddBar adds a new progress bar to the live progress. This does not start the live progress itself, see Start().
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

// SetMainLineAsBar sets the main line as a bar. MainLine will always be the last line.
// This does not start the live progress itself, see Start().
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

// Start starts the live progress. It will render every bars and custom lines added previously or even after.
// It is important to note that Output (default to os.Stdout) should not be used directly (for example with fmt.Print*()) after Start() is called and until Stop() is called.
// See SetOutput() to change the output writer (call it before anything else).
// See ByPass() to get a writer that will bypass the live progress and write directly to the output without disrupting it.
func Start(refreshInterval time.Duration) (err error) {
	if refreshInterval <= time.Millisecond {
		refreshInterval = DefaultRefreshInterval
	}
	liveterm.RefreshInterval = refreshInterval
	liveterm.Output = output
	liveterm.SetMultiLinesUpdateFx(updater)
	return liveterm.Start()
}

// Stop stops the live progress and remove. Set clear to true to clear the liveprogress output.
// After this call, Output can be used directly again.
func Stop(clear bool) {
	// if clear is false, liveterm will call updater one last time (and thus locking the mutex)
	liveterm.Stop(clear)
	RemoveAll()
	termOutput = nil
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

// SetMainLineAsCustomLine sets the main line as a custom line. MainLine will always be the last line.
// This does not start the live progress itself, see Start().
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
