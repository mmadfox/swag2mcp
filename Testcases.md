# swag2mcp — Manual Test Checklist

## 1. Installation & Build

- [ ] `go build ./cmd/swag2mcp/` — builds without errors
- [ ] `go build ./cmd/swag2mcp-mock/` — builds without errors
- [ ] `swag2mcp --help` — shows all 9 subcommands + flags
- [ ] `swag2mcp-mock --help` — shows mockserver subcommand + flags
- [ ] `swag2mcp version` (or `--version`) — prints version string

---

## 2. Workspace Initialization (`swag2mcp init`)

- [ ] `swag2mcp init` — creates `~/.swag2mcp/` with `cache/`, `specs/`, `responses/`, `auth_scripts/`, `swag2mcp.yaml`
- [ ] `swag2mcp init /custom/path` — creates workspace at custom path
- [ ] `swag2mcp init -i` — interactive wizard starts (18 states)
- [ ] `swag2mcp init -f` — force overwrite of existing config
- [ ] `swag2mcp init` on existing workspace without `-f` — shows error / no overwrite
- [ ] `swag2mcp init` — generated `swag2mcp.yaml` is valid YAML
- [ ] `swag2mcp init -i` — complete full wizard flow, verify config is written correctly

---

## 3. Configuration (`swag2mcp.yaml`)

### 3.1 Global Settings

- [ ] `http_client.random: true` — random browser-like headers are applied
- [ ] `http_client.timeout: 30s` — request times out after 30s
- [ ] `http_client.follow_redirects: false` — redirects are NOT followed
- [ ] `http_client.max_redirects: 5` — redirect limit works
- [ ] `http_client.max_response_size: 2048` — response truncated at 2KB
- [ ] `http_client.proxy.url` — requests go through HTTP proxy
- [ ] `http_client.proxy.username/password` — proxy auth works
- [ ] `http_client.proxy.bypass` — bypass list works (e.g. `localhost`)
- [ ] `http_client.headers` — custom headers added to every request
- [ ] `http_client.cookies` — custom cookies sent with every request
- [ ] `http_client.user_agent` — custom UA overrides default
- [ ] `mcp.transport: stdio` — MCP starts on stdio
- [ ] `mcp.transport: sse` — MCP starts SSE server on `:8080`
- [ ] `mcp.transport: streamable-http` — MCP starts streamable HTTP
- [ ] `mcp.auth.token: mytoken` — MCP HTTP endpoint requires Bearer token
- [ ] `$(ENV_VAR)` in any config field — resolved from environment

### 3.2 Spec Configuration

- [ ] `domain: my-api` — spec registered with domain `my-api`
- [ ] `llm_title: "My API"` — title appears in MCP tool descriptions
- [ ] `llm_instruction: "Use this for..."` — instruction appended to LLM prompt
- [ ] `base_url: https://api.example.com` — all requests go to this base
- [ ] `disable: true` — spec is excluded from MCP tools
- [ ] `tags: ["public"]` — spec filtered by `--tags public`
- [ ] `auth.type: bearer` + `auth.config.token: xxx` — auth applied to all endpoints
- [ ] `http_client` per spec — overrides global HTTP settings
- [ ] `base_url` per collection — overrides spec base_url

### 3.3 Collection Configuration

- [ ] `title: "Pets"` — collection appears with correct title
- [ ] `location: ./specs/petstore.yaml` — local file loaded
- [ ] `location: https://example.com/spec.yaml` — remote URL fetched + cached
- [ ] `disable: true` — collection excluded
- [ ] `llm_title` + `llm_instruction` per collection — overrides spec
- [ ] `base_mock_url: localhost:8081` — mock server uses this port
- [ ] `http_client` per collection — overrides spec and global

### 3.4 Config Validation

- [ ] `swag2mcp validate` — valid config reports no issues
- [ ] `swag2mcp validate` — duplicate domain detected
- [ ] `swag2mcp validate` — mock port conflict detected
- [ ] `swag2mcp validate` — unreachable spec location reported
- [ ] `swag2mcp validate` — invalid domain format (e.g. `UPPERCASE`, `spaces`, `>60 chars`)
- [ ] `swag2mcp validate` — invalid title length (<5 or >120 chars)
- [ ] `swag2mcp validate` — invalid instruction length (>500 chars)
- [ ] `swag2mcp validate` — invalid collection location (<5 or >250 chars)
- [ ] `swag2mcp validate` — invalid base_url format
- [ ] `swag2mcp validate -t public,internal` — filter validation by tags

---

## 4. CLI Commands

### 4.1 `swag2mcp add spec`

- [ ] `swag2mcp add spec` — interactive TUI wizard for adding a spec
- [ ] `swag2mcp add spec --yaml "..."` — non-interactive YAML import
- [ ] `swag2mcp add spec --yaml -` — YAML piped from stdin
- [ ] `swag2mcp add spec --example` — example spec added
- [ ] `swag2mcp add spec` with invalid YAML — error message shown
- [ ] `swag2mcp add spec` — config file atomically updated

### 4.2 `swag2mcp add collection`

- [ ] `swag2mcp add collection` — interactive TUI wizard
- [ ] `swag2mcp add collection --yaml "..."` — non-interactive YAML import
- [ ] `swag2mcp add collection --yaml -` — YAML piped from stdin
- [ ] `swag2mcp add collection` — collection added to existing spec
- [ ] `swag2mcp add collection` with no specs in config — error / empty state handled

### 4.3 `swag2mcp delete spec`

- [ ] `swag2mcp delete spec` — interactive selection, spec removed
- [ ] `swag2mcp delete spec` — confirm dialog works (yes/no)
- [ ] `swag2mcp delete spec` — cancel does not modify config
- [ ] `swag2mcp delete spec` with no specs — error / empty state handled

### 4.4 `swag2mcp delete collection`

- [ ] `swag2mcp delete collection` — select spec → select collection → confirm → removed
- [ ] `swag2mcp delete collection` — cancel at any step does not modify config
- [ ] `swag2mcp delete collection` with no collections — error / empty state handled

### 4.5 `swag2mcp ls`

- [ ] `swag2mcp ls` — shows all specs and collections in formatted table
- [ ] `swag2mcp ls -t public` — filters by tag
- [ ] `swag2mcp ls -t public,internal` — multiple tags
- [ ] `swag2mcp ls` with no specs — shows empty table / message
- [ ] `swag2mcp ls` — columns: domain, title, baseURL, tags, auth type, collections

### 4.6 `swag2mcp run` (TUI Explorer)

- [ ] `swag2mcp run` — TUI starts with 4 menu options
- [ ] **Search mode**: enter query → paginated results (10/page) → select → endpoint detail
- [ ] **Search mode**: empty query — shows all / error
- [ ] **Browse mode**: Specs → Collections → Tags → Endpoints → endpoint detail
- [ ] **Browse mode**: empty spec (no collections) — handled
- [ ] **Auth mode**: select spec → confirm → view token/headers/query params
- [ ] **Auth mode**: spec with no auth — shows appropriate message
- [ ] **Save endpoint**: `[S]` saves JSON file to current directory
- [ ] **Navigation**: `[B]ack`, `[M]enu`, `Esc`, `Ctrl+C` all work
- [ ] **Pagination**: `N`/`P` keys navigate pages
- [ ] Schema rendering: properties, types, required fields, enums, examples displayed
- [ ] `swag2mcp run` with no specs — error / empty state handled

### 4.7 `swag2mcp update`

- [ ] `swag2mcp update` — validates config, clears cache, re-caches all specs
- [ ] `swag2mcp update` — orphan auth scripts cleaned
- [ ] `swag2mcp update` with invalid config — validation errors shown, update stops
- [ ] `swag2mcp update` — remote specs re-downloaded to cache

### 4.8 `swag2mcp clean`

- [ ] `swag2mcp clean` — `cache/` contents removed
- [ ] `swag2mcp clean` — `responses/` contents removed
- [ ] `swag2mcp clean` — orphan auth scripts removed
- [ ] `swag2mcp clean` — `specs/` and `auth_scripts/` (non-orphan) preserved

### 4.9 `swag2mcp mcp`

- [ ] `swag2mcp mcp` — starts MCP server on stdio (default)
- [ ] `swag2mcp mcp --transport sse` — starts SSE server
- [ ] `swag2mcp mcp --transport streamable-http` — starts streamable HTTP
- [ ] `swag2mcp mcp --http-addr :9090` — custom address
- [ ] `swag2mcp mcp --http-path /custom-mcp` — custom path
- [ ] `swag2mcp mcp --auth-token secret` — Bearer token auth on HTTP
- [ ] `swag2mcp mcp --disable-llm-auth` — `auth` tool removed from tool list
- [ ] `swag2mcp mcp --dump-dir /tmp/dumps` — HTTP requests dumped to directory
- [ ] `swag2mcp mcp --logfile /tmp/mcp.log` — logs written to file
- [ ] `swag2mcp mcp -t public` — only specs with tag `public` are loaded
- [ ] `swag2mcp mcp` — old responses (>48h) cleaned on startup

---

## 5. MCP Tools

### 5.1 `spec_list`

- [ ] Returns all specs with correct IDs and domains
- [ ] Returns empty list when no specs configured
- [ ] Returns only tag-filtered specs when `--tags` used

### 5.2 `spec_by_id`

- [ ] Returns spec details + collections for valid ID
- [ ] Returns `not_found` error for non-existent ID
- [ ] Returns `not_found` error for empty ID
- [ ] Returns `not_found` error for malformed ID (not 32-char hex)

### 5.3 `collection_by_spec`

- [ ] Returns all collections for valid specId
- [ ] Returns `not_found` for non-existent specId
- [ ] Returns empty list for spec with no collections

### 5.4 `collection_by_id`

- [ ] Returns collection details + tags for valid ID
- [ ] Returns `not_found` for non-existent ID
- [ ] Returns `not_found` for malformed ID

### 5.5 `tag_by_spec`

- [ ] Returns all tags across spec for valid specId
- [ ] Returns `not_found` for non-existent specId
- [ ] Returns empty list for spec with no tags

### 5.6 `tag_by_collection`

- [ ] Returns all tags for valid collectionId
- [ ] Returns `not_found` for non-existent collectionId
- [ ] Returns empty list for collection with no tags

### 5.7 `tag_by_id`

- [ ] Returns tag details for valid ID
- [ ] Returns `not_found` for non-existent ID

### 5.8 `endpoint_by_spec`

- [ ] Returns all endpoints across spec for valid specId
- [ ] Returns `not_found` for non-existent specId
- [ ] Returns empty list for spec with no endpoints

### 5.9 `endpoint_by_collection`

- [ ] Returns all endpoints for valid collectionId
- [ ] Returns `not_found` for non-existent collectionId
- [ ] Returns empty list for collection with no endpoints

### 5.10 `endpoint_by_tag`

- [ ] Returns all endpoints for valid tagId
- [ ] Returns `not_found` for non-existent tagId
- [ ] Returns empty list for tag with no endpoints

### 5.11 `endpoint_by_id`

- [ ] Returns endpoint summary (method, path, summary, deprecated) for valid ID
- [ ] Returns `not_found` for non-existent ID
- [ ] Deprecated endpoint shows `deprecated: true`

### 5.12 `search`

- [ ] `search("pet")` — returns matching endpoints
- [ ] `search("method:GET")` — only GET endpoints
- [ ] `search("tag:auth")` — only auth-tagged endpoints
- [ ] `search("path:/api/v1/users")` — exact path match
- [ ] `search("+method:POST +summary:create")` — boolean AND
- [ ] `search("summary:\"create user\"")` — phrase search
- [ ] `search("sumary~")` — fuzzy search (typo tolerance)
- [ ] `search("cr*")` — wildcard search
- [ ] `search("zzzzz")` — empty results
- [ ] `search("*")` — returns all endpoints
- [ ] `search("pet", limit=1)` — returns exactly 1 result
- [ ] `search("pet", limit=50)` — returns up to 50 results
- [ ] `search("pet", limit=0)` — error (min 1)
- [ ] `search("pet", limit=51)` — error (max 50)

### 5.13 `inspect`

- [ ] Returns full operation object for valid endpointId
- [ ] Parameters (path, query, header) with schemas are present
- [ ] Request body schema is present (for POST/PUT/PATCH)
- [ ] Response schemas with status codes are present
- [ ] Referenced `$ref` schemas are resolved
- [ ] Returns `not_found` for non-existent endpointId

### 5.14 `invoke`

- [ ] `invoke` on GET endpoint — returns response with status code, headers, body
- [ ] `invoke` with path parameters — URL correctly interpolated
- [ ] `invoke` with query parameters — query string correctly built
- [ ] `invoke` with header parameters — headers sent
- [ ] `invoke` with requestBody — JSON body sent
- [ ] `invoke` on POST/PUT/PATCH — request body sent correctly
- [ ] `invoke` on DELETE — request sent (requires explicit user confirmation in LLM)
- [ ] `invoke` with invalid endpointId — `not_found` error
- [ ] `invoke` on non-existent server — `invoke_error` with connection refused
- [ ] `invoke` on 4xx/5xx response — status code and error body returned
- [ ] `invoke` same endpoint twice within 10s — `rate_limit` error with retry-after message
- [ ] `invoke` same endpoint after 10s wait — succeeds
- [ ] `invoke` with response >1KB (default) — body truncated, `FileReference` returned
- [ ] `invoke` with response >configured `max_response_size` — saved to `responses/`
- [ ] `invoke` with response >1MB — truncated at 1MB max

### 5.15 `auth`

- [ ] `auth(specId)` — returns token/headers/query params for valid spec
- [ ] `auth(specId)` with `--disable-llm-auth` — tool not present in list
- [ ] `auth(specId)` for non-existent specId — `not_found` error
- [ ] `auth(specId)` for spec with `auth.type: none` — returns empty / no-auth

---

## 6. Auth Methods

### 6.1 None

- [ ] Requests sent without any auth headers
- [ ] `auth` tool returns empty/no-auth response

### 6.2 Basic

- [ ] `Authorization: Basic base64(user:pass)` header sent
- [ ] Wrong credentials — 401 returned from server
- [ ] `$(ENV_VAR)` in username/password — resolved from environment

### 6.3 Bearer

- [ ] `Authorization: Bearer <token>` header sent
- [ ] Invalid/expired token — 401 returned
- [ ] `$(ENV_VAR)` in token — resolved from environment

### 6.4 Digest

- [ ] Full MD5 digest auth flow: challenge → response with nonce, cnonce, qop
- [ ] Nonce cached for 5 minutes (subsequent requests reuse)
- [ ] Nonce expired — new challenge fetched
- [ ] Wrong credentials — 401 after digest attempt
- [ ] `$(ENV_VAR)` in username/password — resolved

### 6.5 OAuth2 Client Credentials

- [ ] Token obtained from `token_url` using client_id + client_secret
- [ ] Token cached and reused until expiry
- [ ] Expired token — new token fetched automatically
- [ ] `Authorization: Bearer <token>` header sent
- [ ] `scopes` included in token request
- [ ] Invalid credentials — error returned
- [ ] `$(ENV_VAR)` in fields — resolved

### 6.6 OAuth2 Password

- [ ] Token obtained from `token_url` using username + password + client_id
- [ ] `client_secret` optional (public client — Keycloak support)
- [ ] Token cached and reused until expiry
- [ ] Expired token — new token fetched
- [ ] Invalid credentials — error returned
- [ ] `$(ENV_VAR)` in fields — resolved

### 6.7 API Key

- [ ] `in: header` — key placed in request header
- [ ] `in: query` — key placed in query parameter
- [ ] Wrong key — 401 returned
- [ ] `$(ENV_VAR)` in key/value — resolved

### 6.8 Script

- [ ] `{workspace}/auth_scripts/{domain}.sh` executed
- [ ] Script output JSON `{"token":"...","expires_in":N}` parsed correctly
- [ ] Token cached and reused until expiry
- [ ] Script returns non-zero exit — error returned
- [ ] Script returns invalid JSON — error returned
- [ ] Script file does not exist — error returned
- [ ] `.bat` script on Windows (if applicable)

---

## 7. Spec Parsing

### 7.1 OpenAPI 3.x

- [ ] Paths, operations, parameters parsed correctly
- [ ] Request bodies with JSON schema parsed
- [ ] Response schemas with status codes parsed
- [ ] `$ref` references resolved
- [ ] Tags extracted from spec
- [ ] Enums, examples, descriptions preserved

### 7.2 Swagger 2.0

- [ ] Paths, operations, parameters parsed correctly
- [ ] `definitions` resolved for `$ref`
- [ ] Tags extracted
- [ ] All HTTP methods supported

### 7.3 Postman Collections

- [ ] Collection items parsed as endpoints
- [ ] Request methods, URLs, headers extracted
- [ ] Request body (raw, JSON, form-data) parsed
- [ ] Auth defined in collection applied

### 7.4 Invalid / Edge Cases

- [ ] Invalid YAML/JSON spec — error reported
- [ ] Spec with no paths — empty endpoints list
- [ ] Spec with no tags — endpoints grouped under default
- [ ] Spec with circular `$ref` — handled without infinite loop
- [ ] Remote spec URL returns 404 — error reported
- [ ] Remote spec URL times out — error reported

---

## 8. Mock Server (`swag2mcp-mock`)

- [ ] `swag2mcp-mock mockserver` — starts mock servers for all specs
- [ ] `swag2mcp-mock mockserver --tls` — starts with TLS (self-signed)
- [ ] `swag2mcp-mock mockserver --tls-cert cert.pem --tls-key key.pem` — custom TLS cert
- [ ] Mock server responds to requests on configured ports
- [ ] Mock OAuth2 server on port 9090 — returns valid tokens
- [ ] Mock Digest server on port 9091 — handles digest auth flow
- [ ] `base_mock_url` per collection — mock uses correct port
- [ ] Mock server returns realistic responses based on spec
- [ ] Multiple specs — each gets its own mock server
- [ ] `invoke` against mock server — end-to-end flow works

---

## 9. Workspace Management

- [ ] `~/.swag2mcp/` created with all subdirectories
- [ ] `cache/` stores downloaded remote specs
- [ ] `specs/` stores local spec files
- [ ] `responses/` stores large invocation responses
- [ ] `auth_scripts/` stores custom auth scripts
- [ ] `swag2mcp clean` — `cache/` and `responses/` emptied
- [ ] Old responses (>48h) cleaned on `swag2mcp mcp` startup
- [ ] Old responses (<48h) preserved on `swag2mcp mcp` startup
- [ ] Orphan auth scripts (no matching spec domain) cleaned on `update` and `clean`

---

## 10. Error Handling

- [ ] `not_found` error — JSON with code, message, hint
- [ ] `validation_failed` error — actionable message
- [ ] `rate_limit` error — "try again in X seconds" message
- [ ] `invoke_error` error — connection/HTTP error details
- [ ] All errors serialized as valid JSON
- [ ] Error messages include guidance for LLM on what to do next

---

## 11. Cross-Cutting / Integration

- [ ] Full workflow: `init` → `add spec` → `add collection` → `validate` → `mcp` → `spec_list` → `search` → `inspect` → `invoke`
- [ ] Config cascade: global timeout → spec timeout → collection timeout (most specific wins)
- [ ] Tag filtering: `--tags public` on `mcp` — only public-tagged specs loaded
- [ ] `disable: true` on spec — spec excluded from all tools
- [ ] `disable: true` on collection — collection excluded
- [ ] Multiple specs with different auth types — each works independently
- [ ] Multiple collections per spec — all accessible
- [ ] `swag2mcp update` after changing spec file — changes reflected
- [ ] `swag2mcp update` after adding new spec file — new spec available
- [ ] `swag2mcp update` after removing spec file — spec removed from index
- [ ] MCP server restart with `--tags` — only filtered specs available
- [ ] Concurrent `invoke` requests to different endpoints — both succeed
- [ ] Concurrent `invoke` requests to same endpoint — rate limit applies per endpoint
