# Bearer Auth

## Zweck

Bearer-Token-Authentifizierung — die gebräuchlichste Methode für moderne REST-APIs. Das Token wird im `Authorization: Bearer &lt;token&gt;`-Header gesendet.

## Wann verwenden

- Moderne REST-APIs
- JWT (JSON Web Tokens)
- OAuth2-Zugriffstokens (wenn das Token bereits bezogen wurde)
- Jede API, die ein Bearer-Token akzeptiert

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
      type: bearer
      config:
        token: "eyJhbGciOiJIUzI1NiIs..."
```

## Parameter

| Parameter | Erforderlich | Beschreibung |
|-----------|-------------|--------------|
| `token` | Ja | Bearer-Token (JWT, OAuth2-Token usw.) |

## Hinweise

- Das Token ist statisch — wenn es abläuft, müssen Sie es manuell in der Konfiguration aktualisieren
- Für automatische Token-Erneuerung verwenden Sie `oauth2-cc` oder `oauth2-pwd`
- Speichern Sie das Token in einer Umgebungsvariable: `token: "$(API_TOKEN)"`
