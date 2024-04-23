package liveprogress

import (
	"fmt"
	"math"
	"strings"
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

// BarConfig is the configuration of a progress bar. It includes its width and style.
type BarConfig struct {
	// Config
	Width int // Width is the width of the progress bar, if 0 its size will be automatically calculated based on terminal and decoractors width
	// Style
	LeftEnd  rune // LeftEnd is the default character in the left most part of the progress indicator (can be 0 to hide it)
	Fill     rune // Fill is the default character representing completed progress
	Head     rune // Head is the default character that moves when progress is updated
	Empty    rune // Empty is the default character that represents the empty progress
	RightEnd rune // RightEnd is the default character in the right most part of the progress indicator (can be 0 to hide it)
}

// SetStyleASCII sets the progress bar style to a simple ASCII style.
func (bc *BarConfig) SetStyleASCII() {
	bc.LeftEnd = '['
	bc.Fill = '='
	bc.Head = '>'
	bc.Empty = '-'
	bc.RightEnd = ']'
}

// SetStyleUnicodeArrows sets the progress bar style to a unicode arrows style.
func (bc *BarConfig) SetStyleUnicodeArrows() {
	bc.LeftEnd = '◂'
	bc.Fill = '⎯'
	bc.Head = '→'
	bc.Empty = ' '
	bc.RightEnd = '▸'
}

func (bc *BarConfig) validStyle() bool {
	return bc.Fill != 0 && bc.Head != 0 && bc.Empty != 0
}

type barStyleWidth struct {
	LeftEnd  int
	Fill     int
	Head     int
	Empty    int
	RightEnd int
}

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

// Bar is a progress bar that can be added to the live progress. Do not instanciate it directly, use AddBar() instead.
type Bar struct {
	// ui
	config     BarConfig
	styleWidth barStyleWidth
	// progress
	current atomic.Uint64
	total   uint64
	// decorators
	createdAt    time.Time
	prependFuncs []DecoratorFunc
	appendFuncs  []DecoratorFunc
}

func newBar(total uint64, config BarConfig, decorators ...DecoratorAddition) (pb *Bar) {
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
	return
}

// Current returns the current value of the progress bar.
func (pb *Bar) Current() uint64 {
	return pb.current.Load()
}

// CurrentAdd adds a value to the current value of the progress bar.
func (pb *Bar) CurrentAdd(value uint64) {
	pb.current.Add(value)
}

// CurrentSet sets the current value of the progress bar.
func (pb *Bar) CurrentSet(value uint64) {
	pb.current.Store(value)
}

// CurrentIncrement increments the current value of the progress bar by 1.
func (pb *Bar) CurrentIncrement() {
	pb.CurrentAdd(1)
}

// Progress returns the progress of the bar as a float64 between 0 and 1.
func (pb *Bar) Progress() float64 {
	return float64(pb.current.Load()) / float64(pb.total)
}

// String returns the string representation of the progress bar.
func (pb *Bar) String() string {
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
	for i := 0; i < (barWidth-completionActualWidth)/pb.styleWidth.Empty; i++ {
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

// Total returns the total value of the progress bar.
func (pb *Bar) Total() uint64 {
	return pb.total
}

/*
	Decorators
*/

// DecoratorFunc is a function that can be called by a progress bar in order to add custom informations to it.
type DecoratorFunc func(pb *Bar) string

// PrependPercent is a ready to use DecoratorAddition you can use when creating a new progress bar.
func PrependPercent() DecoratorAddition {
	return PrependDecorator(func(pb *Bar) string {
		return fmt.Sprintf("%3d%% ", getPercent(pb))
	})
}

// AppendPercent is a ready to use DecoratorAddition you can use when creating a new progress bar.
func AppendPercent() DecoratorAddition {
	return AppendDecorator(func(pb *Bar) string {
		return fmt.Sprintf(" %3d%%", getPercent(pb))
	})
}

func getPercent(pb *Bar) (percent int) {
	progress := pb.Progress() * 100
	percent = int(math.Round(progress))
	if percent == 100 && progress < 100 {
		// round has made up reach 100 but we don't want to show 100% if not entirely complete
		percent = 99
	}
	return
}

// PrependTimeElapsed is a ready to use DecoratorAddition you can use when creating a new progress bar.
func PrependTimeElapsed() DecoratorAddition {
	return PrependDecorator(func(pb *Bar) string {
		return fmt.Sprintf("%s ", time.Since(pb.createdAt).Round(time.Second))
	})
}

// AppendTimeElapsed is a ready to use DecoratorAddition you can use when creating a new progress bar.
func AppendTimeElapsed() DecoratorAddition {
	return AppendDecorator(func(pb *Bar) string {
		return fmt.Sprintf(" %s", time.Since(pb.createdAt).Round(time.Second))
	})
}

// PrependTimeRemaining is a ready to use DecoratorAddition you can use when creating a new progress bar.
func PrependTimeRemaining() DecoratorAddition {
	return PrependDecorator(func(pb *Bar) string {
		progress := pb.Progress()
		return fmt.Sprintf("~%s ", time.Duration((1-progress)*(float64(time.Since(pb.createdAt))/progress)).Round(time.Second))
	})
}

// AppendTimeRemaining is a ready to use DecoratorAddition you can use when creating a new progress bar.
func AppendTimeRemaining() DecoratorAddition {
	return AppendDecorator(func(pb *Bar) string {
		progress := pb.Progress()
		return fmt.Sprintf(" ~%s", time.Duration((1-progress)*(float64(time.Since(pb.createdAt))/progress)).Round(time.Second))
	})
}
