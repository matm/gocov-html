// Copyright (c) 2013 Mathias Monnerville
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

const (
	ProjectUrl = "https://github.com/matm/gocov-html"
	htmlHeader = `<html>
    <head>
        <title>Coverage Report</title>
        <!-- FIXME: Embedded style sheet -->
        <style type="text/css">
            body { background-color: #fff; }
            table {
                margin-left: auto;
                margin-right: auto;
            }
            td { 
                background-color: #eee; 
                padding: 5px;
            }
            td.percent, td.linecount { text-align: right; }
            div.package, #totalcov, #doctitle { 
                position: fixed;
                color: brown;
                background-color: #eee; 
                font-size: 24px;
                font-weight: bold;
                padding: 10px;
            }
            div.package, #totalcov { 
                float: right; 
                right: 10px;
            }
            #totalcov { 
                top: 60px;
                margin-top: 10px;
                color: green;
                clear: both;
            }
            #doctitle { 
                float: left;
                left: 10px;
                font-size: 20px;
                margin-top: 10px;
                color: black;
                text-align: center;
            }
            #about {
                position: fixed;
                float: left;
                top: 90px;
                left: 10px;
                font-size: 10px;
            }
            table tr:last-child td {
                font-weight: bold;
                color: brown;
            }
            .functitle, .funcname { 
                text-align: center; 
                font-size: 30px; 
                color: brown; 
            }
            .funcname {
                text-align: left;
                margin-left: 140px;
                margin-bottom: 20px;
                padding-top: 20px;
                padding-bottom: 5px;
                border-bottom: 3px solid #ccc;
            }
            table.listing {
                margin-left: 140px;
                border-collapse: collapse;
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
                margin-left: 140px;
            }
            .info code {
                font-weight: bold;
            }
        </style>
    </head>
    <body>
        <div id="doctitle">Coverage<br/>Report</div>
    `

	htmlFooter = `
    </body>
</html>`
)
