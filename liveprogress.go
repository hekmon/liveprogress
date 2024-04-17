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

type CustomLine struct {
	generator func() string
}

func (cl *CustomLine) String() string {
	return cl.generator()
}

func AddBar(total uint64, config BarConfig) (pb *Bar) {
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
	itemsAccess.Lock()
	items = append(items, pb)
	itemsAccess.Unlock()
	return
}

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

func Bypass() io.Writer {
	return liveterm.Bypass()
}

func RemoveAll() {
	itemsAccess.Lock()
	items = make([]fmt.Stringer, 0, 1)
	itemsAccess.Unlock()
}

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

func Start() {
	liveterm.RefreshInterval = RefreshInterval
	liveterm.Output = Output
	liveterm.SetMultiLinesUpdateFx(updater)
	liveterm.Start()
}

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
