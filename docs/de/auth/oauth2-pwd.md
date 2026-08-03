# OAuth2 Password Grant

## Zweck

OAuth2-Resource-Owner-Password-Grant — Authentifizierung mit Benutzername und Passwort eines Benutzers. Geeignet für Erstanbieteranwendungen, bei denen der Benutzer der App seine Anmeldeinformationen anvertraut.

## Wann verwenden

- Erstanbieteranwendungen (mobil, Web)
- Integration mit Keycloak und ähnlichen Identitätsanbietern
- Wenn die API OAuth2 Password Grant unterstützt

## Konfiguration

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: oauth2-pwd
      config:
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        username: "$(USERNAME)"
        password: "$(PASSWORD)"
        token_url: "https://auth.example.com/oauth/token"
        scopes:
          - openid
          - profile
        request_format: form
```

## Parameter

| Parameter | Erforderlich | Beschreibung |
|-----------|-------------|--------------|
| `client_id` | Ja | Client-Identifikator |
| `username` | Ja | Benutzername |
| `password` | Ja | Passwort |
| `token_url` | Ja | Token-Endpunkt-URL |
| `client_secret` | Nein | Client-Geheimnis (optional, für öffentliche Clients) |
| `scopes` | Nein | Liste der Berechtigungen (optional) |
| `request_format` | Nein | Anfrageformat: `form` (Standard, `application/x-www-form-urlencoded`) oder `json` (`application/json`) |

## Hinweise

- `client_secret` ist optional — **öffentliche Clients** werden unterstützt (z. B. Keycloak)
- swag2mcp erneuert das Token automatisch, wenn es abläuft
- Das Token wird bis zum Ablauf zwischengespeichert
- Alle Parameter können in Umgebungsvariablen gespeichert werden:
- `client_id` — unterstützt `$(VAR)`-Syntax
- `client_secret` — unterstützt `$(VAR)`-Syntax (optional)
- `username` — unterstützt `$(VAR)`-Syntax
- `password` — unterstützt `$(VAR)`-Syntax
- `token_url` — unterstützt `$(VAR)`-Syntax
- `request_format: json` sendet die Token-Anfrage als JSON-Body — verwenden Sie dies, wenn der Token-Endpunkt `Content-Type: application/json` erfordert
