# Bauen

## Anforderungen

- Go 1.26+
- Make

## Befehle

```bash
# Bauen
make build

# Mit Version bauen
make build VERSION=1.0.0

# Lint
make lint

# Tests
go test ./...

# Alle Tests
make testall
```

## GoReleaser

Für Veröffentlichungen:

```bash
goreleaser release --snapshot --clean
```

## Plattformen

| Plattform | Architektur |
|-----------|-------------|
| Linux | amd64, arm64 |
| macOS | amd64, arm64 |
| Windows | amd64 |

## Lint

```bash
make lint
```

Verwendet `golangci-lint` mit 80+ Linters. Konfiguration in `.golangci.yml`.
