build:
	go build -o bin ./previewer/main.go

test:
	go test -race ./previewer/...

lint: install-lint-deps
	golangci-lint run ./previewer/...

install-lint-deps:
	rm -rf $(shell go env GOPATH)/bin/golangci-lint
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.30.0
.PHONY: build test lint