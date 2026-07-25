# Structure du projet

```
swag2mcp/
├── cmd/
│   ├── swag2mcp/          # Binaire principal
│   │   └── main.go
│   └── swag2mcp-mock/     # Serveur de simulation
│       └── main.go
├── internal/
│   ├── auth/              # 9 méthodes d'authentification
│   ├── cache/             # Mise en cache des spécifications
│   ├── commands/          # 13 commandes CLI (cobra)
│   ├── config/            # Configuration YAML
│   ├── env/               # Variables d'environnement
│   ├── httpclient/        # Client HTTP
│   ├── id/                # Génération d'ID MD5
│   ├── index/             # Recherche en texte intégral (bluge)
│   ├── model/             # Modèles de données
│   ├── reader/            # Lecture des grandes réponses
│   ├── server/
│   │   ├── mcp/           # Serveur MCP (19 outils)
│   │   └── mockserver/    # Serveur de simulation
│   ├── service/           # Logique métier
│   ├── spec/              # Analyseurs de spécifications
│   ├── tui/               # Interface TUI
│   └── workspace/         # Gestion de l'espace de travail
├── specs/                 # Exemples de spécifications
├── tests/                 # Tests d'intégration
├── docs/                  # Documentation
├── examples/              # Exemples de configuration
└── playground/            # Bac à sable de développement
```

## Packages clés

| Package | Description |
|---------|-------------|
| `auth` | 9 méthodes d'authentification |
| `cache` | Mise en cache sur disque avec TTL |
| `commands` | Commandes CLI Cobra |
| `config` | Configuration YAML avec cascade |
| `httpclient` | Client HTTP configurable |
| `index` | Recherche en texte intégral (bluge) |
| `server/mcp` | Serveur MCP (3 transports) |
| `service` | Logique métier (noyau) |
| `spec` | Analyseurs OpenAPI/Swagger/Postman |
| `tui` | TUI Bubbletea |
| `workspace` | Gestion des fichiers |
