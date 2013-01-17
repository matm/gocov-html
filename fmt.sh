#!/bin/sh

find . -type f -name \*.go -exec gofmt -tabs=false -tabwidth=4 -l -w -s {} \;
