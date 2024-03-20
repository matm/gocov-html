package themes

//go:generate ../../generator

type separatePackageTheme struct{}

func (t separatePackageTheme) Assets() StaticAssets {
	return StaticAssets{
		Stylesheets: []string{"style.css"},
		Index:       "index.html",
	}
}

func (t separatePackageTheme) Name() string {
	return "separate"
}

func (t separatePackageTheme) Description() string {
	return "original golang theme"
}
