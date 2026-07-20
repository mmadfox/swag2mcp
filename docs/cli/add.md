# add

Add an API spec.

## Syntax

```bash
swag2mcp add [location] [flags]
```

## Arguments

| Argument | Description |
|----------|-------------|
| `location` | URL or path to spec file |

## Flags

| Flag | Description |
|------|-------------|
| `-n, --name` | Collection name |
| `-t, --tags` | Collection tags |

## Usage

=== "From URL"
    ```bash
    swag2mcp add https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    ```

=== "From local file"
    ```bash
    swag2mcp add ./specs/my-api.yaml
    ```

=== "With options"
    ```bash
    swag2mcp add https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml \
      --name "forecast" \
      --tags "weather"
    ```

=== "Interactive"
    ```bash
    swag2mcp add
    ```

## Result

The spec is added to the config file and cached.
