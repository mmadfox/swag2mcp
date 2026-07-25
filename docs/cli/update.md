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

=== "Update from same source"
    ```bash
    swag2mcp update abc123...
    ```

=== "Update with new URL"
    ```bash
    swag2mcp update abc123... --location https://new-api.example.com/openapi.json
    ```

## What Happens

1. Fresh spec download
2. Re-index all endpoints
3. Cache update
