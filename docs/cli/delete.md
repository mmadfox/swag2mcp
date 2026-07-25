# delete

## Purpose

Remove a **spec** (API service) or **collection** (spec file) from the configuration. This is the inverse of `add`.

## When to use

- An API is no longer needed
- You want to remove a specific spec file from a spec
- You are cleaning up your workspace

## Syntax

```bash
swag2mcp delete spec [path]
swag2mcp delete collection [path]
```

## Arguments

| Argument | Position | Required | Description |
|----------|----------|----------|-------------|
| `path` | 1 | No | Workspace directory. If omitted, resolves via path resolution rules. |

## Flags

None. Both subcommands are purely interactive.

## How it works

### Delete a spec

Prompts you to select a spec from a list, then asks for confirmation before deleting.

```bash
swag2mcp delete spec
```

### Delete a collection

Prompts you to select a spec, then a collection within that spec, then asks for confirmation.

```bash
swag2mcp delete collection
```

## Finding IDs

The interactive prompts show human-readable names, not IDs. If you need IDs for reference:

```bash
# List all specs with their IDs
swag2mcp ls

# List collections for a specific spec
swag2mcp ls --tags
```

## Post-command verification

```bash
swag2mcp ls [path]
# The deleted spec or collection should no longer appear
```

## Nuances

- **TTY required:** Both commands require an interactive terminal. They will **not** work in CI/CD pipelines, cron jobs, or non-interactive scripts.
- **No `--force` or `--yes`:** There is no way to skip the confirmation prompt. This is intentional to prevent accidental deletions.
- **Auto-init:** If no config file exists, `delete` automatically runs the init wizard first.
- **No YAML mode:** Unlike `add`, there is no `--yaml` flag. Deletion is always interactive.
