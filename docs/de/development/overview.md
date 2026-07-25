# Entwicklungsübersicht

## Über dieses Projekt

swag2mcp ist ein Go-Projekt, das OpenAPI/Swagger/Postman-Spezifikationen mit LLM-Agenten über das Model Context Protocol (MCP) verbindet. Es ist mit Go 1.23+ gebaut und folgt strengen Codierungskonventionen, die von 80+ Linters durchgesetzt werden.

Dieser Abschnitt ist für **Entwickler** geschrieben, die die Codebasis verstehen, dazu beitragen oder swag2mcp mit neuen Auth-Methoden, MCP-Tools oder Integrationen erweitern möchten.

## Entwicklungs-Skills

Das Projekt enthält zwei Entwicklungs-Skills, die die Konventionen und Muster des Projekts kodieren. Sie können sie verwenden oder ignorieren — sie sind Werkzeuge, keine Regeln.

### godeveloper

Der [godeveloper-Skill](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/godeveloper/SKILL.md) definiert jede Code-Konvention im Projekt:

- **Namensgebung** — Pakete, Dateien, Typen, Interfaces, Empfänger, Konstanten
- **Formatierung** — gofmt/gofumpt/goimports/gci, 120-Zeilen-Limit, Import-Reihenfolge
- **Fehlerbehandlung** — `LLMError` mit 8 Fehlercodes, Sentinel-Fehler, Fehler-Wrapping
- **Interfaces** — kleine Interfaces, Komposition, verbraucherseitige Definitionen
- **Nebenläufigkeit** — Mutex-Granularität, Goroutine-Lebensdauern, Kontext-Übergabe
- **Tests** — tabellengesteuerte Tests, `newTestService()`/`seedTestData()`-Helfer, Mock-Generierung
- **Projektmuster** — Service-Schicht, Anfrage-/Antwort-Strukturen, funktionale Optionen, MCP-Handler-Muster

### swag2mcp-cli

Der [swag2mcp-cli-Skill](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md) dokumentiert jeden CLI-Befehl mit Syntax, Flags, Argumenten und Beispielen. Nützlich bei der Arbeit an CLI-Befehlen oder beim Schreiben von Dokumentation.

## Wichtige Architekturentscheidungen

### Service-Schicht-Muster

Jede Funktion folgt dem gleichen dreistufigen Muster:

1. **Validieren** der Anfrage mit `s.validateRequest(req)` (verwendet `go-playground/validator`)
2. **Nachschlagen** von Entitäten aus dem In-Memory-Index (gibt `LLMError` mit `not_found`-Code zurück)
3. **Ausführen** der Geschäftslogik und Rückgabe einer typisierten Antwort oder eines `LLMError`

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

### Anfrage-/Antwort-Strukturen

Jede Methode hat eine dedizierte `{Method}Request`- und `{Method}Response`-Struktur. Anfrage-Strukturen verwenden `validate`-Tags für die Validierung und `jsonschema`-Tags für die Dokumentation:

```go
type SearchRequest struct {
    Query string `json:"query" validate:"required,min=1" jsonschema:"description=Suchabfrage mit Feldfilter-Unterstützung"`
    Limit int    `json:"limit" validate:"required,min=1,max=50" jsonschema:"description=Maximale Ergebnisse"`
}

type SearchResponse struct {
    Results []EndpointSearchItem `json:"results"`
}
```

### Funktionale Optionen

Die Konfiguration verwendet das Muster der funktionalen Optionen:

```go
type Option func(*Service)

func New(opts ...Option) (*Service, error)

func WithDisableLLMAuth(disable bool) Option {
    return func(s *Service) {
        s.disableLLMAuth.Store(disable)
    }
}
```

### MCP-Handler-Muster

Der MCP-Server verwendet ein zusammengesetztes Interface-Muster. Das `Svc`-Interface in `internal/server/mcp/handler.go` wird aus kleineren Interfaces zusammengesetzt (`CatalogReader`, `EndpointExplorer`, `EndpointExecutor`, `SystemInfo`, `ResponseManager`). Jede Handler-Methode delegiert an die Service-Schicht:

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

Alle an den LLM zurückgegebenen Fehler verwenden den Typ `LLMError` mit einem von 8 Codes:

| Code | Wann |
|------|------|
| `validation_failed` | Ungültige Eingabe (falsches ID-Format, fehlende Pflichtfelder) |
| `not_found` | Entität nicht im Index gefunden |
| `rate_limit` | Pro-Endpunkt-10s-Abklingzeit überschritten |
| `invoke_error` | HTTP-Anfrage-/Antwortfehler |
| `config_error` | Fehler beim Laden oder Validieren der Konfiguration |
| `workspace_error` | Fehler bei Arbeitsbereichsverzeichnis- oder Dateioperation |
| `parse_error` | Fehler beim Parsen der Spezifikationsdatei |
| `auth_error` | Fehler beim Abrufen des Authentifizierungstokens |

Nachrichten müssen erklären, was schiefgelaufen ist UND was als nächstes zu tun ist, in einfacher Sprache, die für einen LLM-Verbraucher geeignet ist.

### ID-Generierung

Alle IDs sind deterministische MD5-Hashes:

```go
id.Domain("meteo")                          // 32-stelliger Hex-Wert
id.Collection("meteo", "Forecast")          // 32-stelliger Hex-Wert
id.Tag("meteo", "Forecast", "pets")         // 32-stelliger Hex-Wert
id.Method("meteo", "Forecast", "pets", "GET", "/v2/pet/{petId}")
```

### Konfigurationskaskade

Die Konfiguration kaskadiert durch drei Ebenen: **global → spec → collection**. Jede Ebene überschreibt die vorherige. Alle `http_client`-Einstellungen können auf jeder Ebene überschrieben werden. Header und Cookies werden zusammengeführt; einfache Werte werden ersetzt.

## Kurzreferenz

| Bereich | Konvention |
|---------|------------|
| **Go-Version** | 1.23+ |
| **Formatierer** | gofmt, gofumpt, goimports, gci |
| **Zellenlänge** | 120 Zeichen |
| **Linter** | 80+ in `.golangci.yml` |
| **Fehlertyp** | `LLMError` mit 8 Codes |
| **Mock-Framework** | `go.uber.org/mock` |
| **Test-Helfer** | `newTestService()`, `seedTestData()` |
| **Konfigurationsformat** | YAML mit Kaskade |
| **Auth-Dispatch** | `UnmarshalYAML` liest `type`-Feld |
| **ID-Generierung** | MD5-basiert (`id.Domain()`, `id.Collection()`, usw.) |
| **Ratenlimit** | 10s pro Endpunkt für `invoke` |
| **Antwortgröße** | 1 MB Standard, bei Überschreitung in Datei gespeichert |
| **Abdeckungsziel** | 80%+ für Kernpakete |
| **Bauen** | `make build` |
| **Lint** | `make lint` |
| **Test** | `go test ./...` |
| **Generieren** | `go generate ./...` |
