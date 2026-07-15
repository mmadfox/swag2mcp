# swag2mcp — Manual Test Checklist

## 1. Installation & Build

- [x] `go build ./cmd/swag2mcp/` — builds without errors (integration-test, main_test.go, TestMain)
- [x] `go build ./cmd/swag2mcp-mock/` — builds without errors (integration-test, main_test.go, TestMain)
- [x] `swag2mcp --help` — shows all 9 subcommands + flags (manual)
- [x] `swag2mcp-mock --help` — shows mockserver subcommand + flags (manual)
- [x] `swag2mcp version` (or `--version`) — prints version string (manual)

---

## 2. Workspace Initialization (`swag2mcp init`)

- [x] `swag2mcp init` — creates `~/.swag2mcp/` with subdirectories (integration-test, suite_init_test.go, TestScript_Init_CreatesWorkspace)
- [x] `swag2mcp init /custom/path` — creates workspace at custom path (integration-test, suite_init_test.go, TestScript_Init_CustomPath)
- [x] `swag2mcp init -i` — interactive wizard starts (18 states) (manual — requires TTY)
- [x] `swag2mcp init -f` — force overwrite of existing config (integration-test, suite_init_test.go, TestScript_Init_ForceOverwrite)
- [x] `swag2mcp init` on existing workspace without `-f` — shows error / no overwrite (integration-test, suite_init_test.go, TestScript_Init_ForceOverwrite)
- [x] `swag2mcp init` — generated `swag2mcp.yaml` is valid YAML (integration-test, suite_init_test.go, TestScript_Init_CreatesWorkspace)
- [x] `swag2mcp init -i` — complete full wizard flow, verify config is written correctly (manual — requires TTY)

---

## 3. Configuration (`swag2mcp.yaml`)

### 3.1 Global Settings

- [ ] `http_client.random: true` — random browser-like headers are applied (not covered)
- [ ] `http_client.timeout: 30s` — request times out after 30s (not covered)
- [ ] `http_client.follow_redirects: false` — redirects are NOT followed (not covered)
- [ ] `http_client.max_redirects: 5` — redirect limit works (not covered)
- [x] `http_client.max_response_size: 2048` — response truncated at 2KB (integration-test, suite_response_test.go, TestScript_ResponseSize_Configurable)
- [ ] `http_client.proxy.url` — requests go through HTTP proxy (not covered)
- [ ] `http_client.proxy.username/password` — proxy auth works (not covered)
- [ ] `http_client.proxy.bypass` — bypass list works (e.g. `localhost`) (not covered)
- [ ] `http_client.headers` — custom headers added to every request (not covered)
- [ ] `http_client.cookies` — custom cookies sent with every request (not covered)
- [ ] `http_client.user_agent` — custom UA overrides default (not covered)
- [x] `mcp.transport: stdio` — MCP starts on stdio (integration-test, suite_transport_test.go, TestScript_Transport_Stdio)
- [x] `mcp.transport: sse` — MCP starts SSE server on `:8080` (integration-test, suite_transport_test.go, TestScript_Transport_SSE)
- [x] `mcp.transport: streamable-http` — MCP starts streamable HTTP (integration-test, suite_transport_test.go, TestScript_Transport_StreamableHTTP)
- [x] `mcp.auth.token: mytoken` — MCP HTTP endpoint requires Bearer token (integration-test, suite_transport_test.go, TestScript_Transport_AuthToken)
- [x] `$(ENV_VAR)` in any config field — resolved from environment (integration-test, suite_config_test.go, TestScript_EnvVarResolution)

### 3.2 Spec Configuration

- [x] `domain: my-api` — spec registered with domain `my-api` (integration-test, suite_mcp_tools_test.go, TestScript_MCP_SpecByID)
- [x] `llm_title: "My API"` — title appears in MCP tool descriptions (integration-test, suite_mcp_tools_test.go, TestScript_MCP_SpecList)
- [x] `llm_instruction: "Use this for..."` — instruction appended to LLM prompt (integration-test, suite_mcp_tools_test.go, TestScript_MCP_SpecList)
- [x] `base_url: https://api.example.com` — all requests go to this base (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Invoke)
- [ ] `disable: true` — spec is excluded from MCP tools (not covered)
- [x] `tags: ["public"]` — spec filtered by `--tags public` (integration-test, suite_mcp_tools_test.go, TestScript_MCP_TagFilter)
- [x] `auth.type: bearer` + `auth.config.token: xxx` — auth applied to all endpoints (integration-test, suite_auth_test.go, TestScript_Auth_InvokeWithBearer)
- [x] `http_client` per spec — overrides global HTTP settings (integration-test, suite_config_test.go, TestScript_ConfigCascade)
- [ ] `base_url` per collection — overrides spec base_url (not covered)

### 3.3 Collection Configuration

- [x] `title: "Pets"` — collection appears with correct title (integration-test, suite_mcp_tools_test.go, TestScript_MCP_CollectionByID)
- [x] `location: ./specs/petstore.yaml` — local file loaded (integration-test, suite_mcp_tools_test.go, TestScript_MCP_SpecList)
- [ ] `location: https://example.com/spec.yaml` — remote URL fetched + cached (not covered)
- [ ] `disable: true` — collection excluded (not covered)
- [x] `llm_title` + `llm_instruction` per collection — overrides spec (integration-test, suite_mcp_tools_test.go, TestScript_MCP_CollectionByID)
- [ ] `base_mock_url: localhost:8081` — mock server uses this port (not covered)
- [x] `http_client` per collection — overrides spec and global (integration-test, suite_config_test.go, TestScript_ConfigCascade)

### 3.4 Config Validation

- [x] `swag2mcp validate` — valid config reports no issues (integration-test, suite_config_test.go, TestScript_Validate_ValidConfig)
- [x] `swag2mcp validate` — duplicate domain detected (integration-test, suite_config_test.go, TestScript_Validate_DuplicateDomain)
- [ ] `swag2mcp validate` — mock port conflict detected (not covered)
- [x] `swag2mcp validate` — unreachable spec location reported (integration-test, suite_config_test.go, TestScript_Validate_UnreachableLocation)
- [x] `swag2mcp validate` — invalid domain format (e.g. `UPPERCASE`, `spaces`, `>60 chars`) (integration-test, suite_config_test.go, TestScript_Validate_InvalidDomainFormat)
- [ ] `swag2mcp validate` — invalid title length (<5 or >120 chars) (not covered)
- [ ] `swag2mcp validate` — invalid instruction length (>500 chars) (not covered)
- [ ] `swag2mcp validate` — invalid collection location (<5 or >250 chars) (not covered)
- [ ] `swag2mcp validate` — invalid base_url format (not covered)
- [x] `swag2mcp validate -t public,internal` — filter validation by tags (integration-test, suite_config_test.go, TestScript_Validate_TagFilter)

---

## 4. CLI Commands

### 4.1 `swag2mcp add spec`

- [ ] `swag2mcp add spec` — interactive TUI wizard for adding a spec (manual — requires TTY)
- [x] `swag2mcp add spec --yaml "..."` — non-interactive YAML import (integration-test, suite_config_test.go, TestScript_AddSpec_FromYAML)
- [x] `swag2mcp add spec --yaml -` — YAML piped from stdin (integration-test, suite_config_test.go, TestScript_AddSpec_FromStdin)
- [ ] `swag2mcp add spec --example` — example spec added (not covered)
- [x] `swag2mcp add spec` with invalid YAML — error message shown (integration-test, suite_config_test.go, TestScript_AddSpec_InvalidYAML)
- [x] `swag2mcp add spec` — config file atomically updated (integration-test, suite_config_test.go, TestScript_AddSpec_FromYAML)

### 4.2 `swag2mcp add collection`

- [ ] `swag2mcp add collection` — interactive TUI wizard (manual — requires TTY)
- [x] `swag2mcp add collection --yaml "..."` — non-interactive YAML import (integration-test, suite_config_test.go, TestScript_AddCollection_FromYAML)
- [ ] `swag2mcp add collection --yaml -` — YAML piped from stdin (not covered)
- [x] `swag2mcp add collection` — collection added to existing spec (integration-test, suite_config_test.go, TestScript_AddCollection_FromYAML)
- [ ] `swag2mcp add collection` with no specs in config — error / empty state handled (not covered)

### 4.3 `swag2mcp delete spec`

- [x] `swag2mcp delete spec` — interactive selection, spec removed (integration-test, suite_config_test.go, TestScript_DeleteSpec)
- [x] `swag2mcp delete spec` — confirm dialog works (yes/no) (integration-test, suite_config_test.go, TestScript_DeleteSpec)
- [x] `swag2mcp delete spec` — cancel does not modify config (integration-test, suite_config_test.go, TestScript_DeleteSpec_Cancel)
- [ ] `swag2mcp delete spec` with no specs — error / empty state handled (not covered)

### 4.4 `swag2mcp delete collection`

- [ ] `swag2mcp delete collection` — select spec → select collection → confirm → removed (not covered)
- [ ] `swag2mcp delete collection` — cancel at any step does not modify config (not covered)
- [ ] `swag2mcp delete collection` with no collections — error / empty state handled (not covered)

### 4.5 `swag2mcp ls`

- [x] `swag2mcp ls` — shows all specs and collections in formatted table (integration-test, suite_config_test.go, TestScript_ListSpecs)
- [x] `swag2mcp ls -t public` — filters by tag (integration-test, suite_config_test.go, TestScript_ListSpecs_TagFilter)
- [ ] `swag2mcp ls -t public,internal` — multiple tags (not covered)
- [x] `swag2mcp ls` with no specs — shows empty table / message (integration-test, suite_config_test.go, TestScript_ListSpecs_Empty)
- [x] `swag2mcp ls` — columns: domain, title, baseURL, tags, auth type, collections (integration-test, suite_config_test.go, TestScript_ListSpecs)

### 4.6 `swag2mcp run` (TUI Explorer)

- [ ] `swag2mcp run` — TUI starts with 4 menu options (manual — requires TTY)
- [ ] **Search mode**: enter query → paginated results (10/page) → select → endpoint detail (manual — requires TTY)
- [ ] **Search mode**: empty query — shows all / error (manual — requires TTY)
- [ ] **Browse mode**: Specs → Collections → Tags → Endpoints → endpoint detail (manual — requires TTY)
- [ ] **Browse mode**: empty spec (no collections) — handled (manual — requires TTY)
- [ ] **Auth mode**: select spec → confirm → view token/headers/query params (manual — requires TTY)
- [ ] **Auth mode**: spec with no auth — shows appropriate message (manual — requires TTY)
- [ ] **Save endpoint**: `[S]` saves JSON file to current directory (manual — requires TTY)
- [ ] **Navigation**: `[B]ack`, `[M]enu`, `Esc`, `Ctrl+C` all work (manual — requires TTY)
- [ ] **Pagination**: `N`/`P` keys navigate pages (manual — requires TTY)
- [ ] Schema rendering: properties, types, required fields, enums, examples displayed (manual — requires TTY)
- [ ] `swag2mcp run` with no specs — error / empty state handled (manual — requires TTY)

### 4.7 `swag2mcp update`

- [x] `swag2mcp update` — validates config, clears cache, re-caches all specs (integration-test, suite_config_test.go, TestScript_Update_ReCachesSpecs)
- [ ] `swag2mcp update` — orphan auth scripts cleaned (not covered)
- [x] `swag2mcp update` with invalid config — validation errors shown, update stops (integration-test, suite_config_test.go, TestScript_Update_InvalidConfig)
- [ ] `swag2mcp update` — remote specs re-downloaded to cache (not covered)

### 4.8 `swag2mcp clean`

- [x] `swag2mcp clean` — `cache/` contents removed (integration-test, suite_config_test.go, TestScript_Clean_RemovesCache)
- [x] `swag2mcp clean` — `responses/` contents removed (integration-test, suite_config_test.go, TestScript_Clean_RemovesCache)
- [ ] `swag2mcp clean` — orphan auth scripts removed (not covered)
- [x] `swag2mcp clean` — `specs/` and `auth_scripts/` (non-orphan) preserved (integration-test, suite_config_test.go, TestScript_Clean_PreservesSpecs)

### 4.9 `swag2mcp mcp`

- [x] `swag2mcp mcp` — starts MCP server on stdio (default) (integration-test, suite_transport_test.go, TestScript_Transport_Stdio)
- [x] `swag2mcp mcp --transport sse` — starts SSE server (integration-test, suite_transport_test.go, TestScript_Transport_SSE)
- [x] `swag2mcp mcp --transport streamable-http` — starts streamable HTTP (integration-test, suite_transport_test.go, TestScript_Transport_StreamableHTTP)
- [ ] `swag2mcp mcp --http-addr :9090` — custom address (not covered)
- [ ] `swag2mcp mcp --http-path /custom-mcp` — custom path (not covered)
- [x] `swag2mcp mcp --auth-token secret` — Bearer token auth on HTTP (integration-test, suite_transport_test.go, TestScript_Transport_AuthToken)
- [x] `swag2mcp mcp --disable-llm-auth` — `auth` tool removed from tool list (integration-test, suite_mcp_tools_test.go, TestScript_MCP_SpecList_NoAuthTool)
- [x] `swag2mcp mcp --dump-dir /tmp/dumps` — HTTP requests dumped to directory (integration-test, suite_transport_test.go, TestScript_Transport_DumpDir)
- [ ] `swag2mcp mcp --logfile /tmp/mcp.log` — logs written to file (not covered)
- [x] `swag2mcp mcp -t public` — only specs with tag `public` are loaded (integration-test, suite_mcp_tools_test.go, TestScript_MCP_TagFilter)
- [x] `swag2mcp mcp` — old responses (>48h) cleaned on startup (integration-test, suite_workspace_test.go, TestScript_Workspace_OldResponsesCleaned)

---

## 5. MCP Tools

### 5.1 `spec_list`

- [x] Returns all specs with correct IDs and domains (integration-test, suite_mcp_tools_test.go, TestScript_MCP_SpecList)
- [x] Returns list with at least 1 spec when configured (integration-test, suite_mcp_tools_test.go, TestScript_MCP_SpecList_Empty)
- [x] Returns only tag-filtered specs when `--tags` used (integration-test, suite_mcp_tools_test.go, TestScript_MCP_TagFilter)

### 5.2 `spec_by_id`

- [x] Returns spec details + collections for valid ID (integration-test, suite_mcp_tools_test.go, TestScript_MCP_SpecByID)
- [x] Returns error for non-existent ID (integration-test, suite_mcp_tools_test.go, TestScript_MCP_SpecByID_NotFound)
- [x] Returns error for empty ID (integration-test, suite_errors_test.go, TestScript_Errors_EmptyID)
- [x] Returns error for malformed ID (not 32-char hex) (integration-test, suite_errors_test.go, TestScript_Errors_InvalidID)

### 5.3 `collection_by_spec`

- [x] Returns all collections for valid specId (integration-test, suite_mcp_tools_test.go, TestScript_MCP_CollectionBySpec)
- [ ] Returns `not_found` for non-existent specId (not covered)
- [ ] Returns empty list for spec with no collections (not covered)

### 5.4 `collection_by_id`

- [x] Returns collection details + tags for valid ID (integration-test, suite_mcp_tools_test.go, TestScript_MCP_CollectionByID)
- [ ] Returns `not_found` for non-existent ID (not covered)
- [ ] Returns `not_found` for malformed ID (not covered)

### 5.5 `tag_by_spec`

- [x] Returns all tags across spec for valid specId (integration-test, suite_mcp_tools_test.go, TestScript_MCP_TagBySpec)
- [ ] Returns `not_found` for non-existent specId (not covered)
- [ ] Returns empty list for spec with no tags (not covered)

### 5.6 `tag_by_collection`

- [ ] Returns all tags for valid collectionId (not covered)
- [ ] Returns `not_found` for non-existent collectionId (not covered)
- [ ] Returns empty list for collection with no tags (not covered)

### 5.7 `tag_by_id`

- [ ] Returns tag details for valid ID (not covered)
- [ ] Returns `not_found` for non-existent ID (not covered)

### 5.8 `endpoint_by_spec`

- [x] Returns all endpoints across spec for valid specId (integration-test, suite_mcp_tools_test.go, TestScript_MCP_EndpointBySpec)
- [ ] Returns `not_found` for non-existent specId (not covered)
- [ ] Returns empty list for spec with no endpoints (not covered)

### 5.9 `endpoint_by_collection`

- [ ] Returns all endpoints for valid collectionId (not covered)
- [ ] Returns `not_found` for non-existent collectionId (not covered)
- [ ] Returns empty list for collection with no endpoints (not covered)

### 5.10 `endpoint_by_tag`

- [ ] Returns all endpoints for valid tagId (not covered)
- [ ] Returns `not_found` for non-existent tagId (not covered)
- [ ] Returns empty list for tag with no endpoints (not covered)

### 5.11 `endpoint_by_id`

- [x] Returns endpoint summary (method, path, summary, deprecated) for valid ID (integration-test, suite_mcp_tools_test.go, TestScript_MCP_EndpointByID)
- [ ] Returns `not_found` for non-existent ID (not covered)
- [ ] Deprecated endpoint shows `deprecated: true` (not covered)

### 5.12 `search`

- [x] `search("pet")` — returns matching endpoints (integration-test, suite_search_test.go, TestScript_Search_Basic)
- [x] `search("method:GET")` — only GET endpoints (integration-test, suite_search_test.go, TestScript_Search_ByMethod)
- [x] `search("tag:pets")` — only tagged endpoints (integration-test, suite_search_test.go, TestScript_Search_ByTag)
- [x] `search("path:/pets")` — path match (integration-test, suite_search_test.go, TestScript_Search_ByPath)
- [x] `search("+method:GET +summary:pet")` — boolean AND (integration-test, suite_search_test.go, TestScript_Search_BooleanAND)
- [ ] `search("summary:\"create user\"")` — phrase search (not covered)
- [ ] `search("sumary~")` — fuzzy search (typo tolerance) (not covered)
- [x] `search("list*")` — wildcard search (integration-test, suite_search_test.go, TestScript_Search_Wildcard)
- [x] `search("zzzzz")` — empty results (integration-test, suite_search_test.go, TestScript_Search_EmptyResults)
- [x] `search("*")` — returns all endpoints (integration-test, suite_search_test.go, TestScript_Search_AllEndpoints)
- [x] `search("pet", limit=1)` — returns exactly 1 result (integration-test, suite_search_test.go, TestScript_Search_LimitBounds)
- [x] `search("pet", limit=50)` — returns up to 50 results (integration-test, suite_search_test.go, TestScript_Search_AllEndpoints)
- [ ] `search("pet", limit=0)` — error (min 1) (not covered)
- [ ] `search("pet", limit=51)` — error (max 50) (not covered)

### 5.13 `inspect`

- [x] Returns full operation object for valid endpointId (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Inspect)
- [x] Parameters (path, query, header) with schemas are present (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Inspect)
- [x] Request body schema is present (for POST/PUT/PATCH) (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Inspect)
- [x] Response schemas with status codes are present (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Inspect)
- [x] Referenced `$ref` schemas are resolved (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Inspect)
- [ ] Returns `not_found` for non-existent endpointId (not covered)

### 5.14 `invoke`

- [x] `invoke` on GET endpoint — returns response with status code, headers, body (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Invoke)
- [x] `invoke` with path parameters — URL correctly interpolated (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Invoke_WithPathParams)
- [x] `invoke` with query parameters — query string correctly built (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Invoke_WithQueryParams)
- [ ] `invoke` with header parameters — headers sent (not covered)
- [x] `invoke` with requestBody — JSON body sent (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Invoke_WithRequestBody)
- [x] `invoke` on POST — request body sent correctly (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Invoke_WithRequestBody)
- [ ] `invoke` on DELETE — request sent (requires explicit user confirmation in LLM) (not covered)
- [ ] `invoke` with invalid endpointId — `not_found` error (not covered)
- [x] `invoke` on non-existent server — `invoke_error` with connection refused (integration-test, suite_errors_test.go, TestScript_Errors_InvokeConnectionRefused)
- [x] `invoke` on 5xx response — status code and error body returned (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Invoke_ServerError)
- [x] `invoke` same endpoint twice within 10s — `rate_limit` error (integration-test, suite_ratelimit_test.go, TestScript_RateLimit_BlocksSecondCall)
- [x] `invoke` same endpoint after 10s wait — succeeds (integration-test, suite_ratelimit_test.go, TestScript_RateLimit_RecoversAfterWait)
- [x] `invoke` with response >1KB (default) — body truncated, `FileReference` returned (integration-test, suite_response_test.go, TestScript_ResponseSize_DefaultLimit)
- [x] `invoke` with response >configured `max_response_size` — saved to `responses/` (integration-test, suite_response_test.go, TestScript_ResponseSize_FileReference)
- [ ] `invoke` with response >1MB — truncated at 1MB max (not covered)

### 5.15 `auth`

- [x] `auth(specId)` — returns token/headers/query params for valid spec (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Auth)
- [x] `auth(specId)` with `--disable-llm-auth` — tool not present in list (integration-test, suite_mcp_tools_test.go, TestScript_MCP_SpecList_NoAuthTool)
- [ ] `auth(specId)` for non-existent specId — `not_found` error (not covered)
- [x] `auth(specId)` for spec with `auth.type: none` — returns empty / no-auth (integration-test, suite_auth_test.go, TestScript_Auth_None)

---

## 6. Auth Methods

### 6.1 None

- [x] Requests sent without any auth headers (integration-test, suite_auth_test.go, TestScript_Auth_None)
- [x] `auth` tool returns empty/no-auth response (integration-test, suite_auth_test.go, TestScript_Auth_None)

### 6.2 Basic

- [x] `Authorization: Basic base64(user:pass)` header sent (integration-test, suite_auth_test.go, TestScript_Auth_Basic)
- [ ] Wrong credentials — 401 returned from server (not covered)
- [ ] `$(ENV_VAR)` in username/password — resolved from environment (not covered)

### 6.3 Bearer

- [x] `Authorization: Bearer <token>` header sent (integration-test, suite_auth_test.go, TestScript_Auth_Bearer)
- [ ] Invalid/expired token — 401 returned (not covered)
- [x] `$(ENV_VAR)` in token — resolved from environment (integration-test, suite_auth_test.go, TestScript_Auth_EnvVarResolution)

### 6.4 Digest

- [ ] Full MD5 digest auth flow: challenge → response with nonce, cnonce, qop (not covered)
- [ ] Nonce cached for 5 minutes (subsequent requests reuse) (not covered)
- [ ] Nonce expired — new challenge fetched (not covered)
- [ ] Wrong credentials — 401 after digest attempt (not covered)
- [ ] `$(ENV_VAR)` in username/password — resolved (not covered)

### 6.5 OAuth2 Client Credentials

- [ ] Token obtained from `token_url` using client_id + client_secret (not covered)
- [ ] Token cached and reused until expiry (not covered)
- [ ] Expired token — new token fetched automatically (not covered)
- [ ] `Authorization: Bearer <token>` header sent (not covered)
- [ ] `scopes` included in token request (not covered)
- [ ] Invalid credentials — error returned (not covered)
- [ ] `$(ENV_VAR)` in fields — resolved (not covered)

### 6.6 OAuth2 Password

- [ ] Token obtained from `token_url` using username + password + client_id (not covered)
- [ ] `client_secret` optional (public client — Keycloak support) (not covered)
- [ ] Token cached and reused until expiry (not covered)
- [ ] Expired token — new token fetched (not covered)
- [ ] Invalid credentials — error returned (not covered)
- [ ] `$(ENV_VAR)` in fields — resolved (not covered)

### 6.7 API Key

- [x] `in: header` — key placed in request header (integration-test, suite_auth_test.go, TestScript_Auth_APIKey_Header)
- [x] `in: query` — key placed in query parameter (integration-test, suite_auth_test.go, TestScript_Auth_APIKey_Query)
- [ ] Wrong key — 401 returned (not covered)
- [ ] `$(ENV_VAR)` in key/value — resolved (not covered)

### 6.8 Script

- [ ] `{workspace}/auth_scripts/{domain}.sh` executed (not covered)
- [ ] Script output JSON `{"token":"...","expires_in":N}` parsed correctly (not covered)
- [ ] Token cached and reused until expiry (not covered)
- [ ] Script returns non-zero exit — error returned (not covered)
- [ ] Script returns invalid JSON — error returned (not covered)
- [ ] Script file does not exist — error returned (not covered)
- [ ] `.bat` script on Windows (if applicable) (not covered)

### 6.9 HMAC

- [ ] `Authorization` header with HMAC-SHA256 signature sent (not covered)
- [ ] `X-MBX-APIKEY` header sent with API key (not covered)
- [ ] `timestamp` and `signature` query parameters present (not covered)
- [ ] Signature computed correctly from query string + secret key (not covered)
- [ ] `$(ENV_VAR)` in api_key/secret_key — resolved from environment (not covered)

---

## 7. Spec Parsing

### 7.1 OpenAPI 3.x

- [x] Paths, operations, parameters parsed correctly (integration-test, suite_parsing_test.go, TestScript_Parsing_OpenAPI300)
- [x] Request bodies with JSON schema parsed (integration-test, suite_parsing_test.go, TestScript_Parsing_OpenAPI300)
- [x] Response schemas with status codes parsed (integration-test, suite_parsing_test.go, TestScript_Parsing_OpenAPI300)
- [x] `$ref` references resolved (integration-test, suite_parsing_test.go, TestScript_Parsing_OpenAPI300)
- [x] Tags extracted from spec (integration-test, suite_parsing_test.go, TestScript_Parsing_OpenAPI300)
- [ ] Enums, examples, descriptions preserved (not covered)

### 7.2 Swagger 2.0

- [x] Paths, operations, parameters parsed correctly (integration-test, suite_parsing_test.go, TestScript_Parsing_Swagger20)
- [x] `definitions` resolved for `$ref` (integration-test, suite_parsing_test.go, TestScript_Parsing_Swagger20)
- [x] Tags extracted (integration-test, suite_parsing_test.go, TestScript_Parsing_Swagger20)
- [ ] All HTTP methods supported (not covered)

### 7.3 Postman Collections

- [ ] Collection items parsed as endpoints (not covered)
- [ ] Request methods, URLs, headers extracted (not covered)
- [ ] Request body (raw, JSON, form-data) parsed (not covered)
- [ ] Auth defined in collection applied (not covered)

### 7.4 Invalid / Edge Cases

- [x] Invalid YAML/JSON spec — error reported (integration-test, suite_parsing_test.go, TestScript_Parsing_InvalidSpec)
- [x] Spec with no paths — empty endpoints list (integration-test, suite_parsing_test.go, TestScript_Parsing_EmptySpec)
- [ ] Spec with no tags — endpoints grouped under default (not covered)
- [ ] Spec with circular `$ref` — handled without infinite loop (not covered)
- [ ] Remote spec URL returns 404 — error reported (not covered)
- [ ] Remote spec URL times out — error reported (not covered)

---

## 8. Mock Server (`swag2mcp-mock`)

- [ ] `swag2mcp-mock mockserver` — starts mock servers for all specs (manual — requires mock binary)
- [ ] `swag2mcp-mock mockserver --tls` — starts with TLS (self-signed) (manual — requires mock binary)
- [ ] `swag2mcp-mock mockserver --tls-cert cert.pem --tls-key key.pem` — custom TLS cert (manual — requires mock binary)
- [ ] Mock server responds to requests on configured ports (manual — requires mock binary)
- [ ] Mock OAuth2 server on port 9090 — returns valid tokens (manual — requires mock binary)
- [ ] Mock Digest server on port 9091 — handles digest auth flow (manual — requires mock binary)
- [ ] `base_mock_url` per collection — mock uses correct port (manual — requires mock binary)
- [ ] Mock server returns realistic responses based on spec (manual — requires mock binary)
- [ ] Multiple specs — each gets its own mock server (manual — requires mock binary)
- [ ] `invoke` against mock server — end-to-end flow works (manual — requires mock binary)

---

## 9. Workspace Management

- [x] `~/.swag2mcp/` created with all subdirectories (integration-test, suite_workspace_test.go, TestScript_Workspace_DirectoryStructure)
- [x] `cache/` stores downloaded remote specs (integration-test, suite_workspace_test.go, TestScript_Workspace_DirectoryStructure)
- [x] `specs/` stores local spec files (integration-test, suite_workspace_test.go, TestScript_Workspace_DirectoryStructure)
- [x] `responses/` stores large invocation responses (integration-test, suite_workspace_test.go, TestScript_Workspace_DirectoryStructure)
- [x] `auth_scripts/` stores custom auth scripts (integration-test, suite_workspace_test.go, TestScript_Workspace_DirectoryStructure)
- [x] `swag2mcp clean` — `cache/` and `responses/` emptied (integration-test, suite_workspace_test.go, TestScript_Workspace_CleanRemovesCacheAndResponses)
- [x] Old responses (>48h) cleaned on `swag2mcp mcp` startup (integration-test, suite_workspace_test.go, TestScript_Workspace_OldResponsesCleaned)
- [x] Old responses (<48h) preserved on `swag2mcp mcp` startup (integration-test, suite_workspace_test.go, TestScript_Workspace_RecentResponsesPreserved)
- [ ] Orphan auth scripts (no matching spec domain) cleaned on `update` and `clean` (not covered)

---

## 10. Error Handling

- [x] `not_found` error — JSON with code, message, hint (integration-test, suite_errors_test.go, TestScript_Errors_NotFound)
- [x] `validation_failed` error — actionable message (integration-test, suite_errors_test.go, TestScript_Errors_InvalidID)
- [x] `rate_limit` error — "try again in X seconds" message (integration-test, suite_ratelimit_test.go, TestScript_RateLimit_BlocksSecondCall)
- [x] `invoke_error` error — connection/HTTP error details (integration-test, suite_errors_test.go, TestScript_Errors_InvokeConnectionRefused)
- [x] All errors serialized as valid JSON (integration-test, suite_errors_test.go, TestScript_Errors_NotFound)
- [x] Error messages include guidance for LLM on what to do next (integration-test, suite_errors_test.go, TestScript_Errors_NotFound)

---

## 11. Cross-Cutting / Integration

- [x] Full workflow: `init` → `add spec` → `add collection` → `validate` → `mcp` → `spec_list` → `search` → `inspect` → `invoke` (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Invoke)
- [x] Config cascade: global timeout → spec timeout → collection timeout (most specific wins) (integration-test, suite_config_test.go, TestScript_ConfigCascade)
- [x] Tag filtering: `--tags public` on `mcp` — only public-tagged specs loaded (integration-test, suite_mcp_tools_test.go, TestScript_MCP_TagFilter)
- [ ] `disable: true` on spec — spec excluded from all tools (not covered)
- [ ] `disable: true` on collection — collection excluded (not covered)
- [x] Multiple specs with different auth types — each works independently (integration-test, suite_auth_test.go, TestScript_Auth_InvokeWithBearer)
- [x] Multiple collections per spec — all accessible (integration-test, suite_mcp_tools_test.go, TestScript_MCP_CollectionBySpec)
- [x] `swag2mcp update` after changing spec file — changes reflected (integration-test, suite_config_test.go, TestScript_Update_ReCachesSpecs)
- [ ] `swag2mcp update` after adding new spec file — new spec available (not covered)
- [ ] `swag2mcp update` after removing spec file — spec removed from index (not covered)
- [x] MCP server restart with `--tags` — only filtered specs available (integration-test, suite_mcp_tools_test.go, TestScript_MCP_TagFilter)
- [x] Concurrent `invoke` requests to different endpoints — both succeed (integration-test, suite_ratelimit_test.go, TestScript_RateLimit_DifferentEndpoints)
- [ ] Concurrent `invoke` requests to same endpoint — rate limit applies per endpoint (not covered)
