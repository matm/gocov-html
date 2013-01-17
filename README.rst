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

First generate coverage data with axw's `gocov`, then use it with `gocov-html`::

    $ gocov test pkg > coverage.json
    $ gocov-html coverage.json
