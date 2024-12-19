# Gocov HTML export

[![go-recipes](https://raw.githubusercontent.com/nikolaydubina/go-recipes/main/badge.svg?raw=true)](https://github.com/nikolaydubina/go-recipes)

This is a simple helper tool for generating HTML output from [axw/gocov](https://github.com/axw/gocov/).

`gocov-html` has support for themes, you might want to have a [look at the screenshots](themes/README.md).

## Installation

Binaries for most amd64 systems are built for every release. Please just [grab a binary version of the latest release](https://github.com/matm/gocov-html/releases).

You can also build it from source. In this case, a working Go 1.11+ compiler is required:

```bash
$ go install github.com/matm/gocov-html/cmd/gocov-html@latest
```

A [Dockerfile](Dockerfile) is also provided.

## Features Matrix

Feature|CLI Flag|Version
:---|:---|---:
Use custom CSS file|`-s <filename>`|`1.0.0`
Show program version|`-v`|`1.1.1`
Write CSS of default theme to stdout|`-d`|`1.2.0`
Embbed custom CSS into final HTML document|-|`1.2.0`
List available themes|`-lt`|`1.2.0`
Render with a specific theme|`-t <theme>`|`1.2.0`
New `kit` theme |`-t kit`|`1.3.0`
Put lower coverage functions on top|`-r`|`1.3.1`
Only show functions whose coverage is smaller than a max threshold|`-cmax`|`1.4.0`
Only show functions whose coverage is greater than a min threshold|`-cmin`|`1.4.0`

## Usage

```
Usage of gocov-html:
  -cmax uint
        only show functions whose coverage is less than cmax (default 100)
  -cmin uint
        only show functions whose coverage is more than cmin
  -d    output CSS of default theme
  -lt
        list available themes
  -r    put lower coverage functions on top
  -s string
        path to custom CSS file
  -t string
        theme to use for rendering (default "golang")
  -v    show program version
```

## Examples

Generate code coverage for the `strings` package then generate an HTML report:
```
$ gocov test strings | gocov-html > strings.html
ok      strings 0.700s  coverage: 98.1% of statements
```

Merge several coverage stats for different packages into a single report:
```
$ gocov test fmt math io | gocov-html > report.html
ok      fmt     0.045s  coverage: 95.2% of statements
ok      math    0.006s  coverage: 83.6% of statements
ok      io      0.024s  coverage: 88.2% of statements
```

In this case, the generated report will have an *overview* section with stats per package along with the global coverage percentage. This section may be rendered depending on the theme used. The `golang` (default) theme displays it.

List all available themes:
```
$ gocov-html -lt
golang     -- original golang theme (default)
kit        -- AdminKit theme
```

Generate a report using a specific theme with `-t`:
```
$ gocov test io | gocov-html -t kit > io.html
```

Only show functions whose code coverage is lower than 90% for the `strings` package:
```
$ gocov test strings|./gocov-html -cmax 90 > strings.html
```
In this example, only 5 matches are added to the report.

## Donate

If you like this tool and want to support its development, a donation would be greatly appreciated!

It's not about the amount at all: making a donation boosts the motivation to work on a project. Thank you very much if you can give anything.

Monero address:
`86S43wMDNPgNeUd6MkPEpiPUbBM6dS6DGdXBzc34uSw3Lxyg9p5tjmuGHESwmza3wGKfP2njUQdEd6kE3YPFRuaJFzP4Ger`

![My monero address](res/qr-donate.png)
