# update

## Purpose

Re-validate the configuration, clear the cache, and re-download all spec files. This is a **full refresh** of the workspace — it ensures all cached specs are up to date and the index is rebuilt.

## When to use

- Remote spec files have changed and you want the latest version
- After editing `swag2mcp.yaml` to add or change spec locations
- When troubleshooting stale or corrupted cache
- Before running `mcp` to ensure everything is fresh

## Syntax

```bash
swag2mcp update [path]
```

## Arguments

| Argument | Position | Required | Description |
|----------|----------|----------|-------------|
| `path` | 1 | No | Workspace directory. If omitted, resolves via path resolution rules. |

## Flags

None.

## How it works

The `update` command runs a pipeline of operations:

1. **Load config** — reads `swag2mcp.yaml` from the workspace
2. **Validate** — runs the same checks as `validate` (YAML syntax, structure, spec file reachability, format, auth, HTTP client)
3. **Clean** — removes all contents of `cache/` and `responses/`
4. **Re-cache** — downloads all remote spec files and copies local spec files into the cache
5. **Re-index** — rebuilds the full-text search index for all endpoints
6. **Auth scripts** — creates stub auth scripts for specs using `ScriptAuth`
7. **Orphan cleanup** — removes auth scripts for specs that no longer exist

```bash
swag2mcp update
swag2mcp update ./my-workspace
```

## What happens to disabled collections

Collections with `disable: true` are skipped entirely — they are not cached or indexed.

## Post-command verification

```bash
swag2mcp ls [path]
# All specs should still be listed and reachable
```

## Nuances

- **No auto-init:** If the config file does not exist, `update` returns an error: `"configuration not found at <path>"`. Run `init` first.
- **Network dependency:** All remote spec URLs must be reachable. If any download fails, the entire update fails with a clear error message.
- **Auth script creation:** If a spec uses `ScriptAuth` and the stub script doesn't exist, `update` creates it. If creation fails, the update fails.
- **`update` vs `clean`:** `clean` only removes cache. `update` removes cache **and** re-downloads everything. Use `clean` when you just want to free space; use `update` when you want to refresh.
