package liveprogress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/hekmon/liveterm"
	"github.com/mattn/go-runewidth"
)

var (
	// Config
	Output          = os.Stdout
	RefreshInterval = 100 * time.Millisecond
)

var (
	items       []fmt.Stringer
	itemsAccess sync.Mutex
)

// DecoratorAddition allows to customize a bar when creating it with AddBar().
type DecoratorAddition struct {
	Decorator DecoratorFunc
	Prepend   bool
}

// PreprendDecorator is a wrapper to facilitate the creation of a DecoratorAddition.
func PrependDecorator(decorator DecoratorFunc) DecoratorAddition {
	return DecoratorAddition{
		Decorator: decorator,
		Prepend:   true,
	}
}

// AppendDecorator is a wrapper to facilitate the creation of a DecoratorAddition.
func AppendDecorator(decorator DecoratorFunc) DecoratorAddition {
	return DecoratorAddition{
		Decorator: decorator,
		// Prepend:   false,
	}
}

// AddBar adds a new progress bar to the live progress. This does not start the live progress itself, see Start().
func AddBar(total uint64, config BarConfig, decorators ...DecoratorAddition) (pb *Bar) {
	if total == 0 {
		return
	}
	if !config.validStyle() {
		return
	}
	pb = &Bar{
		// ui
		config: config,
		styleWidth: barStyleWidth{
			LeftEnd:  runewidth.RuneWidth(config.LeftEnd),
			Fill:     runewidth.RuneWidth(config.Fill),
			Head:     runewidth.RuneWidth(config.Head),
			Empty:    runewidth.RuneWidth(config.Empty),
			RightEnd: runewidth.RuneWidth(config.RightEnd),
		},
		// progress
		createdAt: time.Now(),
		total:     total,
	}
	// decorators
	var nbPrepend, nbAppend int
	for _, decorator := range decorators {
		if decorator.Decorator == nil {
			continue
		}
		if decorator.Prepend {
			nbPrepend++
		} else {
			nbAppend++
		}
	}
	pb.prependFuncs = make([]DecoratorFunc, 0, nbPrepend)
	pb.appendFuncs = make([]DecoratorFunc, 0, nbAppend)
	for _, decorator := range decorators {
		if decorator.Decorator == nil {
			continue
		}
		if decorator.Prepend {
			pb.prependFuncs = append(pb.prependFuncs, decorator.Decorator)
		} else {
			pb.appendFuncs = append(pb.appendFuncs, decorator.Decorator)
		}
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
	itemsAccess.Unlock()
}

// RemoveBar removes a bar from the live progress.
// This is needed only if you want to remove a bar while leaving the live progress running, otherwise use Stop(true).
func RemoveBar(pb *Bar) {
	if pb == nil {
		return
	}
	itemsAccess.Lock()
	for index, item := range items {
		if item, ok := item.(*Bar); ok && item == pb {
			items = append(items[:index], items[index+1:]...)
			break
		}
	}
	itemsAccess.Unlock()
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
	lines = make([]string, len(items))
	for index, item := range items {
		lines[index] = item.String()
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
	itemsAccess.Lock()
	for index, item := range items {
		if item, ok := item.(*CustomLine); ok && item == cl {
			items = append(items[:index], items[index+1:]...)
			break
		}
	}
	itemsAccess.Unlock()
}
