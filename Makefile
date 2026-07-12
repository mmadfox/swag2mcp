VERSION ?= dev

.PHONY: lint cover integration-tests build

build:
	go build -ldflags "-X github.com/mmadfox/swag2mcp/internal/commands.Version=$(VERSION)" -o swag2mcp ./cmd/swag2mcp

lint:
	golangci-lint run ./...

cover:
	go test ./... -coverprofile=coverage.out -covermode=atomic && \
	go tool cover -html=coverage.out -o coverage.html && \
	go tool cover -func=coverage.out | tail -1 && \
	rm -f coverage.out

integration-tests:
	go test -v -count=1 -timeout 600s ./tests/...