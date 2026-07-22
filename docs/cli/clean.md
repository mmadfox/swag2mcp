# clean

## Purpose

Remove cached remote specs and saved API invocation responses. This frees up disk space and forces a fresh download of spec files on the next `update` or `mcp` start.

## When to use

- Spec files have changed on the remote server and you want to force a refresh
- You want to free up disk space
- You are troubleshooting stale cache issues
- Before running `update` to ensure a clean re-cache

## Syntax

```bash
swag2mcp clean [path]
```

## Arguments

| Argument | Position | Required | Description |
|----------|----------|----------|-------------|
| `path` | 1 | No | Workspace directory. If omitted, resolves via path resolution rules. |

## Flags

None.

## How it works

```bash
swag2mcp clean
swag2mcp clean ./my-workspace
```

## What is cleaned

| Directory | Contents | Why |
|-----------|----------|-----|
| `cache/` | Downloaded remote spec files | Forces re-download on next access |
| `responses/` | Saved API invocation responses | Frees disk space |

## What is preserved

| Directory | Contents | Why |
|-----------|----------|-----|
| `specs/` | Local spec files | These are your source files, not cache |
| `auth_scripts/` | Authentication scripts | These are user-created, not cache |

## Orphan auth script cleanup

After cleaning, `clean` also removes auth scripts for specs that no longer exist in the configuration. This prevents stale scripts from accumulating.

## Automatic cleanup

When the MCP server starts (`swag2mcp mcp`), responses older than 48 hours are removed automatically. You typically don't need to run `clean` manually for routine maintenance.

## Post-command verification

```bash
ls ~/.swag2mcp/cache
# Should be empty (directory exists but has no files)
```

## Nuances

- **No config required:** `clean` works even without a valid config file. It simply removes the cache and responses directories.
- **Orphan cleanup is best-effort:** If the config file is corrupted or unreadable, orphan auth script cleanup is skipped (not fatal).
- **Directories are preserved:** The `cache/` and `responses/` directories themselves are kept — only their contents are removed.
