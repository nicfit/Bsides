SHELL := /usr/bin/env bash
.DEFAULT_GOAL := build

GOOS ?=
GOARCH ?=

all: clean build test

build:
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/Bsides cmd/*.go
.PHONY: all build

test:
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go test -coverprofile=coverage.out -count=1 ./...
test-coverage: test
	@go tool cover -html=coverage.out
.PHONY: test test-coverage

clean:
	@rm -f bin/Bsides coverage.out *.log
clean-dist: clean
	@rm -rf bin
.PHONY: clean clean-dist

update-deps:
	go get -u ./...
	go mod tidy
.PHONY: update-deps