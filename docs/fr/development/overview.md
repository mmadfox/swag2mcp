# Aperçu du développement

## À propos de ce projet

swag2mcp est un projet Go qui fait le pont entre les spécifications OpenAPI/Swagger/Postman et les agents LLM via le Model Context Protocol (MCP). Il est construit avec Go 1.23+ et suit des conventions de codage strictes appliquées par plus de 80 analyseurs de code.

Cette section est écrite pour les **ingénieurs** qui souhaitent comprendre la base de code, contribuer ou étendre swag2mcp avec de nouvelles méthodes d'authentification, de nouveaux outils MCP ou de nouvelles intégrations.

## Compétences de développement

Le projet est livré avec deux compétences de développement qui encodent les conventions et les modèles du projet. Vous pouvez les utiliser ou les ignorer — ce sont des outils, pas des règles.

### godeveloper

La [compétence godeveloper](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/godeveloper/SKILL.md) définit toutes les conventions de code du projet :

- **Nommage** — packages, fichiers, types, interfaces, récepteurs, constantes
- **Formatage** — gofmt/gofumpt/goimports/gci, limite de 120 lignes, ordre des importations
- **Gestion des erreurs** — `LLMError` avec 8 codes d'erreur, erreurs sentinelles, encapsulation des erreurs
- **Interfaces** — petites interfaces, composition, définitions côté consommateur
- **Concurrence** — granularité des mutex, durée de vie des goroutines, passage de contexte
- **Tests** — tests pilotés par tableaux, helpers `newTestService()`/`seedTestData()`, génération de mocks
- **Modèles du projet** — couche de service, structures requête/réponse, options fonctionnelles, modèle de gestionnaire MCP

### swag2mcp-cli

La [compétence swag2mcp-cli](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md) documente chaque commande CLI avec sa syntaxe, ses indicateurs, ses arguments et ses exemples. Utile lorsque vous travaillez sur des commandes CLI ou que vous rédigez de la documentation.

## Décisions architecturales clés

### Modèle de couche de service

Chaque fonctionnalité suit le même modèle en trois étapes :

1. **Valider** la requête avec `s.validateRequest(req)` (utilise `go-playground/validator`)
2. **Rechercher** les entités dans l'index en mémoire (renvoie `LLMError` avec le code `not_found`)
3. **Exécuter** la logique métier et renvoyer une réponse typée ou une `LLMError`

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

### Structures Requête/Réponse

Chaque méthode a des structures `{Méthode}Request` et `{Méthode}Response` dédiées. Les structures de requête utilisent les balises `validate` pour la validation et les balises `jsonschema` pour la documentation :

```go
type SearchRequest struct {
    Query string `json:"query" validate:"required,min=1" jsonschema:"description=Requête de recherche prenant en charge les filtres de champ"`
    Limit int    `json:"limit" validate:"required,min=1,max=50" jsonschema:"description=Résultats maximum"`
}

type SearchResponse struct {
    Results []EndpointSearchItem `json:"results"`
}
```

### Options fonctionnelles

La configuration utilise le modèle des options fonctionnelles :

```go
type Option func(*Service)

func New(opts ...Option) (*Service, error)

func WithDisableLLMAuth(disable bool) Option {
    return func(s *Service) {
        s.disableLLMAuth.Store(disable)
    }
}
```

### Modèle de gestionnaire MCP

Le serveur MCP utilise un modèle d'interface composée. L'interface `Svc` dans `internal/server/mcp/handler.go` est composée d'interfaces plus petites (`CatalogReader`, `EndpointExplorer`, `EndpointExecutor`, `SystemInfo`, `ResponseManager`). Chaque méthode de gestionnaire délègue à la couche de service :

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

Toutes les erreurs renvoyées au LLM utilisent le type `LLMError` avec l'un des 8 codes :

| Code | Quand |
|------|------|
| `validation_failed` | Entrée invalide (mauvais format d'ID, champs obligatoires manquants) |
| `not_found` | Entité non trouvée dans l'index |
| `rate_limit` | Délai de refroidissement de 10s par point de terminaison dépassé |
| `invoke_error` | Échecs de requête/réponse HTTP |
| `config_error` | Échec de chargement ou de validation de la configuration |
| `workspace_error` | Échec d'opération sur le répertoire ou le fichier de l'espace de travail |
| `parse_error` | Échec d'analyse du fichier de spécification |
| `auth_error` | Échec de récupération du jeton d'authentification |

Les messages doivent expliquer ce qui n'a pas fonctionné ET quoi faire ensuite, dans un langage simple adapté à un consommateur LLM.

### Génération d'ID

Tous les ID sont des hachages MD5 déterministes :

```go
id.Domain("meteo")                          // hexa 32 caractères
id.Collection("meteo", "Forecast")          // hexa 32 caractères
id.Tag("meteo", "Forecast", "pets")         // hexa 32 caractères
id.Method("meteo", "Forecast", "pets", "GET", "/v2/pet/{petId}")
```

### Cascade de configuration

La configuration se propage en cascade à travers trois niveaux : **global → spécification → collection**. Chaque niveau remplace le précédent. Tous les paramètres `http_client` peuvent être remplacés à chaque niveau. Les en-têtes et les cookies sont fusionnés ; les valeurs simples sont remplacées.

## Référence rapide

| Domaine | Convention |
|------|------------|
| **Version Go** | 1.23+ |
| **Formateurs** | gofmt, gofumpt, goimports, gci |
| **Longueur de ligne** | 120 caractères |
| **Analyseurs** | 80+ dans `.golangci.yml` |
| **Type d'erreur** | `LLMError` avec 8 codes |
| **Framework de mock** | `go.uber.org/mock` |
| **Helpers de test** | `newTestService()`, `seedTestData()` |
| **Format de configuration** | YAML avec cascade |
| **Répartition de l'authentification** | `UnmarshalYAML` lit le champ `type` |
| **Génération d'ID** | Basé sur MD5 (`id.Domain()`, `id.Collection()`, etc.) |
| **Limite de débit** | 10s par point de terminaison pour `invoke` |
| **Taille de réponse** | 1 Mo par défaut, enregistré dans un fichier si dépassé |
| **Objectif de couverture** | 80%+ pour les packages principaux |
| **Construction** | `make build` |
| **Analyse** | `make lint` |
| **Test** | `go test ./...` |
| **Génération** | `go generate ./...` |
