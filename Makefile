.PHONY: gocov-html build dist linux darwin windows buildall cleardist clean


BIN=gocov-html
MAIN_CMD=github.com/matm/${BIN}/cmd/${BIN}

GENERATOR_BIN=generator
GENERATOR_CMD=github.com/matm/${BIN}/cmd/${GENERATOR_BIN}

include version.mk
include build.mk
