package themes

//go:generate ../../generator

type defaultTheme struct{}

func (t defaultTheme) Assets() StaticAssets {
	return StaticAssets{
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
