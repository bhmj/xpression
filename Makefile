LDFLAGS ?=-s -w -X main.appVersion=dev-$(shell git rev-parse --short HEAD)-$(shell date +%y-%m-%d)
OUT ?= ./bin
PROJECT ?=$(shell basename $(PWD))
SRC ?= ./cmd/$(PROJECT)
BINARY ?= $(OUT)/$(PROJECT)
PREFIX ?= manual

define USAGE

Usage: make <target>

some of the <targets> are:

  configure     - install tools and dependencies (gocyclo and golangci-lint)
  build         - build xpression evaluator CLI
  run           - run xpression evaluator CLI
  lint          - run linters
  test          - run tests
  cover         - generate coverage report

endef
export USAGE

all: configure build lint test

help:
	echo "$$USAGE"

configure:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

build:
	mkdir -p $(OUT)
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -trimpath -o $(BINARY) $(SRC)

run:
	mkdir -p $(OUT)
	CGO_ENABLED=0 go run -ldflags "$(LDFLAGS)" -trimpath $(SRC)

lint:
	echo "------ golangci-lint"
	golangci-lint run
	echo "------ gocyclo"
	gocyclo -over 18 .

test: 
	go test -cover ./...
	go test -benchmem -bench=.

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: all configure help build run lint test

$(V).SILENT: