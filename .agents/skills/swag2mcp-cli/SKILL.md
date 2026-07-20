---
name: swag2mcp-cli
description: |
  Complete CLI reference for swag2mcp commands, flags, and config.
  Use when the user asks how to use swag2mcp from the terminal,
  needs help with a specific command, or wants to understand
  the configuration file structure.
license: MIT
metadata:
  author: mmadfox
  version: "2.0.0"
---

# swag2mcp-cli — CLI Reference Skill

This document describes every CLI command in swag2mcp. Use it when you need to tell a user how to interact with the tool from the terminal.

## When this skill activates

This skill activates when the user asks to:

- Set up swag2mcp or initialize a workspace
- Connect an API to an AI agent / LLM
- Add, configure, or manage API specifications (OpenAPI/Swagger)
- Start an MCP server for API access
- Explore, search, or invoke API endpoints through MCP tools
- Work with Open-Meteo, Binance, PokéAPI, icanhazdadjoke, or similar APIs

### Example user requests that trigger this skill

| User says | What the skill does |
|-----------|-------------------|
| "Set up swag2mcp" | Downloads and installs swag2mcp, runs `swag2mcp init .` |
| "Initialize a workspace for my APIs" | Runs `swag2mcp init .` in current directory |
| "List my configured APIs" | Runs `swag2mcp ls` |
| "Start the MCP server for my specs" | Runs `swag2mcp mcp` |
| "Show me what APIs are available" | Calls `spec_list` MCP tool |
| "Find an endpoint to get BTC price" | Calls `search` MCP tool |
| "Call the API to get a random dad joke" | Calls `invoke` MCP tool |
| "How do I add authentication to my API?" | Guides through auth config in YAML |
| "Export my workspace as a backup" | Runs `swag2mcp export` |

---

## Workspace creation rules

When the user asks to create a project or workspace, create it in a `.swag2mcp` subdirectory of the current folder:

| User says | Command |
|-----------|---------|
| "Create a project in the current folder" | `mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp` |
| "Create a project called my-api" | `mkdir -p my-api/.swag2mcp && swag2mcp init ./my-api/.swag2mcp` |
| "Set up swag2mcp" (no path) | `mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp` |
| "Initialize a workspace for my APIs" | `mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp` |

Always create the workspace in a `.swag2mcp` subdirectory of the user's current directory unless they specify a custom path. This creates `swag2mcp.yaml`, `cache/`, `specs/`, `responses/`, and `auth_scripts/` inside `.swag2mcp/`.

---

## Installation

### Option 1 — Download from GitHub Releases (recommended)

1. Open https://github.com/mmadfox/swag2mcp/releases/latest
2. Find the archive for the user's system:

   | OS | Architecture | Archive |
   |----|-------------|---------|
   | Linux | x86_64 | `swag2mcp_<version>_linux_amd64.tar.gz` |
   | Linux | ARM64 | `swag2mcp_<version>_linux_arm64.tar.gz` |
   | macOS | Intel | `swag2mcp_<version>_darwin_amd64.tar.gz` |
   | macOS | Apple Silicon | `swag2mcp_<version>_darwin_arm64.tar.gz` |
   | Windows | x86_64 | `swag2mcp_<version>_windows_amd64.zip` |

3. Download and install:

   **Linux / macOS:**
   ```sh
   tar -xzf swag2mcp_<version>_<os>_<arch>.tar.gz
   sudo mv swag2mcp /usr/local/bin/
   swag2mcp --version
   ```

   **Windows (PowerShell):**
   ```powershell
   Expand-Archive swag2mcp_<version>_windows_amd64.zip -DestinationPath .
   move swag2mcp.exe C:\Windows\System32\
   swag2mcp --version
   ```

4. (Optional) Repeat for mock server — download `swag2mcp-mock_<version>_<os>_<arch>.tar.gz`

### Option 2 — Install with Go

If Go is installed:

```sh
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

---

## Quick Start

```sh
# Initialize a workspace
swag2mcp init ~/my-api

# List configured specs
swag2mcp ls
```

---

## Concepts: Spec vs Collection

Before writing config, understand the two levels:

- **Spec** = one API (a domain). Represents a logical service (e.g. "Open-Meteo API", "Binance Market Data"). A spec can have multiple collections.
- **Collection** = one OpenAPI/Swagger file. If an API has multiple spec files (different versions, different microservices), each file is a separate collection under the same spec.

**Important:** The `location` field in a collection must point to an **OpenAPI 3.x, Swagger 2.0, or Postman collection** file — not to the API endpoint itself. For example:
- ✅ `location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo.json` — correct, this is an OpenAPI spec
- ❌ `location: https://api.example.com/v1/users` — wrong, this is a JSON response, not a spec

---

## 1. `swag2mcp init [path]`

Initialize a workspace directory with the default config file.

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--interactive` | `-i` | bool | `false` | Run interactive Bubbletea TUI wizard |
| `--force` | `-f` | bool | `false` | Overwrite existing configuration in non-empty directory |

**Examples:**
```sh
swag2mcp init                    # create ~/.swag2mcp/swag2mcp.yaml
swag2mcp init ./                 # create ./swag2mcp.yaml (in current directory!)
swag2mcp init path/to            # create path/to/swag2mcp.yaml
swag2mcp init -i                 # interactive wizard
swag2mcp init -f                 # force overwrite
```

> **Note:** `swag2mcp init ./` creates `swag2mcp.yaml` directly in the current directory, NOT inside a `.swag2mcp/` subdirectory. For the recommended layout (`.swag2mcp/` subdirectory), use the workspace creation rules in the "Workspace creation rules" section above: `mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp`.

**Behavior:** Creates `cache/`, `specs/`, `responses/`, `auth_scripts/` subdirectories and a `swag2mcp.yaml` config file. After init, prints a hint: `"Next step: edit swag2mcp.yaml or run 'swag2mcp ls' to list configured specs"`.

---

## 2. `swag2mcp add spec [path]`

Add a new API specification to the config.

> **`[path]` is the workspace directory** containing `swag2mcp.yaml`, NOT the path to a spec file or URL. If omitted, swag2mcp resolves the workspace using the path resolution rules (see "Path Resolution" section below).

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--yaml` | `-y` | string | `""` | YAML input inline or `-` for stdin |
| `--example` | `-e` | bool | `false` | Print YAML template and exit |

**Examples:**
```sh
swag2mcp add spec --example

# Inline YAML (simple spec, no special chars in values)
swag2mcp add spec --yaml 'domain: meteo
llm_title: Open-Meteo API
base_url: https://meteo.swagger.io/v2
collections:
  - llm_title: Pets
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo.json'

# Safer: pipe from file or heredoc (avoids shell quoting issues with colons, &, #)
cat spec.yaml | swag2mcp add spec --yaml -

# Heredoc — recommended for complex YAML with special characters
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My API
llm_instruction: "Use this API for X & Y # important"
base_url: https://api.example.com/v1
collections:
  - llm_title: Main
    location: https://raw.githubusercontent.com/org/repo/main/spec.yaml
EOF
```

**YAML format:**
```yaml
domain: meteo
llm_title: Open-Meteo API
llm_instruction: Use this API to manage pets.
base_url: https://meteo.swagger.io/v2
tags: [public, demo]
auth:
  type: bearer
  config:
    token: $(TOKEN)
collections:
  - llm_title: Open-Meteo Swagger
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo.json
```

---

## 3. `swag2mcp add collection [path]`

Add a new collection to an existing spec.

> **`[path]` is the workspace directory** containing `swag2mcp.yaml`, NOT the path to a collection spec file. If omitted, swag2mcp resolves the workspace using the path resolution rules (see "Path Resolution" section below).

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--yaml` | `-y` | string | `""` | YAML input inline or `-` for stdin |
| `--example` | `-e` | bool | `false` | Print YAML template and exit |

**Examples:**
```sh
swag2mcp add collection --example

# Inline YAML
swag2mcp add collection --yaml 'spec_domain: meteo
llm_title: Orders Collection
location: https://meteo.example.com/orders.json'

# Safer: pipe from file or heredoc
cat collection.yaml | swag2mcp add collection --yaml -

swag2mcp add collection --yaml - <<EOF
spec_domain: meteo
llm_title: Orders Collection
location: https://meteo.example.com/orders.json
EOF
```

**YAML format:**
```yaml
spec_domain: meteo
llm_title: Orders Collection
location: https://meteo.example.com/orders.json
```

---

## 4. `swag2mcp delete spec [path]`

Delete a specification interactively. No flags. Prompts for selection and confirmation.

> **Requires a TTY (interactive terminal).** This command will not work in CI/CD pipelines, cron jobs, or non-interactive scripts because it requires user input for selection and confirmation. There is no `--force` or `--yes` flag to skip prompts.

```sh
swag2mcp delete spec
swag2mcp delete spec ./
```

---

## 5. `swag2mcp delete collection [path]`

Delete a collection interactively. No flags. Prompts for spec selection, collection selection, and confirmation.

> **Requires a TTY (interactive terminal).** This command will not work in CI/CD pipelines, cron jobs, or non-interactive scripts because it requires user input for selection and confirmation. There is no `--force` or `--yes` flag to skip prompts.

```sh
swag2mcp delete collection
swag2mcp delete collection ./
```

---

## 6. `swag2mcp ls [path]`

List all specs and their collections.

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--tags` | `-t` | string | `""` | Filter by tags (comma-separated) |

```sh
swag2mcp ls
swag2mcp ls ./
swag2mcp ls path/to
swag2mcp ls --tags=public,internal
```

**Output columns:** domain, title, base URL, auth type, collections.

---

## 7. `swag2mcp run [path]`

Launch the interactive Bubbletea TUI API explorer. No flags.

```sh
swag2mcp run
swag2mcp run ./
swag2mcp run path/to
```

**Behavior:** Opens a full-screen terminal UI for searching, browsing, and invoking endpoints. Supports two modes:
- **Search mode** — search by query, browse results, inspect endpoints
- **Browse mode** — navigate spec → collection → tag → endpoint hierarchy

---

## 8. `swag2mcp validate [path]`

Validate the configuration file and report issues.

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--tags` | `-t` | string | `""` | Filter specs by tags (comma-separated) |

```sh
swag2mcp validate
swag2mcp validate ./
swag2mcp validate path/to
swag2mcp validate --tags=public
```

**Behavior:** Loads config, validates all specs. Checks that:
- Config is valid YAML
- Domains are unique and lowercase
- `location` points to a valid OpenAPI/Swagger/Postman file
- Remote URLs are reachable
- Auth config is valid for the selected type

**What `validate` does NOT check:**
- Authentication endpoints (it validates auth config syntax, but does not test login/token exchange)
- Runtime availability of API endpoints (only spec URL reachability)
- Correctness of `base_url` (it validates format, but does not make test requests)
- Mock server configuration (`base_mock_url` is not verified for connectivity)

Prints `"Configuration is valid."` or detailed error messages.

---

## 9. `swag2mcp clean [path]`

Remove cached remote specs and invocation responses. No flags.

```sh
swag2mcp clean
swag2mcp clean ./
swag2mcp clean path/to
```

**Behavior:** Removes `cache/` and `responses/` directories. Preserves `specs/` and `auth_scripts/` (non-orphan).

---

## 10. `swag2mcp update [path]`

Validate config, clear cache, and re-cache all spec files. No flags.

```sh
swag2mcp update
swag2mcp update ./
swag2mcp update path/to
```

**Behavior:**
1. Loads and validates config
2. Cleans cache
3. Re-caches all spec files
4. Ensures auth scripts exist
5. Removes orphan auth scripts
6. Prints `"✅ N specs processed"`

---

## 11. `swag2mcp mcp [path]`

Start the MCP server (headless mode). This is the primary command for LLM tool access.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--logfile` / `-f` | string | `""` | Log file path (stderr if unset) |
| `--tags` / `-t` | string | `""` | Filter specs by tags (comma-separated) |
| `--disable-llm-auth` | bool | `true` | Remove `auth` tool from MCP tool list |
| `--dump-dir` | string | `""` | Directory to dump HTTP requests for debugging |
| `--transport` | string | `"stdio"` | MCP transport: `stdio`, `sse`, `streamable-http` |
| `--http-addr` | string | `":8080"` | HTTP server address (for sse/streamable-http) |
| `--http-path` | string | `"/mcp"` | HTTP path for MCP handler |
| `--auth-token` | string | `""` | Bearer token for HTTP transport auth |

**Examples:**
```sh
# Default stdio transport (for LLM integration)
swag2mcp mcp

# SSE transport with auth
swag2mcp mcp --transport sse --http-addr :8080 --auth-token my-secret

# With tag filtering and auth enabled
swag2mcp mcp --tags=public --disable-llm-auth=false

# With request dump directory
swag2mcp mcp --dump-dir ./dumps
```

> **Note:** If the workspace was initialized at a custom path (e.g. `swag2mcp init ./my-project`), specify the path when starting the MCP server: `swag2mcp mcp ./my-project`. The IDE configuration must also use the full path to the config file.

> **Path resolution warning:** When `[path]` is omitted, `swag2mcp mcp` searches for `swag2mcp.yaml` in the current directory first, then falls back to `~/.swag2mcp/`. If you run the command from the wrong directory, it may load a different workspace than intended. **Always specify `[path]` explicitly when running as a service or in IDE config.**

### IDE Configuration Examples

**VS Code** (`.vscode/settings.json` or global settings):
```json
{
  "mcp": {
    "servers": {
      "swag2mcp": {
        "command": "swag2mcp",
        "args": ["mcp", "/absolute/path/to/.swag2mcp"]
      }
    }
  }
}
```

**Cursor / Windsurf** (`~/.cursor/mcp.json` or project `.cursor/mcp.json`):
```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absolute/path/to/.swag2mcp"]
    }
  }
}
```

**Claude Desktop** (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):
```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absolute/path/to/.swag2mcp"]
    }
  }
}
```

**JetBrains IDEs** (Settings → Tools → MCP):
- Name: `swag2mcp`
- Command: `swag2mcp`
- Arguments: `mcp /absolute/path/to/.swag2mcp`

> Always use an **absolute path** to the workspace directory in IDE config. Relative paths may fail depending on the IDE's working directory.

**Behavior:**
- Requires existing config (does NOT auto-init)
- Applies MCP settings from YAML config as fallback
- Cleans old responses (>48h) on startup
- Starts MCP server with selected transport
- Prints `"MCP server listening on http://<addr><path>"` on stdout
- `--disable-llm-auth=true` (default) removes the `auth` tool
- HTTP transport provides `GET /health` returning `{"status":"ok","version":"..."}`

### MCP HTTP Transport — Handshake Protocol

When using `--transport streamable-http` or `--transport sse`, the MCP protocol requires a specific handshake sequence. `tools/list` and other tool calls will fail before initialization:

```
Step 1: POST /mcp → {"method":"initialize", ...}
Step 2: POST /mcp → {"method":"notifications/initialized"}
Step 3: POST /mcp → {"method":"tools/list", ...}   ← now works
```

**Health check** (works without initialization):
```sh
curl http://localhost:8080/health
# → {"status":"ok","version":"v1.2.0"}
```

---

## 12. `swag2mcp version`

Print the swag2mcp version. Also available as `--version` flag.

```sh
swag2mcp version
swag2mcp --version
```

---

## 13. `swag2mcp info [path]`

Show detailed configuration and runtime information as JSON.

```sh
swag2mcp info
swag2mcp info ./
swag2mcp info path/to
```

**Output includes:** version, workspace path, specs summary (total/active/disabled/collections/endpoints), HTTP client config, MCP transport, auth methods, mock mode status.

---

## 14. `swag2mcp import [path] [source] [name]`

Import spec files into the workspace. Three modes of operation.

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--spec` | `-s` | stringSlice | `nil` | Import collections from specified specs (comma-separated) |
| `--from-zip` | | string | `""` | Restore workspace from a swag2mcp backup ZIP |

> **`[path]` is the workspace directory** containing `swag2mcp.yaml`. `[source]` is a URL or local path to an OpenAPI/Swagger/Postman spec file. `[name]` is the domain name for the new spec.

### Mode 1 — Single import from URL or file

Import a spec file (URL or local path) and assign it a domain name:
```sh
swag2mcp import https://example.com/spec.yaml myspec
swag2mcp import /path/to/workspace https://example.com/spec.yaml myspec
swag2mcp import ./local-spec.yaml myspec
```

### Mode 2 — Bulk import from existing config

Import all collections for the specified domains from the configured spec URLs:
```sh
swag2mcp import --spec meteo
swag2mcp import /path/to/workspace --spec meteo,store
```

### Mode 3 — Restore from backup

> **The `--from-zip` file must be a ZIP created by `swag2mcp export`.** Arbitrary ZIP files or ZIPs from other tools will not work — the archive has a specific internal structure (`swag2mcp.yaml`, `specs/`, `auth_scripts/`).

```sh
swag2mcp import --from-zip /path/to/backup.zip
swag2mcp import /path/to/workspace /path/to/backup.zip
```

---

## 15. `swag2mcp export [path] [output]`

Export workspace as a portable ZIP backup.

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--spec` | `-s` | stringSlice | `nil` | Export only specified specs (comma-separated) |

**IMPORTANT:** The `[output]` argument must be a **full file path ending in `.zip`** (e.g. `/path/to/backup.zip`). Do NOT pass a directory — the command will not create a ZIP archive if given a directory path.

```sh
# ✅ Correct — full .zip file path
swag2mcp export /path/to/workspace /path/to/backup.zip

# ❌ Wrong — directory, NOT a file path (no ZIP will be created)
swag2mcp export /path/to/workspace /some/directory

# No output arg — creates swag2mcp-backup-<timestamp>.zip in current directory
swag2mcp export /path/to/workspace

# Filter by spec
swag2mcp export --spec meteo
swag2mcp export --spec meteo,store
```

**Behavior:** Creates a ZIP with all spec files, config, and auth scripts. Default output: `swag2mcp-backup-<timestamp>.zip` in current directory.

**Post-command verification:** After running `export`, always verify the ZIP file was created:
```sh
ls -la swag2mcp-backup-*.zip
# or for a custom output path:
ls -la /path/to/backup.zip
```

---

## 16. `swag2mcp-mock mockserver [path]`

Start mock servers for all API specs. **This is a separate binary** — `swag2mcp-mock` is not included in the main `swag2mcp` package and must be installed independently:

```sh
# Option 1: Download from GitHub Releases (separate archive: swag2mcp-mock_<version>_<os>_<arch>.tar.gz)
# Option 2: Install with Go
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Mock servers generate fake responses based on the OpenAPI schema — useful for testing without hitting real APIs.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--tls` | bool | `false` | Enable TLS with self-signed certificate |
| `--tls-cert` | string | `""` | Path to TLS certificate file |
| `--tls-key` | string | `""` | Path to TLS key file |

```sh
swag2mcp-mock
swag2mcp-mock ./
swag2mcp-mock path/to
swag2mcp-mock --tls
swag2mcp-mock --tls-cert cert.pem --tls-key key.pem
```

**Behavior:** Requires `mock_enabled: true` and `base_mock_url` in config. Starts HTTP servers for each spec/collection on configured ports.

**Example config with mock:**
```yaml
mock_enabled: true
specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: https://meteo.swagger.io/v2
    base_mock_url: localhost:8080
    collections:
      - llm_title: Open-Meteo
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo.json
        base_mock_url: localhost:8081
```

**Testing the mock:**
```sh
swag2mcp-mock ./
# In another terminal:
curl http://localhost:8081/pets
# → [{"id":1,"name":"Pet_name","status":"available"}]
```

---

## Post-Command Verification

Always verify that file-creating commands succeeded by checking the result. This prevents silent failures where a command reports success but the expected output is missing.

| Command | What to verify |
|---------|---------------|
| `swag2mcp init` | `ls <path>/swag2mcp.yaml` — config file exists |
| `swag2mcp export` | `ls -la <output>.zip` or `ls -la swag2mcp-backup-*.zip` — ZIP file exists |
| `swag2mcp import` | `swag2mcp ls <path>` — new spec appears in list |
| `swag2mcp add spec` | `swag2mcp ls <path>` — new spec appears in list |
| `swag2mcp add collection` | `swag2mcp ls <path>` — new collection appears under spec |
| `swag2mcp clean` | `ls <path>/cache <path>/responses` — directories removed |
| `swag2mcp update` | `swag2mcp ls <path>` — specs still listed and reachable |
| `swag2mcp delete spec/collection` | `swag2mcp ls <path>` — removed item no longer appears |

**General rule:** After any command that creates, modifies, or deletes files on disk, run a quick check (`ls`, `swag2mcp ls`, or `swag2mcp validate`) to confirm the expected result before reporting success to the user.

---

## Troubleshooting / Common Issues

### `swag2mcp: command not found`

swag2mcp is not installed or not in `$PATH`. Install it (see "Installation" section) or verify the binary location:
```sh
which swag2mcp
# If empty, add the install location to PATH or use the full path:
/usr/local/bin/swag2mcp --version
```

### `Configuration file not found` or wrong workspace loaded

swag2mcp resolves the config file using the "Path Resolution" rules above. If you get this error:
1. Verify you're in the correct directory, OR
2. Pass the workspace path explicitly: `swag2mcp ls /path/to/.swag2mcp`

### Spec URL returns 403 / 404 / CORS error

- **403 / 404:** The `location` URL is incorrect or the file was moved. Verify the URL in a browser.
- **CORS error when using `swag2mcp mcp`:** CORS only affects browser-based requests. The MCP server makes server-side requests and is not affected by CORS. If you see CORS errors, check that you're not confusing `location` (the spec file URL) with `base_url` (the API endpoint).

### YAML parse errors

swag2mcp requires valid YAML with correct indentation. Common mistakes:
- Tabs instead of spaces (YAML requires spaces)
- Missing indentation for nested fields (`collections`, `auth`, `http_client`)
- Unquoted strings with special characters (`:`, `#`, `&`, `{`)

Validate with: `swag2mcp validate /path/to/.swag2mcp`

### MCP server starts but tools/list returns nothing

- Verify the workspace path matches between `swag2mcp init` and `swag2mcp mcp`
- Check that `disable: true` is not set on all specs
- For HTTP transport, ensure the handshake sequence is completed (see "MCP HTTP Transport — Handshake Protocol")

### Port already in use (HTTP transport)

If `swag2mcp mcp --transport sse --http-addr :8080` fails with "address already in use":
```sh
# Find the process using port 8080
lsof -i :8080
# Use a different port
swag2mcp mcp --transport sse --http-addr :8081
```

---

## Path Resolution

Many swag2mcp commands accept an optional `[path]` argument — the workspace directory containing `swag2mcp.yaml`. Understanding how this path is resolved is critical to avoid loading the wrong config.

**Resolution order:**
1. **Explicit `[path]` argument** — if provided, swag2mcp uses it directly
2. **Current directory** — looks for `./swag2mcp.yaml`
3. **Default home directory** — falls back to `~/.swag2mcp/swag2mcp.yaml`

**Common pitfalls:**
- Running `swag2mcp mcp` without `[path]` from a random directory may load `~/.swag2mcp/` instead of your project's `.swag2mcp/`
- `swag2mcp init ./` creates `swag2mcp.yaml` in the current directory, NOT inside a `.swag2mcp/` subdirectory (see "Workspace creation rules" for the recommended pattern)
- **Always pass an explicit `[path]`** when running as a service or configuring an IDE

---

## Config File Location

The config file is `swag2mcp.yaml`. Resolution order:
1. Explicit `[path]` argument
2. Current directory (`./swag2mcp.yaml`)
3. Default: `~/.swag2mcp/swag2mcp.yaml`

---

## Full Configuration Reference

Every field in `swag2mcp.yaml` with type, required status, and description:

### Global Settings

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `mock_enabled` | bool | no | `false` | Enable mock servers for all specs |
| `http_client.timeout` | duration | no | `30s` | HTTP request timeout |
| `http_client.follow_redirects` | bool | no | `true` | Follow HTTP redirects |
| `http_client.max_redirects` | int | no | `10` | Max redirects to follow |
| `http_client.max_response_size` | int | no | `2048` | Max response body size in bytes (truncated, saved to file if exceeded) |
| `http_client.randomize` | bool | no | `false` | Randomize browser-like headers |
| `http_client.proxy.url` | string | no | `""` | HTTP proxy URL |
| `http_client.proxy.username` | string | no | `""` | Proxy username |
| `http_client.proxy.password` | string | no | `""` | Proxy password |
| `http_client.proxy.bypass` | []string | no | `[]` | Hosts to bypass proxy |
| `http_client.headers` | map[string]string | no | `{}` | Custom headers for every request |
| `http_client.cookies` | map[string]string | no | `{}` | Custom cookies for every request |
| `http_client.user_agent` | string | no | `""` | Custom User-Agent override |
| `mcp.transport` | string | no | `"stdio"` | MCP transport: `stdio`, `sse`, `streamable-http` |
| `mcp.addr` | string | no | `":8080"` | HTTP server address |
| `mcp.path` | string | no | `"/mcp"` | HTTP path for MCP handler |
| `mcp.auth_token` | string | no | `""` | Bearer token for HTTP transport auth |

### Spec Settings

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `domain` | string | **yes** | — | Unique identifier (lowercase, digits, hyphens only) |
| `llm_title` | string | **yes** | — | Human-readable name shown in MCP tools |
| `llm_instruction` | string | no | `""` | Instruction for LLM on how to use this API |
| `base_url` | string | **yes** | — | Base URL for all API requests |
| `base_mock_url` | string | no | `""` | Host:port for mock server |
| `disable` | bool | no | `false` | Exclude this spec from MCP tools |
| `tags` | []string | no | `[]` | Tags for filtering (`--tags` flag) |
| `auth` | object | no | `{}` | Authentication config (see Auth Types) |
| `http_client` | object | no | `{}` | Override global HTTP settings for this spec |
| `collections` | []object | **yes** | — | At least 1 collection required |

### Collection Settings

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `llm_title` | string | **yes** | — | Collection name |
| `location` | string | **yes** | — | Path or URL to OpenAPI/Swagger/Postman file |
| `disable` | bool | no | `false` | Exclude this collection |
| `base_mock_url` | string | no | `""` | Host:port for mock server (overrides spec) |
| `http_client` | object | no | `{}` | Override HTTP settings (overrides spec and global) |

---

## Real-World YAML Examples

> **Spec files** are included in the repository at `specs/` and available via raw GitHub URL:
> - `https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo.json` — Open-Meteo API (OpenAPI 3.0)
> - `https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml` — Binance Market Data (OpenAPI 3.0)
> - `https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml` — icanhazdadjoke (OpenAPI 3.0)
> - `https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/pokeapi.yaml` — PokéAPI (OpenAPI 3.0)
>
> Use the raw URL in `location` — works from anywhere without cloning the repo.
>
> Full ready-to-run examples are in the `examples/` directory.

### Example 1: Open-Meteo (public, no auth)

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo API
    llm_instruction: |
      Classic Swagger Open-Meteo API. Use this to manage pets,
      store inventory, and user accounts.
    base_url: https://meteo.swagger.io/v2
    collections:
      - llm_title: Open-Meteo
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo.json
```

### Example 2: Binance Market Data (public, no auth)

```yaml
specs:
  - domain: binance
    llm_title: Binance Market Data API
    llm_instruction: |
      Binance cryptocurrency exchange public market data.
      Use this to get prices, klines, exchange info, and 24hr ticker.
    base_url: https://api.binance.com
    collections:
      - llm_title: Market Data
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml
```

### Example 3: icanhazdadjoke (public, no auth)

```yaml
specs:
  - domain: dadjoke
    llm_title: icanhazdadjoke API
    llm_instruction: |
      The largest selection of dad jokes on the internet.
      Use this to get random jokes, search by term, or find a specific joke by ID.
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

### Example 4: PokéAPI (public, no auth)

```yaml
specs:
  - domain: pokeapi
    llm_title: PokéAPI
    llm_instruction: |
      The RESTful Pokémon API. Use this to get Pokémon data,
      list Pokémon, search by type, and find Pokémon by ID or name.
    base_url: https://pokeapi.co
    collections:
      - llm_title: Pokémon
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/pokeapi.yaml
```