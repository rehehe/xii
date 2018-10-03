###
PWD			= $(shell pwd)
BUILD_DIR	= ${PWD}/build
PLATFORM	= $(shell go env GOOS)

DASHBOARD 	= dashboard
SURVEY 		= survey

###
default: build

build:
	@echo "making:"
	@$(MAKE) --no-print-directory -C dashboard
	@$(MAKE) --no-print-directory -C survey

run-dashboard:
	@$(MAKE) --no-print-directory -C dashboard run

run-survey:
	@$(MAKE) --no-print-directory -C survey run
###
.PHONY: default