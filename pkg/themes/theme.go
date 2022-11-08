package themes

import (
	"github.com/matm/gocov-html/pkg/types"
	"github.com/rotisserie/eris"
)

var themes = []types.Beautifier{
	defaultTheme{},
}

// Theme to use for rendering.
var curTheme types.Beautifier = defaultTheme{}

// List returns all available themes.
func List() []types.Beautifier {
	return themes
}

// Get a theme by name. Returns nil if none found.
func Get(name string) types.Beautifier {
	for _, t := range themes {
		if t.Name() == name {
			return t
		}
	}
	return nil
}

// Use takes the name of the theme that will be used for rendering.
// Returns an error for an unknown theme.
func Use(name string) error {
	p := Get(name)
	if p == nil {
		return eris.Errorf("unknown theme %q", name)
	}
	curTheme = p
	return nil
}

// Current returns the theme to use for rendering HTML.
func Current() types.Beautifier {
	return curTheme
}
