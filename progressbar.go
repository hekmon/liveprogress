package liveprogress

import (
	"fmt"
	"math"
	"strings"
	"sync/atomic"
	"time"

	"github.com/hekmon/liveterm/v2"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/termenv"
)

const (
	DefaultTotal         = 100 // DefaultTotal is the default total value of a progress bar. See WithTotal() to change a bar total at creation.
	minimumProgressWidth = 8
)

// BarOption is a function that can be used to configure a progress bar at creation, see AddBar() or SetMainLineAsBar().
type BarOption func(*Bar)

// WithTotal sets the total value of the progress bar.
func WithTotal(total uint64) BarOption {
	return func(pb *Bar) {
		pb.total = total
	}
}

// WithWidth sets the width of the progress bar.
// By default the width is set to 0: the bar will take the full terminal width (minus decorators).
func WithWidth(width int) BarOption {
	return func(pb *Bar) {
		pb.barWidth = width
	}
}

// WithRunes sets the runes used by the progress bar.
func WithRunes(runes BarRunes) BarOption {
	return func(pb *Bar) {
		if runes.Valid() {
			pb.barRunes = runes
			pb.barRunesMaxLen = runes.maxLen()
			pb.barRunesWidth = runes.width()
		}
	}
}

// WithASCIIRunes sets the style of the progress bar to an ASCII style. This is applied by default.
func WithASCIIRunes() BarOption {
	return WithRunes(BarRunes{
		LeftEnd:  '[',
		Fill:     '=',
		Head:     '>',
		Empty:    '-',
		RightEnd: ']',
	})
}

// WithPlainRunes sets the style of the progress bar to a plain style.
func WithPlainRunes() BarOption {
	return WithRunes(BarRunes{
		LeftEnd:  0,
		Fill:     '█',
		Head:     '▌',
		Empty:    '░',
		RightEnd: 0,
	})
}

// WithLineFillRunes sets the style of the progress bar to an box drawing lines style.
func WithLineFillRunes() BarOption {
	return WithRunes(BarRunes{
		LeftEnd:  0,
		Fill:     '━',
		Head:     '╍',
		Empty:    '┅',
		RightEnd: 0,
	})
}

// WithLineBracketsRunes sets the style of the progress bar to an box drawing lines style.
func WithLineBracketsRunes() BarOption {
	return WithRunes(BarRunes{
		LeftEnd:  '┣',
		Fill:     '━',
		Head:     '╸',
		Empty:    ' ',
		RightEnd: '┫',
	})
}

// WithBarStyle sets the style of the progress bar. See advanced example for style usage.
func WithBarStyle(style termenv.Style) BarOption {
	return func(pb *Bar) {
		pb.barStyle = style
		pb.barStyleLen = len(style.String())
	}
}

// DecoratorFunc is a function that can be used to decorate the progress bar.
type DecoratorFunc func(pb *Bar) string

// WithAppendDecorator adds a decorator function to the end of the progress bar.
func WithAppendDecorator(decorators ...DecoratorFunc) BarOption {
	return func(pb *Bar) {
		pb.appendFuncs = append(pb.appendFuncs, decorators...)
	}
}

// WithPrependDecorator adds a decorator function to the beginning of the progress bar.
func WithPrependDecorator(decorators ...DecoratorFunc) BarOption {
	return func(pb *Bar) {
		pb.prependFuncs = append(pb.prependFuncs, decorators...)
	}
}

// WithPrependPercent adds the percentage of the progress bar to the beginning of the bar.
// Use BaseStyle() if you do not want any particular style.
func WithPrependPercent(style termenv.Style) BarOption {
	return func(pb *Bar) {
		pb.prependFuncs = append(pb.prependFuncs, func(pb *Bar) string {
			return style.Styled(fmt.Sprintf("%3d%% ", getPercent(pb)))
		})
	}
}

// WithAppendPercent adds the percentage of the progress bar to the end of the bar.
// Use BaseStyle() if you do not want any particular style.
func WithAppendPercent(style termenv.Style) BarOption {
	return func(pb *Bar) {
		pb.appendFuncs = append(pb.appendFuncs, func(pb *Bar) string {
			return style.Styled(fmt.Sprintf(" %3d%%", getPercent(pb)))
		})
	}
}

// WithPrependTimeElapsed adds the time elapsed since the creation of the progress bar to the beginning of the bar.
// Use BaseStyle() if you do not want any particular style.
func WithPrependTimeElapsed(style termenv.Style) BarOption {
	return func(pb *Bar) {
		pb.prependFuncs = append(pb.prependFuncs, func(pb *Bar) string {
			return style.Styled(fmt.Sprintf("%s ", time.Since(pb.createdAt).Round(time.Second)))
		})
	}
}

// WithAppendTimeElapsed adds the time elapsed since the creation of the progress bar to the end of the bar.
// Use BaseStyle() if you do not want any particular style.
func WithAppendTimeElapsed(style termenv.Style) BarOption {
	return func(pb *Bar) {
		pb.appendFuncs = append(pb.appendFuncs, func(pb *Bar) string {
			return style.Styled(fmt.Sprintf(" %s", time.Since(pb.createdAt).Round(time.Second)))
		})
	}
}

// WithPrependTimeRemaining adds the time remaining until the end of the progress bar to the beginning of the bar.
// Use BaseStyle() if you do not want any particular style.
func WithPrependTimeRemaining(style termenv.Style) BarOption {
	return func(pb *Bar) {
		pb.prependFuncs = append(pb.prependFuncs, func(pb *Bar) string {
			progress := pb.Progress()
			return style.Styled(fmt.Sprintf("~%s ", time.Duration((1-progress)*(float64(time.Since(pb.createdAt))/progress)).Round(time.Second)))
		})
	}
}

// WithAppendTimeRemaining adds the time remaining until the end of the progress bar to the end of the bar.
// Use BaseStyle() if you do not want any particular style.
func WithAppendTimeRemaining(style termenv.Style) BarOption {
	return func(pb *Bar) {
		pb.appendFuncs = append(pb.appendFuncs, func(pb *Bar) string {
			progress := pb.Progress()
			return style.Styled(fmt.Sprintf(" ~%s", time.Duration((1-progress)*(float64(time.Since(pb.createdAt))/progress)).Round(time.Second)))
		})
	}
}

// BarRunes is the composition of a progress bar.
type BarRunes struct {
	LeftEnd  rune
	Fill     rune
	Head     rune
	Empty    rune
	RightEnd rune
}

// Valid returns true if all the mandatory runes are set (Fill, Head and Empty).
func (b BarRunes) Valid() bool {
	return b.Fill != 0 && b.Head != 0 && b.Empty != 0
}

func (b BarRunes) width() barRunesWidth {
	return barRunesWidth{
		LeftEnd:  runewidth.RuneWidth(b.LeftEnd),
		Fill:     runewidth.RuneWidth(b.Fill),
		Head:     runewidth.RuneWidth(b.Head),
		Empty:    runewidth.RuneWidth(b.Empty),
		RightEnd: runewidth.RuneWidth(b.RightEnd),
	}
}

func (br BarRunes) maxLen() (max int) {
	if len(string(br.LeftEnd)) > max {
		max = len(string(br.LeftEnd))
	}
	if len(string(br.Fill)) > max {
		max = len(string(br.Fill))
	}
	if len(string(br.Head)) > max {
		max = len(string(br.Head))
	}
	if len(string(br.Empty)) > max {
		max = len(string(br.Empty))
	}
	if len(string(br.RightEnd)) > max {
		max = len(string(br.RightEnd))
	}
	return
}

type barRunesWidth struct {
	LeftEnd  int
	Fill     int
	Head     int
	Empty    int
	RightEnd int
}

// Bar is a progress bar that can be added to the live progress. Do not instanciate it directly, use AddBar() instead.
type Bar struct {
	// bar config and properties
	barWidth       int
	barRunes       BarRunes
	barRunesMaxLen int
	barRunesWidth  barRunesWidth
	barStyle       termenv.Style
	barStyleLen    int
	// progress values
	current atomic.Uint64
	total   uint64
	// decorators
	createdAt    time.Time
	prependFuncs []DecoratorFunc
	appendFuncs  []DecoratorFunc
}

func newBar(opts ...BarOption) (b *Bar) {
	// Init base
	b = &Bar{
		total:        DefaultTotal,
		createdAt:    time.Now(),
		prependFuncs: make([]DecoratorFunc, 0, len(opts)),
		appendFuncs:  make([]DecoratorFunc, 0, len(opts)),
	}
	WithASCIIRunes()(b)          // default, can be overridden by opts
	WithBarStyle(BaseStyle())(b) // default, can be overridden by opts
	// Apply user options
	for _, opt := range opts {
		opt(b)
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

// GetCreationTime returns the time at which the progress bar was created.
func (pb *Bar) GetCreationTime() time.Time {
	return pb.createdAt
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
		if pb.barWidth == 0 {
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
		if pb.barWidth == 0 {
			afxWidth += runewidth.StringWidth(afx[index])
		}
	}
	// Progress
	var (
		progressWidth int
		progress      strings.Builder
	)
	switch {
	case pb.barWidth == 0:
		// Calculate the width of the progress bar
		termCols, _ := liveterm.GetTermSize()
		progressWidth = termCols - pfxWidth - afxWidth
		if progressWidth < minimumProgressWidth {
			// this will break line
			progressWidth = minimumProgressWidth
		}
	case pb.barWidth < minimumProgressWidth:
		progressWidth = minimumProgressWidth
	default:
		progressWidth = pb.barWidth
	}
	progress.Grow(pb.barRunesMaxLen * progressWidth) // theorical maximum number of bytes the progress bar can take
	progress.WriteRune(pb.barRunes.LeftEnd)
	barWidth := progressWidth - pb.barRunesWidth.LeftEnd - pb.barRunesWidth.RightEnd
	progressRatio := pb.Progress()
	if progressRatio > 1 {
		progressRatio = 1
	}
	completionWidth := int(math.Round(progressRatio * float64(barWidth)))
	completionActualWidth := 0
	if progressRatio == 1 {
		for i := 0; i < completionWidth/pb.barRunesWidth.Fill; i++ {
			progress.WriteRune(pb.barRunes.Fill)
			completionActualWidth += pb.barRunesWidth.Fill
		}
	} else if completionWidth >= pb.barRunesWidth.Head {
		for i := 0; i < (completionWidth-pb.barRunesWidth.Head)/pb.barRunesWidth.Fill; i++ {
			progress.WriteRune(pb.barRunes.Fill)
			completionActualWidth += pb.barRunesWidth.Fill
		}
		progress.WriteRune(pb.barRunes.Head)
		completionActualWidth += pb.barRunesWidth.Head
	}
	for i := 0; i < (barWidth-completionActualWidth)/pb.barRunesWidth.Empty; i++ {
		progress.WriteRune(pb.barRunes.Empty)
	}
	progress.WriteRune(pb.barRunes.RightEnd)
	// Assemble
	var assembler strings.Builder
	assembler.Grow(pfxLen + pb.barStyleLen + progress.Len() + afxLen)
	for i := 0; i < len(pfx); i++ {
		assembler.WriteString(pfx[i])
	}
	assembler.WriteString(pb.barStyle.Styled(progress.String()))
	for i := 0; i < len(afx); i++ {
		assembler.WriteString(afx[i])
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
