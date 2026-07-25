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

## Hinweise

- `client_secret` ist optional — **öffentliche Clients** werden unterstützt (z. B. Keycloak)
- swag2mcp erneuert das Token automatisch, wenn es abläuft
- Das Token wird bis zum Ablauf zwischengespeichert
- Alle Parameter können in Umgebungsvariablen gespeichert werden
