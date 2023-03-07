#
# Cross building for various operating systems (amd64)
#
WINDOWS_BIN=${BIN}.exe
DISTDIR=dist
BINDIR=${DISTDIR}/${BIN}
BUILDDIR=build
#
BUILD_VERSION=${BIN}-${VERSION}
BUILD_DARWIN_AMD64=${BUILD_VERSION}-darwin-amd64
BUILD_FREEBSD_AMD64=${BUILD_VERSION}-freebsd-amd64
BUILD_OPENBSD_AMD64=${BUILD_VERSION}-openbsd-amd64
BUILD_LINUX_AMD64=${BUILD_VERSION}-linux-amd64
BUILD_WINDOWS_AMD64=${BUILD_VERSION}-windows-amd64
#
GOBUILD64=GOARCH=amd64 go build -ldflags "all=$(GO_LDFLAGS)"

all: build

dist: cleardist buildall zip sourcearchive checksum

checksum:
	@for f in ${DISTDIR}/*; do \
		sha256sum $$f > $$f.sha256; \
		sed -i 's,${DISTDIR}/,,' $$f.sha256; \
	done

zip: linux freebsd openbsd darwin windows
	@rm -rf ${BINDIR}

linux:
	@cp ${BUILDDIR}/${BUILD_VERSION}-linux* ${BINDIR}/${BIN} && \
		(cd ${DISTDIR} && zip -r ${BUILD_LINUX_AMD64}.zip ${BIN})

darwin:
	@cp ${BUILDDIR}/${BUILD_VERSION}-darwin* ${BINDIR}/${BIN} && \
		(cd ${DISTDIR} && zip -r ${BUILD_DARWIN_AMD64}.zip ${BIN})

windows:
	@cp ${BUILDDIR}/${BUILD_VERSION}-windows* ${BINDIR}/${WINDOWS_BIN} && \
		(cd ${DISTDIR} && rm ${BIN}/${BIN} && zip -r ${BUILD_WINDOWS_AMD64}.zip ${BIN})

freebsd:
	@cp ${BUILDDIR}/${BUILD_VERSION}-freebsd* ${BINDIR}/${BIN} && \
		(cd ${DISTDIR} && zip -r ${BUILD_FREEBSD_AMD64}.zip ${BIN})

openbsd:
	@cp ${BUILDDIR}/${BUILD_VERSION}-openbsd* ${BINDIR}/${BIN} && \
		(cd ${DISTDIR} && zip -r ${BUILD_OPENBSD_AMD64}.zip ${BIN})

buildall:
	@echo ">>>>>> OpenBSD build <<<<<<<"
	@GOOS=openbsd ${GOBUILD64} -v -o ${BUILDDIR}/${BUILD_OPENBSD_AMD64} ${MAIN_CMD}
	@echo ">>>>>> FreeBSD build <<<<<<<"
	@GOOS=freebsd ${GOBUILD64} -v -o ${BUILDDIR}/${BUILD_FREEBSD_AMD64} ${MAIN_CMD}
	@echo ">>>>>> Linux build <<<<<<<"
	@GOOS=linux ${GOBUILD64} -v -o ${BUILDDIR}/${BUILD_LINUX_AMD64} ${MAIN_CMD}
	@echo ">>>>>> MacOSX build <<<<<<<"
	@GOOS=darwin ${GOBUILD64} -v -o ${BUILDDIR}/${BUILD_DARWIN_AMD64} ${MAIN_CMD}
	@echo ">>>>>> Windows build <<<<<<<"
	@GOOS=windows ${GOBUILD64} -v -o ${BUILDDIR}/${BUILD_WINDOWS_AMD64} ${MAIN_CMD}

sourcearchive:
	@git archive --format=zip -o ${DISTDIR}/${BUILD_VERSION}.zip ${VERSION}
	@echo ${DISTDIR}/${BUILD_VERSION}.zip
	@git archive -o ${DISTDIR}/${BUILD_VERSION}.tar ${VERSION}
	@gzip ${DISTDIR}/${BUILD_VERSION}.tar
	@echo ${DISTDIR}/${BUILD_VERSION}.tar.gz

cleardist:
	@rm -rf ${DISTDIR} && mkdir -p ${BINDIR} && mkdir -p ${BUILDDIR}

build:
	@go generate ./...
	@go build -ldflags "all=$(GO_LDFLAGS)" ${MAIN_CMD}

test:
	@go test ./...

clean:
	@find pkg -name \*_gen.go -delete
	@rm -rf ${BIN} ${BUILDDIR} ${DISTDIR}
