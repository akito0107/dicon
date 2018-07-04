NAME := dicon
VERSION := $(shell git tag -l | tail -1)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' -X 'main.revision=$(REVISION)'
PACKAGENAME := github.com/akito0107/dicon

.PHONY: setup dep test main clean install lint lint/internal

all: main

main:
	go build -ldflags "$(LDFLAGS)" -o bin/dicon cmd/dicon/main.go

## Install dependencies
setup:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/golang/lint/golint

## install go dependencies
dep:
	dep ensure

test:
	go test -v -cover -race $(PACKAGENAME)

install:
	go install $(PACKAGENAME)/cmd/dicon

lint: lint/main lint/internal

lint/main:
	golint main.go

lint/internal:
	golint internal

## remove build files
clean:
	rm -rf ./bin/*

