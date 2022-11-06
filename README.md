# Gocov HTML export

This is a simple helper tool for generating HTML output from [axw/gocov](https://github.com/axw/gocov/)

Here is a screenshot:

![HTML coverage report screenshot](https://github.com/matm/gocov-html/blob/master/gocovh-html.png)

## Installation

Just type the following to install the program and its dependencies:
```
$ go install github.com/axw/gocov/gocov@latest
$ go install github.com/matm/gocov-html/cmd/gocov-html@latest
```

## Features Matrix

Feature|CLI Flag|Version
:---|:---|---:
Use custom CSS file|`-s <filename>`|`1.0`
Show program version|`-v`|`1.1.1`
Write CSS of default theme to stdout|`-d`|`1.2`
Embbed custom CSS into final HTML document|-|`1.2`
List available themes|`-lt`|`1.2`
Render with a specific theme|`-t <theme>`|`1.2`

## Usage

```bash
$ gocov-html -h
Usage of ./gocov-html:
  -d    output CSS of default theme
  -s string
        path to custom CSS file
  -v    show program version
```

`gocov-html` can read a JSON file or read from standard input:
```
$ gocov test strings | gocov-html > strings.html
ok      strings 0.700s  coverage: 98.1% of statements
```

The generated HTML content comes along with a default embedded CSS. However a custom stylesheet can be used with the `-s` flag:
```
$ gocov test net/http | gocov-html -s mystyle.css > http.html
```

As of version 1.2:
- A `-d` flag is available to write the defaut stylesheet to the standard output. This is provided for convenience and easy editing:
  ```
  $ gocov-html -d > newstyle.css
  ... edit newstyle.css ...
  $ gocov test strings | gocov-html -s newstyle.css > http.html
  ```
- The content of the stylesheet given to `-s` is embedded into the final HTML document
