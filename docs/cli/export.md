# export

## Purpose

Create a portable ZIP backup of the workspace. The archive contains the configuration file, all spec files, and auth scripts — everything needed to restore the workspace on another machine.

## When to use

- You want to back up your workspace before making changes
- You are migrating swag2mcp to another machine
- You want to share your API configuration with a colleague
- You are preparing a reproducible environment

## Syntax

```bash
swag2mcp export [path] [output] [flags]
```

## Arguments

| Argument | Position | Required | Description |
|----------|----------|----------|-------------|
| `path` | 1 | No | Workspace directory. If omitted, resolves via path resolution rules. |
| `output` | 2 | No | Full path for the output ZIP file. If omitted, defaults to `./swag2mcp-backup-<timestamp>.zip`. |

## Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--spec` | `-s` | `stringSlice` | `nil` | Export only specified specs (comma-separated) |

## How it works

### Default export

Creates a ZIP in the current directory with a timestamped name:

```bash
swag2mcp export
# Creates ./swag2mcp-backup-2026-07-22-143022.zip
```

### Custom output path

```bash
swag2mcp export /path/to/workspace /path/to/backup.zip
```

### Export specific specs

```bash
swag2mcp export --spec meteo
swag2mcp export --spec meteo,store
```

## What's in the ZIP

| Entry | Description |
|-------|-------------|
| `swag2mcp.meta` | Metadata about the export |
| `swag2mcp.yaml` | Configuration file |
| `specs/` | All spec files (OpenAPI/Swagger/Postman) |
| `auth_scripts/` | Authentication scripts |
| `cache/` | Empty (cache is not exported) |
| `responses/` | Empty (responses are not exported) |

## Restore

Use `import` to restore from a backup:

```bash
swag2mcp import --from-zip /path/to/backup.zip
```

## Post-command verification

Always verify the ZIP file was created:

```bash
ls -la swag2mcp-backup-*.zip
# or for a custom output path:
ls -la /path/to/backup.zip
```

## Nuances

- **Output must be a file path:** The `[output]` argument must be a full file path ending in `.zip`. Do **not** pass a directory — the command will not create a ZIP if given a directory path.
- **Default filename:** `swag2mcp-backup-<YYYY-MM-DD-HHMMSS>.zip` using UTC timestamp.
- **`--spec` filter:** When set, only the specified specs are included. Other specs are excluded from the archive.
- **No config required:** `export` works even without a valid config file. It exports whatever exists in the workspace.
- **Cache and responses are excluded:** These are transient data that would be stale on restore. Only the config, specs, and auth scripts are preserved.
