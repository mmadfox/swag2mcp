# Collections

A collection is a single OpenAPI/Swagger/Postman file that describes a specific API. It points to a `location` (URL or local file path) and belongs to a spec (domain).

One spec can have multiple collections — for example, the "meteo" spec might have "Forecast", "Air Quality", and "Marine" collections, each pointing to a different spec file.

## How Collections Are Parsed

When swag2mcp loads a collection file, it parses the OpenAPI/Swagger/Postman structure:

::: code-group

```text [OpenAPI 3.x]
Each top-level `tag` becomes a category of endpoints.
```

```text [Swagger 2.0]
Each tag from the tags list becomes a category of endpoints.
```

```text [Postman]
Each top-level folder becomes a category of endpoints.
```

:::

## Example

From the Dad Joke collection file:

```yaml
tags:
  - name: jokes
    description: Everything about dad jokes
```

Parsed category: `jokes`.

## Multiple Collections from One Spec

In YAML config, you can add the same spec file under different domains with different base URLs:

```yaml
specs:
  - domain: meteo-forecast
    llm_title: Open-Meteo Forecast
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml

  - domain: meteo-air-quality
    llm_title: Open-Meteo Air Quality
    base_url: https://air-quality-api.open-meteo.com
    collections:
      - llm_title: Air Quality
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml

  - domain: meteo-marine
    llm_title: Open-Meteo Marine
    base_url: https://marine-api.open-meteo.com
    collections:
      - llm_title: Marine
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

## Management

```bash
# List collections in a spec
swag2mcp ls --spec <spec_id>

# Delete a collection
swag2mcp delete <spec_id>/<collection_id>
```
