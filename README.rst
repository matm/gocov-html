Gocov HTML export
=================

This is a simple helper tool for generating HTML output from `axw/gocov`_

.. _axw/gocov: https://github.com/axw/gocov  

and it forked from `matm/gocov-html`_.

.. _matm/gocov-html: https://github.com/matm/gocov-html 



Installation
------------

Just type the following to install the program and its dependencies::

    $ go get github.com/axw/gocov/gocov
    $ go get github.com/XingyanLee/gocov-html

Usage
-----


first step: Set environment variable `GO_COV_HTML_REORDER`::

    $ export GO_COV_HTML_REORDER=1
    
`gocov-html` can read a JSON file or read from standard input::

    $ gocov convert c.out | gocov-html > http.html

or::

    $ gocov test net/http > http.json
    $ gocov-html http.json > http.html

if you want to use former sorting algorithm,just set `GO_COV_HTML_REORDER` to null or do not set environment variable `GO_COV_HTML_REORDER`.


if you want to compare current and previous code coverage file, the former must be a json file and need be written after the filed '-diff' . In addition,you need set an environment variable `GO_COV_HTML_REORDER_TOPN` whose value is up to you. Its  default value is 10. These two command will show `GO_COV_HTML_REORDER_TOPN` functions whose coverage is less than the former file::
    
    $ export GO_COV_HTML_REORDER_TOPN=10
    $ gocov convert c.out | gocov-html -diff formerfilename.json> http.html
    
    
if you want to see functions whose coverage is less than a certain value in current code coverage file, you just need to set a environment variable `GO_COV_HTML_REORDER_NEWFUNLIMIT` whose value is up to you. this command will highlight the functions whose coverage is less than `GO_COV_HTML_REORDER_NEWFUNLIMIT`::
   
    $ export GO_COV_HTML_REORDER_NEWFUNLIMIT=90
    
    
The generated HTML content comes along with a default embedded CSS. Use the `-s` 
flag to use a custom stylesheet::

    $ gocov test net/http | gocov-html -s mystyle.css > http.html
