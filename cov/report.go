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
	"flag"
	"fmt"
	"github.com/axw/gocov"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"
    "encoding/json"
)

const (
    htmlHeader = `<html>
    <head>
        <title>Coverage</title>
        <!-- FIXME: Embedded style sheet -->
        <style type="text/css">
            body { background-color: black; }
            table {
                margin-left: auto;
                margin-right: auto;
            }
            td { 
                background-color: #ccc; 
                padding: 5px;
            }
            td.percent, td.linecount { text-align: right; }
            div.package, #totalcov { 
                position: fixed;
                color: brown;
                background-color: white; 
                float: left; 
                font-size: 30px;
                font-weight: bold;
                padding: 10px;
            }
            #totalcov { 
                top: 60px;
                clear: both;
                margin-top: 10px;
                color: green;
            }
            table tr:last-child td {
                font-weight: bold;
                color: brown;
            }
        </style>
    </head>
    <body>
    `
    htmlFooter = `
    </body>
</html>`
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
	packages []*gocov.Package
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
		printPackage(w, pkg)
		fmt.Fprintln(w)
	}
}

func printPackage(w io.Writer, pkg *gocov.Package) {
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

    fmt.Fprintf(w, htmlHeader)
    fmt.Fprintf(w, "<div class=\"package\">%s</div>\n", pkg.Name)
    fmt.Fprintf(w, "<div id=\"totalcov\">%s</div>\n", pkg.Name)
    fmt.Fprintf(w, "<table>\n")
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
		fmt.Fprintf(w, "<tr><td>%s/%s</td><td class=\"fn\"><code>%s(...)</code></td><td class=\"percent\">%.2f%%</td><td class=\"linecount\">%d/%d</td></tr>\n",
			pkg.Name, filepath.Base(fn.File), fn.Name, stmtPercent,
			reached, len(fn.Statements))
	}

	var funcPercent float64
	if totalStatements > 0 {
		funcPercent = float64(totalReached) / float64(totalStatements) * 100
	}
	fmt.Fprintf(w, "<tr><td>%s</td><td></td><td class=\"percent\">%.2f%%</td><td class=\"linecount\">%d/%d</td></tr>\n",
		pkg.Name, funcPercent,
		totalReached, totalStatements)
    fmt.Fprintf(w, "</table>\n")
    fmt.Fprintf(w, "<script type=\"text/javascript\">\ndocument.getElementById(\"totalcov\").textContent = \"%.2f%%\"\n</script>", funcPercent)
    fmt.Fprintf(w, htmlFooter)
}

func HTMLReportCoverage() (rc int) {
	report := newReport()
    data, err := ioutil.ReadFile(flag.Arg(0))
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to read coverage file: %s\n", err)
        return 1
    }
    packages, err := unmarshalJson(data)
    if err != nil {
        fmt.Fprintf(
            os.Stderr, "failed to unmarshal coverage data: %s\n", err)
        return 1
    }
    for _, pkg := range packages {
        report.addPackage(pkg)
    }
	fmt.Println()
	printReport(os.Stdout, report)
	return 0
}
