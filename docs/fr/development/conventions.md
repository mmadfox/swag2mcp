# Conventions de code

## Go

- **Go 1.26+**
- **gofmt** / **gofumpt** / **goimports** / **gci**
- **120 caractères** par ligne
- **Clauses de garde** au lieu de if imbriqués
- **Nommage** : `camelCase` pour le privé, `PascalCase` pour l'exporté

## Erreurs

Utilisez `LLMError` pour les erreurs visibles par le LLM :

```go
type LLMError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

Codes d'erreur :
- `validation_failed` — paramètres invalides
- `not_found` — ressource non trouvée
- `rate_limit` — limite de débit dépassée
- `invoke_error` — erreur d'appel API

## Interfaces

- Petites interfaces (1-3 méthodes)
- Composition d'interfaces
- Options fonctionnelles pour la configuration

## Tests

- Tests pilotés par tableaux
- Helpers de test (`newTestService()`, `seedTestData()`)
- Mocks via `go.uber.org/mock`
- 80%+ de couverture pour les packages principaux

## Configuration

- Format YAML
- Cascade : global → spécification → collection
- Validation via `go-playground/validator`
- Variables d'environnement via `$(VAR)`
