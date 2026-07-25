# Tests

## Commandes

```bash
# Tests unitaires
go test ./...

# Package spécifique
go test ./internal/service/...

# Tests d'intégration
make integration-tests

# Couverture
make cover

# Tous les tests
make testall
```

## Structure des tests

```
tests/
├── main_test.go              # Point d'entrée
├── suite_test.go             # Configuration de la suite
├── suite_auth_test.go        # Tests d'authentification
├── suite_config_test.go      # Tests de configuration
├── suite_mcp_tools_test.go   # Tests des outils MCP
├── suite_search_test.go      # Tests de recherche
├── suite_ratelimit_test.go   # Tests de limite de débit
├── suite_response_test.go    # Tests de réponse
├── suite_export_test.go      # Tests d'exportation
├── suite_import_test.go      # Tests d'importation
├── suite_parsing_test.go     # Tests d'analyse
├── suite_transport_test.go   # Tests de transport
├── suite_mock_test.go        # Tests du serveur de simulation
├── suite_workspace_test.go   # Tests de l'espace de travail
├── suite_errors_test.go      # Tests d'erreurs
└── suite_version_test.go     # Tests de version
```

## Couverture

Objectif : 80%+ pour les packages principaux :

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

## Mocks

Utilise `go.uber.org/mock` pour les tests du serveur MCP :

```bash
go generate ./...
```

Génère `internal/server/mcp/mock_svc_test.go` à partir de `handler.go`.

## Tests pilotés par tableaux

```go
func TestQuelqueChose(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"entrée valide", "bonjour", "BONJOUR", false},
        {"entrée vide", "", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FaireQuelqueChose(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.Equal(t, tt.want, got)
        })
    }
}
```
