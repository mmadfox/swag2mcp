# export

Export a workspace.

## Syntax

```bash
swag2mcp export [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `-o, --output` | ZIP file output path |

## Usage

::: code-group

```bash [Default]
swag2mcp export
```
Saves to `swag2mcp-export.zip`.

```bash [Custom path]
swag2mcp export --output ~/backups/swag2mcp-2024-01-01.zip
```

:::

## What's Exported

- Configuration file
- Specs
- Cache
- Auth scripts

## Restore

```bash
swag2mcp import --from-zip backup.zip
```
