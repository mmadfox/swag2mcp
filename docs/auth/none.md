# None

## Purpose

No authentication required. The API is accessible without tokens or keys.

## When to use

- Public APIs (Open-Meteo, icanhazdadjoke, PokéAPI)
- Test and demo environments
- When the API does not require authorization

## Configuration

Set `type: none` or simply omit the `auth` section:

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

## Parameters

None.

## Notes

- If the `auth` section is completely absent from the config, it is equivalent to `type: none`
- No authorization headers are added to requests
