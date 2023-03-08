package themes

//go:generate ../../generator

import (
	"github.com/matm/gocov-html/pkg/types"
)

type kitTheme struct{}

func (t kitTheme) Assets() types.StaticAssets {
	return types.StaticAssets{
		Stylesheets: []string{
			// From the official theme.
			"app.css",
			//"a.css", "b.css",
			// Custom rules.
			"kit.css",
		},
		Scripts: []string{"app.js"},
		Index:   "index.html",
	}
}

func (t kitTheme) Name() string {
	return "kit"
}

func (t kitTheme) Description() string {
	return "AdminKit theme"
}
