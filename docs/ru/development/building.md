# Сборка

## Требования

- Go 1.26+
- Make

## Команды

```bash
# Сборка
make build

# Сборка с версией
make build VERSION=1.0.0

# Линтинг
make lint

# Тесты
go test ./...

# Все тесты
make testall
```

## GoReleaser

Для релизов:

```bash
goreleaser release --snapshot --clean
```

## Платформы

| Платформа | Архитектура |
|-----------|-------------|
| Linux | amd64, arm64 |
| macOS | amd64, arm64 |
| Windows | amd64 |

## Линтинг

```bash
make lint
```

Использует `golangci-lint` с 80+ линтерами. Конфиг в `.golangci.yml`.
