VERSION ?= dev

CORE_PACKAGES = ./internal/auth/... ./internal/cache/... ./internal/config/... \
                ./internal/env/... ./internal/httpclient/... ./internal/id/... \
                ./internal/index/... ./internal/server/mcp/... \
                ./internal/service/... ./internal/spec/... \
                ./internal/types/... ./internal/workspace/...

.PHONY: lint cover cover-core integration-tests build testall docs

build:
	go build -ldflags "-X github.com/mmadfox/swag2mcp/internal/commands.Version=$(VERSION)" -o swag2mcp ./cmd/swag2mcp

docs:
	npm install --silent && npx vitepress build docs

lint:
	golangci-lint run ./...

# Full coverage across all packages (including commands, tui, mockserver)
cover:
	COVER_PKGS=$$(go list ./... | grep -v -E '(/cmd|/tests$$|/tests/|/mocks/|internal/tui|internal/commands/init|internal/commands/add|internal/commands/delete|internal/commands/run)' | tr '\n' ',') ; \
	go test -count=1 \
	  -coverpkg=$$COVER_PKGS \
	  -coverprofile=coverage.out \
	  -covermode=atomic \
	  ./...

integration-tests:
	go test -v -count=1 -timeout 600s ./tests/...

testall: lint integration-tests
	go test ./... -count=1