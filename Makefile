SHELL := /bin/bash

GIT_TAG:=$(shell git tag --contains HEAD|awk 'END {print}')
GIT_REV:=$(shell git rev-parse --short=8 HEAD)
VERSION:=$(if $(GIT_TAG),$(GIT_TAG),$(GIT_REV))

SRC_FILES:=$(shell find . -type f -name '*.go')
TEST_PKGS:=$(shell go list ./... | grep -v /test | grep -v /docs | grep -v /mock_gen)
GO_BUILD_FLAGS := $(BUILD_ARGS) -ldflags '-X "gitlab.szzhijing.com/opensource/vega2-go/utils.REVISION=$(GIT_REV)" -X "gitlab.szzhijing.com/opensource/vega2-go/utils.VERSION=$(VERSION)"'
TEST_LOG_FILE := $(shell echo "test_`date +'%Y%m%d%H%M%S'`.log")

.PHONY: all
all: fmt lint build test

.PHONY: lint
lint:
	@echo "[golangci-lint] Running golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.0
	@golangci-lint run ./...
	@go mod tidy
	@echo "`date +"%Y-%m-%d %H:%M:%S"` ------------------------------------[Done]"

.PHONY: test
test: lint
	@echo "[test] Running go test..."
	@go test $(TEST_PKGS) -race -cover

.PHONY: fmt
fmt:
	@echo "[gofmt] Replace interface{} to any..."
	@gofmt -w -r 'interface{} -> any' . 2>&1
	@echo "[goimports-reviser] Format go project..."
	@#go install golang.org/x/tools/cmd/goimports@latest
	@#goimports -local "gitlab.szzhijing.com" -w . 2>&1
	@go install github.com/incu6us/goimports-reviser/v3@v3.3.0
	@goimports-reviser -project-name "gitlab.szzhijing.com" ./... 2>&1
	@go mod tidy -compat=1.17
	@echo "`date +"%Y-%m-%d %H:%M:%S"` ------------------------------------[Done]"


.PHONY: generate
gen:
	go generate ./...

