# import

## Purpose

Import spec files into the workspace or restore a full workspace from a ZIP backup. Three modes cover different scenarios: adding a single spec, bulk-importing from existing config, or restoring a complete workspace.

## When to use

- You have a spec URL or file and want to add it to the workspace
- You want to download all spec files referenced in the config
- You need to restore a workspace from a ZIP backup created by `export`
- You are migrating swag2mcp to another machine

## Syntax

```bash
swag2mcp import [path] [source] [name] [flags]
```

## Arguments

| Argument | Position | Required | Description |
|----------|----------|----------|-------------|
| `path` | 1 | No | Workspace directory. If omitted, resolves via path resolution rules. |
| `source` | 2 | Varies | URL or local path to a spec file, or path to a ZIP archive |
| `name` | 3 | Varies | Domain name for the new spec |

## Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--spec` | `-s` | `stringSlice` | `nil` | Import collections from specified specs (comma-separated) |
| `--from-zip` | | `string` | `""` | Restore workspace from a swag2mcp backup ZIP |

## How it works

### Mode 1 — Single import from URL or file

Download a spec file and add it to the workspace with a domain name:

```bash
swag2mcp import https://example.com/spec.yaml myspec
swag2mcp import /path/to/workspace https://example.com/spec.yaml myspec
swag2mcp import ./local-spec.yaml myspec
```

The spec file is saved to `specs/` and the config is updated with the new spec entry.

### Mode 2 — Bulk import from existing config

Download all collections for the specified domains from their configured URLs:

```bash
swag2mcp import --spec meteo
swag2mcp import /path/to/workspace --spec meteo,store
```

Each collection's spec file is downloaded and saved to `specs/`. The config is updated to point to the local copies.

### Mode 3 — Restore from ZIP backup

Restore a full workspace from a ZIP archive created by `swag2mcp export`:

```bash
swag2mcp import --from-zip /path/to/backup.zip
swag2mcp import /path/to/workspace /path/to/backup.zip
```

> **The ZIP must be created by `swag2mcp export`.** Arbitrary ZIP files will not work — the archive has a specific internal structure (`swag2mcp.yaml`, `specs/`, `auth_scripts/`).

## Post-command verification

```bash
# Single or bulk import
swag2mcp ls [path]
# The new spec should appear in the list

# ZIP restore
swag2mcp ls [path]
# All specs from the backup should appear
```

## Nuances

- **Bulk mode requires config:** When using `--spec`, the config file must exist. Run `init` first if needed.
- **Single import creates workspace:** If the workspace doesn't exist, it is created automatically.
- **ZIP detection:** A positional argument ending in `.zip` is treated as a ZIP source. The `--from-zip` flag takes priority over positional detection.
- **`--force`:** Available for ZIP restore to overwrite an existing workspace.
- **HTTP client:** The global HTTP client settings from the config are applied during import (timeout, proxy, headers, etc.).
