# Ajouter un nouvel outil MCP

## Étapes

1. **Ajouter une constante de nom d'outil** dans `internal/service/service.go`
2. **Créer les types requête/réponse** dans `internal/service/types.go`
3. **Implémenter le service** dans `internal/service/` (nouveau fichier ou ajouter à un existant)
4. **Créer une définition markdown** dans `internal/service/definitions/` — c'est ce que `MakeToolDefinitions` lit
5. **Ajouter une méthode à l'interface `Svc`** dans `internal/server/mcp/handler.go`
6. **Ajouter un gestionnaire** dans `handler.go`
7. **Enregistrer l'outil** dans `registerTools` dans `mcp.go`
8. **Générer les mocks** : `go generate ./...`
9. **Écrire les tests**

## 1. Constante de nom d'outil

Ajoutez une constante dans `internal/service/service.go` :

```go
const MyNewTool = "my_new_tool"
```

## 2. Types Requête/Réponse

Définissez dans `internal/service/types.go` :

```go
type MyNewToolRequest struct {
    Param1 string `json:"param1" validate:"required" jsonschema:"required,Description de param1"`
}

type MyNewToolResponse struct {
    Result string `json:"result"`
}
```

## 3. Implémentation du service

Créez `internal/service/my_new_tool.go` ou ajoutez à un fichier de service existant. Suivez le modèle de service standard : valider → rechercher → exécuter → renvoyer :

```go
func (s *Service) MyNewTool(ctx context.Context, req MyNewToolRequest) (MyNewToolResponse, error) {
    if err := s.validateRequest(req); err != nil {
        return MyNewToolResponse{}, NewLLMError(validationFailedErrCode, err.Error())
    }
    // logique métier
    return MyNewToolResponse{Result: "ok"}, nil
}
```

## 4. Définition markdown

Créez `internal/service/definitions/my_new_tool.md`. Ce fichier est lu par `MakeToolDefinitions()` et intégré dans le binaire. Le champ `name:` du frontmatter doit correspondre à la constante :

```markdown
---
name: my_new_tool
---

# my_new_tool

Description de l'outil.

## Paramètres

| Paramètre | Type | Description |
|-----------|------|-------------|
| `param1` | string | Description |
```

La fonction `MakeToolDefinitions()` dans `tools.go` lit tous les fichiers `.md` du répertoire `definitions/` intégré, analyse le frontmatter YAML pour le champ `name` et utilise le corps comme description de l'outil. Le fichier `instruction.md` est traité spécialement — il devient l'instruction système pour le LLM.

## 5. Interface Svc

Ajoutez une méthode à l'interface composée `Svc` dans `handler.go` :

```go
type Svc interface {
    // ... méthodes existantes
    MyNewTool(ctx context.Context, req service.MyNewToolRequest) (service.MyNewToolResponse, error)
}
```

## 6. Gestionnaire

Ajoutez une méthode de gestionnaire sur `handler` dans `handler.go`. Le gestionnaire délègue au service et encapsule le résultat dans `StructuredContent` :

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

## 7. Enregistrement

Enregistrez l'outil dans la fonction `registerTools` dans `mcp.go`. Ajoutez une entrée à la map `toolRegistrations` :

```go
service.MyNewTool: {
    addTool[service.MyNewToolRequest](mcpServer, h.handleMyNewTool),
    true, // false si l'outil est modifiable (comme invoke ou auth)
},
```

La signature de la fonction `registerTools` est :

```go
func registerTools(mcpServer *sdkmcp.Server, tools []service.Tool, h handler) {
```

Elle itère sur les définitions d'outils renvoyées par `MakeToolDefinitions()` et enregistre chacune avec son gestionnaire typé. La map `toolRegistrations` relie les constantes de nom d'outil à leurs gestionnaires.
