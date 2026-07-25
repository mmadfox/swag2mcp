# Code-Konventionen

## Go

- **Go 1.26+**
- **gofmt** / **gofumpt** / **goimports** / **gci**
- **120 Zeichen** pro Zeile
- **Guard-Klauseln** statt verschachtelter Ifs
- **Namensgebung**: `camelCase` für private, `PascalCase` für exportierte

## Fehler

Verwenden Sie `LLMError` für LLM-sichtbare Fehler:

```go
type LLMError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

Fehlercodes:
- `validation_failed` — ungültige Parameter
- `not_found` — Ressource nicht gefunden
- `rate_limit` — Ratenlimit überschritten
- `invoke_error` — API-Aufruffehler

## Interfaces

- Kleine Interfaces (1-3 Methoden)
- Interface-Komposition
- Funktionale Optionen für die Konfiguration

## Tests

- Tabellengesteuerte Tests
- Test-Helfer (`newTestService()`, `seedTestData()`)
- Mocks über `go.uber.org/mock`
- 80%+ Abdeckung für Kernpakete

## Konfiguration

- YAML-Format
- Kaskade: global → spec → collection
- Validierung über `go-playground/validator`
- Umgebungsvariablen über `$(VAR)`
