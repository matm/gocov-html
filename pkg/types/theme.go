package types

import "text/template"

// Beautifier defines a theme used for rendering the HTML coverage stats.
type Beautifier interface {
	// Name is the name of the theme.
	Name() string
	// Description is a single line comment about the theme.
	Description() string
	// Template is the content that will be rendered.
	Template() *template.Template
	Data() *TemplateData
}

// TemplateData has all the fields needed by the the HTML template for rendering.
type TemplateData struct {
	// CSS is the stylesheet content that will be embedded in the HTML page.
	CSS string
	// When is the date time of report generation.
	When string
	// Overview holds data used for an additional header in case of multiple Go packages
	// have been analysed. Can be used for a high level summary. Is nil if the report has
	// only one package.
	Overview *ReportPackage
	// Packages is the list of all Go packages analysed.
	Packages ReportPackageList
	// ProjectURL is the project's site on GitHub.
	ProjectURL string
}
