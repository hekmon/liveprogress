package liveprogress

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hekmon/liveterm"
	"github.com/mattn/go-runewidth"
)

const (
	defaultProgressWidth = 70
	minimumProgressWidth = 12
)

var (
	// These values are copied when you call NewBar()
	LeftEnd  rune = '[' // LeftEnd is the default character in the left most part of the progress indicator
	Fill     rune = '=' // Fill is the default character representing completed progress
	Head     rune = '>' // Head is the default character that moves when progress is updated
	Empty    rune = '-' // Empty is the default character that represents the empty progress
	RightEnd rune = ']' // RightEnd is the default character in the right most part of the progress indicator
	Width         = 0   // Width is the default width of the progress bar. 0 for automatic width.
)

func SetProgressStyleASCII() {
	LeftEnd = '['
	Fill = '='
	Head = '>'
	Empty = '-'
	RightEnd = ']'
}

func SetProgressStyleUTF8Arrows() {
	LeftEnd = '◂'
	Fill = '⎯'
	Head = '→'
	Empty = ' '
	RightEnd = '▸'
}

type Bar struct {
	// ui
	fill           rune
	fillWidth      int
	head           rune
	headWidth      int
	empty          rune
	emptyWidth     int
	leftEnd        rune
	leftEndWidth   int
	rightEnd       rune
	rightEndWidth  int
	enclosureWidth int
	width          int
	// progress
	current atomic.Uint64
	total   uint64
	// decorators
	createdAt        time.Time
	prependFuncs     []DecoratorFunc
	appendFuncs      []DecoratorFunc
	decoratorsAccess sync.Mutex
}

// DecoratorFunc is a function that can be prepended and appended to the progress bar
type DecoratorFunc func(pb *Bar) string

func (pb *Bar) Current() uint64 {
	return pb.current.Load()
}

func (pb *Bar) CurrentAdd(value uint64) {
	pb.current.Add(value)
}

func (pb *Bar) CurrentCompareAndSwap(expectedCurrent, newCurrent uint64) bool {
	return pb.current.CompareAndSwap(expectedCurrent, newCurrent)
}

func (pb *Bar) CurrentSet(value uint64) {
	pb.current.Store(value)
}

func (pb *Bar) CurrentSwap(value uint64) (oldValue uint64) {
	return pb.current.Swap(value)
}

func (pb *Bar) CurrentIncrement() {
	pb.CurrentAdd(1)
}

func (pb *Bar) Progress() float64 {
	return float64(pb.current.Load()) / float64(pb.total)
}

func (pb *Bar) String() string {
	defer pb.decoratorsAccess.Unlock()
	pb.decoratorsAccess.Lock()
	// Prepend fx
	pfx := make([]string, len(pb.prependFuncs))
	pfxLen := 0
	pfxWidth := 0
	for index, fx := range pb.prependFuncs {
		pfx[index] = fx(pb)
		pfxLen += len(pfx[index])
		if pb.width == 0 {
			pfxWidth += runewidth.StringWidth(pfx[index])
		}
	}
	// Append fx
	afx := make([]string, len(pb.appendFuncs))
	afxLen := 0
	afxWidth := 0
	for index, fx := range pb.appendFuncs {
		afx[index] = fx(pb)
		afxLen += len(afx[index])
		if pb.width == 0 {
			afxWidth += runewidth.StringWidth(afx[index])
		}
	}
	// Progress
	var (
		progressWidth int
		progress      strings.Builder
	)
	switch {
	case pb.width == 0:
		// Calculate the width of the progress bar
		termCols, _ := liveterm.GetTermSize()
		progressWidth = termCols - pfxWidth - afxWidth
		if progressWidth < minimumProgressWidth {
			// this will break line
			progressWidth = minimumProgressWidth
		}
	case pb.width < minimumProgressWidth:
		progressWidth = minimumProgressWidth
	default:
		progressWidth = pb.width
	}
	progress.Grow(progressWidth)
	progress.WriteRune(pb.leftEnd)
	barWidth := progressWidth - pb.enclosureWidth
	progressRatio := pb.Progress()
	if progressRatio > 1 {
		progressRatio = 1
	}
	completionWidth := int(math.Round(progressRatio * float64(barWidth)))
	completionActualWidth := 0
	if progressRatio == 1 {
		for i := 0; i < completionWidth/pb.fillWidth; i++ {
			progress.WriteRune(pb.fill)
			completionActualWidth += pb.fillWidth
		}
	} else if completionWidth >= pb.headWidth {
		for i := 0; i < (completionWidth-pb.headWidth)/pb.fillWidth; i++ {
			progress.WriteRune(pb.fill)
			completionActualWidth += pb.fillWidth
		}
		progress.WriteRune(pb.head)
		completionActualWidth += pb.headWidth
	}
	for i := 0; i < barWidth-completionActualWidth; i++ {
		progress.WriteRune(pb.empty)
	}
	progress.WriteRune(pb.rightEnd)
	// Assemble
	var assembler strings.Builder
	assembler.Grow(pfxLen + progress.Len() + afxLen)
	for _, line := range pfx {
		assembler.WriteString(line)
	}
	assembler.WriteString(progress.String())
	for _, line := range afx {
		assembler.WriteString(line)
	}
	return assembler.String()
}

func (pb *Bar) Total() uint64 {
	return pb.total
}

/*
	Decorators
*/

func (pb *Bar) PrependFunc(fx DecoratorFunc) {
	pb.decoratorsAccess.Lock()
	pb.prependFuncs = append(pb.prependFuncs, fx)
	pb.decoratorsAccess.Unlock()
}

func (pb *Bar) AppendFunc(fx DecoratorFunc) {
	pb.decoratorsAccess.Lock()
	pb.appendFuncs = append(pb.appendFuncs, fx)
	pb.decoratorsAccess.Unlock()
}

func (pb *Bar) PrependPercent() {
	pb.PrependFunc(func(pb *Bar) string {
		return fmt.Sprintf("%3d%% ", int(math.Round(pb.Progress()*100)))
	})
}

func (pb *Bar) AppendPercent() {
	pb.AppendFunc(func(pb *Bar) string {
		return fmt.Sprintf(" %3d%%", int(math.Round(pb.Progress()*100)))
	})
}

func (pb *Bar) PrependTimeElapsed() {
	pb.PrependFunc(func(pb *Bar) string {
		return fmt.Sprintf("%s ", time.Since(pb.createdAt).Round(time.Second))
	})
}

func (pb *Bar) AppendTimeElapsed() {
	pb.AppendFunc(func(pb *Bar) string {
		return fmt.Sprintf(" %s", time.Since(pb.createdAt).Round(time.Second))
	})
}

func (pb *Bar) PrependTimeRemaining() {
	pb.PrependFunc(func(pb *Bar) string {
		progress := pb.Progress()
		return fmt.Sprintf("%s ", time.Duration((1-progress)*(float64(time.Since(pb.createdAt))/progress)).Round(time.Second))
	})
}

func (pb *Bar) AppendTimeRemaining() {
	pb.AppendFunc(func(pb *Bar) string {
		progress := pb.Progress()
		return fmt.Sprintf(" ~%s", time.Duration((1-progress)*(float64(time.Since(pb.createdAt))/progress)).Round(time.Second))
	})
}
