// Copyright (c) 2013-2022 Mathias Monnerville
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

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/matm/gocov-html/pkg/config"
	"github.com/matm/gocov-html/pkg/cov"
	"github.com/matm/gocov-html/pkg/themes"
)

func main() {
	var r io.Reader
	log.SetFlags(0)

	css := flag.String("s", "", "path to custom CSS file")
	showVersion := flag.Bool("v", false, "show program version")
	showDefaultCSS := flag.Bool("d", false, "output CSS of default theme")
	listThemes := flag.Bool("lt", false, "list available themes")
	theme := flag.String("t", "golang", "theme to use for rendering")

	flag.Parse()

	if *showVersion {
		fmt.Printf("Version:      %s\n", config.Version)
		fmt.Printf("Git revision: %s\n", config.GitRev)
		fmt.Printf("Git branch:   %s\n", config.GitBranch)
		fmt.Printf("Go version:   %s\n", runtime.Version())
		fmt.Printf("Built:        %s\n", config.BuildDate)
		fmt.Printf("OS/Arch:      %s/%s\n", runtime.GOOS, runtime.GOARCH)
		return
	}

	if *listThemes {
		for _, th := range themes.List() {
			fmt.Printf("%-10s -- %s\n", th.Name(), th.Description())
		}
		return
	}

	err := themes.Use(*theme)
	if err != nil {
		log.Fatalf("theme selection: %v", err)
	}

	if *showDefaultCSS {
		fmt.Println(themes.Current().Data().CSS)
		return
	}

	switch flag.NArg() {
	case 0:
		r = os.Stdin
	case 1:
		var err error
		if r, err = os.Open(flag.Arg(0)); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Usage: %s data.json\n", os.Args[0])
	}

	if err := cov.HTMLReportCoverage(r, *css); err != nil {
		log.Fatal(err)
	}
}
