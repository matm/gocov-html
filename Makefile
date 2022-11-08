.PHONY: gocov-html build dist linux darwin windows buildall cleardist clean

include version.mk

BIN=gocov-html
MAIN_CMD=github.com/matm/${BIN}/cmd/${BIN}

include build.mk
