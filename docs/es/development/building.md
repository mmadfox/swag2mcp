# Compilación

## Requisitos

- Go 1.26+
- Make

## Comandos

```bash
# Compilar
make build

# Compilar con versión
make build VERSION=1.0.0

# Lint
make lint

# Pruebas
go test ./...

# Todas las pruebas
make testall
```

## GoReleaser

Para lanzamientos:

```bash
goreleaser release --snapshot --clean
```

## Plataformas

| Plataforma | Arquitectura |
|------------|-------------|
| Linux | amd64, arm64 |
| macOS | amd64, arm64 |
| Windows | amd64 |

## Lint

```bash
make lint
```

Usa `golangci-lint` con más de 80 linters. Configuración en `.golangci.yml`.
