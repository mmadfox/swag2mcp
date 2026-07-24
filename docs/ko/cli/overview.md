# CLI Commands

## Overview

The `swag2mcp` CLI is the single entry point for all operations — from initializing a workspace and managing API specifications to starting an MCP server for LLM integration. It provides **13 commands** that cover the full lifecycle of working with OpenAPI/Swagger/Postman specs.

### What the CLI solves

- **Workspace lifecycle** — create (`init`), inspect (`info`, `ls`), clean (`clean`), update (`update`), and remove (`delete`) workspaces and their contents
- **Spec & collection management** — add (`add`), list (`ls`), and delete (`delete`) API specifications and their collections
- **Running modes** — start the MCP server for LLM tool access (`mcp`) or launch the interactive TUI explorer (`run`)
- **Diagnostics** — validate configuration (`validate`), show version (`version`), display runtime info (`info`)
- **Backup & restore** — full workspace round-trip via ZIP (`export`, `import`)

### Key nuances

- **Path resolution** — commands that accept `[path]` expect a **workspace directory** (not a file path). Resolution order: explicit `[path]` → current directory (`./`) → `~/.swag2mcp/`. The CLI appends `swag2mcp.yaml` automatically. Always pass an explicit path when running as a service or in IDE config to avoid loading the wrong workspace.
- **Spec vs Collection** — a **spec** represents a logical API service (e.g. "Open-Meteo API"), while a **collection** is one OpenAPI/Swagger/Postman file. A spec can have multiple collections.
- **`--version`** is supported both as a flag (`swag2mcp --version`) and as a subcommand (`swag2mcp version`).
- **`add spec` / `add collection`** accept YAML input via `--yaml` (inline string or `-` for stdin). Piping from a file or heredoc avoids shell quoting issues with special characters.
- **`delete`** requires a TTY (interactive terminal). There is no `--force` or `--yes` flag — it always prompts for selection and confirmation.
- **`mcp`** is the primary command for LLM integration. It supports three transports: `stdio` (default), `sse`, and `streamable-http`. The `--disable-llm-auth` flag (default: `true`) removes the `auth` tool from the MCP tool list, preventing the LLM from seeing or requesting tokens. Auth still works — tokens are obtained through the standard config mechanism, not via the LLM. This mode is recommended for **production** (LLM never has access to credentials). For **debugging** or when using short-lived tokens, set `--disable-llm-auth=false` to let the LLM request fresh tokens via the `auth` tool.
- **`validate`** checks YAML syntax, config structure, spec file existence, URL reachability, spec format (OpenAPI/Swagger/Postman), auth settings, and HTTP client correctness. It does **not** test authentication endpoints or API endpoint availability.
- **`export` / `import`** provide a full workspace round-trip — config file, spec files, cache, and auth scripts are all included in the ZIP archive.
- **`clean`** removes `cache/` and `responses/` directories but preserves `specs/` and `auth_scripts/`. Old responses (>48h) are also cleaned automatically on `mcp` startup.

## Commands

| Command | Description |
|---------|-------------|
| [`init`](/cli/init) | Initialize a workspace directory with default config |
| [`add`](/cli/add) | Add a spec or collection to the configuration |
| [`delete`](/cli/delete) | Delete a spec or collection interactively |
| [`ls`](/cli/ls) | List all specs and their collections |
| [`run`](/cli/run) | Launch the interactive TUI API explorer |
| [`validate`](/cli/validate) | Validate configuration and spec files |
| [`clean`](/cli/clean) | Clear cached specs and invocation responses |
| [`update`](/cli/update) | Re-validate, re-cache, and re-index all specs |
| [`mcp`](/cli/mcp) | Start the MCP server for LLM tool access |
| [`version`](/cli/version) | Print the swag2mcp version |
| [`info`](/cli/info) | Show detailed configuration and runtime information |
| [`import`](/cli/import) | Import spec files or restore workspace from ZIP |
| [`export`](/cli/export) | Export workspace as a portable ZIP backup |

## Global Flags

| Flag | Description |
|------|-------------|
| `--version` | Show version (same as `version` subcommand) |
| `--help` | Show help for any command |
