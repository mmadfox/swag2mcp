# Hinzufügen eines neuen MCP-Tools

## Schritte

1. **Tool-Namenskonstante hinzufügen** in `internal/service/service.go`
2. **Anfrage-/Antworttypen erstellen** in `internal/service/types.go`
3. **Service implementieren** in `internal/service/` (neue Datei oder zu bestehender hinzufügen)
4. **Markdown-Definition erstellen** in `internal/service/definitions/` — dies liest `MakeToolDefinitions`
5. **Methode zum `Svc`-Interface hinzufügen** in `internal/server/mcp/handler.go`
6. **Handler hinzufügen** in `handler.go`
7. **Tool registrieren** in `registerTools` in `mcp.go`
8. **Mocks generieren**: `go generate ./...`
9. **Tests schreiben**

## 1. Tool-Namenskonstante

Fügen Sie eine Konstante in `internal/service/service.go` hinzu:

```go
const MyNewTool = "my_new_tool"
```

## 2. Anfrage-/Antworttypen

Definieren Sie in `internal/service/types.go`:

```go
type MyNewToolRequest struct {
    Param1 string `json:"param1" validate:"required" jsonschema:"required,Beschreibung von param1"`
}

type MyNewToolResponse struct {
    Result string `json:"result"`
}
```

## 3. Service-Implementierung

Erstellen Sie `internal/service/my_new_tool.go` oder fügen Sie zu einer bestehenden Service-Datei hinzu. Folgen Sie dem Standard-Service-Muster: validieren → nachschlagen → ausführen → zurückgeben:

```go
func (s *Service) MyNewTool(ctx context.Context, req MyNewToolRequest) (MyNewToolResponse, error) {
    if err := s.validateRequest(req); err != nil {
        return MyNewToolResponse{}, NewLLMError(validationFailedErrCode, err.Error())
    }
    // Geschäftslogik
    return MyNewToolResponse{Result: "ok"}, nil
}
```

## 4. Markdown-Definition

Erstellen Sie `internal/service/definitions/my_new_tool.md`. Diese Datei wird von `MakeToolDefinitions()` gelesen und in die Binärdatei eingebettet. Das `name:`-Feld im Frontmatter muss mit der Konstante übereinstimmen:

```markdown
---
name: my_new_tool
---

# my_new_tool

Beschreibung des Tools.

## Parameter

| Parameter | Typ | Beschreibung |
|-----------|-----|--------------|
| `param1` | string | Beschreibung |
```

Die Funktion `MakeToolDefinitions()` in `tools.go` liest alle `.md`-Dateien aus dem eingebetteten `definitions/`-Verzeichnis, parst das YAML-Frontmatter für das `name`-Feld und verwendet den Body als Tool-Beschreibung. Die Datei `instruction.md` wird speziell behandelt — sie wird zur Systemanweisung für den LLM.

## 5. Svc-Interface

Fügen Sie eine Methode zum zusammengesetzten `Svc`-Interface in `handler.go` hinzu:

```go
type Svc interface {
    // ... bestehende Methoden
    MyNewTool(ctx context.Context, req service.MyNewToolRequest) (service.MyNewToolResponse, error)
}
```

## 6. Handler

Fügen Sie eine Handler-Methode auf `handler` in `handler.go` hinzu. Der Handler delegiert an den Service und verpackt das Ergebnis in `StructuredContent`:

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

## 7. Registrierung

Registrieren Sie das Tool in der `registerTools`-Funktion in `mcp.go`. Fügen Sie einen Eintrag zur `toolRegistrations`-Map hinzu:

```go
service.MyNewTool: {
    addTool[service.MyNewToolRequest](mcpServer, h.handleMyNewTool),
    true, // false, wenn das Tool veränderlich ist (wie invoke oder auth)
},
```

Die Signatur der `registerTools`-Funktion ist:

```go
func registerTools(mcpServer *sdkmcp.Server, tools []service.Tool, h handler) {
```

Sie iteriert über die von `MakeToolDefinitions()` zurückgegebenen Tool-Definitionen und registriert jede mit ihrem typisierten Handler. Die `toolRegistrations`-Map verbindet Tool-Namenskonstanten mit ihren Handlern.
