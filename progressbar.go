package liveprogress

import (
	"fmt"
	"math"
	"strings"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/hekmon/liveterm"
	"github.com/mattn/go-runewidth"
)

const (
	minimumProgressWidth = 8
)

// BarOption is a function that can be used to configure a progress bar.
type BarOption func(*Bar) error

// WithTotal sets the total value of the progress bar.
func WithTotal(total uint64) BarOption {
	return func(pb *Bar) error {
		pb.total = total
		return nil
	}
}

// WithWidth sets the width of the progress bar.
// By default the with is set to 0, this will take the full terminal width (minus decorators).
func WithWidth(width int) BarOption {
	return func(pb *Bar) error {
		pb.width = width
		return nil
	}
}

// WithStyle sets the style of the progress bar.
func WithStyle(style BarStyle) BarOption {
	return func(pb *Bar) error {
		pb.style = style
		return nil
	}
}

// WithASCIIStyle sets the style of the progress bar to an ASCII style.
func WithASCIIStyle() BarOption {
	return func(pb *Bar) error {
		pb.style = BarStyle{
			LeftEnd:  '[',
			Fill:     '=',
			Head:     '>',
			Empty:    '-',
			RightEnd: ']',
		}
		return nil
	}
}

// WithPlainStyle sets the style of the progress bar to a plain style.
func WithPlainStyle() BarOption {
	return func(pb *Bar) error {
		pb.style = BarStyle{
			LeftEnd:  0,
			Fill:     '█',
			Head:     '█',
			Empty:    '░',
			RightEnd: 0,
		}
		return nil
	}
}

// WithUnicodeArrowsStyle sets the style of the progress bar to an Unicode arrows style.
func WithUnicodeArrowsStyle() BarOption {
	return func(pb *Bar) error {
		pb.style = BarStyle{
			LeftEnd:  '◂',
			Fill:     '⎯',
			Head:     '→',
			Empty:    ' ',
			RightEnd: '▸',
		}
		return nil
	}
}

// WithBarColor sets the color of the progress bar.
// It uses lipgloss library to render the color.
// See: https://github.com/charmbracelet/lipgloss
// ANSI 16 colors (4-bit)
// lipgloss.Color("5")  // magenta
// lipgloss.Color("9")  // red
// lipgloss.Color("12") // light blue
// ANSI 256 Colors (8-bit)
// lipgloss.Color("86")  // aqua
// lipgloss.Color("201") // hot pink
// lipgloss.Color("202") // orange
// True Color (16,777,216 colors; 24-bit)
// lipgloss.Color("#0000FF") // good ol' 100% blue
// lipgloss.Color("#04B575") // a green
// lipgloss.Color("#3C3C3C") // a dark gray
// ...as well as a 1-bit ASCII profile, which is black and white only.
// The terminal's color profile will be automatically detected, and colors outside the gamut
// of the current palette will be automatically coerced to their closest available value.
func WithBarColor(color string) BarOption {
	return func(pb *Bar) error {
		pb.barStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(color))
		return nil
	}
}

// DecoratorFunc is a function that can be used to decorate the progress bar.
type DecoratorFunc func(pb *Bar) string

// WithAppendDecorator adds a decorator function to the end of the progress bar.
func WithAppendDecorator(decorators ...DecoratorFunc) BarOption {
	return func(pb *Bar) error {
		pb.appendFuncs = append(pb.appendFuncs, decorators...)
		return nil
	}
}

// WithPrependDecorator adds a decorator function to the beginning of the progress bar.
func WithPrependDecorator(decorators ...DecoratorFunc) BarOption {
	return func(pb *Bar) error {
		pb.prependFuncs = append(pb.prependFuncs, decorators...)
		return nil
	}
}

// WithAppendPercent adds the percentage of the progress bar to the end of the bar.
func WithAppendPercent(color string) BarOption {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	return func(pb *Bar) error {
		pb.appendFuncs = append(pb.appendFuncs, func(pb *Bar) string {
			return style.Render(fmt.Sprintf("%3d%% ", getPercent(pb)))
		})
		return nil
	}
}

// WithPrependPercent adds the percentage of the progress bar to the beginning of the bar.
func WithPrependPercent(color string) BarOption {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	return func(pb *Bar) error {
		pb.prependFuncs = append(pb.prependFuncs, func(pb *Bar) string {
			return style.Render(fmt.Sprintf("%3d%% ", getPercent(pb)))
		})
		return nil
	}
}

// WithAppendTimeElapsed adds the time elapsed since the creation of the progress bar to the end of the bar.
func WithAppendTimeElapsed(color string) BarOption {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	return func(pb *Bar) error {
		pb.appendFuncs = append(pb.appendFuncs, func(pb *Bar) string {
			return style.Render(fmt.Sprintf(" %s", time.Since(pb.createdAt).Round(time.Second)))
		})
		return nil
	}
}

// WithPrependTimeElapsed adds the time elapsed since the creation of the progress bar to the beginning of the bar.
func WithPrependTimeElapsed(color string) BarOption {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	return func(pb *Bar) error {
		pb.prependFuncs = append(pb.prependFuncs, func(pb *Bar) string {
			return style.Render(fmt.Sprintf("%s ", time.Since(pb.createdAt).Round(time.Second)))
		})
		return nil
	}
}

// WithAppendTimeRemaining adds the time remaining until the end of the progress bar to the end of the bar.
func WithAppendTimeRemaining(color string) BarOption {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	return func(pb *Bar) error {
		pb.appendFuncs = append(pb.appendFuncs, func(pb *Bar) string {
			progress := pb.Progress()
			return style.Render(fmt.Sprintf(" ~%s", time.Duration((1-progress)*(float64(time.Since(pb.createdAt))/progress)).Round(time.Second)))
		})
		return nil
	}
}

// WithPrependTimeRemaining adds the time remaining until the end of the progress bar to the beginning of the bar.
func WithPrependTimeRemaining(color string) BarOption {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	return func(pb *Bar) error {
		pb.prependFuncs = append(pb.prependFuncs, func(pb *Bar) string {
			progress := pb.Progress()
			return style.Render(fmt.Sprintf("~%s ", time.Duration((1-progress)*(float64(time.Since(pb.createdAt))/progress)).Round(time.Second)))
		})
		return nil
	}
}

type barStyleWidth struct {
	LeftEnd  int
	Fill     int
	Head     int
	Empty    int
	RightEnd int
}

// BarStyle is the style of a progress bar.
type BarStyle struct {
	LeftEnd  rune
	Fill     rune
	Head     rune
	Empty    rune
	RightEnd rune
}

func (b BarStyle) width() barStyleWidth {
	return barStyleWidth{
		LeftEnd:  runewidth.RuneWidth(b.LeftEnd),
		Fill:     runewidth.RuneWidth(b.Fill),
		Head:     runewidth.RuneWidth(b.Head),
		Empty:    runewidth.RuneWidth(b.Empty),
		RightEnd: runewidth.RuneWidth(b.RightEnd),
	}
}

// Bar is a progress bar that can be added to the live progress. Do not instanciate it directly, use AddBar() instead.
type Bar struct {
	// style
	width      int
	style      BarStyle
	styleWidth barStyleWidth
	barStyle   lipgloss.Style
	// progress
	current atomic.Uint64
	total   uint64
	// decorators
	createdAt    time.Time
	prependFuncs []DecoratorFunc
	appendFuncs  []DecoratorFunc
}

func newBar(opts ...BarOption) *Bar {
	style := BarStyle{
		LeftEnd:  '[',
		Fill:     '=',
		Head:     '>',
		Empty:    '-',
		RightEnd: ']',
	}
	bar := Bar{
		style:      style,
		styleWidth: style.width(),
		barStyle:   lipgloss.NewStyle(),
		createdAt:  time.Now(),
		total:      100,
	}
	for _, opt := range opts {
		opt(&bar)
	}
	return &bar
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
	progress.WriteRune(pb.style.LeftEnd)
	barWidth := progressWidth - pb.styleWidth.LeftEnd - pb.styleWidth.RightEnd
	progressRatio := pb.Progress()
	if progressRatio > 1 {
		progressRatio = 1
	}
	completionWidth := int(math.Round(progressRatio * float64(barWidth)))
	completionActualWidth := 0
	if progressRatio == 1 {
		for i := 0; i < completionWidth/pb.styleWidth.Fill; i++ {
			progress.WriteRune(pb.style.Fill)
			completionActualWidth += pb.styleWidth.Fill
		}
	} else if completionWidth >= pb.styleWidth.Head {
		for i := 0; i < (completionWidth-pb.styleWidth.Head)/pb.styleWidth.Fill; i++ {
			progress.WriteRune(pb.style.Fill)
			completionActualWidth += pb.styleWidth.Fill
		}
		progress.WriteRune(pb.style.Head)
		completionActualWidth += pb.styleWidth.Head
	}
	for i := 0; i < (barWidth-completionActualWidth)/pb.styleWidth.Empty; i++ {
		progress.WriteRune(pb.style.Empty)
	}
	progress.WriteRune(pb.style.RightEnd)
	// Assemble
	var assembler strings.Builder
	assembler.Grow(pfxLen + progress.Len() + afxLen)
	for _, line := range pfx {
		assembler.WriteString(line)
	}
	assembler.WriteString(pb.barStyle.Render(progress.String()))
	for _, line := range afx {
		assembler.WriteString(line)
	}
	return assembler.String()
}

// Total returns the total value of the progress bar.
func (pb *Bar) Total() uint64 {
	return pb.total
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
