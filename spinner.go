package liveprogress

var (
	spinnerStates = []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
)

// Spinner is a custom item that can be added as custom DecoratorFunc.
type Spinner struct {
	lastShown int
}

// Next returns the next spinner state, call it in a loop to animate the spinner.
func (s *Spinner) Next() string {
	s.lastShown++
	if s.lastShown >= len(spinnerStates) {
		s.lastShown = 0
	}
	return string(spinnerStates[s.lastShown])
}
