# Export and Import

## Overview

swag2mcp supports full workspace round-trip via ZIP archives. You can export your entire workspace (config, spec files, auth scripts) to a ZIP file and restore it on another machine.

## Export

Creates a portable ZIP backup of your workspace.

```bash
# Export to default file (swag2mcp-backup-<timestamp>.zip)
swag2mcp export

# Export with custom path
swag2mcp export --output ~/backups/swag2mcp-backup.zip

# Export only specific specs
swag2mcp export --spec meteo
swag2mcp export --spec meteo,store
```

### What's included in the export

| Item | Description |
|------|-------------|
| `swag2mcp.yaml` | Configuration file |
| `specs/` | All spec files (OpenAPI/Swagger/Postman) |
| `auth_scripts/` | Authentication scripts |
| `swag2mcp.meta` | Metadata (version info for compatibility) |

Cache and responses are **not** exported — they are transient and would be stale on restore.

### Default filename

If you don't specify an output path, the file is saved as `swag2mcp-backup-<YYYY-MM-DD-HHMMSS>.zip` in the current directory (UTC timestamp).

## Import

Restore a workspace from a ZIP backup or import spec files.

### Restore from ZIP

```bash
# Restore full workspace
swag2mcp import --from-zip /path/to/backup.zip

# Restore with overwrite
swag2mcp import --from-zip /path/to/backup.zip -f
```

The ZIP must be created by `swag2mcp export` — arbitrary ZIP files will not work.

### Import a single spec file

Download a spec file and add it to the workspace:

```bash
swag2mcp import https://example.com/spec.yaml myspec
swag2mcp import /path/to/workspace https://example.com/spec.yaml myspec
```

### Bulk import from existing config

Download all collection spec files for the specified specs (domains):

```bash
swag2mcp import --spec meteo
swag2mcp import /path/to/workspace --spec meteo,store
```

This downloads each collection's spec file, saves it to `specs/`, and updates the config to point to the local copy.

## Use cases

### Backup

```bash
swag2mcp export --output swag2mcp-$(date +%Y-%m-%d).zip
```

### Transfer to another machine

```bash
# On old machine
swag2mcp export --output swag2mcp.zip

# Copy the ZIP to the new machine, then:
swag2mcp import --from-zip swag2mcp.zip
```

### Share configuration

```bash
swag2mcp init
swag2mcp export --output template.zip
# Share template.zip with a colleague
```

## Post-export verification

Always verify the ZIP file was created:

```bash
ls -la swag2mcp-backup-*.zip
```

## Important notes

- **The output must be a file path ending in `.zip`** — do not pass a directory
- **Cache and responses are excluded** — only the config, specs, and auth scripts are preserved
- **The ZIP is self-contained** — it can be restored on any machine with swag2mcp installed
- **Spec filter** — use `--spec` to export or import only specific specs
