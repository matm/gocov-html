// Copyright (c) 2012 The Gocov Authors.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package cov

import (
	"encoding/json"
	"fmt"
	"github.com/axw/gocov"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"
	"time"
)

func unmarshalJson(data []byte) (packages []*gocov.Package, err error) {
	result := &struct{ Packages []*gocov.Package }{}
	err = json.Unmarshal(data, result)
	if err == nil {
		packages = result.Packages
	}
	return
}

type report struct {
	packages   []*gocov.Package
	stylesheet string // absolute path to CSS
}

type reportFunction struct {
	*gocov.Function
	statementsReached int
}

type reportFunctionList []reportFunction

func (l reportFunctionList) Len() int {
	return len(l)
}

// TODO make sort method configurable?
func (l reportFunctionList) Less(i, j int) bool {
	var left, right float64
	if len(l[i].Statements) > 0 {
		left = float64(l[i].statementsReached) / float64(len(l[i].Statements))
	}
	if len(l[j].Statements) > 0 {
		right = float64(l[j].statementsReached) / float64(len(l[j].Statements))
	}
	if left < right {
		return true
	}
	return left == right && len(l[i].Statements) < len(l[j].Statements)
}

func (l reportFunctionList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type reverse struct {
	sort.Interface
}

func (r reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

// NewReport creates a new report.
func newReport() (r *report) {
	r = &report{}
	return
}

// AddPackage adds a package's coverage information to the report.
func (r *report) addPackage(p *gocov.Package) {
	i := sort.Search(len(r.packages), func(i int) bool {
		return r.packages[i].Name >= p.Name
	})
	if i < len(r.packages) && r.packages[i].Name == p.Name {
		r.packages[i].Accumulate(p)
	} else {
		head := r.packages[:i]
		tail := append([]*gocov.Package{p}, r.packages[i:]...)
		r.packages = append(head, tail...)
	}
}

// Clear clears the coverage information from the report.
func (r *report) clear() {
	r.packages = nil
}

// PrintReport prints a coverage report to the given writer.
func printReport(w io.Writer, r *report) {
	w = tabwriter.NewWriter(w, 0, 8, 0, '\t', 0)
	for _, pkg := range r.packages {
		printPackage(w, r, pkg)
		fmt.Fprintln(w)
	}
}

func printPackage(w io.Writer, r *report, pkg *gocov.Package) {
	functions := make(reportFunctionList, len(pkg.Functions))
	for i, fn := range pkg.Functions {
		reached := 0
		for _, stmt := range fn.Statements {
			if stmt.Reached > 0 {
				reached++
			}
		}
		functions[i] = reportFunction{fn, reached}
	}
	sort.Sort(reverse{functions})

	var longestFunctionName int
	var totalStatements, totalReached int

	css := defaultCSS
	if len(r.stylesheet) > 0 {
		css = fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\" />", r.stylesheet)
	}
	fmt.Fprintf(w, htmlHeader, css)
	fmt.Fprintf(w, "<div id=\"about\">Generated on %s with <a href=\"%s\">gocov-html</a></div>",
		time.Now().Format(time.RFC822Z), ProjectUrl)
	fmt.Fprintf(w, "<div class=\"package\">%s</div>\n", pkg.Name)
	fmt.Fprintf(w, "<div id=\"totalcov\">%s</div>\n", pkg.Name)
	fmt.Fprintf(w, "<div class=\"funcname\">Overview</div>")
	fmt.Fprintf(w, overview, pkg.Name)
	fmt.Fprintf(w, "<table class=\"overview\">\n")
	for _, fn := range functions {
		reached := fn.statementsReached
		totalStatements += len(fn.Statements)
		totalReached += reached
		var stmtPercent float64 = 0
		if len(fn.Statements) > 0 {
			stmtPercent = float64(reached) / float64(len(fn.Statements)) * 100
		}
		if len(fn.Name) > longestFunctionName {
			longestFunctionName = len(fn.Name)
		}
		fmt.Fprintf(w, "<tr id=\"s_fn_%s\"><td><code><a href=\"#fn_%s\">%s(...)</a></code></td><td><code>%s/%s</code></td><td class=\"percent\"><code>%.2f%%</code></td><td class=\"linecount\"><code>%d/%d</code></td></tr>\n",
			fn.Name, fn.Name, fn.Name, pkg.Name, filepath.Base(fn.File), stmtPercent,
			reached, len(fn.Statements))
	}

	var funcPercent float64
	if totalStatements > 0 {
		funcPercent = float64(totalReached) / float64(totalStatements) * 100
	}
	fmt.Fprintf(w, "<tr><td colspan=\"2\"><code>%s</code></td><td class=\"percent\"><code>%.2f%%</code></td><td class=\"linecount\"><code>%d/%d</code></td></tr>\n",
		pkg.Name, funcPercent,
		totalReached, totalStatements)
	fmt.Fprintf(w, "</table>\n")

	// Embbed function source code
	for _, fn := range functions {
		annotateFunctionSource(w, fn.Function)
	}

	fmt.Fprintf(w, "<script type=\"text/javascript\">\ndocument.getElementById(\"totalcov\").textContent = \"%.2f%%\"\n</script>", funcPercent)
	fmt.Fprintf(w, "\n<!-- Can be parsed by external script\nPACKAGE:%s DONE:%.2f\n-->\n",
		pkg.Name, funcPercent)
	fmt.Fprintf(w, htmlFooter)
}

func exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		return false, err
	}
	return true, nil
}

// HTMLReportCoverage outputs an HTML report on stdout by
// parsing JSON data generated by axw/gocov. The css parameter
// is an absolute path to a custom stylesheet. Use an empty
// string to use the default stylesheet available.
func HTMLReportCoverage(r io.Reader, css string) error {
	report := newReport()

	// Custom stylesheet?
	stylesheet := ""
	if len(css) > 0 {
		if _, err := exists(css); err != nil {
			return err
		}
		stylesheet = css
	}
	report.stylesheet = stylesheet

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read coverage data: %s\n", err)
	}

	packages, err := unmarshalJson(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal coverage data: %s\n", err)
	}

	for _, pkg := range packages {
		report.addPackage(pkg)
	}
	fmt.Println()
	printReport(os.Stdout, report)
	return nil
}
