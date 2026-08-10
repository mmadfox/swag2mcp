# OpenID Connect (OIDC Discovery)

## Zweck

OpenID Connect Discovery-basierte Authentifizierung. swag2mcp ruft das OIDC-Discovery-Dokument vom Issuer ab, ermittelt den Token-Endpunkt und erhält ein Bearer-Token über den **Client-Credentials**-Grant. Das Token wird bis zum Ablauf zwischengespeichert.

## Wann verwenden

- APIs, die ein OIDC-Discovery-Dokument bereitstellen (`.well-known/openid-configuration`)
- Server-zu-Server-Integrationen mit `client_id` + `client_secret`
- Wenn die API OIDC verwendet und Sie eine automatische Token-Beschaffung wünschen

## Konfiguration

```yaml
specs:
  - domain: my-api
    llm_title: My API
    base_url: https://api.example.com
    collections:
      - llm_title: Main
        location: https://example.com/spec.yaml
    auth:
      type: openid-connect
      config:
        issuer: "https://auth.example.com"
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        scopes:
          - openid
          - profile
```

## Parameter

| Parameter | Erforderlich | Beschreibung |
|-----------|----------|-------------|
| `issuer` | Ja | OIDC-Issuer-URL. swag2mcp hängt `/.well-known/openid-configuration` an, um den Token-Endpunkt zu ermitteln |
| `client_id` | Ja | Client-Identifikator |
| `client_secret` | Ja | Client-Geheimnis |
| `scopes` | Nein | Liste der Berechtigungen (optional) |

## Hinweise

- swag2mcp fordert automatisch ein neues Token an, wenn das aktuelle abläuft
- Das Token wird bis zum Ablauf (`expires_in`) zwischengespeichert
- Wenn der Server kein `expires_in` liefert, gilt das Token 1 Stunde
- `issuer`, `client_id` und `client_secret` unterstützen `$(VAR)` für Umgebungsvariablen
- Das Token wird über den **Client-Credentials**-Grant bezogen (Server-zu-Server, ohne Benutzer)
- Das Discovery-Dokument muss ein `token_endpoint`-Feld enthalten
