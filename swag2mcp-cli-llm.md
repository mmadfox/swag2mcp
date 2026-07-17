# swag2mcp CLI — LLM Agent Guide

This document describes every CLI command in swag2mcp. Use it when you need to tell a user how to interact with the tool from the terminal.

---

## Quick Start

```sh
# Initialize a workspace
swag2mcp init ~/my-api

# Add a spec
swag2mcp add spec --yaml 'domain: petstore
llm_title: Petstore API
base_url: https://petstore.swagger.io/v2
collections:
  - llm_title: Pets
    location: https://petstore.swagger.io/v2/swagger.json'

# List specs
swag2mcp ls

# Start MCP server (for LLM tools)
swag2mcp mcp
```

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
swag2mcp init ./                 # create ./swag2mcp.yaml
swag2mcp init path/to            # create path/to/swag2mcp.yaml
swag2mcp init -i                 # interactive wizard
swag2mcp init -f                 # force overwrite
```

**Behavior:** Creates `cache/`, `specs/`, `responses/`, `auth_scripts/` subdirectories and a `swag2mcp.yaml` config file. Without `-i`, runs a guided prompt. With `-i`, launches full TUI.

---

## 2. `swag2mcp add spec [path]`

Add a new API specification to the config.

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--yaml` | `-y` | string | `""` | YAML input inline or `-` for stdin |
| `--example` | `-e` | bool | `false` | Print YAML template and exit |

**Examples:**
```sh
swag2mcp add spec --example
swag2mcp add spec --yaml 'domain: petstore
llm_title: Petstore API
base_url: https://petstore.swagger.io/v2
collections:
  - llm_title: Pets
    location: https://petstore.swagger.io/v2/swagger.json'
cat spec.yaml | swag2mcp add spec --yaml -
```

**YAML format:**
```yaml
domain: petstore
llm_title: Petstore API
llm_instruction: Use this API to manage pets.
base_url: https://petstore.swagger.io/v2
tags: [public, demo]
auth:
  type: bearer
  config:
    token: $(TOKEN)
collections:
  - llm_title: Petstore Swagger
    location: https://petstore.swagger.io/v2/swagger.json
```

---

## 3. `swag2mcp add collection [path]`

Add a new collection to an existing spec.

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--yaml` | `-y` | string | `""` | YAML input inline or `-` for stdin |
| `--example` | `-e` | bool | `false` | Print YAML template and exit |

**Examples:**
```sh
swag2mcp add collection --example
swag2mcp add collection --yaml 'spec_domain: petstore
llm_title: Orders Collection
location: https://petstore.example.com/orders.json'
cat collection.yaml | swag2mcp add collection --yaml -
```

**YAML format:**
```yaml
spec_domain: petstore
llm_title: Orders Collection
location: https://petstore.example.com/orders.json
```

---

## 4. `swag2mcp delete spec [path]`

Delete a specification interactively. No flags. Prompts for selection and confirmation.

```sh
swag2mcp delete spec
swag2mcp delete spec ./
```

---

## 5. `swag2mcp delete collection [path]`

Delete a collection interactively. No flags. Prompts for spec selection, collection selection, and confirmation.

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

**Behavior:** Loads config, creates cache, validates all specs. Prints `"Configuration is valid."` or error details.

---

## 9. `swag2mcp clean [path]`

Remove cached remote specs and invocation responses. No flags.

```sh
swag2mcp clean
swag2mcp clean ./
swag2mcp clean path/to
```

**Behavior:** Removes `cache/` and `responses/` directories. Also removes orphan auth scripts.

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
6. Prints `"✅ Cache updated (N specs cached)"`

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

**Behavior:**
- Requires existing config (does NOT auto-init)
- Applies MCP settings from YAML config as fallback
- Cleans old responses (>48h) on startup
- Starts MCP server with selected transport
- `--disable-llm-auth=true` (default) removes the `auth` tool

---

## 12. `swag2mcp version`

Print the swag2mcp version. No flags, no args.

```sh
swag2mcp version
```

---

## 13. `swag2mcp info [path]`

Show detailed configuration and runtime information as JSON.

```sh
swag2mcp info
swag2mcp info ./
swag2mcp info path/to
```

**Output includes:** version, workspace path, uptime, specs summary, HTTP client config, MCP transport, auth methods, mock mode status.

---

## 14. `swag2mcp import [path] [source] [name]`

Import spec files into the workspace. Three modes of operation.

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--spec` | `-s` | stringSlice | `nil` | Import collections from specified specs (comma-separated) |
| `--from-zip` | | string | `""` | Restore workspace from a swag2mcp backup ZIP |

**Mode 1 — Single import:**
```sh
swag2mcp import https://example.com/spec.yaml myspec
swag2mcp import /path/to/workspace https://example.com/spec.yaml myspec
swag2mcp import ./local-spec.yaml myspec
```

**Mode 2 — Bulk import from config:**
```sh
swag2mcp import --spec petstore
swag2mcp import /path/to/workspace --spec petstore,store
```

**Mode 3 — Restore from backup:**
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

```sh
swag2mcp export
swag2mcp export /path/to/workspace
swag2mcp export /path/to/workspace /path/to/backup.zip
swag2mcp export --spec petstore
swag2mcp export --spec petstore,store
```

**Behavior:** Creates a ZIP with all spec files, config, and auth scripts. Default output: `swag2mcp-backup-<timestamp>.zip` in current directory.

---

## 16. `swag2mcp-mock mockserver [path]`

Start mock servers for all API specs (separate binary).

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
```

**Behavior:** Requires `mock_enabled: true` and `base_mock_url` in config. Starts HTTP servers for each spec/collection.

---

## Config File Location

The config file is `swag2mcp.yaml`. Resolution order:
1. Explicit `[path]` argument
2. Current directory (`./swag2mcp.yaml`)
3. Default: `~/.swag2mcp/swag2mcp.yaml`

---

## Config File Structure

```yaml
mock_enabled: false
http_client:
  timeout: 30s
  follow_redirects: true
  max_response_size: 2048
mcp:
  transport: stdio
  addr: ":8080"
  path: "/mcp"
  auth_token: ""
specs:
  - domain: "petstore"
    llm_title: "Petstore API"
    llm_instruction: "Use this API to manage pets."
    tags: ["public"]
    base_url: "https://petstore.swagger.io/v2"
    base_mock_url: "localhost:8080"
    auth:
      type: api-key
      config:
        key: "X-API-Key"
        value: "abc123"
        in: header
    collections:
      - llm_title: "Pet Operations"
        location: "https://petstore.swagger.io/v2/swagger.json"
        base_mock_url: "localhost:8081"
```

---

## Auth Types

| Type | Config Fields |
|------|---------------|
| `none` | — |
| `basic` | `username`, `password` |
| `bearer` | `token` |
| `digest` | `username`, `password` |
| `oauth2-cc` | `token_url`, `client_id`, `client_secret`, `scopes` |
| `oauth2-pwd` | `token_url`, `client_id`, `client_secret`, `username`, `password`, `scopes` |
| `api-key` | `key`, `value`, `in` (header/query/cookie) |
| `script` | `path`, `args`, `env` |

All config values support `$(ENV_VAR)` syntax for environment variable resolution.

---

## MCP Tools (when running `swag2mcp mcp`)

| Tool | Description |
|------|-------------|
| `spec_list` | List all API specifications |
| `spec_by_id` | Get spec details by ID |
| `collection_by_spec` | List collections in a spec |
| `collection_by_id` | Get collection details by ID |
| `tag_by_spec` | List all tags across a spec |
| `tag_by_collection` | List tags in a collection |
| `tag_by_id` | Get tag details by ID |
| `endpoint_by_spec` | List all endpoints in a spec |
| `endpoint_by_collection` | List endpoints in a collection |
| `endpoint_by_tag` | List endpoints in a tag |
| `endpoint_by_id` | Get endpoint summary by ID |
| `search` | Full-text search across all endpoints |
| `inspect` | Get full OpenAPI operation details |
| `invoke` | Execute a real API call |
| `auth` | Get auth token/headers for a spec (disabled with `--disable-llm-auth`) |
