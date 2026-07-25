# Clé API

## Objectif

Authentification via une clé API. La clé peut être envoyée comme en-tête HTTP ou comme paramètre de requête URL.

## Quand l'utiliser

- Services qui utilisent des clés API
- Services météo, géodonnées, API de traduction
- Quand l'API attend une clé dans un en-tête (`X-API-Key`) ou un paramètre de requête (`?api_key=...`)

## Configuration

### Clé dans l'en-tête

```yaml
specs:
  - domain: jokes
    llm_title: API Dad Joke
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Blagues
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: api-key
      config:
        key: "X-API-Key"
        in: header
        value: "$(CLE_API)"
```

### Clé dans le paramètre de requête

```yaml
specs:
  - domain: jokes
    llm_title: API Dad Joke
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Blagues
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: api-key
      config:
        key: "api_key"
        in: query
        value: "$(CLE_API)"
```

## Paramètres

| Paramètre | Requis | Description |
|-----------|--------|-------------|
| `key` | Oui | Nom de l'en-tête ou du paramètre de requête |
| `in` | Oui | Où placer la clé : `header` ou `query` |
| `value` | Oui | La valeur de la clé |

## Notes

- En mode `header`, la clé est ajoutée comme en-tête HTTP
- En mode `query`, la clé est ajoutée comme paramètre URL
- Stockez la valeur dans une variable d'environnement : `value: "$(MA_CLE_API)"`
