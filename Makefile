.PHONY: gocov-html build clean

include version.mk

BIN=gocov-html
MAIN_CMD=github.com/matm/${BIN}/cmd/${BIN}

all: build

build:
	@go build -ldflags "all=$(GO_LDFLAGS)" ${MAIN_CMD}

test:
	@go test ./...

clean:
	@rm -rf ${BIN}
