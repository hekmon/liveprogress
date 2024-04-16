package termprogress

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

func AddBar(total uint64) (pb *ProgressBar) {
	if total == 0 {
		return
	}
	pb = &ProgressBar{
		// ui
		fill:          Fill,
		fillWidth:     runewidth.RuneWidth(Fill),
		head:          Head,
		headWidth:     runewidth.RuneWidth(Head),
		empty:         Empty,
		emptyWidth:    runewidth.RuneWidth(Empty),
		leftEnd:       LeftEnd,
		leftEndWidth:  runewidth.RuneWidth(LeftEnd),
		rightEnd:      RightEnd,
		rightEndWidth: runewidth.RuneWidth(RightEnd),
		width:         Width,
		// progress
		createdAt: time.Now(),
		total:     total,
	}
	pb.enclosureWidth = pb.leftEndWidth + pb.rightEndWidth
	itemsAccess.Lock()
	items = append(items, pb)
	itemsAccess.Unlock()
	return
}

func AddCustom(custom fmt.Stringer) {
	if custom == nil {
		return
	}
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

func RemoveBar(pb *ProgressBar) {
	if pb == nil {
		return
	}
	itemsAccess.Lock()
	for index, item := range items {
		if item, ok := item.(*ProgressBar); ok && item == pb {
			items = append(items[:index], items[index+1:]...)
			break
		}
	}
	itemsAccess.Unlock()
}

func RemoveCustom(item fmt.Stringer) {
	if item == nil {
		return
	}
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

type StringGenerator func() string

type stringerWrapper struct {
	sg StringGenerator
}

func (sw *stringerWrapper) String() string {
	return sw.sg()
}

func NewCustomLine(sg StringGenerator) fmt.Stringer {
	return &stringerWrapper{sg: sg}
}
