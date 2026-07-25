# API-Schlüssel

## Zweck

Authentifizierung über einen API-Schlüssel. Der Schlüssel kann als HTTP-Header oder als URL-Abfrageparameter gesendet werden.

## Wann verwenden

- Dienste, die API-Schlüssel verwenden
- Wetterdienste, Geodaten, Übersetzungs-APIs
- Wenn die API einen Schlüssel in einem Header (`X-API-Key`) oder Abfrageparameter (`?api_key=...`) erwartet

## Konfiguration

### Schlüssel im Header

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: api-key
      config:
        key: "X-API-Key"
        in: header
        value: "$(API_KEY)"
```

### Schlüssel im Abfrageparameter

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: api-key
      config:
        key: "api_key"
        in: query
        value: "$(API_KEY)"
```

## Parameter

| Parameter | Erforderlich | Beschreibung |
|-----------|-------------|--------------|
| `key` | Ja | Name des Headers oder Abfrageparameters |
| `in` | Ja | Wo der Schlüssel platziert wird: `header` oder `query` |
| `value` | Ja | Der Schlüsselwert |

## Hinweise

- Im `header`-Modus wird der Schlüssel als HTTP-Header hinzugefügt
- Im `query`-Modus wird der Schlüssel als URL-Parameter hinzugefügt
- Speichern Sie den Wert in einer Umgebungsvariable: `value: "$(MY_API_KEY)"`
