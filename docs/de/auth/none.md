# Keine (None)

## Zweck

Keine Authentifizierung erforderlich. Die API ist ohne Tokens oder Schlüssel zugänglich.

## Wann verwenden

- Öffentliche APIs (Open-Meteo, icanhazdadjoke, PokéAPI)
- Test- und Demo-Umgebungen
- Wenn die API keine Autorisierung erfordert

## Konfiguration

Setzen Sie `type: none` oder lassen Sie den `auth`-Abschnitt einfach weg:

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: none
```

## Parameter

Keine.

## Hinweise

- Wenn der `auth`-Abschnitt vollständig in der Konfiguration fehlt, entspricht dies `type: none`
- Es werden keine Autorisierungs-Header zu Anfragen hinzugefügt
