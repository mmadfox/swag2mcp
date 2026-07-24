# Konfigurationskaskade

swag2mcp verwendet eine dreistufige Konfigurationskaskade. Jede Ebene überschreibt die vorherige. Dies ermöglicht es Ihnen, sinnvolle Standardwerte global festzulegen und Einstellungen für bestimmte Specs oder Collections fein abzustimmen.

## Ebenen

```
Global (http_client, mcp, mock_enabled, disable_ratelimiter, rate_limit_interval)
    ↓ überschreibt
Spec (specs[].http_client, specs[].auth, specs[].base_url, specs[].disable, specs[].tags)
    ↓ überschreibt
Collection (specs[].collections[].http_client, specs[].collections[].base_url, specs[].collections[].disable)
```

## Was was überschreibt

| Parameter | Global | Spec | Collection |
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

Alle `http_client`-Einstellungen können auf jeder Ebene überschrieben werden. Collection-Ebene-Einstellungen haben volle Priorität vor Spec und Global.

## Kaskadenbeispiel

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
      timeout: 60s  # überschreibt globales timeout
      headers:
        "X-API-Version": "2"  # zu globalen Headern hinzugefügt
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          timeout: 120s  # überschreibt Spec-timeout
          headers:
            "X-Custom": "value"  # zu Spec + globalen Headern hinzugefügt
```

## Effektive Einstellungen für die "Forecast"-Collection

```
timeout: 120s (von Collection, überschreibt Spec 60s und Global 30s)
max_response_size: 1048576 (von Global)
headers:
  - User-Agent: swag2mcp/1.0 (von Global)
  - X-API-Version: 2 (von Spec)
  - X-Custom: value (von Collection)
```

## Wie die Zusammenführung funktioniert

### HTTP-Client-Einstellungen

Einfache Werte (`timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`) werden auf jeder Ebene **ersetzt**. Wenn eine Spec `timeout: 60s` setzt, ersetzt dies vollständig das globale `30s`.

### Header

Header werden über Ebenen hinweg **zusammengeführt**. Die Header aller drei Ebenen werden kombiniert. Wenn derselbe Header-Schlüssel auf mehreren Ebenen erscheint, gewinnt die niedrigste Ebene.

### Cookies

Cookies werden über Ebenen hinweg **zusammengeführt**. Wenn derselbe Cookie-Name auf mehreren Ebenen erscheint, gewinnt die niedrigste Ebene.

### Proxy

Der Proxy wird auf jeder Ebene **ersetzt**. Wenn eine Spec einen Proxy setzt, ersetzt dies vollständig den globalen Proxy für diese Spec.
