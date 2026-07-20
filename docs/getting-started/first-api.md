# Adding Your First API

## Step-by-Step Guide

### Step 1: Find an OpenAPI Spec

swag2mcp supports three formats:

- **OpenAPI 3.x** (JSON/YAML)
- **Swagger 2.0** (JSON/YAML)
- **Postman Collection v2.1** (JSON)

### Step 2: Add the Spec

=== "From URL"
    ```bash
    swag2mcp add https://petstore.swagger.io/v2/swagger.json
    ```

=== "From local file"
    ```bash
    swag2mcp add ./specs/petstore.yaml
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
  petstore (https://petstore.swagger.io/v2/swagger.json)
  └── pet (3 endpoints)
  └── store (4 endpoints)
  └── user (8 endpoints)
```

### Step 4: Start MCP Server

```bash
swag2mcp mcp
```

### Step 5: Explore via LLM

Your LLM agent can now:

- **Find endpoints**: "Show all endpoints in petstore"
- **Inspect details**: "What does POST /pet accept?"
- **Invoke API**: "Create a new pet named Rex"

## Popular API Examples

| API | Command |
|-----|---------|
| Open-Meteo (weather) | `swag2mcp add https://.../meteo/forecast.yml` |
| PokeAPI | `swag2mcp add https://.../pokeapi.yaml` |
| Dad Jokes | `swag2mcp add https://.../dadjoke.yaml` |
| Rick & Morty | `swag2mcp add https://.../rick-and-morty.json` |
| Binance | `swag2mcp add https://.../binance.yaml` |

!!! tip
    Sample specs are in the `specs/` directory of the repository.
