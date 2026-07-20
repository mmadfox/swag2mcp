# update

Update a spec.

## Syntax

```bash
swag2mcp update <id> [flags]
```

## Arguments

| Argument | Description |
|----------|-------------|
| `id` | Spec ID to update |

## Flags

| Flag | Description |
|------|-------------|
| `-l, --location` | New URL or path |

## Usage

::: code-group

```bash [Update from same source]
swag2mcp update abc123...
```

```bash [Update with new URL]
swag2mcp update abc123... --location https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/pokeapi.yaml
```

:::

## What Happens

1. Fresh spec download
2. Re-index all endpoints
3. Cache update
