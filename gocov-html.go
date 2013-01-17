package main

// gocov-html
// https://github.com/axw/gocov/blob/master/gocov/annotate.go

import (
    "encoding/json"
    "flag"
    "fmt"
    "github.com/axw/gocov"
    "io/ioutil"
    "log"
    "os"
)

func unmarshalJson(data []byte) (packages []*gocov.Package, err error) {
    result := &struct{ Packages []*gocov.Package }{}
    err = json.Unmarshal(data, result)
    if err == nil {
        packages = result.Packages
    }
    return
}

func main() {
    flag.Parse()
    if flag.NArg() != 1 {
        fmt.Fprintf(os.Stderr, fmt.Sprintf("Usage: %s data.json\n", os.Args[0]))
        os.Exit(2)
    }

    data, err := ioutil.ReadFile(flag.Arg(0))
    if err != nil {
        log.Fatal(err.Error())
    }

    packages, err := unmarshalJson(data)
    if err != nil {
        log.Fatal(err.Error())
    }

    for _, p := range packages {
        fmt.Println(p.Name)
    }
}
