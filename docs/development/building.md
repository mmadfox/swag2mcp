# Building

## Requirements

- Go 1.26+
- Make

## Commands

```bash
# Build
make build

# Build with version
make build VERSION=1.0.0

# Lint
make lint

# Tests
go test ./...

# All tests
make testall
```

## GoReleaser

For releases:

```bash
goreleaser release --snapshot --clean
```

## Platforms

| Platform | Architecture |
|----------|-------------|
| Linux | amd64, arm64 |
| macOS | amd64, arm64 |
| Windows | amd64 |

## Lint

```bash
make lint
```

Uses `golangci-lint` with 80+ linters. Config in `.golangci.yml`.
