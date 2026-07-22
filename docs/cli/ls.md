# ls

## Purpose

List all configured **specs** and their **collections** in a human-readable format. This is the primary way to inspect what APIs are available in your workspace.

## When to use

- You want to see what APIs are configured
- You need to find a spec or collection ID
- You want to check how many endpoints each collection has
- You want to filter specs by tags

## Syntax

```bash
swag2mcp ls [path] [flags]
```

## Arguments

| Argument | Position | Required | Description |
|----------|----------|----------|-------------|
| `path` | 1 | No | Workspace directory. If omitted, resolves via path resolution rules. |

## Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--tags` | `-t` | `string` | `""` | Filter specs by tags (comma-separated) |

## How it works

### List all specs

Shows every spec with its domain, collections, and endpoint counts:

```bash
swag2mcp ls
```

Example output:

```
Specifications:
  dadjoke (https://icanhazdadjoke.com)
    jokes (3 endpoints)
  meteo (https://meteo.swagger.io/v2)
    forecast (5 endpoints)
    current (8 endpoints)
  binance (https://api.binance.com)
    market-data (12 endpoints)
```

### Filter by tags

Show only specs that have the specified tags:

```bash
swag2mcp ls --tags=public
swag2mcp ls --tags=public,internal
```

## Post-command verification

Use `ls` after `add`, `delete`, `update`, or `import` to confirm the workspace state matches your expectations.

## Nuances

- **Auto-init:** If no config file exists, `ls` automatically runs the init wizard first.
- **Tag filtering:** Tags are comma-separated. Specs matching **any** of the specified tags are shown (OR logic).
- **Output format:** The output is plain text, not JSON. For machine-readable output, use `info`.
