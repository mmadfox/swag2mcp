# Serveur MCP

Le serveur MCP est le principal point d'interaction pour les agents LLM. Il expose toutes les API configurées sous forme d'outils MCP que le LLM peut appeler.

## Configuration

```yaml
mcp:
  transport: stdio
  addr: ":8080"
  path: "/mcp"
  auth:
    token: ""
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
