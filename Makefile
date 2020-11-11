lint: install-lint-deps
	golangci-lint run ./previewer/... ./internal/...

unit-test-fast:
	go test -race -count 100 -timeout 30s -short ./internal/...

unit-test-slow:
	go test -race -timeout 150s -run Slow ./internal/...

integration-test:
	go test -v ./cmd/...

build:
	go build -o bin ./cmd/main.go


docker-build:
	sudo docker build -t previewer .

docker-run:
	sudo docker run -p 8080:8080 previewer

install-lint-deps:
	rm -rf $(shell go env GOPATH)/bin/golangci-lint
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin
.PHONY: build test lint