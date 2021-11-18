#
# vim:ft=make
#

MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
.DELETE_ON_ERROR:
.ONESHELL:

GIT_REF := $(shell git rev-parse --short HEAD)
GIT_TAG := $(shell git name-rev --tags --name-only $(GIT_REF))

.PHONY: all
all: ./bin/go-stacks.darwin ./bin/go-stacks.linux

./bin/go-stacks.%: $(shell find ./ -name '*.go')
	GOOS=$* go build -o $@ -ldflags "-X github.com/mhristof/go-stacks/cmd.version=$(GIT_TAG)+$(GIT_REF)" main.go

.PHONY: fast-test
fast-test:  ## Run fast tests
	go test ./... -tags fast

.PHONY: test
test:	## Run all tests
	go test ./...

.PHONY: clean
clean:
	rm -rf bin/go-stacks.*

.PHONY: help
help:           ## Show this help.
	@grep '.*:.*##' Makefile | grep -v grep  | sort | sed 's/:.*## /:/g' | column -t -s:
