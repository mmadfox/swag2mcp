---
name: swag2mcp-format
description: |
  [WHAT] Response formatting rules for swag2mcp MCP tools.
  [WHEN] Use when displaying results from ANY swag2mcp MCP tool:
         spec_list, spec_by_id, collection_by_*, tag_by_*,
         endpoint_by_*, search, inspect, invoke, auth, info.
         Automatically triggered on every swag2mcp tool response.
  [WHY] Ensures consistent human-readable markdown, enforces pagination,
        colorizes HTTP methods, structures schemas/headers/errors.
license: MIT
metadata:
  author: mmadfox
  version: "1.0.0"
---

# swag2mcp-format — Response Formatting Skill

This document defines how to format **every swag2mcp MCP tool response** into human-readable markdown. Each response type has a dedicated format optimized for clarity and usability.

## When to Use

TRIGGERS (use when user says or tool returns):
- "show me specs" / "list APIs" → spec_list
- "show spec details" → spec_by_id
- "list collections" → collection_by_spec / collection_by_id
- "show tags" → tag_by_*
- "list endpoints" / "show endpoints" → endpoint_by_*
- "search for ..." / "find endpoint" → search
- "inspect endpoint" / "show me the full spec" → inspect
- "call API" / "invoke" / "make request" → invoke
- "get auth token" → auth
- "show system info" / "status" → info
- ANY response from swag2mcp MCP tools

## When NOT to Use

- User asks to display their OWN data structure (not from swag2mcp tools)
- User explicitly requests raw/unformatted JSON output
- User asks for a custom format different from these rules

---

## 1. General Rules

### 1.1 HTTP Method Colors

Always colorize HTTP methods in tables and headers:

| Method | Markdown |
|--------|----------|
| GET    | `🟢 GET` |
| POST   | `🔵 POST` |
| PUT    | `🟠 PUT` |
| PATCH  | `🟣 PATCH` |
| DELETE | `🔴 DELETE` |
| HEAD   | `⚪ HEAD` |
| OPTIONS| `⚪ OPTIONS` |

### 1.2 Sorting

- Endpoints within a list: sort by **HTTP method** (GET < POST < PUT < PATCH < DELETE < HEAD < OPTIONS), then by **path** alphabetically.
- Collections, tags, specs: sort by **title** alphabetically (case-insensitive).
- Search results: keep the order returned by the search engine (relevance-based).

### 1.3 ID Display

- Always show IDs in the **last column** of tables.
- Truncate to first 8 characters for readability: `a1b2c3d4...`.
- Full ID is available via `endpoint_by_id` or `inspect` if needed.

### 1.4 Pagination — Never Dump All Results

**Never display more than 10 items at once.** Always paginate:

- Show the first 5-10 results, then ask: `_Showing 5 of 14 results. Load 5 more?_`
- Wait for user confirmation before showing the next batch.
- This applies to ALL list responses: `search`, `endpoint_by_*`, `tag_by_*`, `collection_by_*`, `spec_list`.
- Exception: `tag_by_id`, `endpoint_by_id`, `spec_by_id` (single-item responses) — show fully.

**Example interaction:**

```
**Search results for "pet" (14):**

| Method | Path | Summary | Tag | ID |
|--------|------|---------|-----|----|
| 🟢 GET | /pet/{petId} | Find pet by ID | pets | e1f2... |
| 🟢 GET | /pet/findByStatus | Find pets by status | pets | i5j6... |
| 🟢 GET | /pet/findByTags | Find pets by tags | pets | m9n0... |
| 🔵 POST | /pet | Add a new pet | pets | q3r4... |
| 🔴 DELETE | /pet/{petId} | Delete a pet | pets | u7v8... |

_Showing 5 of 14 results. Load 5 more?_
```

When user confirms, show the next batch with the same prompt at the bottom.

### 1.5 Empty States

- If a list is empty: `_No items found._`
- If a field is empty/omitted: `—` (em dash).

### 1.5 Deprecated Endpoints

Append `⚠️ Deprecated` after the summary for deprecated endpoints.

---

## 2. List Responses (Markdown Tables)

These responses contain flat lists of items. Always use a markdown table.

### 2.1 `spec_list` → `SpecsResponse`

```
**Available Specifications (2):**

| ID | Domain |
|----|--------|
| a1b2c3d4... | Petstore API |
| e5f6g7h8... | Weather API |
```

**Format:**
- Header: `**Available Specifications ({count}):**`
- Table columns: `ID`, `Domain`
- Sort by domain alphabetically.

### 2.2 `collection_by_spec` → `CollectionsResponse`

```
**Spec:** Petstore API (`a1b2c3d4...`)

**Collections (3):**

| ID | Title | Tags | Methods |
|----|-------|------|---------|
| c1d2e3f4... | Pet Operations | 2 | 8 |
| g5h6i7j8... | Store | 1 | 4 |
| k9l0m1n2... | User Management | 3 | 6 |
```

**Format:**
- Context line: `**Spec:** {Domain} ({truncated ID})`
- Header: `**Collections ({count}):**`
- Table columns: `ID`, `Title`, `Tags`, `Methods`
- If `LLMTitle` is present, show it in parentheses after Title: `Pet Operations (pets)`.

### 2.3 `collection_by_id` → `CollectionByIDResponse`

```
**Spec:** Petstore API (`a1b2c3d4...`)

**Collection:** Pet Operations (`c1d2e3f4...`) — 8 methods

**Tags (2):**

| ID | Title | Methods |
|----|-------|---------|
| t1u2v3w4... | pets | 5 |
| x5y6z7a8... | store | 3 |
```

**Format:**
- Context line: `**Spec:** {Domain} ({truncated ID})`
- Collection line: `**Collection:** {Title} ({truncated ID}) — {countMethods} methods`
- Header: `**Tags ({count}):**`
- Table columns: `ID`, `Title`, `Methods`

### 2.4 `tag_by_collection` → `TagsByCollectionResponse`

```
**Spec:** Petstore API (`a1b2c3d4...`)

**Collection:** Pet Operations (`c1d2e3f4...`) — 8 methods

**Tags (2):**

| ID | Title | Methods |
|----|-------|---------|
| t1u2v3w4... | pets | 5 |
| x5y6z7a8... | store | 3 |
```

**Format:** Same as `collection_by_id` — identical structure.

### 2.5 `tag_by_spec` → `TagsBySpecResponse`

```
**Tags across Petstore API (5):**

| ID | Title | Methods |
|----|-------|---------|
| t1u2v3w4... | pets | 5 |
| x5y6z7a8... | store | 3 |
| b9c0d1e2... | user | 6 |
```

**Format:**
- Header: `**Tags across {Domain} ({count}):**`
- Table columns: `ID`, `Title`, `Methods`

### 2.6 `tag_by_id` → `TagByIDResponse`

```
**Tag:** pets (`t1u2v3w4...`) — 5 methods
```

**Format:**
- Single line: `**Tag:** {Title} ({truncated ID}) — {countMethods} methods`

### 2.7 `endpoint_by_tag` → `EndpointsByTagResponse`

```
**Spec:** Petstore API (`a1b2c3d4...`)
**Collection:** Pet Operations (`c1d2e3f4...`)
**Tag:** pets (`t1u2v3w4...`) — 5 methods

**Endpoints (5):**

| Method | Path | Summary | ID |
|--------|------|---------|----|
| 🟢 GET | /pet/{petId} | Find pet by ID | e1f2g3h4... |
| 🟢 GET | /pet/findByStatus | Find pets by status | i5j6k7l8... |
| 🔵 POST | /pet | Add a new pet | m9n0o1p2... |
| 🟠 PUT | /pet | Update an existing pet | q3r4s5t6... |
| 🔴 DELETE | /pet/{petId} | Delete a pet | u7v8w9x0... |
```

**Format:**
- Context lines (3 lines): Spec, Collection, Tag
- Header: `**Endpoints ({count}):**`
- Table columns: `Method` (colored), `Path`, `Summary`, `ID`
- Sort by method order, then path.

### 2.8 `endpoint_by_collection` → `EndpointsByCollectionResponse`

```
**Spec:** Petstore API (`a1b2c3d4...`)
**Collection:** Pet Operations (`c1d2e3f4...`) — 8 methods

**Endpoints (8):**

| Method | Path | Summary | Tag | ID |
|--------|------|---------|-----|----|
| 🟢 GET | /pet/{petId} | Find pet by ID | pets | e1f2g3h4... |
| 🟢 GET | /pet/findByStatus | Find pets by status | pets | i5j6k7l8... |
| 🔵 POST | /pet | Add a new pet | pets | m9n0o1p2... |
| 🟠 PUT | /pet | Update an existing pet | pets | q3r4s5t6... |
| 🔴 DELETE | /pet/{petId} | Delete a pet | pets | u7v8w9x0... |
```

**Format:**
- Context lines (2 lines): Spec, Collection
- Header: `**Endpoints ({count}):**`
- Table columns: `Method` (colored), `Path`, `Summary`, `Tag`, `ID`
- Group by tag visually (sort by tag, then method, then path).

### 2.9 `endpoint_by_spec` → `EndpointsBySpecResponse`

```
**All endpoints in Petstore API (14):**

| Method | Path | Summary | Tag | Collection | Spec | ID |
|--------|------|---------|-----|-----------|------|----|
| 🟢 GET | /pet/{petId} | Find pet by ID | pets | Pet Operations | Petstore | e1f2... |
| 🔵 POST | /pet | Add a new pet | pets | Pet Operations | Petstore | m9n0... |
| 🔵 POST | /store/order | Place order | store | Store | Petstore | a1b2... |
```

**Format:**
- Header: `**All endpoints in {Domain} ({count}):**`
- Table columns: `Method` (colored), `Path`, `Summary`, `Tag`, `Collection`, `Spec`, `ID`
- ID column: truncate to first 4 chars for space.
- Sort by spec, then collection, then tag, then method, then path.

### 2.10 `search` → `SearchResponse`

```
**Search results for "pet" (5):**

| Method | Path | Summary | Tag | Collection | Spec | ID |
|--------|------|---------|-----|-----------|------|----|
| 🟢 GET | /pet/{petId} | Find pet by ID | pets | Pet Operations | Petstore | e1f2... |
| 🔵 POST | /pet | Add a new pet | pets | Pet Operations | Petstore | m9n0... |
| 🟠 PUT | /pet | Update an existing pet | pets | Pet Operations | Petstore | q3r4... |
```

**Format:**
- Header: `**Search results for "{query}" ({count}):**`
- Table columns: `Method` (colored), `Path`, `Summary`, `Tag`, `Collection`, `Spec`, `ID`
- ID column: truncate to first 4 chars.
- Keep search engine's relevance order.

---

## 3. Detail Responses (Sections with Headers)

### 3.1 `spec_by_id` → `SpecByIDResponse`

```
## Spec: Petstore API

| Property | Value |
|----------|-------|
| ID | a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6 |
| Domain | Petstore API |

## Collections (3)

| ID | Title | Tags | Methods |
|----|-------|------|---------|
| c1d2e3f4... | Pet Operations | 2 | 8 |
| g5h6i7j8... | Store | 1 | 4 |
| k9l0m1n2... | User Management | 3 | 6 |
```

**Format:**
- `## Spec: {Domain}` — level-2 heading
- Key-value table for spec properties (ID, Domain)
- `## Collections ({count})` — level-2 heading
- Table columns: `ID`, `Title`, `Tags`, `Methods`

### 3.2 `endpoint_by_id` → `EndpointByIDResponse`

```
## Spec: Petstore API

| Property | Value |
|----------|-------|
| ID | a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6 |
| Domain | Petstore API |

## Collection: Pet Operations

| Property | Value |
|----------|-------|
| ID | c1d2e3f4... |
| Methods | 8 |

## Tag: pets

| Property | Value |
|----------|-------|
| ID | t1u2v3w4... |
| Methods | 5 |

## Endpoint

| Method | Path | Summary | ID |
|--------|------|---------|----|
| 🟢 GET | /pet/{petId} | Find pet by ID | e1f2g3h4... |
```

**Format:**
- Four sections: `## Spec`, `## Collection`, `## Tag`, `## Endpoint`
- Each section has a key-value table (2 columns) or a single-row table for the endpoint itself.
- Full ID (not truncated) for the endpoint.

---

## 4. Inspect Response (Sections with Code Blocks)

### 4.1 `inspect` → `InspectResponse`

```
## 🟢 GET /pet/{petId}

**Spec:** Petstore API
**Full URL:** `https://petstore.swagger.io/v2/pet/{petId}`
**Summary:** Find pet by ID
**Description:** Returns a single pet by its ID

### Parameters

| Name | In | Type | Required | Description |
|------|----|------|----------|-------------|
| petId | path | integer | ✅ Yes | ID of pet to return |
| api_key | header | string | ❌ No | API key for authentication |

### Request Body

_No request body required._

### Responses

**200 — successful operation**

```json
{
  "type": "object",
  "properties": {
    "id": { "type": "integer", "format": "int64" },
    "name": { "type": "string" },
    "category": {
      "type": "object",
      "properties": {
        "id": { "type": "integer" },
        "name": { "type": "string" }
      }
    },
    "status": {
      "type": "string",
      "enum": ["available", "pending", "sold"]
    }
  },
  "required": ["name", "status"]
}
```

**400 — Invalid pet ID**

```json
{
  "type": "object",
  "properties": {
    "code": { "type": "integer" },
    "message": { "type": "string" }
  }
}
```

**404 — Pet not found**

_No response body defined._
```

**Format:**
- Title: `## {method emoji} {Method} {Path}`
- Overview block (key-value lines):
  - `**Spec:** {Domain}`
  - `**Full URL:** \`{fullUrl}\``
  - `**Summary:** {summary}` (if present)
  - `**Description:** {description}` (if present)
  - If deprecated: `⚠️ **Deprecated**`
- `### Parameters` section:
  - If empty: `_No parameters._`
  - Table columns: `Name`, `In`, `Type`, `Required`, `Description`
  - `Required`: `✅ Yes` or `❌ No`
  - `Type`: from schema type (e.g., `string`, `integer`, `array<Pet>`)
- `### Request Body` section:
  - If empty: `_No request body required._`
  - Show description, required status
  - Schema as ` ```json ` block (pretty-printed, simplified — remove `$ref`, resolve inline)
- `### Responses` section:
  - For each status code: `**{code} — {description}**`
  - Schema as ` ```json ` block
  - If no content: `_No response body defined._`
  - If multiple content types, show the first one (prefer `application/json`)

### Schema Simplification Rules for Code Blocks

When displaying schemas in ` ```json ` blocks:
1. Resolve `$ref` to inline definitions (show the actual structure, not references).
2. Remove `$ref`, `oneOf`, `anyOf`, `allOf` — show only the resolved structure.
3. Keep `type`, `properties`, `items`, `required`, `enum`, `format`, `description`, `example`, `default`.
4. For `array` types, show `"type": "array"` with `items` containing the element schema.
5. If a property has `enum`, show the enum values inline.
6. If a property has `example`, append it as a comment: `// example: "doggie"`.

---

## 5. Invoke Response (Sections)

### 5.1 `invoke` → `InvokeResponse` (normal)

```
**Status:** 200 OK

### Headers

| Header | Value |
|--------|-------|
| content-type | application/json |
| x-request-id | abc-123-def |

### Body

```json
{
  "id": 1,
  "name": "Buddy",
  "status": "available"
}
```
```

**Format:**
- Status line: `**Status:** {code} {text}` (e.g., `200 OK`, `404 Not Found`)
  - Colorize: 2xx 🟢, 3xx 🟡, 4xx 🟠, 5xx 🔴
- `### Headers` section:
  - Table columns: `Header`, `Value`
  - Sort alphabetically by header name
  - Omit `Content-Length`, `Date`, `Connection`, `Keep-Alive` (noise reduction)
- `### Body` section:
  - If JSON: ` ```json ` block with pretty-printed JSON
  - If string: ` ``` ` block with raw text
  - If empty: `_Empty response body._`

### 5.2 `invoke` → `InvokeResponse` (large response with FileReference)

```
**Status:** 200 OK

### Headers

| Header | Value |
|--------|-------|
| content-type | application/json |

### Body

📁 **Response body (2.5 KB) exceeds the maximum size limit (1 KB).**
The full response has been saved to disk.

| Property | Value |
|----------|-------|
| File | `~/.swag2mcp/responses/petstore-get-pet-findByStatus-abc123.json` |
| Size | 2.5 KB |
| Max Size | 1 KB |
| Open | `open ~/.swag2mcp/responses/petstore-get-pet-findByStatus-abc123.json` |
```

**Format:**
- Status line: same as normal
- Headers table: same as normal
- Body section:
  - 📁 icon + message from `FileRef.Message`
  - Key-value table with: `File`, `Size`, `Max Size`, `Open`

---

## 6. Info Response (Two-Column Layout)

### 6.1 `info` → `InfoResponse`

```
## System

| Property | Value |
|----------|-------|
| Version | 1.2.3 |
| Workspace | ~/.swag2mcp |
| Uptime | 2h15m30s |

## Specs

| Property | Count |
|----------|-------|
| Total | 3 |
| Active | 2 |
| Disabled | 1 |
| Collections | 7 |
| Endpoints | 42 |

## HTTP Client

| Property | Value |
|----------|-------|
| Randomize | ❌ No |
| User Agent | swag2mcp/1.0 |
| Timeout | 30s |
| Follow Redirects | ✅ Yes |
| Max Redirects | 5 |
| Max Response Size | 1 KB |

| Proxy | |
|-------|---|
| URL | http://proxy:8080 |
| Username | user |
| Bypass | localhost, 127.0.0.1 |

| Headers | |
|---------|---|
| X-Custom | value1 |

| Cookies | |
|---------|---|
| session_id | Domain: .example.com, Path: /, Secure: ✅, HTTPOnly: ✅ |

## MCP

| Property | Value |
|----------|-------|
| Transport | stdio |
| Auth Enabled | ✅ Yes |

## Auth

| Property | Value |
|----------|-------|
| Methods | bearer, oauth2-cc |

## Mock

| Property | Value |
|----------|-------|
| Enabled | ❌ No |
```

**Format:**
- Each top-level field becomes a `## {Section}` heading.
- Simple key-value pairs: 2-column table (`Property`, `Value`).
- Booleans: `✅ Yes` / `❌ No`.
- Nested objects (Proxy, Headers, Cookies): separate 2-column tables with the nested name as a header row.
- `MaxResponseSize`: convert bytes to human-readable (e.g., `1048` → `1 KB`).
- `Uptime`: already human-readable, show as-is.
- Omit empty sections entirely (e.g., if no auth methods, skip `## Auth`).

---

## 7. Auth Response (Key-Value)

### 7.1 `auth` → `AuthResponse`

```
**Token:** Bearer eyJhbGciOiJIUzI1NiIs...

**Headers:**

| Header | Value |
|--------|-------|
| Authorization | Bearer eyJhbGciOiJIUzI1NiIs... |
| X-API-Key | abc123 |

**Query Parameters:**

| Parameter | Value |
|-----------|-------|
| api_key | abc123 |
```

**Format:**
- Token line: `**Token:** {value}` (truncate long tokens to 40 chars + `...`)
- `**Headers:**` section with table (if present)
- `**Query Parameters:**` section with table (if present)
- If `disableLLMAuth` is active and response is empty: `_Auth is disabled._`

---

## 8. Error Responses

### 8.1 `LLMError`

```
❌ [validation_failed] The endpoint ID is invalid — it must be a 32-character hex string.
Use the search tool to find the correct endpoint ID.
```

**Format:**
- `❌ [{code}] {message}`
- If `hint` is present, append it on a new line in a ` ``` ` block:

```
❌ [invoke_error] The API request failed — the server may be unreachable or returned an error.

Technical details:
```
connection refused: dial tcp 127.0.0.1:8080: connect: connection refused
```
```

**Error code prefixes:**

| Code | Icon | Tone |
|------|------|------|
| `validation_failed` | ❌ | "Fix your input and try again" |
| `not_found` | 🔍 | "Search for the correct ID" |
| `rate_limit` | ⏳ | "Wait before retrying" |
| `invoke_error` | 🌐 | "The server may be down" |

---

## 9. Quick Reference

| Tool | Response Type | Format | Section |
|------|-------------|--------|---------|
| `spec_list` | `SpecsResponse` | Table (ID, Domain) | 2.1 |
| `spec_by_id` | `SpecByIDResponse` | Sections (## Spec, ## Collections) | 3.1 |
| `collection_by_spec` | `CollectionsResponse` | Context + Table | 2.2 |
| `collection_by_id` | `CollectionByIDResponse` | Context + Table | 2.3 |
| `tag_by_collection` | `TagsByCollectionResponse` | Context + Table | 2.4 |
| `tag_by_spec` | `TagsBySpecResponse` | Table | 2.5 |
| `tag_by_id` | `TagByIDResponse` | Single line | 2.6 |
| `endpoint_by_tag` | `EndpointsByTagResponse` | Context + Table | 2.7 |
| `endpoint_by_collection` | `EndpointsByCollectionResponse` | Context + Table | 2.8 |
| `endpoint_by_spec` | `EndpointsBySpecResponse` | Table | 2.9 |
| `endpoint_by_id` | `EndpointByIDResponse` | Sections (## Spec, ## Collection, ## Tag, ## Endpoint) | 3.2 |
| `search` | `SearchResponse` | Table | 2.10 |
| `inspect` | `InspectResponse` | Sections with code blocks | 4.1 |
| `invoke` | `InvokeResponse` | Sections (Status, Headers, Body) | 5.1 / 5.2 |
| `auth` | `AuthResponse` | Key-Value | 7.1 |
| `info` | `InfoResponse` | Two-column sections | 6.1 |
| error | `LLMError` | `❌ [code] message` | 8.1 |
