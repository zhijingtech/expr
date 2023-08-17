.PHONY: all
all: fmt lint test

.PHONY: fmt
fmt:
	@echo "[gofmt] Replace interface{} to any..."
	@gofmt -w -r 'interface{} -> any' . 2>&1
	@echo "[goimports-reviser] Format go project..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@goimports -local "gitlab.szzhijing.com" -w . 2>&1
	@go install github.com/incu6us/goimports-reviser/v3@v3.3.0
	@goimports-reviser -project-name "gitlab.szzhijing.com" ./... 2>&1
	@go mod tidy -compat=1.17
	@echo "`date +"%Y-%m-%d %H:%M:%S"` ------------------------------------[Done]"

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
	@go test ./... -race -cover

.PHONY: gen
gen:
	go generate ./...