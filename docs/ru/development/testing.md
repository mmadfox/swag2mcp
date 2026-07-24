# Тестирование

## Команды

```bash
# Модульные тесты
go test ./...

# Конкретный пакет
go test ./internal/service/...

# Интеграционные тесты
make integration-tests

# Покрытие
make cover

# Все тесты
make testall
```

## Структура тестов

```
tests/
├── main_test.go              # Точка входа
├── suite_test.go             # Настройка набора тестов
├── suite_auth_test.go        # Тесты аутентификации
├── suite_config_test.go      # Тесты конфигурации
├── suite_mcp_tools_test.go   # Тесты MCP-инструментов
├── suite_search_test.go      # Тесты поиска
├── suite_ratelimit_test.go   # Тесты ограничения частоты
├── suite_response_test.go    # Тесты ответов
├── suite_export_test.go      # Тесты экспорта
├── suite_import_test.go      # Тесты импорта
├── suite_parsing_test.go     # Тесты парсинга
├── suite_transport_test.go   # Тесты транспорта
├── suite_mock_test.go        # Тесты мок-сервера
├── suite_workspace_test.go   # Тесты рабочей области
├── suite_errors_test.go      # Тесты ошибок
└── suite_version_test.go     # Тесты версии
```

## Покрытие

Цель: 80%+ для основных пакетов:

- `auth`
- `cache`
- `config`
- `env`
- `httpclient`
- `id`
- `index`
- `server/mcp`
- `service`
- `spec`
- `workspace`

## Моки

Использует `go.uber.org/mock` для тестов MCP-сервера:

```bash
go generate ./...
```

Генерирует `internal/server/mcp/mock_svc_test.go` из `handler.go`.

## Table-Driven тесты

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "hello", "HELLO", false},
        {"empty input", "", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := DoSomething(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.Equal(t, tt.want, got)
        })
    }
}
```
