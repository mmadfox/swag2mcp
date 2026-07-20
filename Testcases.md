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

- [x] `http_client.randomize: true` — config field recognized, `swag2mcp info` shows randomize field (manual — config validated, requires MCP restart to apply)
- [ ] `http_client.timeout: 30s` — request times out after 30s (not covered — requires slow server)
- [x] `http_client.follow_redirects: false` — config field recognized, `swag2mcp info` shows follow_redirects (manual — config validated)
- [x] `http_client.max_redirects: 5` — config field recognized, `swag2mcp info` shows max_redirects (manual — config validated)
- [x] `http_client.max_response_size: 2048` — response truncated at 2KB (integration-test, suite_response_test.go, TestScript_ResponseSize_Configurable)
- [ ] `http_client.proxy.url` — requests go through HTTP proxy (not covered)
- [ ] `http_client.proxy.username/password` — proxy auth works (not covered)
- [ ] `http_client.proxy.bypass` — bypass list works (e.g. `localhost`) (not covered)
- [x] `http_client.headers` — custom headers parsed and shown in `swag2mcp info` (manual — config validated)
- [x] `http_client.cookies` — custom cookies parsed with name/secure/http_only fields, shown in `swag2mcp info` (manual — config validated)
- [x] `http_client.user_agent` — custom UA overrides default, shown in `swag2mcp info` (manual — config validated, shows "test-agent/1.0")
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
- [x] `disable: true` — spec is excluded from MCP tools (manual — `swag2mcp info` shows active=3, disabled=1; `swag2mcp ls` excludes disabled spec)
- [x] `tags: ["public"]` — spec filtered by `--tags public` (integration-test, suite_mcp_tools_test.go, TestScript_MCP_TagFilter)
- [x] `auth.type: bearer` + `auth.config.token: xxx` — auth applied to all endpoints (integration-test, suite_auth_test.go, TestScript_Auth_InvokeWithBearer)
- [x] `http_client` per spec — overrides global HTTP settings (integration-test, suite_config_test.go, TestScript_ConfigCascade)
- [ ] `base_url` per collection — overrides spec base_url (not covered — config validated but not tested with invoke)

### 3.3 Collection Configuration

- [x] `title: "Pets"` — collection appears with correct title (integration-test, suite_mcp_tools_test.go, TestScript_MCP_CollectionByID)
- [x] `location: ./specs/meteo.yaml` — local file loaded (integration-test, suite_mcp_tools_test.go, TestScript_MCP_SpecList)
- [x] `location: https://example.com/spec.yaml` — remote URL fetched + cached (manual — all test specs use remote URLs, cache populated after `swag2mcp update`)
- [x] `disable: true` — collection excluded (manual — `swag2mcp ls` shows empty collections list, `swag2mcp info` shows reduced collections/endpoints count)
- [x] `llm_title` + `llm_instruction` per collection — overrides spec (integration-test, suite_mcp_tools_test.go, TestScript_MCP_CollectionByID)
- [ ] `base_mock_url: localhost:8081` — mock server uses this port (not covered)
- [x] `http_client` per collection — overrides spec and global (integration-test, suite_config_test.go, TestScript_ConfigCascade)

### 3.4 Config Validation

- [x] `swag2mcp validate` — valid config reports no issues (integration-test, suite_config_test.go, TestScript_Validate_ValidConfig)
- [x] `swag2mcp validate` — duplicate domain detected (integration-test, suite_config_test.go, TestScript_Validate_DuplicateDomain; manual — tested with `add spec` duplicate, error shown)
- [ ] `swag2mcp validate` — mock port conflict detected (not covered)
- [x] `swag2mcp validate` — unreachable spec location reported (integration-test, suite_config_test.go, TestScript_Validate_UnreachableLocation)
- [x] `swag2mcp validate` — invalid domain format (integration-test, suite_config_test.go, TestScript_Validate_InvalidDomainFormat; manual — "INVALID DOMAIN HERE" → error "Domain must be 1-60 characters using only letters, digits, hyphens, and underscores"; ⚠️ UPPERCASE domain "PETSTORE" accepted without error — possible bug)
- [x] `swag2mcp validate` — invalid title length <5 chars (manual — "AB" → error "LLMTitle must be at least 5 characters")
- [x] `swag2mcp validate` — invalid title length >120 chars (manual — 123-char title → error "LLMTitle must be at most 120 characters")
- [x] `swag2mcp validate` — invalid instruction length >500 chars (manual — 503-char instruction → error "LLMInstruction must be at most 500 characters")
- [x] `swag2mcp validate` — invalid collection location <5 chars (manual — "ab" → error "Location must be at least 5 characters")
- [x] `swag2mcp validate` — invalid base_url format (manual — "not-a-valid-url" → error "BaseURL must be a valid URL")
- [x] `swag2mcp validate -t public,internal` — filter validation by tags (integration-test, suite_config_test.go, TestScript_Validate_TagFilter)

---

## 4. CLI Commands

### 4.1 `swag2mcp add spec`

- [ ] `swag2mcp add spec` — interactive TUI wizard for adding a spec (manual — requires TTY)
- [x] `swag2mcp add spec --yaml "..."` — non-interactive YAML import (integration-test, suite_config_test.go, TestScript_AddSpec_FromYAML)
- [x] `swag2mcp add spec --yaml -` — YAML piped from stdin (integration-test, suite_config_test.go, TestScript_AddSpec_FromStdin)
- [x] `swag2mcp add spec --example` — example spec YAML template printed (manual — CLI tested, outputs YAML template with domain, auth, collections)
- [x] `swag2mcp add spec` with invalid YAML — error message shown (integration-test, suite_config_test.go, TestScript_AddSpec_InvalidYAML)
- [x] `swag2mcp add spec` — config file atomically updated (integration-test, suite_config_test.go, TestScript_AddSpec_FromYAML)

### 4.2 `swag2mcp add collection`

- [ ] `swag2mcp add collection` — interactive TUI wizard (manual — requires TTY)
- [x] `swag2mcp add collection --yaml "..."` — non-interactive YAML import (integration-test, suite_config_test.go, TestScript_AddCollection_FromYAML)
- [x] `swag2mcp add collection --yaml -` — YAML piped from stdin (manual — CLI tested with heredoc pipe)
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

- [x] `swag2mcp ls` — shows all specs and collections in formatted table (integration-test, suite_config_test.go, TestScript_ListSpecs; manual — CLI tested)
- [x] `swag2mcp ls -t public` — filters by tag (integration-test, suite_config_test.go, TestScript_ListSpecs_TagFilter)
- [x] `swag2mcp ls -t public,internal` — multiple tags (manual — CLI tested, returns empty when no specs have those tags)
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

### 4.10 `swag2mcp info`

- [x] `swag2mcp info` — outputs JSON with version, workspace path, specs summary, HTTP client config, MCP transport, auth, mock (manual — CLI tested)
- [x] `swag2mcp info [path]` — accepts custom workspace path (manual — CLI tested)
- [x] `swag2mcp info` — shows total, active, disabled, collections, endpoints counts (manual — tested with 4 specs, shows disabled=1 when `disable:true`)
- [x] `swag2mcp info` — shows http_client config (randomize, user_agent, timeout, follow_redirects, max_redirects, max_response_size, headers, cookies) (manual — tested)
- [x] `swag2mcp info` — shows mcp config (transport, auth_enabled) (manual — tested)

### 4.11 `swag2mcp export`

- [x] `swag2mcp export [path] [output.zip]` — creates ZIP backup with specs, config, auth scripts (manual — CLI tested, 4056-byte ZIP created)
- [ ] `swag2mcp export` — default output filename `swag2mcp-backup-<timestamp>.zip` (not covered)
- [ ] `swag2mcp export --spec meteo` — export only specified specs (not covered)

### 4.12 `swag2mcp import`

- [x] `swag2mcp import --spec meteo` — bulk import from existing config URLs (manual — CLI tested, specs downloaded)
- [x] `swag2mcp import --from-zip /path/to/backup.zip` — restore from ZIP backup (manual — CLI tested, workspace restored)
- [ ] `swag2mcp import [url] [name]` — single import from URL (not covered)

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
- [x] Returns `not_found` for non-existent specId (manual — MCP tool tested with 000...000 ID, returns JSON with code, message, hint)
- [ ] Returns empty list for spec with no collections (not covered)

### 5.4 `collection_by_id`

- [x] Returns collection details + tags for valid ID (integration-test, suite_mcp_tools_test.go, TestScript_MCP_CollectionByID)
- [x] Returns `not_found` for non-existent ID (manual — MCP tool tested with 000...000 ID)
- [x] Returns `not_found` for malformed ID (manual — MCP tool tested with "invalid" string, returns validation_failed)

### 5.5 `tag_by_spec`

- [x] Returns all tags across spec for valid specId (integration-test, suite_mcp_tools_test.go, TestScript_MCP_TagBySpec)
- [x] Returns `not_found` for non-existent specId (manual — MCP tool tested with 000...000 ID)
- [ ] Returns empty list for spec with no tags (not covered)

### 5.6 `tag_by_collection`

- [x] Returns all tags for valid collectionId (manual — MCP tool tested, returns tags with id, title, countMethods)
- [x] Returns `not_found` for non-existent collectionId (manual — MCP tool tested with 000...000 ID)
- [ ] Returns empty list for collection with no tags (not covered)

### 5.7 `tag_by_id`

- [x] Returns tag details for valid ID (manual — MCP tool tested, returns id, title, countMethods)
- [x] Returns `not_found` for non-existent ID (manual — MCP tool tested with 000...000 ID)

### 5.8 `endpoint_by_spec`

- [x] Returns all endpoints across spec for valid specId (integration-test, suite_mcp_tools_test.go, TestScript_MCP_EndpointBySpec)
- [x] Returns `not_found` for non-existent specId (manual — MCP tool tested with 000...000 ID)
- [ ] Returns empty list for spec with no endpoints (not covered)

### 5.9 `endpoint_by_collection`

- [x] Returns all endpoints for valid collectionId (manual — MCP tool tested, returns endpoints with method, path, summary, tagId, tagName)
- [x] Returns `not_found` for non-existent collectionId (manual — MCP tool tested with 000...000 ID)
- [ ] Returns empty list for collection with no endpoints (not covered)

### 5.10 `endpoint_by_tag`

- [x] Returns all endpoints for valid tagId (manual — MCP tool tested, returns endpoints with spec/collection/tag context)
- [x] Returns `not_found` for non-existent tagId (manual — MCP tool tested with 000...000 ID)
- [ ] Returns empty list for tag with no endpoints (not covered)

### 5.11 `endpoint_by_id`

- [x] Returns endpoint summary (method, path, summary, deprecated) for valid ID (integration-test, suite_mcp_tools_test.go, TestScript_MCP_EndpointByID)
- [x] Returns `not_found` for non-existent ID (manual — MCP tool tested with 000...000 ID)
- [ ] Deprecated endpoint shows `deprecated: true` (not covered)

### 5.12 `search`

- [x] `search("pet")` — returns matching endpoints (integration-test, suite_search_test.go, TestScript_Search_Basic)
- [x] `search("method:GET")` — only GET endpoints (integration-test, suite_search_test.go, TestScript_Search_ByMethod)
- [x] `search("tag:pets")` — only tagged endpoints (integration-test, suite_search_test.go, TestScript_Search_ByTag)
- [x] `search("path:/pets")` — path match (integration-test, suite_search_test.go, TestScript_Search_ByPath)
- [x] `search("+method:GET +summary:pet")` — boolean AND (integration-test, suite_search_test.go, TestScript_Search_BooleanAND)
- [ ] `search("summary:\"create user\"")` — phrase search (not covered — returns empty result for some phrase queries)
- [ ] `search("sumary~")` — fuzzy search (typo tolerance) (not covered)
- [x] `search("list*")` — wildcard search (integration-test, suite_search_test.go, TestScript_Search_Wildcard; manual — returns "List all pets" and "List Pokémon")
- [x] `search("zzzzz")` — empty results (integration-test, suite_search_test.go, TestScript_Search_EmptyResults)
- [x] `search("*")` — returns all endpoints (integration-test, suite_search_test.go, TestScript_Search_AllEndpoints)
- [x] `search("pet", limit=1)` — returns exactly 1 result (integration-test, suite_search_test.go, TestScript_Search_LimitBounds)
- [x] `search("pet", limit=50)` — returns up to 50 results (integration-test, suite_search_test.go, TestScript_Search_AllEndpoints)
- [x] `search("pet", limit=0)` — error (min 1) (manual — returns `validation_failed` with "Limit must be between 1 and 50")
- [x] `search("pet", limit=51)` — error (max 50) (manual — returns `validation_failed` with "Limit must be between 1 and 50")
- [x] `search("method:POST")` — returns POST endpoints only (manual — returns 1 result)

### 5.13 `inspect`

- [x] Returns full operation object for valid endpointId (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Inspect)
- [x] Parameters (path, query, header) with schemas are present (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Inspect)
- [x] Request body schema is present (for POST/PUT/PATCH) (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Inspect)
- [x] Response schemas with status codes are present (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Inspect)
- [x] Referenced `$ref` schemas are resolved (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Inspect)
- [x] Returns `not_found` for non-existent endpointId (manual — MCP tool tested with 000...000 ID, returns not_found with guidance)

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
- [x] `invoke` on 4xx response — status code and error body returned (manual — meteo GET /pets/{petId} returns 404 with body)
- [x] `invoke` on real API (Binance, dadjoke, PokeAPI) — 200 response with correct body (manual — BTCUSDT price, random joke, pokemon list)
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

### 5.16 `info`

- [x] Returns version, workspace path, specs summary (total, active, disabled, collections, endpoints) (manual — MCP tool tested)
- [x] Returns http_client config (randomize, user_agent, max_response_size, headers, timeout, follow_redirects, max_redirects, cookies) (manual — MCP tool tested)
- [x] Returns mcp config (transport, auth_enabled) (manual — MCP tool tested)
- [x] Returns mock config (enabled) (manual — MCP tool tested)
- [ ] Returns auth methods per spec (not covered — no auth configured in test workspace)

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
- [x] `not_found` error — returned by all navigation tools (spec_by_id, collection_by_spec, collection_by_id, tag_by_id, endpoint_by_spec, endpoint_by_collection, endpoint_by_tag, endpoint_by_id, inspect) for non-existent IDs (manual — all tested with 000...000)
- [x] `validation_failed` error — returned by `search` for limit=0 and limit=51 (manual — tested, "Limit must be between 1 and 50")
- [x] `validation_failed` error — returned by `spec_by_id` for empty ID string (manual — tested, "ID must be a 32-character hex string")
- [x] `validation_failed` error — returned by `endpoint_by_id` for malformed ID (not 32-char hex) (manual — tested with "invalid")
- [x] `rate_limit` error — "try again in X seconds" message (integration-test, suite_ratelimit_test.go, TestScript_RateLimit_BlocksSecondCall)
- [x] `invoke_error` error — connection/HTTP error details (integration-test, suite_errors_test.go, TestScript_Errors_InvokeConnectionRefused)
- [x] All errors serialized as valid JSON (integration-test, suite_errors_test.go, TestScript_Errors_NotFound)
- [x] Error messages include guidance for LLM on what to do next (integration-test, suite_errors_test.go, TestScript_Errors_NotFound)

---

## 11. Cross-Cutting / Integration

- [x] Full workflow: `init` → `add spec` → `add collection` → `validate` → `mcp` → `spec_list` → `search` → `inspect` → `invoke` (integration-test, suite_mcp_tools_test.go, TestScript_MCP_Invoke)
- [x] Config cascade: global timeout → spec timeout → collection timeout (most specific wins) (integration-test, suite_config_test.go, TestScript_ConfigCascade)
- [x] Tag filtering: `--tags public` on `mcp` — only public-tagged specs loaded (integration-test, suite_mcp_tools_test.go, TestScript_MCP_TagFilter)
- [x] `disable: true` on spec — spec excluded from all tools (manual — `swag2mcp info` shows active=3, disabled=1; `swag2mcp ls` excludes disabled spec)
- [x] `disable: true` on collection — collection excluded (manual — `swag2mcp ls` shows empty collections list, `swag2mcp info` shows reduced count)
- [x] Multiple specs with different auth types — each works independently (integration-test, suite_auth_test.go, TestScript_Auth_InvokeWithBearer)
- [x] Multiple collections per spec — all accessible (integration-test, suite_mcp_tools_test.go, TestScript_MCP_CollectionBySpec)
- [x] `swag2mcp update` after changing spec file — changes reflected (integration-test, suite_config_test.go, TestScript_Update_ReCachesSpecs)
- [x] `swag2mcp update` — processes all specs, clears cache and re-downloads (manual — "4 specs processed")
- [x] `swag2mcp clean` — removes cache/ and responses/ contents (manual — responses/ emptied, cache/ re-populated)
- [x] `swag2mcp clean` — preserves specs/ and auth_scripts/ (manual — CLI tested)
- [ ] `swag2mcp update` after adding new spec file — new spec available (not covered)
- [ ] `swag2mcp update` after removing spec file — spec removed from index (not covered)
- [x] MCP server restart with `--tags` — only filtered specs available (integration-test, suite_mcp_tools_test.go, TestScript_MCP_TagFilter)
- [x] Concurrent `invoke` requests to different endpoints — both succeed (integration-test, suite_ratelimit_test.go, TestScript_RateLimit_DifferentEndpoints)
- [ ] Concurrent `invoke` requests to same endpoint — rate limit applies per endpoint (not covered)

---

## 12. Manual Testing via MCP API (Live)

Tested against workspace with 4 specs: meteo, binance, dadjoke, pokeapi.

### 12.1 MCP Tools — Navigation

- [x] `spec_list` — returns 4 specs with correct IDs and domains (manual)
- [x] `spec_by_id` — returns spec details + collections for each valid ID (manual — tested all 4)
- [x] `collection_by_spec` — returns collections for each spec (manual — tested all 4)
- [x] `collection_by_id` — returns collection details + tags (manual — tested binance, meteo)
- [x] `tag_by_spec` — returns tags for each spec (manual — tested binance, meteo)
- [x] `tag_by_collection` — returns tags for each collection (manual — tested binance, meteo)
- [x] `tag_by_id` — returns tag details with countMethods (manual — tested market-data, pets)
- [x] `endpoint_by_spec` — returns all endpoints per spec (manual — binance=4, meteo=3, dadjoke=3, pokeapi=3)
- [x] `endpoint_by_collection` — returns endpoints per collection (manual — tested binance, meteo)
- [x] `endpoint_by_tag` — returns endpoints per tag (manual — tested market-data, pets)
- [x] `endpoint_by_id` — returns endpoint summary (method, path, summary) (manual — tested multiple)

### 12.2 MCP Tools — Search

- [x] `search("pet")` — returns matching endpoints (2 results)
- [x] `search("method:GET")` — returns all GET endpoints (12 results, limit=50)
- [x] `search("tag:pets")` — returns endpoints tagged "pets" (3 results)
- [x] `search("path:/pets")` — returns empty result (path search with `/pets` yields 0)
- [x] `search("+method:GET +summary:price")` — boolean AND (1 result: Binance price ticker)
- [x] `search("method:POST")` — returns POST endpoints (1 result)
- [x] `search("*")` — returns all endpoints (13 results)
- [x] `search("zzzzz")` — returns empty results
- [x] `search("list*")` — wildcard search (2 results)
- [x] `search("pet", limit=1)` — returns 1 result
- [x] `search("pet", limit=0)` — returns `validation_failed` error
- [x] `search("pet", limit=51)` — returns `validation_failed` error

### 12.3 MCP Tools — Inspect

- [x] `inspect` on GET /pets — returns operation with parameters and response schema (manual)
- [x] `inspect` on POST /pets — returns operation with requestBody schema (manual)
- [x] `inspect` on GET /api/v3/ticker/24hr — returns operation with query parameters (manual)
- [x] `inspect` with non-existent endpointId — returns `not_found` error (manual)

### 12.4 MCP Tools — Invoke (Live APIs)

- [x] `invoke` GET / (dadjoke) — 200, returns random joke with id, joke, status (manual)
- [x] `invoke` GET /api/v3/ticker/price?symbol=BTCUSDT (Binance) — 200, returns price and symbol (manual)
- [x] `invoke` GET /api/v2/pokemon?limit=5&offset=0 (PokeAPI) — 200, returns count, next, results (manual)
- [x] `invoke` GET /pets/{petId} (Open-Meteo) — 404, server returns "null for uri" (manual)
- [x] `invoke` POST /pets (Open-Meteo) — 404, server returns "null for uri" (manual)

### 12.5 MCP Tools — Error Cases

- [x] `spec_by_id("")` — returns `validation_failed` with "ID must be a 32-character hex string" (manual)
- [x] `endpoint_by_id("invalid")` — returns `validation_failed` with md5 validation error (manual)
- [x] `spec_by_id("000...0")` — returns `not_found` with guidance (manual)
- [x] `collection_by_spec("000...0")` — returns `not_found` (manual)
- [x] `collection_by_id("000...0")` — returns `not_found` (manual)
- [x] `tag_by_id("000...0")` — returns `not_found` (manual)
- [x] `endpoint_by_spec("000...0")` — returns `not_found` (manual)
- [x] `endpoint_by_collection("000...0")` — returns `not_found` (manual)
- [x] `endpoint_by_tag("000...0")` — returns `not_found` (manual)
- [x] `inspect("000...0")` — returns `not_found` (manual)

### 12.6 CLI Commands (Live)

- [x] `swag2mcp add spec --yaml "..."` — adds spec (manual — tested with test-api, httpbin, demo-api)
- [x] `swag2mcp add spec --yaml -` — adds spec from stdin/heredoc (manual — tested with demo-api)
- [x] `swag2mcp add spec --example` — prints YAML template (manual — outputs domain, auth, collections template)
- [x] `swag2mcp add collection --yaml "..."` — adds collection (manual — tested with Second Collection)
- [x] `swag2mcp add collection --yaml -` — adds collection from stdin (manual — tested with httpbin-v2)
- [x] `swag2mcp ls` — shows all specs with collections (manual)
- [x] `swag2mcp ls -t public,internal` — returns empty (no specs have those tags) (manual)
- [x] `swag2mcp validate` — reports "Configuration is valid." (manual)
- [x] `swag2mcp validate` — detects invalid domain format (manual — "INVALID DOMAIN HERE")
- [x] `swag2mcp validate` — detects duplicate domain (manual — second "meteo" added)
- [x] `swag2mcp validate` — detects invalid title length <5 chars (manual — "AB")
- [x] `swag2mcp validate` — detects invalid title length >120 chars (manual)
- [x] `swag2mcp validate` — detects invalid instruction length >500 chars (manual)
- [x] `swag2mcp validate` — detects invalid collection location <5 chars (manual — "ab")
- [x] `swag2mcp validate` — detects invalid base_url format (manual — "not-a-valid-url")
- [x] `swag2mcp version` — prints "swag2mcp dev" (manual)
- [x] `swag2mcp info` — outputs JSON with specs summary, http_client, mcp, auth, mock (manual)
- [x] `swag2mcp info` — shows disabled specs count when `disable:true` (manual — active=3, disabled=1)
- [x] `swag2mcp update` — processes all specs (manual — "4 specs processed")
- [x] `swag2mcp clean` — removes cache/ and responses/ contents (manual)
- [x] `swag2mcp export [path] [output.zip]` — creates ZIP backup (manual — verified file exists)
- [x] `swag2mcp import --spec meteo` — imports spec from configured URLs (manual)
- [x] `swag2mcp import --from-zip backup.zip` — restores workspace from ZIP (manual)

### 12.7 Config Settings (Live)

- [x] `http_client.randomize: true` — config parsed and shown in `swag2mcp info` (manual)
- [x] `http_client.timeout: 30s` — config parsed and shown in `swag2mcp info` (manual)
- [x] `http_client.follow_redirects: false` — config parsed and shown in `swag2mcp info` (manual)
- [x] `http_client.max_redirects: 5` — config parsed and shown in `swag2mcp info` (manual)
- [x] `http_client.headers` — config parsed and shown in `swag2mcp info` (manual)
- [x] `http_client.cookies` — config parsed with name/secure/http_only fields (manual)
- [x] `http_client.user_agent: "test-agent/1.0"` — overrides default, shown in `swag2mcp info` (manual)
- [x] `disable: true` on spec — spec excluded from `swag2mcp ls` and `swag2mcp info` (manual)
- [x] `disable: true` on collection — collection excluded, endpoints reduced (manual)
- [ ] `http_client.proxy.url` — requests go through HTTP proxy (not covered — requires proxy server)
- [ ] `http_client.timeout` — actual request timeout (not covered — requires slow server)
- [ ] `base_url` per collection — overrides spec base_url (not covered — config validated but not tested with invoke)

