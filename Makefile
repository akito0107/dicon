NAME := dicon
VERSION := 0.0.1
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' -X 'main.revision=$(REVISION)'
PACKAGENAME := github.com/akito0107/dicon

.PHONY: setup dep test test/internal main clean install lint lint/internal

all: main

main:
	go build -ldflags "$(LDFLAGS)" -o bin/dicon main.go

## Install dependencies
setup:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/golang/lint/golint

## install go dependencies
dep:
	dep ensure

test: test/internal 

install:
	go install

test/internal:
	go test -v -cover -race $(PACKAGENAME)/internal

lint: lint/main lint/internal

lint/main:
	golint main.go

lint/internal:
	golint internal

## remove build files
clean:
	rm -rf ./bin/*

