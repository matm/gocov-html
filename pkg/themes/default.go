package themes

import (
	"text/template"
	"time"

	"github.com/matm/gocov-html/pkg/types"
)

type defaultTheme struct{}

func (t defaultTheme) Data() *types.TemplateData {
	css := `<style type="text/css">
    body {
        background-color: #fff;
        font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
    }
    table {
        margin-left: 10px;
        border-collapse: collapse;
    }
    td {
        background-color: #fff;
        padding: 2px;
    }
    table.overview td {
        padding-right: 20px;
    }
    td.percent, td.linecount { text-align: right; }
    div.package, #totalcov {
        color: #fff;
        background-color: #375eab;
        font-size: 16px;
        font-weight: bold;
        padding: 10px;
        border-radius: 5px 5px 5px 5px;
    }
    div.package, #totalcov {
        float: right;
        right: 10px;
    }
    #totalcov {
        top: 10px;
        position: relative;
        background-color: #fff;
        color: #000;
        border: 1px solid #375eab;
        clear: both;
    }
    #summaryWrapper {
        position: fixed;
        top: 10px;
        float: right;
        right: 10px;

    }
    span.packageTotal {
        float: right;
        color: #000;
    }
    #doctitle {
        background-color: #fff;
        font-size: 24px;
        margin-top: 20px;
        margin-left: 10px;
        color: #375eab;
        font-weight: bold;
    }
    #about {
        margin-left: 18px;
        font-size: 10px;
    }
    table tr:last-child td {
        font-weight: bold;
    }
    .functitle, .funcname {
        text-align: center;
        font-size: 20px;
        font-weight: bold;
        color: #375eab;
    }
    .funcname {
        text-align: left;
        margin-top: 20px;
        margin-left: 10px;
        margin-bottom: 20px;
        padding: 2px 5px 5px;
        background-color: #e0ebf5;
    }
    table.listing {
        margin-left: 10px;
    }
    table.listing td {
        padding: 0px;
        font-size: 12px;
        background-color: #eee;
        vertical-align: top;
        padding-left: 10px;
        border-bottom: 1px solid #fff;
    }
    table.listing td:first-child {
        text-align: right;
        font-weight: bold;
        vertical-align: center;
    }
    table.listing tr.miss td {
        background-color: #FFBBB8;
    }
    table.listing tr:last-child td {
        font-weight: normal;
        color: #000;
    }
    table.listing tr:last-child td:first-child {
        font-weight: bold;
    }
    .info {
        margin-left: 10px;
    }
    .info code {
    }
    pre { margin: 1px; }
    pre.cmd {
        background-color: #e9e9e9;
        border-radius: 5px 5px 5px 5px;
        padding: 10px;
        margin: 20px;
        line-height: 18px;
        font-size: 14px;
    }
    a {
        text-decoration: none;
        color: #375eab;
    }
    a:hover { text-decoration: underline; }
    p { margin-left: 10px; }
</style>`
	return &types.TemplateData{
		CSS:        css,
		When:       time.Now().Format(time.RFC822Z),
		ProjectURL: types.ProjectURL,
	}
}

func (t defaultTheme) Name() string {
	return "golang"
}

func (t defaultTheme) Description() string {
	return "original golang theme (default)"
}

func (t defaultTheme) Template() *template.Template {
	tmpl := `{{define "theme"}}
<html>
	<head>
		<title>Coverage Report</title>
		<meta charset="utf-8" />
		{{.CSS}}
	</head>
	<body>
		<div id="doctitle">Coverage Report</div>
        {{if not .Packages}}
		<p>no test files in package.</p>"
        {{else}}
        <div id="about">Generated on {{.When}} with <a href="{{.ProjectURL}}">gocov-html</a></div>
        {{/* Report overview/summary available? */}}
        {{if .Overview}}
        <div class="funcname">Report Overview</div>
            <table class="overview">
            {{range $k,$rp := .Packages}}
            <tr id="s_pkg_{{$rp.Pkg.Name}}">
                <td><code><a href="#pkg_{{$rp.Pkg.Name}}">{{$rp.Pkg.Name}}</a></code></td>
                <td class="percent"><code>{{printf "%.2f%%" $rp.PercentageReached}}</code></td>
                <td class="linecount"><code>{{printf "%d" $rp.ReachedStatements}}/{{printf "%d" $rp.TotalStatements}}</code></td>
            </tr>
            {{end}}
            </table>
        </div>
        {{end}}
        {{range $k,$rp := .Packages}}
        <div id="pkg_{{$rp.Pkg.Name}}" class="funcname">
            Package Overview: {{$rp.Pkg.Name}}
            <span class="packageTotal">{{printf "%.2f%%" $rp.PercentageReached}}</span>
        </div>
        <p>
        This is a coverage report created after analysis of the <code>{{$rp.Pkg.Name}}</code> package. It 
        has been generated with the following command:
        </p>
        <pre class="cmd">gocov test {{$rp.Pkg.Name}} | gocov-html</pre>
        <p>Here are the stats. Please select a function name to view its implementation and see what's left for testing.</p>

        <table class="overview">
        {{range $k,$f := $rp.Functions}}
            <tr id="s_fn_{{$f.Name}}">
                <td>
                    <code><a href="#fn_{{$f.Name}}">{{$f.Name}}(...)</a></code>
                </td>
                <td>
                    <code>{{$rp.Pkg.Name}}/{{$f.ShortFileName}}</code>
                </td>
                <td class="percent">
                    <code>{{printf "%.2f%%" $f.CoveragePercent}}</code>
                </td>
                <td class="linecount">
                    <code>{{$f.StatementsReached}}/{{len $f.Statements}}</code>
                </td>
            </tr>
        {{end}}
        </table>

        {{/* Functions source code here */}}
        {{range $k,$f := $rp.Functions}}
        <div class="funcname" id="fn_{{$f.Name}}">func {{$f.Name}}</div>
        <div class="info">
            <a href="#s_fn_{{$f.Name}}">Back</a>
            <p>In <code>{{$f.File}}</code>:</p>
        </div>
        <table class="listing">
            {{range $p,$info := $f.Lines}}
            <tr{{if $info.Missed}} class="miss"{{end}}>
                <td>{{$info.LineNumber}}</td>
                <td>
                    <code><pre>{{$info.Code}}</pre></code>
                </td>
            </tr>
            {{end}}
        </table>
        {{end}} {{/* range function lines */}}

        <!--    Can be parsed by external script
                PACKAGE:{{$rp.Pkg.Name}} DONE:{{printf "%.2f" $rp.PercentageReached}}
        -->
        {{end}} {{/* range Packages end */}}

        <div id="summaryWrapper">
        {{if not .Overview}}
            {{$rp := index .Packages 0}}
            <div class="package">{{$rp.Pkg.Name}}</div>
            <div id="totalcov">{{printf "%.2f%%" $rp.PercentageReached}}</div>
        {{else}}
            <div class="package">{{.Overview.Pkg.Name}}</div>
            <div id="totalcov">{{printf "%.2f%%" .Overview.PercentageReached}}</div>
        {{end}} {{/* if overview end */}}
        </div>
        {{end}} {{/* range if end */}}
	</body>
</html>
{{end}}`
	p := template.Must(template.New("theme").Parse(tmpl))
	return p
}
