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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/axw/gocov"
	"github.com/matm/gocov-html/pkg/themes"
	"github.com/matm/gocov-html/pkg/types"
	"github.com/rotisserie/eris"
)

// ReportOptions holds various options used when generating the final
// HTML report.
type ReportOptions struct {
	// LowCoverageOnTop puts low coverage functions first.
	LowCoverageOnTop bool
	// Stylesheet is the path to a custom CSS file.
	Stylesheet string
	// CoverageMin filters out all functions whose code coverage is smaller than it is.
	CoverageMin uint8
	// CoverageMax filters out all functions whose code coverage is greater than it is.
	CoverageMax uint8
	// command list arguments.
	cliArgs string
}

type report struct {
	ReportOptions
	packages []*gocov.Package
}

func unmarshalJSON(data []byte) (packages []*gocov.Package, err error) {
	result := &struct{ Packages []*gocov.Package }{}
	err = json.Unmarshal(data, result)
	if err == nil {
		packages = result.Packages
	}
	return
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

func buildReportPackage(pkg *gocov.Package, r *report) types.ReportPackage {
	rv := types.ReportPackage{
		Pkg:       pkg,
		Functions: make(types.ReportFunctionList, 0),
	}
	for _, fn := range pkg.Functions {
		reached := 0
		for _, stmt := range fn.Statements {
			if stmt.Reached > 0 {
				reached++
			}
		}
		rf := types.ReportFunction{Function: fn, StatementsReached: reached}
		covp := rf.CoveragePercent()
		if covp > float64(r.CoverageMin) && covp < float64(r.CoverageMax) {
			rv.Functions = append(rv.Functions, rf)
		}
		rv.TotalStatements += len(fn.Statements)
		rv.ReachedStatements += reached
	}
	if r.LowCoverageOnTop {
		sort.Sort(rv.Functions)
	} else {
		sort.Sort(reverse{rv.Functions})
	}
	return rv
}

// printReport prints a coverage report to the given writer.
func printReport(w io.Writer, r *report) error {
	theme := themes.Current()
	data := theme.Data()

	// Base64 decoding of style data and script.
	s, err := base64.StdEncoding.DecodeString(data.Style)
	if err != nil {
		return eris.Wrap(err, "decode style")
	}
	css := string(s)
	// Decode the script also.
	sc, err := base64.StdEncoding.DecodeString(data.Script)
	if err != nil {
		return eris.Wrap(err, "decode script")
	}

	if len(r.Stylesheet) > 0 {
		// Inline CSS.
		f, err := os.Open(r.Stylesheet)
		if err != nil {
			return eris.Wrap(err, "print report")
		}
		style, err := ioutil.ReadAll(f)
		if err != nil {
			return eris.Wrap(err, "read style")
		}
		css = string(style)
	}
	reportPackages := make(types.ReportPackageList, len(r.packages))
	pkgNames := make([]string, len(r.packages))
	for i, pkg := range r.packages {
		reportPackages[i] = buildReportPackage(pkg, r)
		pkgNames[i] = pkg.Name
	}

	data.Script = string(sc)
	data.Style = css
	data.Packages = reportPackages
	data.Command = fmt.Sprintf("gocov test %s | gocov-html %s",
		strings.Join(pkgNames, " "),
		strings.Join(os.Args[1:], " "),
	)

	if len(reportPackages) > 1 {
		rv := types.ReportPackage{
			Pkg: &gocov.Package{Name: "Report Total"},
		}
		for _, rp := range reportPackages {
			rv.ReachedStatements += rp.ReachedStatements
			rv.TotalStatements += rp.TotalStatements
		}
		data.Overview = &rv
	}
	err = theme.Template().Execute(w, data)
	return eris.Wrap(err, "execute template")
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
func HTMLReportCoverage(r io.Reader, opts ReportOptions) error {
	t0 := time.Now()
	report := newReport()
	report.ReportOptions = opts

	// Custom stylesheet?
	stylesheet := ""
	if opts.Stylesheet != "" {
		if _, err := exists(opts.Stylesheet); err != nil {
			return eris.Wrap(err, "stylesheet")
		}
		stylesheet = opts.Stylesheet
	}
	report.Stylesheet = stylesheet

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return eris.Wrap(err, "read coverage data")
	}

	packages, err := unmarshalJSON(data)
	if err != nil {
		return eris.Wrap(err, "unmarshal coverage data")
	}

	for _, pkg := range packages {
		report.addPackage(pkg)
	}
	fmt.Println()
	err = printReport(os.Stdout, report)
	fmt.Fprintf(os.Stderr, "Took %v\n", time.Since(t0))
	return eris.Wrap(err, "HTML report")
}
