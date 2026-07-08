# swag2mcp

> ⚠️ **En développement** — L'API peut changer, contributions bienvenues.

**swag2mcp** est un outil CLI et serveur MCP (Model Context Protocol) qui connecte les spécifications OpenAPI/Swagger/Postman avec des agents LLM (Opencode, Crush, Copilot, Cursor, etc.).

Il indexe vos spécifications API dans un moteur de recherche en texte intégral, les expose via 14 outils MCP et permet aux LLM de découvrir, inspecter et invoquer de véritables points d'API — sans écrire une seule ligne de code d'intégration.

---

## Table des matières

- [Démarrage rapide](#démarrage-rapide)
- [Configuration](#configuration)
- [Commandes CLI](#commandes-cli)
- [Serveur MCP](#serveur-mcp)
- [Recherche](#recherche)
- [Espace de travail (Workspace)](#espace-de-travail-workspace)
- [Cache](#cache)
- [Développement](#développement)

---

## Démarrage rapide

```bash
# Installation
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest

# Initialiser l'espace de travail
swag2mcp init

# Démarrer le serveur MCP (pour les agents LLM)
swag2mcp mcp

# Ou utiliser l'explorateur interactif
swag2mcp run
```

---

## Configuration

### Schéma YAML

```yaml
http_client:                        # optionnel, paramètres HTTP globaux
  headers:                          # optionnel
    X-API-Version: "2"
  cookies: []                       # optionnel
  user_agent: ""                    # optionnel
  timeout: 0s                       # optionnel
  follow_redirects: true            # optionnel
  max_redirects: 10                 # optionnel
  max_response_size: 1048           # optionnel, octets (défaut 1Ko, max 1Mo)

specs:
  - domain: petstore                    # obligatoire, 1-60 car., [a-zA-Z0-9_-]
    llm_title: Petstore API             # obligatoire, 5-120 car.
    llm_instruction: |                  # optionnel, max 500 car.
      Utilisez cette API pour gérer les animaux, commandes et utilisateurs.
    base_url: https://petstore.swagger.io/v2  # obligatoire, URL valide
    disable: false                      # optionnel
    tags: [public, demo]                # optionnel, pour le filtrage
    http_client:                        # optionnel, remplace le global
      headers:
        X-API-Version: "2"
    auth:                               # optionnel
      type: bearer                      # voir Méthodes d'authentification
      config:
        token: $(TOKEN_AUTH)
    collections:
      - llm_title: Petstore Swagger     # optionnel, max 120 car.
        llm_instruction: |             # optionnel, max 360 car.
          Points d'accès principaux de Petstore
        title: ""                      # optionnel, auto-rempli depuis la spec
        location: https://petstore.swagger.io/v2/swagger.json  # obligatoire, 5-250 car.
        disable: false                  # optionnel
        base_url: ""                    # optionnel, remplace base_url de la spec
        http_client: {}                 # optionnel, remplace la spec
```

### Tags — Filtrage des spécifications par projet

Les tags permettent de regrouper les spécifications par projet, environnement ou équipe. Au démarrage du serveur MCP, utilisez `--tags` pour charger uniquement les spécifications correspondantes :

```bash
# Démarrer le serveur avec uniquement les spécifications publiques
swag2mcp mcp --tags=public

# Démarrer avec plusieurs tags
swag2mcp mcp --tags=public,internal

# Exécuter plusieurs serveurs pour différents projets
swag2mcp mcp --tags=project-alpha --logfile=/tmp/swag2mcp-alpha.log
swag2mcp mcp --tags=project-beta  --logfile=/tmp/swag2mcp-beta.log
```

Cela permet d'exécuter des serveurs MCP séparés pour différents projets à partir d'un seul fichier de configuration.

### Méthodes d'authentification

| Type | Champs | Exemple de configuration |
|------|--------|--------------------------|
| `none` | — | `type: none` |
| `basic` | `username`, `password` | `username: $(USER)`, `password: $(PASS)` |
| `bearer` | `token` | `token: $(TOKEN)` |
| `digest` | `username`, `password` | `username: admin`, `password: secret` |
| `api-key` | `key`, `value`, `in` (header/query) | `key: X-API-Key`, `value: $(KEY)`, `in: header` |
| `oauth2-cc` | `client_id`, `client_secret`, `token_url`, `scopes` | `client_id: $(ID)`, `token_url: https://auth.example.com/token` |
| `oauth2-pwd` | `username`, `password`, `client_id`, `client_secret`, `token_url`, `scopes` | `username: $(USER)`, `token_url: https://auth.example.com/token` |
| `script` | `source` | `source: chemin/vers/auth.sh` |

Tous les champs de chaîne prennent en charge la syntaxe `$(ENV_VAR)` — résolue à l'exécution depuis les variables d'environnement.

---

## Commandes CLI

Toutes les commandes qui acceptent `[path]` utilisent la même résolution de chemin :

```
swag2mcp <command>          → ~/.swag2mcp/
swag2mcp <command> ./       → {cwd}/.swag2mcp/
swag2mcp <command> path/to  → {cwd}/path/to/.swag2mcp/
```

### `init [path]`

Initialiser l'espace de travail et la configuration.

| Option | Court | Défaut | Description |
|--------|-------|--------|-------------|
| `--interactive` | `-i` | `false` | Assistant interactif |
| `--force` | `-f` | `false` | Écraser la configuration existante |

```bash
swag2mcp init              # créer ~/.swag2mcp/swag2mcp.yaml
swag2mcp init ./           # créer ./.swag2mcp/swag2mcp.yaml
swag2mcp init -i           # assistant interactif
```

### `add spec [path]` / `add collection [path]`

Ajouter une spécification ou une collection à la configuration.

| Option | Court | Défaut | Description |
|--------|-------|--------|-------------|
| `--yaml` | `-y` | `""` | Entrée YAML (`-` pour stdin) |
| `--example` | `-e` | `false` | Afficher un exemple YAML |

```bash
swag2mcp add spec
swag2mcp add spec --yaml 'domain: petstore\nllm_title: Petstore API\nbase_url: https://...'
cat spec.yaml | swag2mcp add spec --yaml -
swag2mcp add spec --example
```

### `delete spec [path]` / `delete collection [path]`

Supprimer une spécification ou une collection. Sélection interactive.

```bash
swag2mcp delete spec
swag2mcp delete collection
```

### `ls [path]`

Lister les spécifications et collections.

| Option | Court | Défaut | Description |
|--------|-------|--------|-------------|
| `--tags` | `-t` | `""` | Filtrer par tags (séparés par des virgules) |

```bash
swag2mcp ls
swag2mcp ls --tags=public,internal
```

### `run [path]`

Explorateur API interactif (TUI). Rechercher, parcourir, inspecter et sauvegarder des points d'accès.

```bash
swag2mcp run
```

### `validate [path]`

Valider la configuration et vérifier que tous les emplacements de collections sont accessibles.

| Option | Court | Défaut | Description |
|--------|-------|--------|-------------|
| `--tags` | `-t` | `""` | Filtrer les spécifications par tags |

```bash
swag2mcp validate
swag2mcp validate --tags=public
```

### `clean [path]`

Supprimer tout le contenu des répertoires `cache/` et `responses/`.

```bash
swag2mcp clean
```

### `update [path]`

Valider la configuration, vider le cache, re-mettre en cache tous les fichiers de spécification.

```bash
swag2mcp update
```

### `mcp [path]`

Démarrer le serveur MCP en mode headless (transport stdio). C'est la commande de production principale pour l'intégration LLM.

| Option | Court | Défaut | Description |
|--------|-------|--------|-------------|
| `--logfile` | `-f` | `""` | Chemin du fichier journal |
| `--tags` | `-t` | `""` | Filtrer les spécifications par tags |
| `--disable-llm-auth` | | `true` | `true` — authentification en arrière-plan (LLM ne voit jamais les tokens). `false` — LLM peut demander des tokens via l'outil `auth` |
| `--dump-dir` | | `""` | Répertoire pour vider les requêtes HTTP (débogage) |

```bash
swag2mcp mcp
swag2mcp mcp --tags=public --logfile=/var/log/swag2mcp.log
swag2mcp mcp --disable-llm-auth=false
swag2mcp mcp --dump-dir=/tmp/dump
```

---

## Serveur MCP

Le serveur MCP expose 14 outils via le transport stdio. Les agents LLM (Opencode, Crush, Copilot, Cursor, etc.) se connectent automatiquement une fois configurés.

### Hiérarchie des outils

```
spec_list                       — lister toutes les spécifications disponibles
  └─ spec_by_id                 — détails d'une spécification par ID
       └─ collection_by_spec    — collections dans une spécification
            └─ tag_by_collection     — tags dans une collection
                 └─ endpoint_by_tag  — points d'accès sous un tag
                      └─ inspect          — opération OpenAPI complète
                           └─ invoke       — exécuter un appel API

search                          — recherche en texte intégral sur tous les points d'accès
```

### Référence des outils

| Outil | Arguments | Retourne | Description |
|-------|-----------|----------|-------------|
| `spec_list` | — | `Spec[]` | Toutes les spécifications disponibles |
| `spec_by_id` | `id` | Spec + Collections | Détails d'une spécification |
| `collection_by_spec` | `specId` | Collections | Collections dans une spécification |
| `collection_by_id` | `id` | Collection + Tags | Détails d'une collection |
| `tag_by_collection` | `collectionId` | Tags | Tags dans une collection |
| `tag_by_spec` | `specId` | Tags | Tous les tags d'une spécification |
| `tag_by_id` | `id` | Tag | Métadonnées d'un tag |
| `endpoint_by_tag` | `tagId` | Endpoints | Points d'accès sous un tag |
| `endpoint_by_collection` | `collectionId` | Endpoints | Tous les points d'accès d'une collection |
| `endpoint_by_spec` | `specId` | Endpoints | Tous les points d'accès d'une spécification |
| `endpoint_by_id` | `id` | Endpoint | Résumé rapide d'un point d'accès |
| `search` | `query`, `limit` | Endpoints | Recherche en texte intégral |
| `inspect` | `endpointId` | Full Operation | Objet d'opération OpenAPI complet |
| `invoke` | `endpointId`, `parameters`, `requestBody` | Response | Exécute un véritable appel API |
| `auth` | `specId` | Token | Obtenir un token d'authentification pour une spécification |

---

## Recherche

### Syntaxe des requêtes

| Fonctionnalité | Syntaxe | Exemple |
|----------------|---------|---------|
| Terme | `terme` | `animaux` |
| Phrase | `"phrase"` | `"ajouter un animal"` |
| Champ : method | `method:terme` | `method:post` |
| Champ : tag | `tag:terme` | `tag:auth` |
| Champ : path | `path:terme` | `path:/users` |
| Champ : summary | `summary:terme` | `summary:login` |
| Requis (AND) | `+terme` | `+method:post +tag:user` |
| Exclu (NOT) | `-terme` | `-deprecated` |
| Wildcard | `*` | `path:*/v2/*` |
| Flou | `terme~` | `watex~` |
| Regex | `/motif/` | `/user(s\|sessions)/` |
| Pondération | `terme^N` | `tag:pet^5` |
| Tout | `*` | `*` |

### Exemples

```
# Trouver les points d'accès POST dans le tag auth
+method:post +tag:auth

# Rechercher les points d'accès liés à la connexion
summary:"login"~

# Trouver tous les chemins liés aux utilisateurs, exclure les obsolètes
path:*/users/* -deprecated

# Requête complexe
+method:get +tag:pet summary:"find by status"
```

### Champs indexés

| Champ | Type | Contenu |
|-------|------|---------|
| `method` | text | Méthode HTTP (en minuscules) |
| `tag` | text | Nom du tag (en minuscules) |
| `path` | text | Chemin API (en minuscules) |
| `summary` | text (analysé) | Résumé/description du point d'accès (en minuscules) |
| `_all` | text (analysé) | method + path + tag + summary |

---

## Espace de travail (Workspace)

### Structure des répertoires

```
~/.swag2mcp/                    # ou {projet}/.swag2mcp/
├── swag2mcp.yaml               # Fichier de configuration
├── cache/                      # Spécifications distantes mises en cache
│   ├── {hash}.spec             # Contenu du fichier de spécification
│   └── {hash}.meta             # Métadonnées JSON
├── specs/                      # Fichiers de spécification locaux (gérés par l'utilisateur)
├── responses/                  # Fichiers de réponses d'appels
└── auth_scripts/               # Scripts d'authentification
```

### Résolution de chemin

```
swag2mcp <command>          → ~/.swag2mcp/
swag2mcp <command> ./       → {cwd}/.swag2mcp/
swag2mcp <command> path/to  → {cwd}/path/to/.swag2mcp/
```

### .gitignore

Seules les données temporaires doivent être ignorées :

```
.swag2mcp/cache/*
.swag2mcp/responses/*
```

La configuration `.swag2mcp/swag2mcp.yaml` et les fichiers de spécification dans `.swag2mcp/specs/` **doivent être dans le dépôt**.

### Recommandation

Conservez tous les fichiers de spécification dans `.swag2mcp/specs/` — c'est le seul moyen de garantir qu'ils sont utilisés directement sans être copiés dans le cache.

---

## Cache

### Règles

| Source | Comportement |
|--------|--------------|
| URL HTTP/HTTPS | Toujours mis en cache. TTL : aléatoire 1-48h. |
| Chemin local dans `specs/` | Utilisé directement, non mis en cache. |
| Chemin local hors `specs/` | Copié dans le cache au premier accès. |
| URL `file://` | Traité comme un chemin local. |

### Clé de cache

Hash SHA-256 de l'emplacement normalisé (16 premiers octets = 32 caractères hexadécimaux).

### Logique de succès du cache

1. Lire le fichier `.meta` — expiré ou manquant → échec
2. Pour les sources locales : `ModTime` modifié → échec
3. Fichier `.spec` manquant → échec
4. Sinon → succès

---

## Développement

```bash
# Compilation
go build ./cmd/swag2mcp/

# Tests
go test ./...

# Linter
make lint

# Exécution
go run ./cmd/swag2mcp/main.go
```
