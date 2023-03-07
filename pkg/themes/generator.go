//go:build ignore

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"

	"github.com/matm/gocov-html/pkg/types"
)

const tmpl = `
func (t defaultTheme) Data() *types.TemplateData {
	css := ""
	return &types.TemplateData{
		CSS:        css,
		When:       time.Now().Format(time.RFC822Z),
		ProjectURL: types.ProjectURL,
	}
}

func (t defaultTheme) Template() *template.Template {
	tmpl := ""
	p := template.Must(template.New("theme").Parse(tmpl))
	return p
}
`

func main() {
	name := os.Getenv("GOFILE")
	fset := token.NewFileSet()
	token, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	var assets types.StaticAssets
	var theme string

	ast.Inspect(token, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if ok {
			switch fn.Name.Name {
			case "Name":
				theme = fn.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.BasicLit).Value
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
							assets.Stylesheets = append(assets.Stylesheets, sheet)
						}
					case "Index":
						tmplName := kv.Value.(*ast.BasicLit).Value
						assets.Index = tmplName
						// TODO: case for "Scripts".
					}
				}
			}
			return false
		}
		return true
	})
	if err := render(theme, assets); err != nil {
		log.Fatal(err)
	}
}

func render(theme string, assets types.StaticAssets) error {
	fmt.Printf("%s: %+v\n", theme, assets)
	return nil
}
