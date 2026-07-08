# swag2mcp

**swag2mcp** is a CLI tool and MCP (Model Context Protocol) server that bridges OpenAPI/Swagger/Postman API specifications with LLM agents (Opencode, Crush, Copilot, Cursor, etc.).

It indexes your API specs into a full-text search engine, exposes them through 14 MCP tools, and lets LLMs discover, inspect, and invoke real API endpoints — all without writing a single line of integration code.

---

## Table of Contents

- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [CLI Commands](#cli-commands)
- [MCP Server](#mcp-server)
- [Search](#search)
- [Workspace](#workspace)
- [Caching](#caching)
- [Development](#development)

---

## Quick Start

```bash
# Install
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest

# Initialize workspace
swag2mcp init

# Start MCP server (for LLM agents)
swag2mcp mcp

# Or use interactive explorer
swag2mcp run
```

---

## Configuration

### YAML Schema

```yaml
http_client:                        # optional, global HTTP defaults
  headers:                          # optional
    X-API-Version: "2"
  cookies: []                       # optional
  user_agent: ""                    # optional
  timeout: 0s                       # optional
  follow_redirects: true            # optional
  max_redirects: 10                 # optional
  max_response_size: 1048           # optional, bytes (default 1KB, max 1MB)

specs:
  - domain: petstore                    # required, 1-60 chars, [a-zA-Z0-9_-]
    llm_title: Petstore API             # required, 5-120 chars
    llm_instruction: |                  # optional, max 500 chars
      Use this API to manage pets, orders, and users.
    base_url: https://petstore.swagger.io/v2  # required, valid URL
    disable: false                      # optional
    tags: [public, demo]                # optional, for filtering
    http_client:                        # optional, overrides global
      headers:
        X-API-Version: "2"
    auth:                               # optional
      type: bearer                      # see Auth Methods below
      config:
        token: $(TOKEN_AUTH)
    collections:
      - llm_title: Petstore Swagger     # optional, max 120 chars
        llm_instruction: |             # optional, max 360 chars
          Main petstore endpoints
        title: ""                      # optional, auto-populated from spec
        location: https://petstore.swagger.io/v2/swagger.json  # required, 5-250 chars
        disable: false                  # optional
        base_url: ""                    # optional, overrides spec base_url
        http_client: {}                 # optional, overrides spec
```

### Tags — Filtering Specs by Project

Tags let you organize specs by project, environment, or team. When starting the MCP server, use `--tags` to load only matching specs:

```bash
# Start server with only public specs
swag2mcp mcp --tags=public

# Start server with multiple tags
swag2mcp mcp --tags=public,internal

# Run multiple servers for different projects
swag2mcp mcp --tags=project-alpha --logfile=/tmp/swag2mcp-alpha.log
swag2mcp mcp --tags=project-beta  --logfile=/tmp/swag2mcp-beta.log
```

This allows running separate MCP servers for different projects from a single config file.

### Auth Methods

| Type | Fields | Config Example |
|------|--------|----------------|
| `none` | — | `type: none` |
| `basic` | `username`, `password` | `username: $(USER)`, `password: $(PASS)` |
| `bearer` | `token` | `token: $(TOKEN)` |
| `digest` | `username`, `password` | `username: admin`, `password: secret` |
| `api-key` | `key`, `value`, `in` (header/query) | `key: X-API-Key`, `value: $(KEY)`, `in: header` |
| `oauth2-cc` | `client_id`, `client_secret`, `token_url`, `scopes` | `client_id: $(ID)`, `token_url: https://auth.example.com/token` |
| `oauth2-pwd` | `username`, `password`, `client_id`, `client_secret`, `token_url`, `scopes` | `username: $(USER)`, `token_url: https://auth.example.com/token` |
| `script` | `source` | `source: path/to/auth.sh` |

All string fields support `$(ENV_VAR)` syntax — resolved at runtime from environment variables.

---

## CLI Commands

All commands that accept `[path]` use the same path resolution:

```
swag2mcp <command>          → ~/.swag2mcp/
swag2mcp <command> ./       → {cwd}/.swag2mcp/
swag2mcp <command> path/to  → {cwd}/path/to/.swag2mcp/
```

### `init [path]`

Initialize workspace and configuration.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--interactive` | `-i` | `false` | Run interactive wizard |
| `--force` | `-f` | `false` | Overwrite existing config |

```bash
swag2mcp init              # create ~/.swag2mcp/swag2mcp.yaml
swag2mcp init ./           # create ./.swag2mcp/swag2mcp.yaml
swag2mcp init -i           # interactive wizard
```

### `add spec [path]` / `add collection [path]`

Add a specification or collection to the config.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--yaml` | `-y` | `""` | YAML input (use `-` for stdin) |
| `--example` | `-e` | `false` | Show YAML example |

```bash
swag2mcp add spec
swag2mcp add spec --yaml 'domain: petstore\nllm_title: Petstore API\nbase_url: https://...'
cat spec.yaml | swag2mcp add spec --yaml -
swag2mcp add spec --example
```

### `delete spec [path]` / `delete collection [path]`

Delete a specification or collection. Interactive prompts for selection.

```bash
swag2mcp delete spec
swag2mcp delete collection
```

### `ls [path]`

List specifications and collections.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--tags` | `-t` | `""` | Filter by tags (comma-separated) |

```bash
swag2mcp ls
swag2mcp ls --tags=public,internal
```

### `run [path]`

Interactive API explorer (TUI). Search, browse, inspect, and save endpoints.

```bash
swag2mcp run
```

### `validate [path]`

Validate configuration and check that all collection locations are accessible.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--tags` | `-t` | `""` | Filter specs by tags |

```bash
swag2mcp validate
swag2mcp validate --tags=public
```

### `clean [path]`

Remove all contents of `cache/` and `responses/` directories.

```bash
swag2mcp clean
```

### `update [path]`

Validate config, clear cache, re-cache all spec files.

```bash
swag2mcp update
```

### `mcp [path]`

Start the MCP server in headless mode (stdio transport). This is the primary production command for LLM integration.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--logfile` | `-f` | `""` | Log file path |
| `--tags` | `-t` | `""` | Filter specs by tags |
| `--disable-llm-auth` | | `true` | `true` — auth happens under the hood (LLM never sees tokens). `false` — LLM can request tokens via the `auth` tool |
| `--dump-dir` | | `""` | Directory to dump HTTP requests for debugging |

```bash
swag2mcp mcp
swag2mcp mcp --tags=public --logfile=/var/log/swag2mcp.log
swag2mcp mcp --disable-llm-auth=false
swag2mcp mcp --dump-dir=/tmp/dump
```

---

## MCP Server

The MCP server exposes 14 tools over stdio transport. LLM agents (Opencode, Crush, Copilot, Cursor, etc.) connect automatically when configured.

### Tool Hierarchy

```
spec_list                       — list all available specs
  └─ spec_by_id                 — spec details by ID
       └─ collection_by_spec    — collections in a spec
            └─ tag_by_collection     — tags in a collection
                 └─ endpoint_by_tag  — endpoints under a tag
                      └─ inspect          — full OpenAPI operation
                           └─ invoke       — execute API call

search                          — full-text search across all endpoints
```

### Tool Reference

| Tool | Args | Returns | Description |
|------|------|---------|-------------|
| `spec_list` | — | `Spec[]` | All available specs |
| `spec_by_id` | `id` | Spec + Collections | Spec details |
| `collection_by_spec` | `specId` | Collections | Collections in a spec |
| `collection_by_id` | `id` | Collection + Tags | Collection details |
| `tag_by_collection` | `collectionId` | Tags | Tags in a collection |
| `tag_by_spec` | `specId` | Tags | All tags across a spec |
| `tag_by_id` | `id` | Tag | Single tag metadata |
| `endpoint_by_tag` | `tagId` | Endpoints | Endpoints under a tag |
| `endpoint_by_collection` | `collectionId` | Endpoints | All endpoints in a collection |
| `endpoint_by_spec` | `specId` | Endpoints | All endpoints across a spec |
| `endpoint_by_id` | `id` | Endpoint | Quick endpoint summary |
| `search` | `query`, `limit` | Endpoints | Full-text search |
| `inspect` | `endpointId` | Full Operation | Complete OpenAPI operation object |
| `invoke` | `endpointId`, `parameters`, `requestBody` | Response | Executes real API call |
| `auth` | `specId` | Token | Get auth token for a spec |

---

## Search

### Query Syntax

| Feature | Syntax | Example |
|---------|--------|---------|
| Term | `term` | `pets` |
| Phrase | `"phrase"` | `"add a new pet"` |
| Field: method | `method:term` | `method:post` |
| Field: tag | `tag:term` | `tag:auth` |
| Field: path | `path:term` | `path:/users` |
| Field: summary | `summary:term` | `summary:login` |
| Required (AND) | `+term` | `+method:post +tag:user` |
| Excluded (NOT) | `-term` | `-deprecated` |
| Wildcard | `*` | `path:*/v2/*` |
| Fuzzy | `term~` | `watex~` |
| Regex | `/pattern/` | `/user(s\|sessions)/` |
| Boost | `term^N` | `tag:pet^5` |
| Match all | `*` | `*` |

### Examples

```
# Find POST endpoints in auth tag
+method:post +tag:auth

# Search for login-related endpoints
summary:"login"~

# Find all user-related paths, exclude deprecated
path:*/users/* -deprecated

# Complex query
+method:get +tag:pet summary:"find by status"
```

### Indexed Fields

| Field | Type | Content |
|-------|------|---------|
| `method` | text | HTTP method (lowercased) |
| `tag` | text | Tag name (lowercased) |
| `path` | text | API path (lowercased) |
| `summary` | text (analyzed) | Endpoint summary/description (lowercased) |
| `_all` | text (analyzed) | method + path + tag + summary |

---

## Workspace

### Directory Structure

```
~/.swag2mcp/                    # or {project}/.swag2mcp/
├── swag2mcp.yaml               # Configuration file
├── cache/                      # Cached remote specs
│   ├── {hash}.spec             # Spec file content
│   └── {hash}.meta             # JSON metadata
├── specs/                      # Local spec files (user-managed)
├── responses/                  # Invocation response files
└── auth_scripts/               # Authentication scripts
```

### Path Resolution

```
swag2mcp <command>          → ~/.swag2mcp/
swag2mcp <command> ./       → {cwd}/.swag2mcp/
swag2mcp <command> path/to  → {cwd}/path/to/.swag2mcp/
```

### .gitignore

Only temporary data should be gitignored:

```
.swag2mcp/cache/*
.swag2mcp/responses/*
```

The config `.swag2mcp/swag2mcp.yaml` and spec files in `.swag2mcp/specs/` **must be in the repository**.

### Recommendation

Keep all spec files in `.swag2mcp/specs/` — this is the only way to ensure they are used directly without being copied to cache.

---

### Rules

| Source | Behavior |
|--------|----------|
| HTTP/HTTPS URL | Always cached. TTL: random 1-48h. |
| Local path inside `specs/` | Used directly, not cached. |
| Local path outside `specs/` | Copied to cache on first access. |
| `file://` URL | Treated as local path. |

### Cache Key

SHA-256 hash of the normalized location (first 16 bytes = 32 hex chars).

### Cache Hit Logic

1. Read `.meta` file — expired or missing → miss
2. For local sources: `ModTime` changed → miss
3. `.spec` file missing → miss
4. Otherwise → hit

---

## Development

```bash
# Build
go build ./cmd/swag2mcp/

# Test
go test ./...

# Lint
make lint

# Run
go run ./cmd/swag2mcp/main.go
```
