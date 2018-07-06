NAME := dicon/examples/helloworld
VERSION := $(shell git tag -l | tail -1)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' -X 'main.revision=$(REVISION)'
PACKAGENAME := github.com/akito0107/dicon/examples/helloworld

.PHONY: setup dep test main clean install lint lint/internal

all: main

main:
	go build -ldflags "$(LDFLAGS)"

## Install dependencies
setup:
	go get -u github.com/akito0107/dicon/cmd/dicon

## generate DI Container
dicon:
	dicon generate --pkg di,component
