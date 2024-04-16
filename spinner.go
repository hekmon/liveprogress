package liveprogress

var (
	spinnerStates = []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
)

// Spinner is a custom item that can be added as a custom line or part of a custom DecoratorFunc.
type Spinner struct {
	lastShow int
}

func (s *Spinner) Next() string {
	s.lastShow++
	if s.lastShow >= len(spinnerStates) {
		s.lastShow = 0
	}
	return string(spinnerStates[s.lastShow])
}
