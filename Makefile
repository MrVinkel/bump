.DEFAULT_GOAL := build
SHELL := /bin/bash

ifdef RELEASE_TAG
VERSION := $(RELEASE_TAG)
else
VERSION := dev
endif
VERSION_PACKAGE := github.com/mrvinkel/bump/cmd/bump
ARCH:=amd64 386
OS:=linux windows

.PHONY: setup
setup: ## Install dependencies
	@go install gotest.tools/gotestsum@latest
	@go install github.com/boumenot/gocover-cobertura@latest
	@go mod download

.PHONY: build
build: ## Build application for current OS/ARCH
	@$(eval VERSIONFLAGS=-X '$(VERSION_PACKAGE).BumpVersion=$(VERSION)')
	@go build -o ./bin/bump -ldflags="$(VERSIONFLAGS)" ./cmd

.PHONY: all
all:  ## Build for all OS/ARCHS

define build-os-arch
.PHONY: build-$(1)-$(2)
build-$(1)-$(2):
	@echo Building bump-$(1)-$(2) $(VERSION)
	@$(eval VERSIONFLAGS=-X '$(VERSION_PACKAGE).BumpVersion=$(VERSION)')
	@CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build -o ./bin/bump-$(1)-$(2)  -ldflags="-w -s $(VERSIONFLAGS)" ./cmd
all: build-$(1)-$(2)
endef
$(foreach o,$(OS), $(foreach a,$(ARCH), $(eval $(call build-os-arch,$(o),$(a)))))

.PHONY: test
test: build ## Run tests
	@gotestsum --format pkgname -- -coverprofile=bin/cobertura-coverage.txt -covermode count ./...
	@gocover-cobertura < bin/cobertura-coverage.txt > bin/cobertura-coverage.xml

.PHONY: lint
lint: build ## Lint code
	@golangci-lint run

.PHONY: tidy
tidy: ## go mod tidy
	@go mod tidy

.PHONY: clean
clean: ## Clean bin
	@go clean -testcache
	@rm -rf bin

PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
