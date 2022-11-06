package types

import "text/template"

// Beautifier defines a theme used for rendering the HTML coverage stats.
type Beautifier interface {
	// Name is the name of the theme.
	Name() string
	// Description is a single line comment about the theme.
	Description() string
	// Template is the content that will be rendered.
	Template() (*template.Template, error)
}
