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
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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
	statementsReached    int
	diff float64
	newFunctionFlag bool
}

type reportFunctionList []reportFunction

func (l reportFunctionList) Len() int {
	return len(l)
}

var IsMissingLine bool
var IsCompareOldFile bool
var topN = 10
var newFunctionLimit int
func init() {

	if os.Getenv("GO_COV_HTML_REORDER") != "" {
		IsMissingLine = true
	}
	if os.Getenv("GO_COV_HTML_REORDER_TOPN") != ""{
		str := os.Getenv("GO_COV_HTML_REORDER_TOPN")
		topN = validInput(str)

	}
    if os.Getenv("GO_COV_HTML_REORDER_NEWFUNLIMIT") != ""{
    	str := os.Getenv("GO_COV_HTML_REORDER_NEWFUNLIMIT")
		newFunctionLimit = validInput(str)
	}
}

// valid In
func validInput(str  string) (int) {
	value, err := strconv.Atoi(str)
	if err != nil{
		log.Fatal("environment variable setting error")
	}
	return value
}

func (l reportFunctionList) Less(i, j int) bool {
	if IsMissingLine {
		var left, right int
		if len(l[i].Statements) > 0 {
			left = int(l[i].statementsReached) - int(len(l[i].Statements))
		}
		if len(l[j].Statements) > 0 {
			right = int(l[j].statementsReached) - int(len(l[j].Statements))
		}
		if left > right {
			return true
		}
		return left == right && len(l[i].Statements) > len(l[j].Statements)
	} else {
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

type reportPackageList []reportPackage
type Foo interface {
}
type reportPackage struct {
	pkg               *gocov.Package
	functions         reportFunctionList
	totalStatements   int
	reachedStatements int
}

func (rp *reportPackage) percentageReached() float64 {
	var rv float64
	if rp.totalStatements > 0 {
		rv = float64(rp.reachedStatements) / float64(rp.totalStatements) * 100
	}
	return rv
}

// search package in old file
func boolOldPackages(r *report, p *gocov.Package) (int, error) {
	i := sort.Search(len(r.packages), func(i int) bool {
		return r.packages[i].Name == p.Name
	})
	if i < len(r.packages) && r.packages[i].Name == p.Name {
		return i, nil
	}
	return 0, fmt.Errorf("package Name not in old file: %q  ", p.Name)
}

// search function in package after it's package has found in old file
func addOldFunction(p *gocov.Package, function *gocov.Function) (float64,bool,error) {
	for j, fn := range p.Functions {
		if fn.Name == function.Name {
			if len(p.Functions[j].Statements) != 0 {
				reached := 0
				for _, stmt := range p.Functions[j].Statements {
					if stmt.Reached > 0 {
						reached++
					}
				}
				oldReach := float64(reached) / float64(len(p.Functions[j].Statements)) * 100
				return oldReach,false, nil
			}
		}
	}
	return 0,true, fmt.Errorf("Function Name not in old file: %q ", function.Name)
}

type diffSort struct {
	name string
	diff float64
}
// use slice,top N
func diffSortFunctions(reportPackages reportPackageList) []diffSort{
	var ds []diffSort
	for _, rp := range reportPackages {
		for _ ,fn := range rp.functions{
			name := rp.pkg.Name + "/"+fn.Name
			df := diffSort{
				name : name,
				diff : fn.diff,
			}
			i := sort.Search(len(ds), func(i int) bool {
				return ds[i].diff <= df.diff
			})
			head := ds[:i]
			tail := append([]diffSort{df}, ds[i:]...)
			ds = append(head,tail...)
		}
	}
	if topN > len(ds){
		topN = len(ds)
	}
	ds = ds[:topN]
	return ds
}


func buildReportPackage(pkg *gocov.Package, oldr *report) reportPackage {
	rv := reportPackage{
		pkg:       pkg,
		functions: make(reportFunctionList, len(pkg.Functions)),
	}
	if IsCompareOldFile {
		olditem, err := boolOldPackages(oldr, pkg)
		var p *gocov.Package
		if err == nil {
			p = oldr.packages[olditem]
		}
		for i, fn := range pkg.Functions {
			reached := 0
			for _, stmt := range fn.Statements {
				if stmt.Reached > 0 {
					reached++
				}
			}
			var oldReached float64
			var newFunctionFlag bool
			if err == nil {
				oldReached,newFunctionFlag,_ = addOldFunction(p, fn)
			}
			diff := oldReached - float64(reached) / float64(len(fn.Statements)) * 100
			if diff <= 0{
				diff = 0
			}
			rv.functions[i] = reportFunction{fn, reached, diff,newFunctionFlag }
			rv.totalStatements += len(fn.Statements)
			rv.reachedStatements += reached
		}
	} else {
		for i, fn := range pkg.Functions {
			reached := 0
			for _, stmt := range fn.Statements {
				if stmt.Reached > 0 {
					reached++
				}
			}
			rv.functions[i] = reportFunction{fn, reached, 0,true}
			rv.totalStatements += len(fn.Statements)
			rv.reachedStatements += reached
		}
	}
	sort.Sort(reverse{rv.functions})
	return rv
}

// PrintReport prints a coverage report to the given writer.
func printReport(w io.Writer, r *report, r2 *report) {
	css := defaultCSS
	if len(r.stylesheet) > 0 {
		css = fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\" />", r.stylesheet)
	}
	fmt.Fprintf(w, htmlHeader, css)

	reportPackages := make(reportPackageList, len(r.packages))
	for i, pkg := range r.packages {
		reportPackages[i] = buildReportPackage(pkg, r2)
	}

	if len(reportPackages) == 0 {
		fmt.Fprintf(w, "<p>no test files in package.</p>")
		fmt.Fprintf(w, htmlFooter)
		return

	}
	summaryPackage := reportPackages[0]
	fmt.Fprintf(w, "<div id=\"about\">Generated on %s with <a href=\"%s\">gocov-html</a></div>",
		time.Now().Format(time.RFC822Z), ProjectUrl)
	if len(reportPackages) > 1 {
		summaryPackage = printReportOverview(w, reportPackages)
	}
	w = tabwriter.NewWriter(w, 0, 8, 0, '\t', 0)

	ds := diffSortFunctions(reportPackages)
	for _, rp := range reportPackages {
		printPackage(w, r, rp , ds)
		fmt.Fprintln(w)
	}

	printReportSummary(w, summaryPackage)

	fmt.Fprintf(w, htmlFooter)
}

func printReportSummary(w io.Writer, rp reportPackage) {
	fmt.Fprintf(w, "<div id=\"summaryWrapper\">")
	fmt.Fprintf(w, "<div class=\"package\">%s</div>\n", rp.pkg.Name)
	fmt.Fprintf(w, "<div id=\"totalcov\">%.2f%%</div>\n", rp.percentageReached())
	fmt.Fprintf(w, "</div>")
}

func printReportOverview(w io.Writer, reportPackages reportPackageList) reportPackage {
	rv := reportPackage{
		pkg: &gocov.Package{Name: "Report Total"},
	}
	fmt.Fprintf(w, "<div class=\"funcname\">Report Overview</div>")
	fmt.Fprintf(w, "<table class=\"overview\">\n")
	for _, rp := range reportPackages {
		rv.reachedStatements += rp.reachedStatements
		rv.totalStatements += rp.totalStatements
		fmt.Fprintf(w, "<tr id=\"s_pkg_%s\"><td><code><a href=\"#pkg_%s\">%s</a></code></td><td class=\"percent\"><code>%.2f%%</code></td><td class=\"linecount\"><code>%d/%d</code></td></tr>\n",
			rp.pkg.Name, rp.pkg.Name, rp.pkg.Name, rp.percentageReached(), rp.reachedStatements, rp.totalStatements)
	}

	fmt.Fprintf(w, "<tr><td><code>%s</code></td><td class=\"percent\"><code>%.2f%%</code></td><td class=\"linecount\"><code>%d/%d</code></td></tr>\n",
		"Report Total", rv.percentageReached(),
		rv.reachedStatements, rv.totalStatements)
	fmt.Fprintf(w, "</table>\n")

	return rv
}

func printPackage(w io.Writer, r *report, rp reportPackage,ds []diffSort) {
	fmt.Fprintf(w, "<div id=\"pkg_%s\" class=\"funcname\">Package Overview: %s <span class=\"packageTotal\">%.2f%%</span></div>", rp.pkg.Name, rp.pkg.Name, rp.percentageReached())
	fmt.Fprintf(w, overview, rp.pkg.Name, rp.pkg.Name)
	fmt.Fprintf(w, "<table class=\"overview\">\n")
	for _, fn := range rp.functions {
		// missNumber
		missingNumber := len(fn.Statements) - fn.statementsReached

		reached := fn.statementsReached
		var stmtPercent float64 = 0
		if len(fn.Statements) > 0 {
			stmtPercent = float64(reached) / float64(len(fn.Statements)) * 100
		}
		//log.Print(fn.Name,fn.newFunctionFlag,stmtPercent)

		if newFunctionLimit!=0 && fn.newFunctionFlag && stmtPercent<float64(newFunctionLimit){
			fmt.Fprintf(w, "<tr style = \"color:#ff0000 \" " )
		}else{
			fmt.Fprintf(w, "<tr " )
		}
        fmt.Fprintf(w,"id=\"s_fn_%s\"><td><code><a href=\"#fn_%s\">%s(...)</a></code></td><td><code>%s/%s</code></td><td class=\"percent\"><code>%.2f%%</code></td><td class=\"linecount\"><code>%d/%d</code></td>",
			fn.Name, fn.Name, fn.Name, rp.pkg.Name, filepath.Base(fn.File), stmtPercent,
			reached, len(fn.Statements))
		if IsMissingLine {
			fmt.Fprintf(w, "<td class=\"linecount\"><code>%d</code></td>", missingNumber)
		}
		if IsCompareOldFile {
			name := rp.pkg.Name + "/"+fn.Name
			if fn.diff > 0 {
				for _ , df := range(ds){
					if name == df.name{
						fmt.Fprintf(w, "<td><code text-align=left >&darr; %.2f%%</code></td>", fn.diff)
						break
					}
				}
			}
		}
		fmt.Fprintln(w, "</tr>\n")
	}
	allMissingNumber := rp.totalStatements - rp.reachedStatements

	fmt.Fprintf(w, "<tr><td colspan=\"2\"><code>%s</code></td><td class=\"percent\"><code>%.2f%%</code></td><td class=\"linecount\"><code>%d/%d</code></td>",
		rp.pkg.Name, rp.percentageReached(),
		rp.reachedStatements, rp.totalStatements)
	if IsMissingLine {
		fmt.Fprintf(w, "<td class=\"linecount\"><code>%d</code></td>", allMissingNumber)
	}
	fmt.Fprintln(w, "</tr>\n")
	fmt.Fprintf(w, "</table>\n")

	// Embbed function source code
	for _, fn := range rp.functions {
		annotateFunctionSource(w, fn.Function)
	}

	fmt.Fprintf(w, "\n<!-- Can be parsed by external script\nPACKAGE:%s DONE:%.2f\n-->\n",
		rp.pkg.Name, rp.percentageReached())
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
func HTMLReportCoverage(r io.Reader, css string, oldFileName string,) error {
	report := newReport()
	oldreport:= newReport()
	if oldFileName != "" {
		IsCompareOldFile = true
		var oldR io.Reader
		var errOldfile error
		if oldR, errOldfile = os.Open(oldFileName); errOldfile != nil {
			return fmt.Errorf("Usage: %s data.json\n", oldFileName)
		}
		err := oldreport.readFile(oldR)
		if err != nil {
			return err
		}
	}
	err := report.readFile(r)
	if err != nil {
		return err
	}
	// Custom stylesheet
	stylesheet := ""
	if len(css) > 0 {
		if _, err := exists(css); err != nil {
			return err
		}
		stylesheet = css
	}
	report.stylesheet = stylesheet
	printReport(os.Stdout, report, oldreport)
	return nil
}

func (re *report) readFile(r io.Reader) (error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to unmarshal coverage data: %s\n", err)
	}
	packages, err := unmarshalJson(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal coverage data: %s\n", err)
	}
	for _, pkg := range packages {
		re.addPackage(pkg)
	}
	return nil
}
