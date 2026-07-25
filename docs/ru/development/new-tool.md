# Добавление нового MCP-инструмента

## Шаги

1. **Добавьте константу имени инструмента** в `internal/service/service.go`
2. **Создайте типы запроса/ответа** в `internal/service/types.go`
3. **Реализуйте сервис** в `internal/service/` (новый файл или добавьте в существующий)
4. **Создайте markdown-определение** в `internal/service/definitions/` — это то, что читает `MakeToolDefinitions`
5. **Добавьте метод в интерфейс `Svc`** в `internal/server/mcp/handler.go`
6. **Добавьте обработчик** в `handler.go`
7. **Зарегистрируйте инструмент** в `registerTools` в `mcp.go`
8. **Сгенерируйте моки**: `go generate ./...`
9. **Напишите тесты**

## 1. Константа имени инструмента

Добавьте константу в `internal/service/service.go`:

```go
const MyNewTool = "my_new_tool"
```

## 2. Типы запроса/ответа

Определите в `internal/service/types.go`:

```go
type MyNewToolRequest struct {
    Param1 string `json:"param1" validate:"required" jsonschema:"required,Description of param1"`
}

type MyNewToolResponse struct {
    Result string `json:"result"`
}
```

## 3. Реализация сервиса

Создайте `internal/service/my_new_tool.go` или добавьте в существующий файл сервиса. Следуйте стандартному шаблону сервиса: валидация → поиск → выполнение → возврат:

```go
func (s *Service) MyNewTool(ctx context.Context, req MyNewToolRequest) (MyNewToolResponse, error) {
    if err := s.validateRequest(req); err != nil {
        return MyNewToolResponse{}, NewLLMError(validationFailedErrCode, err.Error())
    }
    // бизнес-логика
    return MyNewToolResponse{Result: "ok"}, nil
}
```

## 4. Markdown-определение

Создайте `internal/service/definitions/my_new_tool.md`. Этот файл читается `MakeToolDefinitions()` и встраивается в бинарник. Поле `name:` в frontmatter должно совпадать с константой:

```markdown
---
name: my_new_tool
---

# my_new_tool

Описание инструмента.

## Параметры

| Параметр | Тип | Описание |
|----------|-----|----------|
| `param1` | string | Описание |
```

Функция `MakeToolDefinitions()` в `tools.go` читает все `.md` файлы из встроенной директории `definitions/`, парсит YAML frontmatter для поля `name` и использует тело как описание инструмента. Файл `instruction.md` обрабатывается особым образом — он становится системной инструкцией для LLM.

## 5. Интерфейс Svc

Добавьте метод в составной интерфейс `Svc` в `handler.go`:

```go
type Svc interface {
    // ... существующие методы
    MyNewTool(ctx context.Context, req service.MyNewToolRequest) (service.MyNewToolResponse, error)
}
```

## 6. Обработчик

Добавьте метод обработчика на `handler` в `handler.go`. Обработчик делегирует сервису и оборачивает результат в `StructuredContent`:

```go
func (h *handler) handleMyNewTool(
    ctx context.Context,
    _ *sdkmcp.CallToolRequest,
    req service.MyNewToolRequest,
) (*sdkmcp.CallToolResult, any, error) {
    resp, err := h.service.MyNewTool(ctx, req)
    if err != nil {
        return nil, nil, err
    }
    return &sdkmcp.CallToolResult{
        StructuredContent: resp,
    }, nil, nil
}
```

## 7. Регистрация

Зарегистрируйте инструмент в функции `registerTools` в `mcp.go`. Добавьте запись в карту `toolRegistrations`:

```go
service.MyNewTool: {
    addTool[service.MyNewToolRequest](mcpServer, h.handleMyNewTool),
    true, // false, если инструмент изменяемый (как invoke или auth)
},
```

Сигнатура функции `registerTools`:

```go
func registerTools(mcpServer *sdkmcp.Server, tools []service.Tool, h handler) {
```

Она итерирует определения инструментов, возвращённые `MakeToolDefinitions()`, и регистрирует каждый с его типизированным обработчиком. Карта `toolRegistrations` соединяет константы имён инструментов с их обработчиками.
