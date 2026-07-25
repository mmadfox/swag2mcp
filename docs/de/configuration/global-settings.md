# Globale Einstellungen

Globale Einstellungen sind die Konfigurationsblöcke der obersten Ebene in `swag2mcp.yaml`. Sie gelten für alle Specs, es sei denn, sie werden auf Spec- oder Collection-Ebene überschrieben.

## Struktur

```yaml
http_client:
  # HTTP-Client-Einstellungen für alle API-Aufrufe

mcp:
  # MCP-Server-Einstellungen

mock_enabled: false
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

disable_ratelimiter: false
rate_limit_interval: 10s
```

## HTTP-Client

Steuert, wie swag2mcp HTTP-Anfragen an APIs stellt: Timeout, Antwortgrößenlimit, Proxy, Header, Cookies, Weiterleitungen und User-Agent. Diese Einstellungen kaskadieren zu Specs und Collections hinunter.

Siehe [HTTP-Client](./http-client) für alle Parameter und Beispiele.

## MCP-Server

Steuert, wie der MCP-Server mit LLM-Agenten kommuniziert: Transporttyp (stdio, SSE, Streamable HTTP), Adresse, Pfad und optionales Bearer-Token-Auth.

Siehe [MCP-Server](./mcp-server) für alle Parameter, Transports und Start-Flags.

## Mock-Server

Der Mock-Server generiert gefälschte API-Antworten basierend auf OpenAPI-Schemata. Nützlich zum Testen ohne echte API-Aufrufe.

```yaml
mock_enabled: true
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092
```

### mock_enabled

- **Typ:** `bool`
- **Standard:** `false`
- **Wirkung:** Wenn `true`, startet swag2mcp Mock-Server für alle Specs, die `base_mock_url` konfiguriert haben. Jede Collection muss `base_mock_url` gesetzt haben.
- **Wann aktivieren:** Sie möchten Ihre API-Integration testen, ohne echte HTTP-Aufrufe zu tätigen. Mock-Server geben gefälschte Daten basierend auf dem OpenAPI-Schema zurück.

### mock_auth

Port-Konfiguration für Mock-Authentifizierungsserver. Diese werden beim Testen von Auth-Methoden (OAuth2, Digest, HMAC) mit dem Mock-Server verwendet.

| Feld | Typ | Standard | Beschreibung |
|------|-----|----------|--------------|
| `oauth2_port` | int | `9090` | Port für den Mock-OAuth2-Token-Server (1024-65535) |
| `digest_port` | int | `9091` | Port für den Mock-Digest-Auth-Server (1024-65535) |
| `hmac_port` | int | `9092` | Port für den Mock-HMAC-Auth-Server (1024-65535) |

## Ratenbegrenzer

Der Ratenbegrenzer verhindert, dass der LLM denselben API-Endpunkt zu häufig aufruft. Standardmäßig kann jeder Endpunkt einmal alle 10 Sekunden aufgerufen werden.

```yaml
disable_ratelimiter: false
rate_limit_interval: 10s
```

### disable_ratelimiter

- **Typ:** `bool`
- **Standard:** `false`
- **Wirkung:** Wenn `true`, wird der Pro-Endpunkt-Ratenbegrenzer vollständig deaktiviert. Der LLM kann denselben Endpunkt wiederholt ohne Wartezeit aufrufen.
- **Wann aktivieren:** Testen, Debuggen oder wenn Sie denselben Endpunkt mehrmals schnell hintereinander aufrufen müssen.
- **Wann deaktiviert lassen (empfohlen):** Produktion. Der Ratenbegrenzer verhindert versehentlichen Missbrauch und respektiert API-Ratenlimits.

### rate_limit_interval

- **Typ:** Dauer (Go-Format: `10s`, `30s`, `1m`)
- **Standard:** `10s`
- **Wirkung:** Legt fest, wie lange der LLM zwischen Aufrufen desselben Endpunkts warten muss.
- **Wann ändern:** Erhöhen für APIs mit strengen Ratenlimits. Verringern für interne APIs, bei denen Sie die Last kontrollieren.
- **Bereich:** Jede gültige Dauer (z. B. `5s`, `30s`, `1m`, `2m`).

## Kaskade

Globale Einstellungen können auf Spec- und Collection-Ebene überschrieben werden. Alle `http_client`-Einstellungen (Timeout, Proxy, User-Agent, Weiterleitungen, Antwortgröße, Randomizer, Header, Cookies) können auf beiden Ebenen überschrieben werden.

```
Global (http_client, mock_enabled, disable_ratelimiter, rate_limit_interval)
    ↓ überschreibt (nur http_client)
Spec (specs[].http_client)
    ↓ überschreibt (nur http_client)
Collection (specs[].collections[].http_client)
```

Siehe [Konfigurationskaskade](./cascade) für Details.
