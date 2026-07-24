# OAuth2 Client Credentials

## Zweck

OAuth2-Client-Credentials-Grant — Authentifizierung für Server-zu-Server-Kommunikation. Die Anwendung erhält ein Token mit ihrer client_id und client_secret, ohne Benutzerbeteiligung.

## Wann verwenden

- Microservices und Server-zu-Server-Integrationen
- Maschine-zu-Maschine-Kommunikation
- Wenn die API OAuth2 verwendet und Sie eine client_id + client_secret haben

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
      type: oauth2-cc
      config:
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        token_url: "https://auth.example.com/oauth/token"
        scopes:
          - read
          - write
```

## Parameter

| Parameter | Erforderlich | Beschreibung |
|-----------|-------------|--------------|
| `client_id` | Ja | Client-Identifikator |
| `client_secret` | Ja | Client-Geheimnis |
| `token_url` | Ja | Token-Endpunkt-URL |
| `scopes` | Nein | Liste der Berechtigungen (optional) |

## Hinweise

- swag2mcp fordert automatisch ein neues Token an, wenn das aktuelle abläuft
- Das Token wird bis zu seiner Ablaufzeit (`expires_in`) zwischengespeichert
- Wenn der Server kein `expires_in` angibt, gilt das Token für 1 Stunde als gültig
- Alle Parameter können in Umgebungsvariablen gespeichert werden
