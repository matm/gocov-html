package themes

import "text/template"

type defaultTheme struct{}

func (t defaultTheme) Name() string {
	return "golang"
}

func (t defaultTheme) Description() string {
	return "original golang theme"
}

func (t defaultTheme) Template() (*template.Template, error) {
	p, err := template.New("theme").Parse(`{{define "theme"}}Hello{{end}}`)
	return p, err
}
