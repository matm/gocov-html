package themes

//go:generate ../../generator

import (
	"github.com/matm/gocov-html/pkg/types"
)

type kitTheme struct{}

func (t kitTheme) Assets() types.StaticAssets {
	return types.StaticAssets{
		Stylesheets: []string{"app.css"},
		Scripts:     []string{"app.js"},
		Index:       "index.html",
	}
}

func (t kitTheme) Name() string {
	return "kit"
}

func (t kitTheme) Description() string {
	return "AdminKit theme"
}
