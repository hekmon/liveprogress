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

func SetASCIIStyle() {
	LeftEnd = '['
	Fill = '='
	Head = '>'
	Empty = '-'
	RightEnd = ']'
}

func SetUTF8ArrowsStyle() {
	LeftEnd = '◂'
	Fill = '⎯'
	Head = '→'
	Empty = ' '
	RightEnd = '▸'
}

type Bar struct {
	// ui
	fill          rune
	fillWidth     int
	head          rune
	headWidth     int
	empty         rune
	emptyWidth    int
	leftEnd       rune
	leftEndWidth  int
	rightEnd      rune
	rightEndWidth int
	width         int
	// progress
	current atomic.Uint64
	total   uint64
	// decorators
	createdAt    time.Time
	prependFuncs []DecoratorFunc
	appendFuncs  []DecoratorFunc
}

// DecoratorFunc is a function that can be prepended and appended to the progress bar
type DecoratorFunc func(b *Bar) string

func (b *Bar) Add(value uint64) {
	b.current.Add(value)
}

func (b *Bar) Current() uint64 {
	return b.current.Load()
}

func (b *Bar) Inc() {
	b.Add(1)
}

func (b *Bar) Progress() float64 {
	return float64(b.current.Load()) / float64(b.total)
}

func (b *Bar) Set(value uint64) {
	b.current.Store(value)
}

func (b *Bar) String() string {
	// Prepend fx
	pfx := make([]string, len(b.prependFuncs))
	pfxLen := 0
	pfxWidth := 0
	for index, fx := range b.prependFuncs {
		pfx[index] = fx(b)
		pfxLen += len(pfx[index])
		if b.width == 0 {
			pfxWidth += runewidth.StringWidth(pfx[index])
		}
	}
	// Append fx
	afx := make([]string, len(b.appendFuncs))
	afxLen := 0
	afxWidth := 0
	for index, fx := range b.appendFuncs {
		afx[index] = fx(b)
		afxLen += len(afx[index])
		if b.width == 0 {
			afxWidth += runewidth.StringWidth(afx[index])
		}
	}
	// Progress
	var (
		progressWidth int
		progress      strings.Builder
	)
	switch {
	case b.width == 0:
		// Calculate the width of the progress bar
		termCols, _ := liveterm.GetTermSize()
		progressWidth = termCols - pfxWidth - afxWidth
		if progressWidth < minimumProgressWidth {
			// this will break line
			progressWidth = defaultProgressWidth
		}
	case b.width < minimumProgressWidth:
		progressWidth = minimumProgressWidth
	default:
		progressWidth = b.width
	}
	enclosureWidth := b.leftEndWidth + b.rightEndWidth
	barWidth := progressWidth - enclosureWidth
	progressRatio := b.Progress()
	if progressRatio > 1 {
		progressRatio = 1
	}
	completionWidth := int(math.Round(progressRatio * float64(barWidth)))
	progress.Grow(progressWidth)
	progress.WriteRune(LeftEnd)
	completionActualWidth := 0
	if completionWidth >= b.headWidth {
		for i := 0; i < (completionWidth-b.headWidth)/b.fillWidth; i++ {
			progress.WriteRune(b.fill)
			completionActualWidth += b.fillWidth
		}
		progress.WriteRune(b.head)
		completionActualWidth += b.headWidth
	}
	for i := 0; i < barWidth-completionActualWidth; i++ {
		progress.WriteRune(b.empty)
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

func (b *Bar) Total() uint64 {
	return b.total
}
