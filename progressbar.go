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
	minimumProgressWidth = 8
)

var (
	// DefaultConfig can be used when creating a new progress bar
	DefaultConfig BarConfig
)

func init() {
	DefaultConfig.SetStyleASCII()
}

type BarConfig struct {
	// Config
	Width int // Width is the width of the progress bar, if 0 its size will be automatically calculated based on the terminal size and the decoractors
	// Style
	LeftEnd  rune // LeftEnd is the default character in the left most part of the progress indicator
	Fill     rune // Fill is the default character representing completed progress
	Head     rune // Head is the default character that moves when progress is updated
	Empty    rune // Empty is the default character that represents the empty progress
	RightEnd rune // RightEnd is the default character in the right most part of the progress indicator
}

func (bc *BarConfig) SetStyleASCII() {
	bc.LeftEnd = '['
	bc.Fill = '='
	bc.Head = '>'
	bc.Empty = '-'
	bc.RightEnd = ']'
}

func (bc *BarConfig) SetStyleUnicodeArrows() {
	bc.LeftEnd = '◂'
	bc.Fill = '⎯'
	bc.Head = '→'
	bc.Empty = ' '
	bc.RightEnd = '▸'
}

func (bc *BarConfig) validStyle() bool {
	return bc.LeftEnd != 0 && bc.Fill != 0 && bc.Head != 0 && bc.Empty != 0 && bc.RightEnd != 0
}

type barStyleWidth struct {
	LeftEnd  int
	Fill     int
	Head     int
	Empty    int
	RightEnd int
}

type Bar struct {
	// ui
	config     BarConfig
	styleWidth barStyleWidth
	// progress
	current atomic.Uint64
	total   uint64
	// decorators
	createdAt        time.Time
	prependFuncs     []DecoratorFunc
	appendFuncs      []DecoratorFunc
	decoratorsAccess sync.Mutex
}

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
		if pb.config.Width == 0 {
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
		if pb.config.Width == 0 {
			afxWidth += runewidth.StringWidth(afx[index])
		}
	}
	// Progress
	var (
		progressWidth int
		progress      strings.Builder
	)
	switch {
	case pb.config.Width == 0:
		// Calculate the width of the progress bar
		termCols, _ := liveterm.GetTermSize()
		progressWidth = termCols - pfxWidth - afxWidth
		if progressWidth < minimumProgressWidth {
			// this will break line
			progressWidth = minimumProgressWidth
		}
	case pb.config.Width < minimumProgressWidth:
		progressWidth = minimumProgressWidth
	default:
		progressWidth = pb.config.Width
	}
	progress.Grow(progressWidth)
	progress.WriteRune(pb.config.LeftEnd)
	barWidth := progressWidth - pb.styleWidth.LeftEnd - pb.styleWidth.RightEnd
	progressRatio := pb.Progress()
	if progressRatio > 1 {
		progressRatio = 1
	}
	completionWidth := int(math.Round(progressRatio * float64(barWidth)))
	completionActualWidth := 0
	if progressRatio == 1 {
		for i := 0; i < completionWidth/pb.styleWidth.Fill; i++ {
			progress.WriteRune(pb.config.Fill)
			completionActualWidth += pb.styleWidth.Fill
		}
	} else if completionWidth >= pb.styleWidth.Head {
		for i := 0; i < (completionWidth-pb.styleWidth.Head)/pb.styleWidth.Fill; i++ {
			progress.WriteRune(pb.config.Fill)
			completionActualWidth += pb.styleWidth.Fill
		}
		progress.WriteRune(pb.config.Head)
		completionActualWidth += pb.styleWidth.Head
	}
	for i := 0; i < barWidth-completionActualWidth; i++ {
		progress.WriteRune(pb.config.Empty)
	}
	progress.WriteRune(pb.config.RightEnd)
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

// DecoratorFunc is a function that can be prepended and appended to the progress bar
type DecoratorFunc func(pb *Bar) string

func PrependPercent() (da DecoratorAddition) {
	da.Decorator = func(pb *Bar) string {
		return fmt.Sprintf("%3d%% ", int(math.Round(pb.Progress()*100)))
	}
	da.Prepend = true
	return
}

func AppendPercent() (da DecoratorAddition) {
	da.Decorator = func(pb *Bar) string {
		return fmt.Sprintf(" %3d%%", int(math.Round(pb.Progress()*100)))
	}
	return
}

func PrependTimeElapsed() (da DecoratorAddition) {
	da.Decorator = func(pb *Bar) string {
		return fmt.Sprintf("%s ", time.Since(pb.createdAt).Round(time.Second))
	}
	da.Prepend = true
	return
}

func AppendTimeElapsed() (da DecoratorAddition) {
	da.Decorator = func(pb *Bar) string {
		return fmt.Sprintf(" %s", time.Since(pb.createdAt).Round(time.Second))
	}
	return
}

func PrependTimeRemaining() (da DecoratorAddition) {
	da.Decorator = func(pb *Bar) string {
		progress := pb.Progress()
		return fmt.Sprintf("~%s ", time.Duration((1-progress)*(float64(time.Since(pb.createdAt))/progress)).Round(time.Second))
	}
	da.Prepend = true
	return
}

func AppendTimeRemaining() (da DecoratorAddition) {
	da.Decorator = func(pb *Bar) string {
		progress := pb.Progress()
		return fmt.Sprintf(" ~%s", time.Duration((1-progress)*(float64(time.Since(pb.createdAt))/progress)).Round(time.Second))
	}
	return
}
