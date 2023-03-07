//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/matm/gocov-html/pkg/types"
)

const tmpl = `
package themes

import (
	"text/template"
	"time"
	"github.com/matm/gocov-html/pkg/types"
)

func (t defaultTheme) Data() *types.TemplateData {
	css := {{.Style}}
	return &types.TemplateData{
		CSS:        css,
		When:       time.Now().Format(time.RFC822Z),
		ProjectURL: types.ProjectURL,
	}
}

func (t defaultTheme) Template() *template.Template {
	tmpl := {{.Template}}
	p := template.Must(template.New("theme").Parse(tmpl))
	return p
}
`

func inspect(name string, theme *string, assets *types.StaticAssets) error {
	fset := token.NewFileSet()
	token, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	ast.Inspect(token, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if ok {
			switch fn.Name.Name {
			case "Name":
				*theme = fn.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.BasicLit).Value
				*theme = strings.Replace(*theme, `"`, "", -1)
			case "Assets":
				es := fn.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CompositeLit).Elts
				for _, e := range es {
					kv := e.(*ast.KeyValueExpr)
					id := kv.Key.(*ast.Ident)
					switch id.Name {
					case "Stylesheets":
						elems := kv.Value.(*ast.CompositeLit).Elts
						for _, elem := range elems {
							sheet := elem.(*ast.BasicLit).Value
							assets.Stylesheets = append(assets.Stylesheets, strings.Replace(sheet, `"`, "", -1))
						}
					case "Index":
						tmplName := kv.Value.(*ast.BasicLit).Value
						assets.Index = strings.Replace(tmplName, `"`, "", -1)
						// TODO: case for "Scripts".
					}
				}
			}
			return false
		}
		return true
	})
	return nil
}

func render(name, theme string, assets types.StaticAssets) error {
	baseThemeDir := path.Join("..", "..", "themes", theme)
	out := strings.Replace(name, ".go", "_gen.go", 1)
	outFile, err := os.Create(out)
	if err != nil {
		return err
	}
	defer outFile.Close()
	index, err := ioutil.ReadFile(path.Join(baseThemeDir, assets.Index))
	if err != nil {
		return err
	}
	// Contains all stylesheets' data.
	var allStyles bytes.Buffer
	for _, css := range assets.Stylesheets {
		style, err := ioutil.ReadFile(path.Join(baseThemeDir, css))
		if err != nil {
			return err
		}
		fmt.Fprintf(&allStyles, "`%s`\n", style)
	}
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return err
	}
	type data struct {
		Style    string
		Template string
	}
	err = t.Execute(outFile, &data{
		Style:    allStyles.String(),
		Template: "`" + string(index) + "`"},
	)
	return err
}

func main() {
	name := os.Getenv("GOFILE")
	assets := new(types.StaticAssets)
	theme := new(string)
	err := inspect(name, theme, assets)
	if err != nil {
		log.Fatal(err)
	}
	if err := render(name, *theme, *assets); err != nil {
		log.Fatal(err)
	}
}
