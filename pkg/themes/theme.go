package themes

import "github.com/matm/gocov-html/pkg/types"

var themes = []types.Beautifier{
	defaultTheme{},
}

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
