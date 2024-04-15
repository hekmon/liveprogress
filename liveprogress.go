package termprogress

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
	itemsAccess sync.Mutex
)

func AddBar() (bar *Bar) {
	bar = &Bar{}
	itemsAccess.Lock()
	items = append(items, bar)
	itemsAccess.Unlock()
	return
}

func AddCustom(custom fmt.Stringer) {
	itemsAccess.Lock()
	items = append(items, custom)
	itemsAccess.Unlock()
}

func Bypass() io.Writer {
	return liveterm.Bypass()
}

func RemoveAll() {
	itemsAccess.Lock()
	items = make([]fmt.Stringer, 0, 1)
	itemsAccess.Unlock()
}

func RemoveBar(bar *Bar) {
	if bar == nil {
		return
	}
	itemsAccess.Lock()
	for index, item := range items {
		if b, ok := item.(*Bar); ok && b == bar {
			items = append(items[:index], items[index+1:]...)
			break
		}
	}
	itemsAccess.Unlock()
}

func RemoveCustom(item fmt.Stringer) {
	itemsAccess.Lock()
	for index, registered := range items {
		if item == registered {
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

type customLine struct {
	fx func() string
}

func (c *customLine) String() string {
	return c.fx()
}

func NewCustomLine(fx func() string) fmt.Stringer {
	return &customLine{fx: fx}
}
