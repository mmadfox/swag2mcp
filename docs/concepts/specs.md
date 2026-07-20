# Specs

A spec is a logical container representing an API domain or service (e.g., YouTube, Binance, Open-Meteo). Each spec has a unique `domain`, a `base_url`, optional `auth`, and contains one or more collections.

Collections point to OpenAPI/Swagger/Postman files — the spec itself is not a file, it's the grouping around them.

## Supported Formats

Collections support three file formats:

| Format | Extensions | Versions |
|--------|------------|----------|
| OpenAPI 3.x | `.json`, `.yaml`, `.yml` | 3.0.0 – 3.1.1 |
| Swagger 2.0 | `.json`, `.yaml`, `.yml` | 2.0 |
| Postman Collection | `.json` | v2.1 |

## Sources

Collection files can be:

- **URL**: `https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml`
- **Local file**: `./specs/my-api.yaml`

## Identification

Each spec gets a unique MD5 hash based on its domain:

```go
id = md5(domain)
```

## Management

```bash
# Add a spec
swag2mcp add https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml

# List all specs
swag2mcp ls

# Delete a spec
swag2mcp delete <id>

# Update a spec
swag2mcp update <id>
```

## Spec Configuration

In YAML config:

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    disable: false
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```
