package themes

//go:generate ../../generator

import (
	"github.com/matm/gocov-html/pkg/types"
)

type defaultTheme struct{}

func (t defaultTheme) Assets() types.StaticAssets {
	return types.StaticAssets{
		Stylesheets: []string{"style.css"},
		Index:       "index.html",
	}
}

func (t defaultTheme) Name() string {
	return "golang"
}

func (t defaultTheme) Description() string {
	return "original golang theme (default)"
}
