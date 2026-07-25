# Construction

## Prérequis

- Go 1.26+
- Make

## Commandes

```bash
# Construction
make build

# Construction avec version
make build VERSION=1.0.0

# Analyse
make lint

# Tests
go test ./...

# Tous les tests
make testall
```

## GoReleaser

Pour les versions :

```bash
goreleaser release --snapshot --clean
```

## Plateformes

| Plateforme | Architecture |
|----------|-------------|
| Linux | amd64, arm64 |
| macOS | amd64, arm64 |
| Windows | amd64 |

## Analyse

```bash
make lint
```

Utilise `golangci-lint` avec plus de 80 analyseurs. Configuration dans `.golangci.yml`.
