# Обзор разработки

## О проекте

swag2mcp — это Go-проект, который соединяет спецификации OpenAPI/Swagger/Postman с LLM-агентами через протокол Model Context Protocol (MCP). Он написан на Go 1.23+ и следует строгим соглашениям по кодированию, обеспечиваемым 80+ линтерами.

Этот раздел написан для **инженеров**, которые хотят понять кодовую базу, внести свой вклад или расширить swag2mcp новыми методами аутентификации, MCP-инструментами или интеграциями.

## Навыки разработки

Проект поставляется с двумя навыками разработки, которые кодируют соглашения и шаблоны проекта. Вы можете использовать их или игнорировать — это инструменты, а не правила.

### godeveloper

Навык [godeveloper](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/godeveloper/SKILL.md) определяет все соглашения по коду в проекте:

- **Именование** — пакеты, файлы, типы, интерфейсы, приёмники, константы
- **Форматирование** — gofmt/gofumpt/goimports/gci, лимит 120 строк, порядок импортов
- **Обработка ошибок** — `LLMError` с 8 кодами ошибок, sentinel-ошибки, обёртывание ошибок
- **Интерфейсы** — маленькие интерфейсы, композиция, определение на стороне потребителя
- **Конкурентность** — гранулярность мьютексов, время жизни горутин, передача контекста
- **Тестирование** — table-driven тесты, хелперы `newTestService()`/`seedTestData()`, генерация моков
- **Шаблоны проекта** — сервисный слой, структуры запросов/ответов, функциональные опции, шаблон MCP-обработчика

### swag2mcp-cli

Навык [swag2mcp-cli](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md) документирует каждую CLI-команду с синтаксисом, флагами, аргументами и примерами. Полезен при работе над CLI-командами или написании документации.

## Ключевые архитектурные решения

### Шаблон сервисного слоя

Каждая функция следует одному и тому же трёхшаговому шаблону:

1. **Валидация** запроса с помощью `s.validateRequest(req)` (использует `go-playground/validator`)
2. **Поиск** сущностей в in-memory индексе (возвращает `LLMError` с кодом `not_found`)
3. **Выполнение** бизнес-логики и возврат типизированного ответа или `LLMError`

```go
func (s *Service) Search(ctx context.Context, req SearchRequest) (SearchResponse, error) {
    if err := s.validateRequest(req); err != nil {
        return SearchResponse{}, NewLLMError(validationFailedErrCode, err.Error())
    }
    results, err := s.index.Search(req.Query, req.Limit)
    if err != nil {
        return SearchResponse{}, NewLLMError(invokeErrorCode, err.Error())
    }
    return SearchResponse{Results: results}, nil
}
```

### Структуры запросов/ответов

Каждый метод имеет выделенные структуры `{Method}Request` и `{Method}Response`. Структуры запросов используют теги `validate` для валидации и теги `jsonschema` для документации:

```go
type SearchRequest struct {
    Query string `json:"query" validate:"required,min=1" jsonschema:"description=Search query supporting field filters"`
    Limit int    `json:"limit" validate:"required,min=1,max=50" jsonschema:"description=Maximum results"`
}

type SearchResponse struct {
    Results []EndpointSearchItem `json:"results"`
}
```

### Функциональные опции

Конфигурация использует шаблон функциональных опций:

```go
type Option func(*Service)

func New(opts ...Option) (*Service, error)

func WithDisableLLMAuth(disable bool) Option {
    return func(s *Service) {
        s.disableLLMAuth.Store(disable)
    }
}
```

### Шаблон MCP-обработчика

MCP-сервер использует шаблон композиции интерфейсов. Интерфейс `Svc` в `internal/server/mcp/handler.go` составлен из меньших интерфейсов (`CatalogReader`, `EndpointExplorer`, `EndpointExecutor`, `SystemInfo`, `ResponseManager`). Каждый метод обработчика делегирует сервисному слою:

```go
type handler struct {
    service Svc
}

func (h *handler) handleSearch(ctx context.Context, _ *sdkmcp.CallToolRequest, req service.SearchRequest) (*sdkmcp.CallToolResult, any, error) {
    resp, err := h.service.Search(ctx, req)
    if err != nil {
        return nil, nil, err
    }
    return &sdkmcp.CallToolResult{StructuredContent: resp}, nil, nil
}
```

### LLMError

Все ошибки, возвращаемые LLM, используют тип `LLMError` с одним из 8 кодов:

| Код | Когда |
|-----|-------|
| `validation_failed` | Некорректный ввод (неверный формат ID, отсутствуют обязательные поля) |
| `not_found` | Сущность не найдена в индексе |
| `rate_limit` | Превышена 10-секундная задержка на эндпоинт |
| `invoke_error` | Ошибки HTTP-запроса/ответа |
| `config_error` | Ошибка загрузки или валидации конфигурации |
| `workspace_error` | Ошибка операции с директорией или файлом рабочей области |
| `parse_error` | Ошибка парсинга файла спецификации |
| `auth_error` | Ошибка получения токена аутентификации |

Сообщения должны объяснять, что пошло не так И что делать дальше, простым языком, подходящим для LLM-потребителя.

### Генерация ID

Все ID — детерминированные MD5-хеши:

```go
id.Domain("meteo")                          // 32-char hex
id.Collection("meteo", "Forecast")          // 32-char hex
id.Tag("meteo", "Forecast", "pets")         // 32-char hex
id.Method("meteo", "Forecast", "pets", "GET", "/v2/pet/{petId}")
```

### Каскад конфигурации

Конфигурация каскадируется через три уровня: **глобальный → спецификация → коллекция**. Каждый уровень переопределяет предыдущий. Все настройки `http_client` могут быть переопределены на каждом уровне. Заголовки и cookies объединяются; простые значения заменяются.

## Быстрая справка

| Область | Соглашение |
|---------|------------|
| **Версия Go** | 1.23+ |
| **Форматтеры** | gofmt, gofumpt, goimports, gci |
| **Длина строки** | 120 символов |
| **Линтеры** | 80+ в `.golangci.yml` |
| **Тип ошибки** | `LLMError` с 8 кодами |
| **Фреймворк моков** | `go.uber.org/mock` |
| **Хелперы тестов** | `newTestService()`, `seedTestData()` |
| **Формат конфига** | YAML с каскадом |
| **Диспетчеризация аутентификации** | `UnmarshalYAML` читает поле `type` |
| **Генерация ID** | На основе MD5 (`id.Domain()`, `id.Collection()` и т.д.) |
| **Лимит частоты** | 10 секунд на эндпоинт для `invoke` |
| **Размер ответа** | 1 МБ по умолчанию, сохраняется в файл при превышении |
| **Цель покрытия** | 80%+ для основных пакетов |
| **Сборка** | `make build` |
| **Линтинг** | `make lint` |
| **Тесты** | `go test ./...` |
| **Генерация** | `go generate ./...` |
