# swag2mcp — Agent Guide

## Commands

```sh
make lint          # golangci-lint (120 line limit, strict linters)
make cover         # go test ./... + coverage HTML report
go test ./pkg/...  # single package
go test ./...      # all packages
go generate ./...  # installs mockgen + generates mocks (internal/server/mcp)
```

Order: `make lint && go test ./...`

## Architecture

- **Entrypoint**: `cmd/swag2mcp/main.go` — cobra CLI with 9 subcommands
- **Core**: `internal/service/` — business logic (Bootstrap, Invoke, Search, Inspect, Auth, Specs, etc.)
- **MCP server**: `internal/server/mcp/` — 14 MCP tools, uses `go.uber.org/mock` for tests
- **TUI**: `internal/tui/` — Bubbletea explorer + wizards (see `internal/tui/AGENTS.md`)
- **Config**: YAML, cascade: global → spec → collection (`internal/config/`)
- **Auth**: 8 methods in `internal/auth/` (none, basic, bearer, digest, oauth2-cc, oauth2-pwd, api-key, script)
- **Search**: bluge full-text engine (`internal/index/`)
- **IDs**: MD5-based (`internal/id/`)

## Key conventions

- **Errors for LLM**: use `LLMError` with codes `validation_failed`, `not_found`, `rate_limit`, `invoke_error` — messages must explain what to do next
- **Rate limit**: 10s per endpoint for `invoke` (`internal/service/ratelimit.go`)
- **Response size**: default 2KB max, saved to `{workspace}/responses/` when exceeded; `FileReference` returned
- **`--disable-llm-auth`**: removes `auth` tool from MCP tool list entirely (not just empty response)
- **`OAuth2PasswordAuthClient.ClientSecret`**: optional (public client support for Keycloak)
- **Workspace**: `~/.swag2mcp` with `cache/`, `specs/`, `responses/`, `auth_scripts/`; old responses cleaned after 48h on `mcp` start
- **Config validation**: `internal/config/validator.go` — human-readable messages for LLM
- **Test helpers**: `newTestService()` + `seedTestData()` in `internal/service/service_test.go`; `newTestIndex()` in `internal/index/index_test.go`
- **MCP handler tests**: use `go.uber.org/mock/gomock` — run `go generate ./...` after changing `svc` interface
- **Lint exclusions**: TUI, commands, config, service paths have relaxed rules in `.golangci.yml`
