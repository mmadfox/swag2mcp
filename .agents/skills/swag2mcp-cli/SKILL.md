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

# Add a spec (use local spec file from the repository)
swag2mcp add spec --yaml 'domain: petstore
llm_title: Petstore API
base_url: https://petstore.swagger.io/v2
collections:
  - llm_title: Pets
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json'

# List specs
swag2mcp ls

# Start MCP server (for LLM tools)
swag2mcp mcp
```

---

## Concepts: Spec vs Collection

Before writing config, understand the two levels:

- **Spec** = one API (a domain). Represents a logical service (e.g. "Petstore API", "Binance Market Data"). A spec can have multiple collections.
- **Collection** = one OpenAPI/Swagger file. If an API has multiple spec files (different versions, different microservices), each file is a separate collection under the same spec.

**Important:** The `location` field in a collection must point to an **OpenAPI 3.x, Swagger 2.0, or Postman collection** file — not to the API endpoint itself. For example:
- ✅ `location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json` — correct, this is an OpenAPI spec
- ❌ `location: https://restcountries.com/v3.1/all` — wrong, this is a JSON response, not a spec

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

**Behavior:** Creates `cache/`, `specs/`, `responses/`, `auth_scripts/` subdirectories and a `swag2mcp.yaml` config file. After init, prints a hint: `"Next step: edit swag2mcp.yaml or run 'swag2mcp add spec --yaml \"...\"'"`.

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
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json'
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
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json
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

Start mock servers for all API specs (separate binary). Mock servers generate fake responses based on the OpenAPI schema — useful for testing without hitting real APIs.

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
  - domain: petstore
    llm_title: Petstore API
    base_url: https://petstore.swagger.io/v2
    base_mock_url: localhost:8080
    collections:
      - llm_title: Petstore
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json
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
> - `https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json` — Petstore API (OpenAPI 3.0)
> - `https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml` — Binance Market Data (OpenAPI 3.0)
> - `https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/countries.yaml` — REST Countries (OpenAPI 3.0)
>
> Use the raw URL in `location` — works from anywhere without cloning the repo.
>
> Full ready-to-run examples are in the `examples/` directory.

### Example 1: Petstore (public, no auth)

```yaml
specs:
  - domain: petstore
    llm_title: Petstore API
    llm_instruction: |
      Classic Swagger Petstore API. Use this to manage pets,
      store inventory, and user accounts.
    base_url: https://petstore.swagger.io/v2
    collections:
      - llm_title: Petstore
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json
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

### Example 3: REST Countries (public, no auth)

```yaml
specs:
  - domain: countries
    llm_title: REST Countries API
    llm_instruction: |
      Get information about countries — names, capitals, currencies, regions.
    base_url: https://restcountries.com
    collections:
      - llm_title: Countries
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/countries.yaml
```

### Example 4: GitHub API (Bearer token)

```yaml
specs:
  - domain: github
    llm_title: GitHub REST API
    llm_instruction: |
      GitHub API for managing repositories, issues, pull requests,
      and user data. Requires authentication.
    base_url: https://api.github.com
    auth:
      type: bearer
      config:
        token: $(GITHUB_TOKEN)
    collections:
      - llm_title: GitHub API v3
        location: https://raw.githubusercontent.com/github/rest-api-description/main/descriptions/api.github.com/api.github.com.yaml
```

### Example 5: Multiple collections per spec

```yaml
specs:
  - domain: my-service
    llm_title: My Microservice API
    base_url: https://api.my-service.com
    collections:
      - llm_title: Users
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json
      - llm_title: Orders
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json
      - llm_title: Payments
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json
```

### Example 6: With mock server

```yaml
mock_enabled: true
specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://petstore.swagger.io/v2
    base_mock_url: localhost:8080
    collections:
      - llm_title: Petstore
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json
        base_mock_url: localhost:8081
```

---

## Common Recipes / Patterns

### Public API without key
```yaml
- domain: my-api
  base_url: https://api.example.com
  collections:
    - llm_title: Main
      location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json
```
No `auth` block needed. Works for Petstore, Binance market data, REST Countries, etc.

### API with Bearer token
```yaml
- domain: my-api
  base_url: https://api.example.com
  auth:
    type: bearer
    config:
      token: $(MY_TOKEN)
  collections:
    - llm_title: Main
      location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json
```
Works for GitHub, OpenAI, most modern APIs. Use `$(ENV_VAR)` to keep secrets out of config.

### API with API Key in header
```yaml
- domain: my-api
  base_url: https://api.example.com
  auth:
    type: api-key
    config:
      key: X-API-Key
      value: $(API_KEY)
      in: header
  collections:
    - llm_title: Main
      location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json
```

### API with Basic Auth
```yaml
- domain: my-api
  base_url: https://api.example.com
  auth:
    type: basic
    config:
      username: $(USER)
      password: $(PASS)
  collections:
    - llm_title: Main
      location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/petstore.json
```

### Disable a spec (exclude from MCP)
```yaml
- domain: old-api
  disable: true
  base_url: https://old.example.com
  collections:
    - llm_title: Old
      location: ./specs/old.yaml
```

### Tag-based filtering
```yaml
- domain: public-api
  tags: [public]
  ...
- domain: internal-api
  tags: [internal]
  ...
```
Then: `swag2mcp mcp --tags=public` — only public specs loaded.

---

## Config Cascade (How Settings Override)

Settings cascade from global → spec → collection. More specific wins:

```
Global http_client (timeout: 30s, headers: {X-Global: "true"})
  └── Spec http_client (timeout: 10s, headers: {X-Spec: "spec-only"})
       └── Collection http_client (headers: {X-Collection: "collection-only"})
```

Result for this collection: `timeout=10s` (from spec), `X-Global=true` (from global), `X-Spec=spec-only` (from spec), `X-Collection=collection-only` (from collection).

---

## Environment Variables

All config values support `$(ENV_VAR)` syntax for environment variable resolution:

```yaml
auth:
  type: bearer
  config:
    token: $(GITHUB_PAT)   # read from GITHUB_PAT env var
```

This works in any string field: `token`, `username`, `password`, `key`, `value`, `base_url`, etc.

---

## Auth Types

| Type | Config Fields | Example |
|------|---------------|---------|
| `none` | — | No auth |
| `basic` | `username`, `password` | `username: user\npassword: $(PASS)` |
| `bearer` | `token` | `token: $(TOKEN)` |
| `digest` | `username`, `password` | `username: user\npassword: $(PASS)` |
| `oauth2-cc` | `token_url`, `client_id`, `client_secret`, `scopes` | `token_url: https://auth.example.com/token\nclient_id: $(CLIENT_ID)` |
| `oauth2-pwd` | `token_url`, `client_id`, `client_secret`, `username`, `password`, `scopes` | `client_secret` is optional (public client support for Keycloak) |
| `api-key` | `key`, `value`, `in` (header/query/cookie) | `key: X-API-Key\nvalue: $(KEY)\nin: header` |
| `script` | `path`, `args`, `env` | `path: ./auth_scripts/my-script.sh` |
| `hmac` | `api_key`, `secret_key` | `api_key: $(API_KEY)\nsecret_key: $(SECRET)` |

All config values support `$(ENV_VAR)` syntax for environment variable resolution.

---

## Error Messages & What They Mean

| Error | Cause | What to do |
|-------|-------|------------|
| `location does not appear to be a valid OpenAPI/Swagger spec` | URL/file is not an OpenAPI spec | Ensure `location` points to an OpenAPI 3.x, Swagger 2.0, or Postman collection file |
| `unexpected status 404` | Spec URL not found | Check the URL. For GitHub raw files, use `raw.githubusercontent.com` not `github.com` |
| `method is invalid during session initialization` | MCP HTTP request before `initialize` | Send `initialize` + `notifications/initialized` first (see MCP HTTP Handshake) |
| `rate_limit` | Same endpoint invoked twice within 10s | Wait 10 seconds before calling again |
| `not_found` | Spec/collection/endpoint ID doesn't exist | Use `search` to find the correct ID |
| `validation_failed` | Config validation error | Read the specific error message and fix the field |
| `invoke_error` | HTTP/connection error during API call | Check the base_url and network connectivity |
| `duplicate domain` | Two specs with the same `domain` | Use unique domain names (lowercase, no spaces) |

---

## MCP Tools (when running `swag2mcp mcp`)

| Tool | Description | When to use |
|------|-------------|-------------|
| `spec_list` | List all API specifications | Start here to discover available APIs |
| `spec_by_id` | Get spec details by ID | After `spec_list`, to see collections |
| `collection_by_spec` | List collections in a spec | Navigate spec → collection |
| `collection_by_id` | Get collection details by ID | See tags and method count |
| `tag_by_spec` | List all tags across a spec | See all endpoint categories |
| `tag_by_collection` | List tags in a collection | See categories within a collection |
| `tag_by_id` | Get tag details by ID | Verify tag exists, see method count |
| `endpoint_by_spec` | List all endpoints in a spec | Full inventory of a spec |
| `endpoint_by_collection` | List endpoints in a collection | All endpoints in one collection |
| `endpoint_by_tag` | List endpoints in a tag | Filtered view by category |
| `endpoint_by_id` | Get endpoint summary by ID | Quick overview (method, path, summary) |
| **`search`** | **Full-text search across all endpoints** | **Use this to find endpoints — do NOT manually traverse spec → collection → tag** |
| `inspect` | Get full OpenAPI operation details | **Always use before `invoke`** to see parameters, request body, responses |
| `invoke` | Execute a real API call | Only after `inspect` and user confirmation (for destructive operations) |
| `auth` | Get auth token/headers for a spec | Disabled with `--disable-llm-auth` (default) |

### MCP Tools — Usage Rules

1. **`search` is the ONLY tool for finding endpoints** when you don't have an ID. Do NOT manually traverse spec → collection → tag → endpoint. One `search` call replaces dozens of manual navigation steps.

2. **Always `inspect` before `invoke`** — you need to know the required parameters, request body schema, and response format.

3. **Never invoke destructive operations (POST/PUT/PATCH/DELETE) without explicit user confirmation.**

### Search Query Examples

| Query | Finds |
|-------|-------|
| `pet` | Any endpoint mentioning "pet" |
| `method:GET` | All GET endpoints |
| `tag:store` | Endpoints tagged "store" |
| `path:/api/v1/users` | Endpoints matching this path |
| `+method:POST +summary:user` | POST endpoints with "user" in summary |
| `"create user"` | Exact phrase match |
| `list*` | Wildcard — "listPets", "listUsers", etc. |
| `sumary~` | Fuzzy match (typo tolerance) |
| `*` | All endpoints |
