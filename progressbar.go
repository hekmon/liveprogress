package liveprogress

import (
	"fmt"
	"math"
	"strings"
	"sync/atomic"
	"time"

	"github.com/hekmon/liveterm/v2"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/ansi"
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
// Auto width can be aligned with others auto width bars with WithSameAutoSize().
func WithWidth(width int) BarOption {
	return func(pb *Bar) {
		pb.barWidth = width
	}
}

// WithInternalPadding sets the padding to be internal instead of external for left and right decorators.
// Only usefull if WithSameAutoSize() has been set too.
func WithSameAutoSizeInternalPadding(left, right bool) BarOption {
	return func(pb *Bar) {
		pb.internalPaddingLeft = left
		pb.internalPaddingRight = right
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
		LeftEnd:  '[', // https://www.compart.com/unicode/U+005B
		Fill:     '=', // https://www.compart.com/unicode/U+003D
		Head:     '>', // https://www.compart.com/unicode/U+003E
		Empty:    '-', // https://www.compart.com/unicode/U+002D
		RightEnd: ']', // https://www.compart.com/unicode/U+005D
	})
}

// WithPlainRunes sets the style of the progress bar to a plain style.
func WithPlainRunes() BarOption {
	return WithRunes(BarRunes{
		LeftEnd:  0,
		Fill:     '█', // https://www.compart.com/unicode/U+2588
		Head:     '█', // https://www.compart.com/unicode/U+2588
		Empty:    '░', // https://www.compart.com/unicode/U+2591
		RightEnd: 0,
	})
}

// WithLineFillRunes sets the style of the progress bar to an box drawing lines style.
func WithLineFillRunes() BarOption {
	return WithRunes(BarRunes{
		LeftEnd:  0,
		Fill:     '━', // https://www.compart.com/unicode/U+2501
		Head:     '━', // https://www.compart.com/unicode/U+2501
		Empty:    '┅', // https://www.compart.com/unicode/U+2505
		RightEnd: 0,
	})
}

func WithMultiplyRunes() BarOption {
	return WithRunes(BarRunes{
		LeftEnd:  '❮', // https://www.compart.com/unicode/U+276E
		Fill:     '×', // https://www.compart.com/unicode/U+00D7
		Head:     '×', // https://www.compart.com/unicode/U+00D7
		Empty:    ' ', // https://www.compart.com/unicode/U+0020
		RightEnd: '❯', // https://www.compart.com/unicode/U+276F
	})
}

// WithBarStyle sets the style of the progress bar. See advanced example for style usage.
func WithBarStyle(style termenv.Style) BarOption {
	return func(pb *Bar) {
		pb.barStyle = style
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
	return WithPrependDecorator(func(pb *Bar) string {
		return style.Styled(getPercent(pb.Progress())) + " "
	})
}

// WithAppendPercent adds the percentage of the progress bar to the end of the bar.
// Use BaseStyle() if you do not want any particular style.
func WithAppendPercent(style termenv.Style) BarOption {
	return WithAppendDecorator(func(pb *Bar) string {
		return " " + style.Styled(getPercent(pb.Progress()))
	})
}

func getPercent(progress float64) (percent string) {
	progress *= 100
	percentInt := int(math.Round(progress))
	if percentInt == 100 && progress < 100 {
		// round has made us reach 100 but we don't want to show 100% if not entirely complete
		percentInt = 99
	}
	return fmt.Sprintf("%3d%%", percentInt)
}

// WithPrependTimeElapsed adds the time elapsed since the creation of the progress bar to the beginning of the bar.
// Use BaseStyle() if you do not want any particular style.
func WithPrependTimeElapsed(style termenv.Style) BarOption {
	return WithPrependDecorator(func(pb *Bar) string {
		return style.Styled(getTimeElapsed(pb.GetCreationTime())) + " "
	})
}

// WithAppendTimeElapsed adds the time elapsed since the creation of the progress bar to the end of the bar.
// Use BaseStyle() if you do not want any particular style.
func WithAppendTimeElapsed(style termenv.Style) BarOption {
	return WithAppendDecorator(func(pb *Bar) string {
		return " " + style.Styled(getTimeElapsed(pb.GetCreationTime()))
	})
}

func getTimeElapsed(start time.Time) string {
	return time.Since(start).Round(time.Second).String()
}

// WithPrependTimeRemaining adds the time remaining until the end of the progress bar to the beginning of the bar.
// Use BaseStyle() if you do not want any particular style.
func WithPrependTimeRemaining(style termenv.Style) BarOption {
	return WithPrependDecorator(func(pb *Bar) string {
		return style.Styled(getRemainingTime(pb.GetCreationTime(), pb.Progress())) + " "
	})
}

// WithAppendTimeRemaining adds the time remaining until the end of the progress bar to the end of the bar.
// Use BaseStyle() if you do not want any particular style.
func WithAppendTimeRemaining(style termenv.Style) BarOption {
	return WithAppendDecorator(func(pb *Bar) string {
		return " " + style.Styled(getRemainingTime(pb.GetCreationTime(), pb.Progress()))
	})
}

func getRemainingTime(start time.Time, progress float64) string {
	if progress == 0 {
		return "∞"
	}
	return "~" + time.Duration((1-progress)*(float64(time.Since(start))/progress)).Round(time.Second).String()
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
	barWidth             int
	internalPaddingLeft  bool
	internalPaddingRight bool
	barRunes             BarRunes
	barRunesMaxLen       int
	barRunesWidth        barRunesWidth
	barStyle             termenv.Style
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

// CurrentIncrement increments the current value of the progress bar by 1.
func (pb *Bar) CurrentIncrement() {
	pb.CurrentAdd(1)
}

// CurrentSet sets the current value of the progress bar.
func (pb *Bar) CurrentSet(value uint64) {
	pb.current.Store(value)
}

// GetCreationTime returns the time at which the progress bar was created.
func (pb *Bar) GetCreationTime() time.Time {
	return pb.createdAt
}

// Progress returns the progress of the bar as a float64 between 0 and 1.
func (pb *Bar) Progress() float64 {
	return float64(pb.current.Load()) / float64(pb.total)
}

// String returns a naive (does not support the AutoSizeSameSize) string representation of the progress bar.
func (pb *Bar) String() (line string) {
	// Generate line parts
	lineWidth, _ := liveterm.GetTermSize()
	pfx, pfxWidth := pb.renderPfx()
	afx, afxWidth := pb.renderAfx()
	bar := pb.renderProgressBar(lineWidth, pfxWidth, afxWidth, 0)
	// Assemble
	var assembler strings.Builder
	assembler.Grow(len(pfx) + len(bar) + len(afx))
	assembler.WriteString(pfx)
	assembler.WriteString(bar)
	assembler.WriteString(afx)
	line = assembler.String()
	return
}

func (pb *Bar) renderAutoSize(pfx, afx string, lineWidth, pfxWidth, pfxPadding, afxWidth, afxPadding int) string {
	// prepare
	var builder strings.Builder
	if pfxPadding < 0 {
		pfxPadding = 0
	}
	if afxPadding < 0 {
		afxPadding = 0
	}
	barWidth := lineWidth - pfxPadding - pfxWidth - afxWidth - afxPadding
	if barWidth < minimumProgressWidth {
		// will overflow
		barWidth = minimumProgressWidth
	}
	// pfx
	if !pb.internalPaddingLeft {
		builder.WriteString(strings.Repeat(" ", pfxPadding))
	}
	builder.WriteString(pfx)
	if pb.internalPaddingLeft {
		builder.WriteString(strings.Repeat(" ", pfxPadding))
	}
	// bar
	builder.WriteString(pb.renderProgressBar(lineWidth, pfxWidth, afxWidth, barWidth))
	// afx
	if pb.internalPaddingRight {
		builder.WriteString(strings.Repeat(" ", afxPadding))
	}
	builder.WriteString(afx)
	// if !pb.internalPaddingRight {
	// 	builder.WriteString(strings.Repeat(" ", afxPadding))
	// }
	// done
	return builder.String()
}

func (pb *Bar) renderPfx() (pfx string, pfxWidth int) {
	var builder strings.Builder
	for _, fx := range pb.prependFuncs {
		builder.WriteString(fx(pb))
	}
	pfx = builder.String()
	// no need to compute width if bar has a fixed size
	if pb.barWidth == 0 {
		pfxWidth = ansi.PrintableRuneWidth(pfx)
	}
	return
}

func (pb *Bar) renderAfx() (afx string, afxWidth int) {
	var builder strings.Builder
	for _, fx := range pb.appendFuncs {
		builder.WriteString(fx(pb))
	}
	afx = builder.String()
	// no need to compute width if bar has a fixed size
	if pb.barWidth == 0 {
		afxWidth = ansi.PrintableRuneWidth(afx)
	}
	return
}

func (pb *Bar) renderProgressBar(lineWidth, pfxWidth, afxWidth, overwriteBarWidth int) (bar string) {
	var (
		progressWidth int
		progress      strings.Builder
	)
	switch {
	case pb.barWidth == 0:
		if BarsAutoSizeSameSize && overwriteBarWidth != 0 {
			// Use the provided overwrite
			progressWidth = overwriteBarWidth
		} else {
			// Calculate the width of the progress bar
			progressWidth = lineWidth - pfxWidth - afxWidth
			if progressWidth < minimumProgressWidth {
				// this will break line
				progressWidth = minimumProgressWidth
			}
		}
	case pb.barWidth < minimumProgressWidth:
		progressWidth = minimumProgressWidth
	default:
		progressWidth = pb.barWidth
	}
	progress.Grow(pb.barRunesMaxLen * progressWidth) // theorical maximum number of bytes the progress bar can take
	progress.WriteRune(pb.barRunes.LeftEnd)
	barWithinWidth := progressWidth - pb.barRunesWidth.LeftEnd - pb.barRunesWidth.RightEnd
	progressRatio := pb.Progress()
	if progressRatio > 1 {
		progressRatio = 1
	}
	completionWidth := int(math.Round(progressRatio * float64(barWithinWidth)))
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
	for i := 0; i < (barWithinWidth-completionActualWidth)/pb.barRunesWidth.Empty; i++ {
		progress.WriteRune(pb.barRunes.Empty)
	}
	progress.WriteRune(pb.barRunes.RightEnd)
	return pb.barStyle.Styled(progress.String())
}

// Total returns the total value of the progress bar.
func (pb *Bar) Total() uint64 {
	return pb.total
}
