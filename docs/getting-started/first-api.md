# Adding Your First API

## Step-by-Step Guide

### Step 1: Find an OpenAPI Spec

swag2mcp supports three formats:

- **OpenAPI 3.x** (JSON/YAML)
- **Swagger 2.0** (JSON/YAML)
- **Postman Collection v2.1** (JSON)

Sample specs are in the `specs/` directory of the repository.

### Step 2: Add the Spec

=== "From URL"
    ```bash
    swag2mcp add https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    ```

=== "From local file"
    ```bash
    swag2mcp add ./specs/dadjoke.yaml
    ```

=== "Interactive"
    ```bash
    swag2mcp add
    ```

### Step 3: Verify

```bash
swag2mcp ls
```

Output:

```
Specifications:
  dadjoke (https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml)
  └── jokes (3 endpoints)
```

### Step 4: Start MCP Server

```bash
swag2mcp mcp
```

### Step 5: Explore via LLM

Your LLM agent can now:

- **Find endpoints**: "Show all endpoints in dadjoke"
- **Inspect details**: "What does GET / random-joke return?"
- **Invoke API**: "Get a random dad joke"

## Popular API Examples

| API | Command |
|-----|---------|
| Open-Meteo (weather) | `swag2mcp add https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml` |
| PokeAPI | `swag2mcp add https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/pokeapi.yaml` |
| Dad Jokes | `swag2mcp add https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml` |
| Rick & Morty | `swag2mcp add https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/rick-and-morty.json` |
| Binance | `swag2mcp add https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml` |

!!! tip
    Sample specs are in the `specs/` directory of the repository.
