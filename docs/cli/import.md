# import

## Purpose

Import spec files into the workspace `specs/` directory for local use, or restore a full workspace from a ZIP backup. Three modes cover different scenarios: adding a single spec, bulk-importing from existing config, or restoring a complete workspace.

## When to use

- You have a spec URL or file and want to save it locally in the workspace
- You want to download all collection spec files from the config and make the workspace self-contained
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
| `name` | 3 | Varies | Filename to save as (e.g. `example-api.yaml`). Derived from URL if omitted. |

## Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--spec` | `-s` | `string` | `""` | Download collection spec files from the config. Use without value for all specs, or specify domains like `--spec meteo,github` |
| `--force` | `-f` | `bool` | `false` | Overwrite existing spec files without error |
| `--from-zip` | | `string` | `""` | Restore workspace from a swag2mcp backup ZIP |

## How it works

### Mode 1 — Single import from URL or file

Download a spec file and save it to `specs/`:

```bash
swag2mcp import https://example.com/spec.yaml example-api.yaml
swag2mcp import /path/to/workspace https://example.com/spec.yaml example-api.yaml
swag2mcp import ./local-spec.yaml example-api.yaml
```

If `name` is omitted, it is derived from the URL filename:
```bash
swag2mcp import https://example.com/specs/petstore.yaml
# → saved as petstore.yaml
```

Overwrite an existing file with `--force`:
```bash
swag2mcp import --force https://example.com/spec.yaml example-api.yaml
```

After import, the output shows the workspace path, the saved file, and a YAML template to add to `swag2mcp.yaml`:

```
✅ Imported to /path/to/workspace
   specs/example-api.yaml

   Add to swag2mcp.yaml:
     specs:
       - domain: <your-domain>
         collections:
           - location: specs/example-api.yaml
```

### Mode 2 — Bulk import from existing config (`--spec`)

Download all collection spec files for the specified domains from their configured `location` URLs, save them to `specs/`, and update the config to point to the local copies:

```bash
swag2mcp import --spec                # all specs
swag2mcp import --spec meteo           # specific spec
swag2mcp import --spec meteo,github    # multiple specs
swag2mcp import /path/to/workspace --spec meteo
```

If a specified domain does not exist in the config, the command returns an error:
```
Error: import_no_match
  Spec "nonexistent" not found in config.
```

This makes the workspace self-contained — no remote spec URLs are needed after import.

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
- **HTTP client:** The global HTTP client settings from the config are applied during import (timeout, proxy, headers, etc.).
