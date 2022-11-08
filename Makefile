.PHONY: gocov-html build dist linux darwin windows buildall cleardist clean


BIN=gocov-html
MAIN_CMD=github.com/matm/${BIN}/cmd/${BIN}

include version.mk
include build.mk
