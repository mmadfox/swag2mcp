# Fichier de configuration

swag2mcp utilise un fichier de configuration YAML. Créé par `swag2mcp init`.

## Emplacement

- **Linux/macOS** : `~/.swag2mcp/swag2mcp.yaml`
- **Windows** : `%USERPROFILE%\.swag2mcp\swag2mcp.yaml`

## Structure de base

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

## Exemple complet

```yaml
# ── Client HTTP global ──────────────────────────────────
http_client:
  timeout: 30s
  max_response_size: 1048576
  user_agent: "swag2mcp-global/1.0"
  follow_redirects: true
  max_redirects: 10
  random: false
  proxy:
    url: ""
    username: ""
    password: ""
    bypass: []
  headers:
    "Accept": "application/json"
  cookies:
    - name: "session"
      value: "abc123"

# ── Serveur MCP ──────────────────────────────────────────
mcp:
  transport: stdio
  addr: ":8080"
  path: "/mcp"
  auth:
    token: ""

# ── Serveur de simulation ───────────────────────────────
mock_enabled: false
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

# ── Limiteur de débit ────────────────────────────────────
disable_ratelimiter: false
rate_limit_interval: 10s

# ── Spécifications ──────────────────────────────────────
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    llm_instruction: "Utilisez cette API pour les prévisions météorologiques et les données climatiques"
    base_url: https://api.open-meteo.com
    disable: false
    tags: ["weather", "climate"]
    http_client:
      timeout: 10s
    auth:
      type: bearer
      config:
        token: "$(TOKEN)"
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: false
        http_client:
          timeout: 5s

  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Variables d'environnement

Utilisez la syntaxe `$(NOM_VAR)` pour référencer des variables d'environnement. swag2mcp les résout au démarrage.

```yaml
specs:
  - domain: meteo
    auth:
      type: bearer
      config:
        token: "$(API_TOKEN)"

mcp:
  auth:
    token: "$(MCP_TOKEN)"
```

`$(VAR)` est résolu dans :
- Les champs de configuration d'authentification : `token`, `username`, `password`, `client_id`, `client_secret`, `api_key`, `secret_key`, `domain`
- Le jeton d'authentification du serveur MCP : `mcp.auth.token`
- Les en-têtes et valeurs de cookies du client HTTP

`$(VAR)` n'est **pas** résolu dans les URL de base ou les emplacements des collections.

## Validation

```bash
# Valider l'espace de travail par défaut (~/.swag2mcp)
swag2mcp validate

# Valider un espace de travail de projet personnalisé
swag2mcp validate ./mon-projet
```

Si l'espace de travail ne se trouve pas dans le répertoire personnel (par exemple, dans un dépôt de projet), spécifiez toujours le chemin lors de l'exécution de `validate`, `update`, `mcp` ou toute autre commande. Sinon, swag2mcp utilisera l'espace de travail par défaut `~/.swag2mcp`.
