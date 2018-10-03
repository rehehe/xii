###
PLATFORM	= $(shell go env GOOS)
GOARCH		= $(shell go env GOARCH)
GOPATH		= $(shell go env GOPATH)
GOBIN		= $(GOROOT)/bin/go
GOBUILD		= ${GOBIN} build
GOMODULE	= GO111MODULE=on
GOCGO		= CGO_ENABLED=0
GOSUFFIX	= -a -installsuffix cgo
GOVARS		= ${GOCGO} ${GOMODULE} GOARCH=${GOARCH}

PWD			= $(shell pwd)
BUILD_DIR	= ${PWD}/../build
CMD_DIR 	= cmd/${BINARY}/main.go

BINARY 		= survey

###
default: clean build

###
build: ${PLATFORM}

linux: dep
	@echo "building linux binary of ${BINARY}"
	GOOS=linux   ${GOVARS} ${GOBUILD} ${GOSUFFIX} -o ${BUILD_DIR}/${BINARY}-linux       ${CMD_DIR}

darwin: dep
	@echo "building darwin binary of ${BINARY}"
	GOOS=darwin  ${GOVARS} ${GOBUILD} ${GOSUFFIX} -o ${BUILD_DIR}/${BINARY}-darwin      ${CMD_DIR}

windows: dep
	@echo "building windows binary of ${BINARY}"
	GOOS=windows ${GOVARS} ${GOBUILD} ${GOSUFFIX} -o ${BUILD_DIR}/${BINARY}-windows.exe ${CMD_DIR}

###
clean:
	@echo "cleaning"
	@rm -rf ${BUILD_DIR}/${BINARY}*

###
dep:
	@echo "get dependencies of ${BINARY}"
	${GOMODULE} ${GOBIN} get -v

###
run:
	@echo "runing ${BINARY}:"
	${BUILD_DIR}/${BINARY}-${PLATFORM}

###
.PHONY: default build clean linux darwin windows dep