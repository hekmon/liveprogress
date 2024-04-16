package termprogress

import (
	"math"
	"strings"
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

type ProgressBar struct {
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
	createdAt    time.Time
	prependFuncs []DecoratorFunc
	appendFuncs  []DecoratorFunc
}

// DecoratorFunc is a function that can be prepended and appended to the progress bar
type DecoratorFunc func(pb *ProgressBar) string

func (pb *ProgressBar) Add(value uint64) {
	pb.current.Add(value)
}

func (pb *ProgressBar) CompareAndSwap(expectedCurrent, newCurrent uint64) bool {
	return pb.current.CompareAndSwap(expectedCurrent, newCurrent)
}

func (pb *ProgressBar) Current() uint64 {
	return pb.current.Load()
}

func (pb *ProgressBar) Inc() {
	pb.Add(1)
}

func (pb *ProgressBar) Progress() float64 {
	return float64(pb.current.Load()) / float64(pb.total)
}

func (pb *ProgressBar) Set(value uint64) {
	pb.current.Store(value)
}

func (pb *ProgressBar) String() string {
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
	progress.WriteRune(LeftEnd)
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
	progress.WriteRune(RightEnd)
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

func (pb *ProgressBar) Total() uint64 {
	return pb.total
}
