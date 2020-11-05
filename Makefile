build:
	go build -o bin ./cmd/main.go

test:
	go test -race ./cmd/... ./internal/...

lint: install-lint-deps
	golangci-lint run ./previewer/... ./internal/...

install-lint-deps:
	rm -rf $(shell go env GOPATH)/bin/golangci-lint
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.30.0
.PHONY: build test lint