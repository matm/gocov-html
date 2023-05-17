# Gocov HTML export

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

## Usage

```
$ gocov-html -h
Usage of ./gocov-html:
  -d    output CSS of default theme
  -lt
        list available themes
  -s string
        path to custom CSS file
  -t string
        theme to use for rendering (default "golang")
  -v    show program version
```

`gocov-html` can read a JSON file or read from standard input:
```
$ gocov test strings | gocov-html > strings.html
ok      strings 0.700s  coverage: 98.1% of statements
```

Several packages can be tested at once and added to a single report. Let's test the `fmt`, `math` and `io` packages:
```
$ gocov test fmt math io | gocov-html > report.html
ok      fmt     0.045s  coverage: 95.2% of statements
ok      math    0.006s  coverage: 83.6% of statements
ok      io      0.024s  coverage: 88.2% of statements
```

In this case, the generated report will have an *overview* section with stats per package along with the global coverage percentage. This section may be rendered depending on the theme used. The `golang` (default) theme displays it.

The generated HTML content comes along with a default embedded CSS. However a custom stylesheet can be used with the `-s` flag:
```
$ gocov test net/http | gocov-html -s mystyle.css > http.html
```

As of version 1.2,
- A `-d` flag is available to write the defaut stylesheet to the standard output. This is provided for convenience and easy editing:
  ```
  $ gocov-html -d > newstyle.css
  ... edit newstyle.css ...
  $ gocov test strings | gocov-html -s newstyle.css > http.html
  ```
- The content of the stylesheet given to `-s` is embedded into the final HTML document
- Theming capabilities are available (to go further than just using a CSS file) through the use of Go templates.
  - Use the `-lt` flag to list available themes:
    ```
    $ gocov-html -lt
    golang     -- original golang theme (default)
    ```
  - Generate a report using a specific theme with `-t`:
    ```
    $ gocov test io | gocov-html -t golang > io.html
    ```

## Donate

If you like this tool and want to support its development, a donation would be greatly appreciated!

It's not about the amount at all: making a donation boosts the motivation to work on a project. Thank you very much if you can give anything.

Monero address:
`86S43wMDNPgNeUd6MkPEpiPUbBM6dS6DGdXBzc34uSw3Lxyg9p5tjmuGHESwmza3wGKfP2njUQdEd6kE3YPFRuaJFzP4Ger`

![My monero address](res/qr-donate.png)
