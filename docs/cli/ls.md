# ls

List specs and collections.

## Syntax

```bash
swag2mcp ls [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `-s, --spec` | Spec ID to view collections |

## Usage

::: code-group

```bash [All specs]
swag2mcp ls
```
```
Specifications:
  dadjoke (https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml)
    jokes (3 endpoints)
  meteo (https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml)
    forecast (5 endpoints)
```

```bash [Spec collections]
swag2mcp ls --spec abc123...
```
```
Collections for dadjoke:
  jokes (3 endpoints, 1 tag)
```

:::
