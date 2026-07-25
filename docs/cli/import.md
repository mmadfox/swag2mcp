# import

Import a workspace.

## Syntax

```bash
swag2mcp import [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `--from-zip` | ZIP file path to import |
| `-f, --force` | Overwrite existing config |

## Usage

=== "From ZIP"
    ```bash
    swag2mcp import --from-zip workspace.zip
    ```

=== "With overwrite"
    ```bash
    swag2mcp import --from-zip workspace.zip -f
    ```

## What's Imported

- Configuration file
- Specs
- Cache
- Auth scripts
