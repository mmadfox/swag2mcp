# Cascade de configuration

swag2mcp utilise une cascade de configuration à trois niveaux. Chaque niveau remplace le précédent. Cela vous permet de définir des valeurs par défaut sensées globalement et d'affiner les paramètres pour des spécifications ou collections spécifiques.

## Niveaux

```
Global (http_client, mcp, mock_enabled, disable_ratelimiter, rate_limit_interval)
    ↓ remplace
Spécification (specs[].http_client, specs[].auth, specs[].base_url, specs[].disable, specs[].tags)
    ↓ remplace
Collection (specs[].collections[].http_client, specs[].collections[].base_url, specs[].collections[].disable)
```

## Ce qui remplace quoi

| Paramètre | Global | Spécification | Collection |
|-----------|--------|------|------------|
| `http_client.timeout` | ✅ | ✅ | ✅ |
| `http_client.max_response_size` | ✅ | ✅ | ✅ |
| `http_client.user_agent` | ✅ | ✅ | ✅ |
| `http_client.follow_redirects` | ✅ | ✅ | ✅ |
| `http_client.max_redirects` | ✅ | ✅ | ✅ |
| `http_client.proxy` | ✅ | ✅ | ✅ |
| `http_client.random` | ✅ | ✅ | ✅ |
| `http_client.headers` | ✅ | ✅ | ✅ |
| `http_client.cookies` | ✅ | ✅ | ✅ |
| `base_url` | ❌ | ✅ | ✅ |
| `auth` | ❌ | ✅ | ❌ |
| `disable` | ❌ | ✅ | ✅ |
| `tags` | ❌ | ✅ | ❌ |
| `mock_enabled` | ✅ | ❌ | ❌ |
| `disable_ratelimiter` | ✅ | ❌ | ❌ |
| `rate_limit_interval` | ✅ | ❌ | ❌ |

Tous les paramètres `http_client` peuvent être remplacés à chaque niveau. Les paramètres au niveau de la collection prévalent entièrement sur ceux de la spécification et du global.

## Exemple de cascade

```yaml
http_client:
  timeout: 30s
  max_response_size: 1048576
  headers:
    "User-Agent": "swag2mcp/1.0"

specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    http_client:
      timeout: 60s  # remplace le délai d'attente global
      headers:
        "X-API-Version": "2"  # ajouté aux en-têtes globaux
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          timeout: 120s  # remplace le délai d'attente de la spécification
          headers:
            "X-Custom": "valeur"  # ajouté aux en-têtes de la spécification + globaux
```

## Paramètres effectifs pour la collection « Forecast »

```
timeout: 120s (de la collection, remplace 60s de la spécification et 30s du global)
max_response_size: 1048576 (du global)
headers:
  - User-Agent: swag2mcp/1.0 (du global)
  - X-API-Version: 2 (de la spécification)
  - X-Custom: valeur (de la collection)
```

## Fonctionnement de la fusion

### Paramètres du client HTTP

Les valeurs simples (`timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`) sont **remplacées** à chaque niveau. Si une spécification définit `timeout: 60s`, elle remplace complètement le `30s` global.

### En-têtes

Les en-têtes sont **fusionnés** entre les niveaux. Les en-têtes des trois niveaux sont combinés. Si la même clé d'en-tête apparaît à plusieurs niveaux, le niveau le plus bas prévaut.

### Cookies

Les cookies sont **fusionnés** entre les niveaux. Si le même nom de cookie apparaît à plusieurs niveaux, le niveau le plus bas prévaut.

### Proxy

Le proxy est **remplacé** à chaque niveau. Si une spécification définit un proxy, il remplace complètement le proxy global pour cette spécification.
