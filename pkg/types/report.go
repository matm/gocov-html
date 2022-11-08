package types

import (
	"go/token"
	"html"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/axw/gocov"
)

const ProjectUrl = "https://github.com/matm/gocov-html"

type ReportPackageList []ReportPackage

type ReportPackage struct {
	Pkg               *gocov.Package
	Functions         ReportFunctionList
	TotalStatements   int
	ReachedStatements int
}

func (rp *ReportPackage) PercentageReached() float64 {
	var rv float64
	if rp.TotalStatements > 0 {
		rv = float64(rp.ReachedStatements) / float64(rp.TotalStatements) * 100
	}
	return rv
}

type ReportFunction struct {
	*gocov.Function
	StatementsReached int
}

type FunctionLine struct {
	Missed     bool
	LineNumber int
	Code       string
}

func (f ReportFunction) CoveragePercent() float64 {
	reached := f.StatementsReached
	var stmtPercent float64 = 0
	if len(f.Statements) > 0 {
		stmtPercent = float64(reached) / float64(len(f.Statements)) * 100
	} else if len(f.Statements) == 0 {
		stmtPercent = 100
	}
	return stmtPercent
}

func (f ReportFunction) ShortFileName() string {
	return filepath.Base(f.File)
}

const (
	hitPrefix  = "    "
	missPrefix = "MISS"
)

type annotator struct {
	fset  *token.FileSet
	files map[string]*token.File
}

func (fn ReportFunction) Lines() []FunctionLine {
	a := &annotator{}
	a.fset = token.NewFileSet()
	a.files = make(map[string]*token.File)

	// Load the file for line information. Probably overkill, maybe
	// just compute the lines from offsets in here.
	setContent := false
	file := a.files[fn.File]
	if file == nil {
		info, err := os.Stat(fn.File)
		if err != nil {
			panic(err)
		}
		file = a.fset.AddFile(fn.File, a.fset.Base(), int(info.Size()))
		setContent = true
	}

	data, err := ioutil.ReadFile(fn.File)
	if err != nil {
		panic(err)
	}

	if setContent {
		// This processes the content and records line number info.
		file.SetLinesForContent(data)
	}

	statements := fn.Statements[:]
	lineno := file.Line(file.Pos(fn.Start))
	lines := strings.Split(string(data)[fn.Start:fn.End], "\n")
	fls := make([]FunctionLine, len(lines))

	for i, line := range lines {
		lineno := lineno + i
		statementFound := false
		hit := false
		for j := 0; j < len(statements); j++ {
			start := file.Line(file.Pos(statements[j].Start))
			if start == lineno {
				statementFound = true
				if !hit && statements[j].Reached > 0 {
					hit = true
				}
				statements = append(statements[:j], statements[j+1:]...)
			}
		}
		hitmiss := hitPrefix
		if statementFound && !hit {
			hitmiss = missPrefix
		}
		/*
			tr := "<tr"
			if hitmiss == missPrefix {
				tr += ` class="miss">`
			} else {
				tr += ">"
			}
		*/
		fls[i] = FunctionLine{
			Missed:     hitmiss == missPrefix,
			LineNumber: lineno,
			Code:       html.EscapeString(strings.Replace(line, "\t", "        ", -1)),
		}
	}
	return fls
}

type ReportFunctionList []ReportFunction

func (l ReportFunctionList) Len() int {
	return len(l)
}

// TODO make sort method configurable?
func (l ReportFunctionList) Less(i, j int) bool {
	var left, right float64
	if len(l[i].Statements) > 0 {
		left = float64(l[i].StatementsReached) / float64(len(l[i].Statements))
	}
	if len(l[j].Statements) > 0 {
		right = float64(l[j].StatementsReached) / float64(len(l[j].Statements))
	}
	if left < right {
		return true
	}
	return left == right && len(l[i].Statements) < len(l[j].Statements)
}

func (l ReportFunctionList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
