VERSION ?= dev

CORE_PACKAGES = ./internal/auth/... ./internal/cache/... ./internal/config/... \
                ./internal/env/... ./internal/httpclient/... ./internal/id/... \
                ./internal/index/... ./internal/server/mcp/... \
                ./internal/service/... ./internal/spec/... \
                ./internal/types/... ./internal/workspace/...

.PHONY: lint cover cover-core integration-tests build testall

build:
	go build -ldflags "-X github.com/mmadfox/swag2mcp/internal/commands.Version=$(VERSION)" -o swag2mcp ./cmd/swag2mcp

lint:
	golangci-lint run ./...

# Full coverage across all packages (including commands, tui, mockserver)
cover:
	go test ./... -coverprofile=coverage.out -covermode=atomic && \
	go tool cover -html=coverage.out -o coverage.html && \
	go tool cover -func=coverage.out | tail -1 && \
	rm -f coverage.out

# Coverage for core packages only — quick quality check
cover-core:
	go test $(CORE_PACKAGES) -coverprofile=coverage.out -covermode=atomic && \
	go tool cover -html=coverage.out -o coverage.html && \
	go tool cover -func=coverage.out | tail -1 && \
	rm -f coverage.out

integration-tests:
	go test -v -count=1 -timeout 600s ./tests/...

testall: lint integration-tests
	go test ./... -count=1