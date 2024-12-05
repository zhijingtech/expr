.PHONY: all
all: fmt lint test

.PHONY: init
init:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/incu6us/goimports-reviser/v3@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: fmt
fmt:
	gofmt -w -r 'interface{} -> any' . 2>&1
	echo "[goimports-reviser] Format go project..."
	goimports -local "github.com/zhijingtech/expr" -w . 2>&1
	goimports-reviser -project-name "github.com/zhijingtech/expr" ./... 2>&1

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test: lint
	go test ./... -race -cover

.PHONY: gen
gen:
	go generate ./...