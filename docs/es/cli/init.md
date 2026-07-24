# init

## Purpose

The `init` command creates a **workspace** — a directory with a `swag2mcp.yaml` config file and subdirectories for cache, specs, responses, and auth scripts. This is the first command to run when setting up swag2mcp.

## When to use

- You are setting up swag2mcp for the first time
- You want to create a new workspace in a specific directory
- You need to re-initialize a corrupted or missing workspace

## Syntax

```bash
swag2mcp init [path] [flags]
```

## Arguments

| Argument | Position | Required | Description |
|----------|----------|----------|-------------|
| `path` | 1 | No | Workspace directory. If omitted, defaults to `~/.swag2mcp`. |

## Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--interactive` | `-i` | `bool` | `false` | Run the interactive TUI wizard |
| `--force` | `-f` | `bool` | `false` | Overwrite existing configuration in a non-empty directory |

## How it works

### Non-interactive mode (default)

Creates a minimal `swag2mcp.yaml` with no specs. You edit the file manually afterwards.

```bash
swag2mcp init
# Creates ~/.swag2mcp/swag2mcp.yaml

swag2mcp init ./my-project
# Creates ./my-project/swag2mcp.yaml

swag2mcp init /absolute/path
# Creates /absolute/path/swag2mcp.yaml
```

### Interactive mode (`-i`)

Launches an 18-step TUI wizard that guides you through:

1. Choosing the workspace directory
2. Adding specs with domain, title, base URL
3. Configuring collections with location URLs
4. Setting up authentication (all 9 methods)
5. Configuring HTTP client settings (timeout, proxy, headers, etc.)

```bash
swag2mcp init -i
```

### Force mode (`--force`)

By default, `init` refuses to run in a non-empty directory. Use `--force` to overwrite:

```bash
swag2mcp init -f
swag2mcp init ./existing-dir -f
```

## What gets created

```
~/.swag2mcp/
├── swag2mcp.yaml       # Configuration file
├── cache/               # Downloaded remote spec files
├── specs/               # Local spec files
├── responses/           # Saved API invocation responses
└── auth_scripts/        # Authentication scripts (for ScriptAuth type)
```

## Post-command verification

```bash
ls ~/.swag2mcp/swag2mcp.yaml
# If the file exists, init succeeded
```

## Nuances

- **Path resolution:** `[path]` is a **workspace directory**, not a file path. The CLI appends `swag2mcp.yaml` automatically. Resolution order: explicit `[path]` → current directory (`./`) → `~/.swag2mcp/`.
- **Non-empty directory check:** Without `--force`, `init` returns an error if the target directory exists and is not empty. This prevents accidental overwrites.
- **Auth script stubs:** If any spec uses `ScriptAuth`, `init` creates stub script files (`.sh` on Unix, `.bat` on Windows) in `auth_scripts/`.
- **Output:** On success, prints the config path and a hint: `"Next step: edit swag2mcp.yaml or run 'swag2mcp ls' to list configured specs"`.
