SHELL := /bin/sh

MAIN_PKG ?= .
ENTRYPOINT ?= ./main.go

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
LDFLAGS ?= -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)

.DEFAULT_GOAL := help

.PHONY: help clean build test run install all release

help: ## Show available targets.
	@awk 'BEGIN {FS = ":.*## "}; /^[a-zA-Z0-9_.-]+:.*## / {printf "%-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

clean: ## Clean Go build/test caches.
	@go clean

build: ## Build the project binary.
	@go build -ldflags "$(LDFLAGS)" $(MAIN_PKG)

test: ## Run all tests with race detection and coverage output.
	@go test -race -covermode=atomic -coverprofile=coverage.out ./...

run: ## Run the main entrypoint. Pass args via: make run cmd="--help"
	@go run -ldflags "$(LDFLAGS)" $(ENTRYPOINT) $(cmd)

install: ## Install the project with linker metadata.
	@go install -ldflags "$(LDFLAGS)" $(MAIN_PKG)

all: clean install ## Clean and install.

release: ## Create a snapshot release with goreleaser.
	@goreleaser release --snapshot
