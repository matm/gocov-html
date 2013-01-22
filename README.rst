Gocov HTML export
=================

This is a simple helper tool for generating HTML output from `axw/gocov`_.

.. _axw/gocov: https://github.com/axw/gocov

Installation
------------

Just type the following to install the program and its dependencies::

    $ go get github.com/axw/gocov/gocov
    $ go get github.com/matm/gocov-html

Usage
-----

`gocov-html` can read a JSON file or read from the standard input::

    $ gocov test net/http | gocov-html > http.html

or::

    $ gocov test net/http > http.json
    $ gocov-html http.json > http.html
