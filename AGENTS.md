# swag2mcp — Agent Guide

## Commands

```sh
make lint              # golangci-lint (120 line limit, strict linters)
make cover             # go test ./... + coverage HTML report (excludes commands, tui, mockserver)
make integration-tests # go test -v -count=1 -timeout 600s ./tests/...
make testall           # lint + integration-tests + go test ./...
go test ./pkg/...      # single package
go test ./...          # all packages
go generate ./...      # installs mockgen + generates mocks (internal/server/mcp)
```

Order: `make lint && go test ./...`

**Note:** `make cover-core` is declared in the Makefile but has no recipe — do not use it.

## Architecture

- **Entrypoint**: `cmd/swag2mcp/main.go` — cobra CLI with **13 subcommands** (init, add, delete, ls, run, validate, clean, update, mcp, version, info, import, export)
- **Mock binary**: `cmd/swag2mcp-mock/main.go` — separate binary with `mockserver` subcommand (`--tls`, `--tls-cert`, `--tls-key`)
- **Core**: `internal/service/` — business logic (Bootstrap, Invoke, Search, Inspect, Auth, Specs, etc.)
- **MCP server**: `internal/server/mcp/` — **16 MCP tools**, uses `go.uber.org/mock` for tests; HTTP transport provides `GET /health` endpoint
- **TUI**: `internal/tui/` — Bubbletea explorer + wizards (see `internal/tui/AGENTS.md` — 382 lines)
- **Config**: YAML, cascade: global → spec → collection (`internal/config/`)
- **Auth**: **9 methods** in `internal/auth/` (none, basic, bearer, digest, hmac, oauth2-cc, oauth2-pwd, api-key, script)
- **Search**: bluge full-text engine (`internal/index/`)
- **IDs**: MD5-based (`internal/id/`)

## CLI reference

Full CLI documentation is in `.agents/skills/swag2mcp-cli/SKILL.md` (763 lines). Key points:

- **`--version`** flag is supported (same as `version` subcommand)
- **`mcp`** prints `"MCP server listening on http://<addr><path>"` on stdout on startup
- **`validate`** checks that `location` is a valid OpenAPI/Swagger/Postman file (not just any URL)
- **`info`** shows `max_response_size` in human-readable format (e.g. `"1 KB"`)
- **`export`/`import --from-zip`** supports full workspace round-trip

## MCP Tools (16 total)

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
| `info` | Get swag2mcp runtime info |

## Key conventions

- **Errors for LLM**: use `LLMError` with codes `validation_failed`, `not_found`, `rate_limit`, `invoke_error` — messages must explain what to do next
- **Rate limit**: 10s per endpoint for `invoke` (`internal/service/ratelimit.go`). Second call within 10s is silently blocked (no error returned to client)
- **Response size**: default 2KB max, saved to `{workspace}/responses/` when exceeded; `FileReference` returned
- **`--disable-llm-auth`**: removes `auth` tool from MCP tool list entirely (not just empty response)
- **`OAuth2PasswordAuthClient.ClientSecret`**: optional (public client support for Keycloak)
- **Workspace**: `~/.swag2mcp` with `cache/`, `specs/`, `responses/`, `auth_scripts/`; old responses cleaned after 48h on `mcp` start
- **Config validation**: `internal/config/validator.go` — human-readable messages for LLM; also checks that `location` is a valid OpenAPI/Swagger/Postman file
- **Test helpers**: `newTestService()` + `seedTestData()` in `internal/service/service_test.go`; `newTestIndex()` in `internal/index/index_test.go`
- **MCP handler tests**: use `go.uber.org/mock/gomock` — run `go generate ./...` after changing `svc` interface
- **Lint exclusions**: TUI, commands, config, service paths have relaxed rules in `.golangci.yml`
- **`.golangci.yml`**: uses version 2 format, 80+ linters enabled, extensive path-based exclusions

## Coverage strategy

- **Core packages** (`auth`, `cache`, `config`, `env`, `httpclient`, `id`, `index`, `server/mcp`, `service`, `spec`, `types`, `workspace`) — target 80%+
- **Integration packages** (`commands`, `tui`, `mockserver`) — inherently hard to unit test (cobra RunE, Bubbletea models, real HTTP servers). Coverage is informational only
- **Coveralls** runs in `informational` mode (`fail-on-error: false`) — never blocks PRs on coverage drop
- **`.coveralls.yml`** excludes `*_test.go`, `mocks/`, `mock_*`, `*.pb.go` from the report
- **Commands**: `make cover` (all packages, excludes commands/tui/mockserver), `make integration-tests` (tests/ package)
