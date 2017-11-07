NAME := dicon
VERSION := 0.0.1
REVISON := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' -X 'main.revision=$(REVISION)'
PACKAGENAME := github.com/akito0107/dicon

.PHONY: setup dep test/internal main clean install

all: main

main:
	go build -ldflags "$(LDFLAGS)" -o bin/dicon main.go

## Install dependencies
setup:
	go get -u github.com/golang/dep/cmd/dep

## install go dependencies
dep:
	dep ensure

test: test/internal 

install:
	go install

test/internal:
	go test -race $(PACKAGENAME)/internal

## remove build files
clean:
	rm -rf ./bin/*

