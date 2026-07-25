# Basic Auth

## Zweck

HTTP-Basic-Authentifizierung — die einfachste Methode zur Authentifizierung mit Benutzername und Passwort.

## Wann verwenden

- Legacy-APIs, die nur Basic Auth unterstützen
- Einfache Authentifizierung ohne komplexe Tokens
- Interne Dienste

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
      type: basic
      config:
        username: "admin"
        password: "$(PASSWORD)"
```

## Parameter

| Parameter | Erforderlich | Beschreibung |
|-----------|-------------|--------------|
| `username` | Ja | Benutzername |
| `password` | Ja | Passwort |

## Hinweise

- Das Passwort wird im `Authorization: Basic ...`-Header Base64-kodiert gesendet — dies ist **keine Verschlüsselung**. Verwenden Sie immer HTTPS.
- Speichern Sie das Passwort in einer Umgebungsvariable: `password: "$(MY_PASSWORD)"`
