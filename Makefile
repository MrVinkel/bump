.DEFAULT_GOAL := build
SHELL := /bin/bash

ifdef CI_COMMIT_TAG
VERSION := $(CI_COMMIT_TAG)
else
VERSION := dev
endif
VERSION_PACKAGE := goto/cmd/gt
ARCH:=amd64 386
OS:=linux windows

.PHONY: setup
setup: ## Install dependencies
	@go install gotest.tools/gotestsum@latest
	@go install github.com/boumenot/gocover-cobertura@latest
	@go mod download

.PHONY: build
build: ## Build application for current OS/ARCH
	@$(eval VERSIONFLAGS=-X '$(VERSION_PACKAGE).Version=$(VERSION)')
	@go build -o ./bin/bump -ldflags="$(VERSIONFLAGS)" ./cmd

.PHONY: build-all
build-all:  ## Build for all OS/ARCHS

define build-os-arch
.PHONY: build-$(1)-$(2)
build-$(1)-$(2):
	@echo Building bump-$(1)-$(2) $(VERSION)
	@$(eval VERSIONFLAGS=-X '$(VERSION_PACKAGE).Version=$(VERSION)')
	@CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build -o ./bin/bump-$(1)-$(2)  -ldflags="-w -s $(VERSIONFLAGS)" ./cmd
build-all: build-$(1)-$(2)
endef
$(foreach o,$(OS), $(foreach a,$(ARCH), $(eval $(call build-os-arch,$(o),$(a)))))

.PHONY: test
test: build ## Run tests
	@gotestsum --format pkgname -- -coverprofile=bin/cobertura-coverage.txt -covermode count ./...
	@gocover-cobertura < bin/cobertura-coverage.txt > bin/cobertura-coverage.xml

.PHONY: lint
lint: build ## Lint code
	@golangci-lint run


.PHONY: clean
clean: ## Clean bin
	@go clean -testcache
	@rm -rf bin

PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
