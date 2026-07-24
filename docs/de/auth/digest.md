# Digest Auth

## Zweck

HTTP-Digest-Access-Authentifizierung — eine sicherere Alternative zu Basic Auth. Das Passwort wird nicht im Klartext gesendet; stattdessen werden MD5-Hashes verwendet.

## Wann verwenden

- Legacy-APIs, die nur Digest unterstützen
- Wenn Sie Authentifizierung benötigen, ohne das Passwort im Klartext zu senden
- Interne Unternehmenssysteme

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
      type: digest
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

- swag2mcp sendet zuerst eine Anfrage ohne Authentifizierung, erhält eine Challenge vom Server (HTTP 401), berechnet die Antwort und wiederholt den Vorgang mit dem `Authorization: Digest ...`-Header
- Die Challenge wird 5 Minuten lang zwischengespeichert — nachfolgende Anfragen benötigen keinen zusätzlichen Roundtrip
- Speichern Sie das Passwort in einer Umgebungsvariable: `password: "$(API_PASSWORD)"`
