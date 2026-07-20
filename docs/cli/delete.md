# delete

Delete a spec or collection.

## Syntax

```bash
swag2mcp delete <id> [flags]
```

## Arguments

| Argument | Description |
|----------|-------------|
| `id` | Spec ID or `spec_id/collection_id` |

## Usage

=== "Delete spec"
    ```bash
    swag2mcp delete abc123def456...
    ```

=== "Delete collection"
    ```bash
    swag2mcp delete abc123def456.../collection123...
    ```

=== "Interactive"
    ```bash
    swag2mcp delete
    ```

## Finding IDs

```bash
# Find spec ID
swag2mcp ls

# Find collection ID
swag2mcp ls --spec <spec_id>
```
