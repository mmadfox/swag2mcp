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

::: code-group

```bash [From URL]
swag2mcp add https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

```bash [From local file]
swag2mcp add ./specs/my-api.yaml
```

```bash [With options]
swag2mcp add https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml \
  --name "forecast" \
  --tags "weather"
```

```bash [Interactive]
swag2mcp add
```

:::

## Result

The spec is added to the config file and cached.
