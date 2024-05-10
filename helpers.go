package liveprogress

import (
	"github.com/hekmon/liveterm/v2"
	"github.com/muesli/termenv"
)

// BaseStyle returns a base termenv style with its terminal profile correctly set.
// You can use it to create your own styles by modifying the returned style and use it in decorators.
// You should call this function after Start() if you have changed default Output value.
func BaseStyle() termenv.Style {
	return liveterm.GetTermProfile().String()
}

// GetTermProfile returns the termenv profile used by liveprogress (actually by liveterm).
// It can be used to create styles and colors that will be compatible with the terminal. See BaseStyle() for a more high level helper.
// You should call this function after Start() if you have changed default Output value.
func GetTermProfile() termenv.Profile {
	return liveterm.GetTermProfile()
}

// HasDarkBackground returns whether terminal uses a dark-ish background.
// You should call this function after Start() if you have changed default Output value.
func HasDarkBackground() bool {
	return liveterm.HasDarkBackground()
}

// Hyperlink creates a hyperlink that can be printed to the terminal.
func Hyperlink(link string, name string) string {
	return liveterm.Hyperlink(link, name)
}

// Notify triggers a notification.
// You should call this function after Start() if you have changed default Output value.
func Notify(title, body string) {
	liveterm.Notify(title, body)
}
