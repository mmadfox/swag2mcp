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

::: code-group

```bash [From ZIP]
swag2mcp import --from-zip workspace.zip
```

```bash [With overwrite]
swag2mcp import --from-zip workspace.zip -f
```

:::

## What's Imported

- Configuration file
- Specs
- Cache
- Auth scripts
