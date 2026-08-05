# Serveur MCP

Le serveur MCP est le principal point d'interaction pour les agents LLM. Il expose toutes les API configurées sous forme d'outils MCP que le LLM peut appeler.

## Configuration

```yaml
mcp:
  transport: stdio
```

## Transports

Trois types de transport sont disponibles :

| Transport | Description | Quand l'utiliser |
|-----------|-------------|-------------|
| `stdio` | Entrée/sortie standard | Clients LLM locaux (VS Code, Cursor, Claude Desktop) |
| `sse` | Événements envoyés par le serveur | Clients distants, communication basée sur HTTP |
| `streamable-http` | HTTP avec streaming | Clients web, clients MCP modernes |

### stdio (par défaut)

Le client LLM exécute swag2mcp comme un processus enfant. La communication se fait par l'entrée et la sortie standard. Aucun port réseau n'est nécessaire.

```yaml
mcp:
  transport: stdio
```

```bash
swag2mcp mcp
```

### SSE

Transport par événements envoyés par le serveur pour une communication basée sur HTTP. Le serveur MCP écoute sur un port HTTP et le client LLM se connecte à distance.

```yaml
mcp:
  transport: sse
  addr: "127.0.0.1:8080"
  path: "/mcp"
```

```bash
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080
```

### Streamable HTTP

Transport HTTP moderne qui prend en charge les réponses en streaming. Similaire à SSE mais utilise un protocole différent.

```yaml
mcp:
  transport: streamable-http
  addr: "127.0.0.1:8080"
  path: "/mcp"
```

```bash
swag2mcp mcp --transport streamable-http --http-addr 0.0.0.0:8080
```

## Paramètres

### transport

- **Type :** `string`
- **Valeur par défaut :** `"stdio"`
- **Options :** `stdio`, `sse`, `streamable-http`
- **Effet :** Détermine comment le serveur MCP communique avec le client LLM.

### addr

- **Type :** `string`
- **Valeur par défaut :** `":8080"`
- **Description :** Adresse d'écoute pour les transports SSE et Streamable HTTP. Format : `host:port`.
- **Exemples :** `":8080"`, `"127.0.0.1:8080"`, `"0.0.0.0:9000"`

### path

- **Type :** `string`
- **Valeur par défaut :** `"/mcp"`
- **Description :** Chemin URL pour le point de terminaison MCP. Le client LLM envoie les requêtes à `http://&lt;addr&gt;&lt;path&gt;`.
- **Exemples :** `"/mcp"`, `"/api/mcp"`, `"/v1/mcp"`

### auth.token

- **Type :** `string`
- **Valeur par défaut :** `""` (pas d'authentification)
- **Description :** Jeton Bearer pour l'authentification du transport HTTP. Lorsqu'il est défini, le client LLM doit inclure `Authorization: Bearer &lt;token&gt;` dans chaque requête.
- **Remarque :** Prend en charge la résolution `$(ENV_VAR)`.

### auth.type

- **Type:** `string`
- **Default:** `""` (no JWT auth)
- **Options:** `jwks`, `oidc`, `introspection`
- **Description:** JWT authentication type for HTTP transport. When set, enables dynamic token verification using JWKS, OIDC Discovery, or token introspection.

### auth.jwks_url

- **Type:** `string`
- **Default:** `""`
- **Description:** URL of the JWKS (JSON Web Key Set) endpoint. Required when `auth.type` is `jwks` or resolved via OIDC discovery.

### auth.issuer

- **Type:** `string`
- **Default:** `""`
- **Description:** Expected JWT issuer (`iss` claim). If set, tokens with a different issuer are rejected.

### auth.audience

- **Type:** `string`
- **Default:** `""`
- **Description:** Expected JWT audience (`aud` claim). If set, tokens without this audience are rejected.

### auth.introspection_url

- **Type:** `string`
- **Default:** `""`
- **Description:** Token introspection endpoint URL. Required when `auth.type` is `introspection`.

### auth.client_id

- **Type:** `string`
- **Default:** `""`
- **Description:** Client ID for introspection auth. Required when `auth.type` is `introspection`.

### auth.client_secret

- **Type:** `string`
- **Default:** `""`
- **Description:** Client secret for introspection auth. Supports `$(ENV_VAR)` resolution.

## Authentification HTTP

Protégez le point de terminaison HTTP MCP avec un jeton bearer :

```yaml
mcp:
  auth:
    token: "mon-jeton-secret"
```

Ou via l'indicateur CLI :

```bash
swag2mcp mcp --auth-token "mon-jeton-secret"
```

### With JWT authentication (JWKS)

Protect the MCP HTTP endpoint with JWT verification via a JWKS endpoint:

```yaml
mcp:
  auth:
    type: jwks
    jwks_url: "https://auth.example.com/.well-known/jwks.json"
    issuer: "https://auth.example.com/"
    audience: "swag2mcp"
```

```bash
swag2mcp mcp --transport sse --http-addr 0.0.0.0:8080 \
  --auth-type jwks \
  --auth-jwks-url "https://auth.example.com/.well-known/jwks.json" \
  --auth-issuer "https://auth.example.com/" \
  --auth-audience "swag2mcp"
```

### With JWT authentication (OIDC Discovery)

```yaml
mcp:
  auth:
    type: oidc
    issuer: "https://auth.example.com/"
    audience: "swag2mcp"
```

### With JWT authentication (Token Introspection)

```yaml
mcp:
  auth:
    type: introspection
    introspection_url: "https://auth.example.com/introspect"
    client_id: "my-client"
    client_secret: "$(MCP_CLIENT_SECRET)"
```

```bash
swag2mcp mcp --transport sse --http-addr 0.0.0.0:8080 \
  --auth-type introspection \
  --auth-introspection-url "https://auth.example.com/introspect" \
  --auth-client-id "my-client" \
  --auth-client-secret "$(MCP_CLIENT_SECRET)"
```

## Vérification de santé

Le serveur MCP fournit un point de terminaison de vérification de santé qui fonctionne sans initialisation MCP :

```bash
curl http://127.0.0.1:8080/health
# {"status":"ok","version":"v1.2.0"}
```

## Indicateurs de démarrage

Les indicateurs CLI remplacent la configuration YAML. Si un indicateur n'est pas défini, la valeur de la section `mcp` dans le YAML est utilisée comme solution de repli.

| Indicateur | Type | Valeur par défaut | Description |
|------|------|---------|-------------|
| `--transport` | string | `"stdio"` | Type de transport : `stdio`, `sse`, `streamable-http` |
| `--http-addr` | string | `":8080"` | Adresse du serveur HTTP (pour SSE et Streamable HTTP) |
| `--http-path` | string | `"/mcp"` | Chemin URL pour le gestionnaire MCP |
| `--auth-token` | string | `""` | Jeton Bearer pour l'authentification du transport HTTP |
| `--logfile` | string | `""` | Chemin du fichier journal (journalise sur stderr si non défini) |
| `--disable-llm-auth` | bool | `true` | Supprime l'outil `auth` de la liste des outils MCP |
| `--dump-dir` | string | `""` | Répertoire pour vider les requêtes HTTP pour le débogage |
| `--tags` | string | `""` | Filtrer les spécifications par balises (séparées par des virgules) |
| `--auth-type` | string | `""` | JWT auth type: `jwks`, `oidc`, `introspection` |
| `--auth-jwks-url` | string | `""` | JWKS URL for JWT auth |
| `--auth-issuer` | string | `""` | JWT issuer for token validation |
| `--auth-audience` | string | `""` | JWT audience for token validation |
| `--auth-introspection-url` | string | `""` | Token introspection URL |
| `--auth-client-id` | string | `""` | Client ID for introspection auth |
| `--auth-client-secret` | string | `""` | Client secret for introspection auth |
