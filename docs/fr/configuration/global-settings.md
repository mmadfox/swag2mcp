# Paramètres globaux

Les paramètres globaux sont les blocs de configuration de premier niveau dans `swag2mcp.yaml`. Ils s'appliquent à toutes les spécifications, sauf s'ils sont remplacés au niveau de la spécification ou de la collection.

## Structure

```yaml
http_client:
  # Paramètres du client HTTP pour tous les appels API

mcp:
  # Paramètres du serveur MCP

mock_enabled: false
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

disable_ratelimiter: false
rate_limit_interval: 10s
```

## Client HTTP

Contrôle la manière dont swag2mcp effectue les requêtes HTTP vers les API : délai d'attente, limite de taille de réponse, proxy, en-têtes, cookies, redirections et user-agent. Ces paramètres se propagent en cascade vers les spécifications et les collections.

Consultez [Client HTTP](./http-client) pour tous les paramètres et exemples.

## Serveur MCP

Contrôle la manière dont le serveur MCP communique avec les agents LLM : type de transport (stdio, SSE, Streamable HTTP), adresse, chemin et authentification par jeton bearer optionnelle.

Consultez [Serveur MCP](./mcp-server) pour tous les paramètres, transports et indicateurs de démarrage.

## Serveur de simulation

Le serveur de simulation génère des réponses API factices basées sur les schémas OpenAPI. Utile pour les tests sans solliciter les vraies API.

```yaml
mock_enabled: true
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092
```

### mock_enabled

- **Type :** `bool`
- **Valeur par défaut :** `false`
- **Effet :** Lorsqu'il est `true`, swag2mcp démarre des serveurs de simulation pour toutes les spécifications qui ont `base_mock_url` configuré. Chaque collection doit avoir `base_mock_url` défini.
- **Quand l'activer :** Vous souhaitez tester votre intégration API sans effectuer de véritables appels HTTP. Les serveurs de simulation renvoient des données factices basées sur le schéma OpenAPI.

### mock_auth

Configuration des ports pour les serveurs d'authentification de simulation. Ils sont utilisés lors des tests des méthodes d'authentification (OAuth2, Digest, HMAC) avec le serveur de simulation.

| Champ | Type | Valeur par défaut | Description |
|-------|------|---------|-------------|
| `oauth2_port` | int | `9090` | Port du serveur de jeton OAuth2 de simulation (1024-65535) |
| `digest_port` | int | `9091` | Port du serveur d'authentification Digest de simulation (1024-65535) |
| `hmac_port` | int | `9092` | Port du serveur d'authentification HMAC de simulation (1024-65535) |

## Limiteur de débit

Le limiteur de débit empêche le LLM d'appeler le même point de terminaison API trop fréquemment. Par défaut, chaque point de terminaison peut être appelé une fois toutes les 10 secondes.

```yaml
disable_ratelimiter: false
rate_limit_interval: 10s
```

### disable_ratelimiter

- **Type :** `bool`
- **Valeur par défaut :** `false`
- **Effet :** Lorsqu'il est `true`, le limiteur de débit par point de terminaison est complètement désactivé. Le LLM peut appeler le même point de terminaison de manière répétée sans attendre.
- **Quand l'activer :** Tests, débogage, ou lorsque vous devez appeler le même point de terminaison plusieurs fois rapidement.
- **Quand le laisser désactivé (recommandé) :** Production. Le limiteur de débit empêche les abus accidentels et respecte les limites de débit des API.

### rate_limit_interval

- **Type :** durée (format Go : `10s`, `30s`, `1m`)
- **Valeur par défaut :** `10s`
- **Effet :** Définit le temps d'attente obligatoire du LLM entre les appels au même point de terminaison.
- **Quand le modifier :** Augmentez pour les API avec des limites de débit strictes. Diminuez pour les API internes dont vous contrôlez la charge.
- **Plage :** Toute durée valide (par exemple, `5s`, `30s`, `1m`, `2m`).

## Cascade

Les paramètres globaux peuvent être remplacés au niveau de la spécification et de la collection. Tous les paramètres `http_client` (délai d'attente, proxy, user-agent, redirections, taille de réponse, randomiseur, en-têtes, cookies) peuvent être remplacés aux niveaux spécification et collection.

```
Global (http_client, mock_enabled, disable_ratelimiter, rate_limit_interval)
    ↓ remplace (http_client uniquement)
Spécification (specs[].http_client)
    ↓ remplace (http_client uniquement)
Collection (specs[].collections[].http_client)
```

Consultez [Cascade de configuration](./cascade) pour plus de détails.
