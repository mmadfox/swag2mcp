# API Key

## Purpose

Authentication via an API key. The key can be sent as an HTTP header or as a URL query parameter.

## When to use

- Services that use API keys
- Weather services, geodata, translation APIs
- When the API expects a key in a header (`X-API-Key`) or query parameter (`?api_key=...`)

## Configuration

### Key in header

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

### Key in query parameter

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

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `key` | Yes | Name of the header or query parameter |
| `in` | Yes | Where to place the key: `header` or `query` |
| `value` | Yes | The key value |

## Notes

- In `header` mode, the key is added as an HTTP header
- In `query` mode, the key is added as a URL parameter
- Store the value in an environment variable: `value: "$(MY_API_KEY)"`
