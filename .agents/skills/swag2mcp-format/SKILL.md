---
name: swag2mcp-format
description: |
  Compact response formatting rules for swag2mcp MCP tools.
  Use when displaying results from ANY swag2mcp MCP tool:
  spec_list, spec_by_id, collection_by_*, tag_by_*,
  endpoint_by_*, search, inspect, invoke, auth, info.
  Automatically triggered on every swag2mcp tool response.
  Ensures consistent, ergonomic, human-readable markdown with
  tight tables, inline headers, and compact schemas.
license: MIT
metadata:
  author: mmadfox
  version: "2.0.0"
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

### 1.1 HTTP Methods

Display methods as plain uppercase text in table cells and headers. Do not use emoji or color markers.

| Method | Display |
|--------|---------|
| GET    | GET |
| POST   | POST |
| PUT    | PUT |
| PATCH  | PATCH |
| DELETE | DELETE |
| HEAD   | HEAD |
| OPTIONS| OPTIONS |

### 1.2 Sorting

- Endpoints within a list: sort by **HTTP method** (GET < POST < PUT < PATCH < DELETE < HEAD < OPTIONS), then by **path** alphabetically.
- Collections, tags, specs: sort by **title** alphabetically (case-insensitive).
- Search results: keep the order returned by the search engine (relevance-based).

### 1.3 ID Display

- Show IDs in the **last column** of tables.
- Truncate to first **6 characters** and wrap in backticks: `` `a1b2c3` ``.
- Full ID is available via `endpoint_by_id` or `inspect` if needed.

### 1.4 Pagination

**Never display more than 5 items at once.** Paginate all list responses (`search`, `endpoint_by_*`, `tag_by_*`, `collection_by_*`, `spec_list`). Exceptions: single-item responses (`tag_by_id`, `endpoint_by_id`, `spec_by_id`).

- Show first 5 items, then append: `▸ {shown}/{total} · reply “more” to load next 5`
- Wait for user confirmation before showing the next batch.

**Example:**

```
**Search: pet** (14)

| Method | Path | Summary | ID |
|--------|------|---------|----|
| GET | /pet/{petId} | Find pet by ID | e1f2g3 |
| GET | /pet/findByStatus | Find pets by status | i5j6k7 |
| GET | /pet/findByTags | Find pets by tags | pets | m9n0o1 |
| POST | /pet | Add a new pet | q3r4s5 |
| DELETE | /pet/{petId} | Delete a pet | u7v8w9 |

▸ 5/14 · reply “more” to load next 5
```

When user confirms, show the next batch with the same footer.

### 1.5 Empty States

- Empty list: `—`
- Empty/omitted field: `—`

### 1.6 Deprecated Endpoints

Append `[deprecated]` after the summary: `Update user [deprecated]`.

---

## 2. List Responses (Markdown Tables)

These responses contain flat lists of items. Always use a markdown table.

### 2.1 `spec_list` → `SpecsResponse`

```
**Specs (2)**

| Domain | ID |
|--------|----|
| Petstore API | `a1b2c3` |
| Weather API | `e5f6g7` |
```

**Format:**
- Header: `**Specs ({count})**`
- Table columns: `Domain`, `ID`
- Empty: `**Specs (0)** —`
- Sort by domain alphabetically.

### 2.2 `collection_by_spec` → `CollectionsResponse`

```
**Petstore API** · Collections (3)

| Collection | Tags | Methods | ID |
|------------|------|---------|----|
| Pet Operations | 2 | 8 | `c1d2e3` |
| Store | 1 | 4 | `g5h6i7` |
| User Management | 3 | 6 | `k9l0m1` |
```

**Format:**
- Header: `**{Domain}** · Collections ({count})`
- Table columns: `Collection`, `Tags`, `Methods`, `ID`
- If `LLMTitle` differs from title, show in parentheses: `Pet Operations (pets)`.
- Sort collections alphabetically.

### 2.3 `collection_by_id` → `CollectionByIDResponse`

```
**Petstore API › Pet Operations** · 8 methods · `c1d2e3...`

**Tags (2)**

| Tag | Methods | ID |
|-----|---------|----|
| pets | 5 | `t1u2v3` |
| store | 3 | `x5y6z7` |
```

**Format:**
- One-line header: `**{Domain} › {Collection}** · {countMethods} methods · `{id}`
- Table columns: `Tag`, `Methods`, `ID`
- Sort tags alphabetically.

### 2.4 `tag_by_collection` → `TagsByCollectionResponse`

Same format as `collection_by_id`:

```
**Petstore API › Pet Operations** · 8 methods · `c1d2e3...`

**Tags (2)**

| Tag | Methods | ID |
|-----|---------|----|
| pets | 5 | `t1u2v3` |
| store | 3 | `x5y6z7` |
```

### 2.5 `tag_by_spec` → `TagsBySpecResponse`

```
**Petstore API** · Tags (5)

| Tag | Methods | ID |
|-----|---------|----|
| pets | 5 | `t1u2v3` |
| store | 3 | `x5y6z7` |
| user | 6 | `b9c0d1` |
```

**Format:**
- Header: `**{Domain}** · Tags ({count})`
- Table columns: `Tag`, `Methods`, `ID`
- Sort tags alphabetically.

### 2.6 `tag_by_id` → `TagByIDResponse`

```
**Tag:** pets · 5 methods · `t1u2v3...`
```

**Format:**
- Single line: `**Tag:** {Title} · {countMethods} methods · `{id}``

### 2.7 `endpoint_by_tag` → `EndpointsByTagResponse`

```
**Petstore API › Pet Operations › pets** · 5 endpoints

| Method | Path | Summary | ID |
|--------|------|---------|----|
| GET | /pet/{petId} | Find pet by ID | e1f2g3 |
| GET | /pet/findByStatus | Find pets by status | i5j6k7 |
| POST | /pet | Add a new pet | m9n0o1 |
| PUT | /pet | Update an existing pet | q3r4s5 |
| DELETE | /pet/{petId} | Delete a pet | u7v8w9 |

▸ 5/5
```

**Format:**
- One-line header: `**{Domain} › {Collection} › {Tag}** · {count} endpoints`
- Table columns: `Method`, `Path`, `Summary`, `ID`
- Sort by method order (GET < POST < PUT < PATCH < DELETE < HEAD < OPTIONS), then path.

### 2.8 `endpoint_by_collection` → `EndpointsByCollectionResponse`

```
**Petstore API › Pet Operations** · 8 endpoints

| Method | Path | Summary | Tag · ID |
|--------|------|---------|----------|
| GET | /pet/{petId} | Find pet by ID | pets `e1f2g3` |
| GET | /pet/findByStatus | Find pets by status | pets `i5j6k7` |
| POST | /pet | Add a new pet | pets `m9n0o1` |
| POST | /store/order | Place order | store `a1b2c3` |
| DELETE | /pet/{petId} | Delete a pet | pets `u7v8w9` |

▸ 5/8 · reply “more” to load next 5
```

**Format:**
- One-line header: `**{Domain} › {Collection}** · {count} endpoints`
- Table columns: `Method`, `Path`, `Summary`, `Tag · ID`
- Sort by tag, then method order, then path.

### 2.9 `endpoint_by_spec` → `EndpointsBySpecResponse`

```
**Petstore API** · 14 endpoints

| Method | Path | Summary | Collection · Tag · ID |
|--------|------|---------|-----------------------|
| GET | /pet/{petId} | Find pet by ID | Pet Operations · pets `e1f2g3` |
| POST | /pet | Add a new pet | Pet Operations · pets `m9n0o1` |
| POST | /store/order | Place order | Store · store `a1b2c3` |
| DELETE | /pet/{petId} | Delete a pet | Pet Operations · pets `u7v8w9` |

▸ 5/14 · reply “more” to load next 5
```

**Format:**
- One-line header: `**{Domain}** · {count} endpoints`
- Table columns: `Method`, `Path`, `Summary`, `Collection · Tag · ID`
- Sort by collection, tag, method order, then path.

### 2.10 `search` → `SearchResponse`

```
**Search: pet** · 14 results

| Method | Path | Summary | Spec · Collection · Tag · ID |
|--------|------|---------|------------------------------|
| GET | /pet/{petId} | Find pet by ID | Petstore · Pet Operations · pets `e1f2g3` |
| POST | /pet | Add a new pet | Petstore · Pet Operations · pets `m9n0o1` |
| POST | /store/order | Place order | Petstore · Store · store `a1b2c3` |

▸ 5/14 · reply “more” to load next 5
```

**Format:**
- Header: `**Search: {query}** · {count} results`
- Table columns: `Method`, `Path`, `Summary`, `Spec · Collection · Tag · ID`
- Keep search engine's relevance order.
- Paginate as in 1.4.

---

## 3. Detail Responses (Sections with Headers)

### 3.1 `spec_by_id` → `SpecByIDResponse`

```
**Petstore API** · `a1b2c3...`

**Collections (3)**

| Collection | Tags | Methods | ID |
|------------|------|---------|----|
| Pet Operations | 2 | 8 | `c1d2e3` |
| Store | 1 | 4 | `g5h6i7` |
| User Management | 3 | 6 | `k9l0m1` |
```

**Format:**
- One-line header: `**{Domain}** · `{id}`
- Table columns: `Collection`, `Tags`, `Methods`, `ID`
- Sort collections alphabetically.

### 3.2 `endpoint_by_id` → `EndpointByIDResponse`

```
**Petstore API › Pet Operations › pets** · `e1f2g3...`

| Method | Path | Summary |
|--------|------|---------|
| GET | /pet/{petId} | Find pet by ID |
```

**Format:**
- One-line header: `**{Domain} › {Collection} › {Tag}** · `{id}`
- Single-row table: `Method`, `Path`, `Summary`.

---

## 4. Inspect Response

### 4.1 `inspect` → `InspectResponse`

```
**GET /pet/{petId}** · Petstore API

`https://petstore.swagger.io/v2/pet/{petId}`

Find pet by ID. Returns a single pet by its ID.

**Parameters**

| Name | In | Type | Req | Description |
|------|----|------|-----|-------------|
| petId | path | integer | yes | ID of pet to return |
| api_key | header | string | no | API key for authentication |

**Request body** — no

**Responses**

| Code | Description | Schema |
|------|-------------|--------|
| 200 | successful operation | `{ id, name, category, status }` |
| 400 | Invalid pet ID | `{ code, message }` |
| 404 | Pet not found | — |
```

**Format:**
- One-line title: `**{Method} {Path}** · {Domain}`
- Full URL on the next line in backticks.
- Summary + description in one compact line (omit if empty).
- If deprecated: append ` [deprecated]` to the title line.
- `**Parameters**` section:
  - If empty: `**Parameters** — none`
  - Table columns: `Name`, `In`, `Type`, `Req`, `Description`
  - `Req`: `yes` / `no`
  - `Type`: schema type (e.g., `string`, `integer`, `array<Pet>`)
- `**Request body** — {yes|no}` line. If yes, append the description inline.
- `**Responses**` section:
  - Table columns: `Code`, `Description`, `Schema`
  - `Schema`: top-level field list in braces. Example: `{ id, name, category }`.
  - For arrays: `{ Pet[] }` or `[Pet]`.
  - If no content: `—`.
  - If the user asks for the full schema, switch to a ````json```` code block.

---

## 5. Invoke Response

**Default rule:** show only the response body. Do not show HTTP status or headers unless the response is not 2xx or the user explicitly asks for full details.

### 5.1 `invoke` → `InvokeResponse` (success, 2xx)

Default output:

```
Курс биткоина (BTCUSDT): $63,999.00
```

JSON body:

```json
{
  "price": 63999.00,
  "symbol": "BTCUSDT"
}
```

Empty body (e.g. 204):

```
No Content
```

**Format:**
- Show only the body.
- JSON: ` ```json ` block with pretty-printed JSON.
- Text/HTML/XML: ` ``` ` block with raw text.
- Empty: `No Content`.

### 5.2 `invoke` → `InvokeResponse` (non-2xx)

```
**404 Not Found**

`content-type: application/json`

```json
{
  "error": "not found"
}
```
```

**Format:**
- Status line: `**{code} {text}**` (e.g., `**404 Not Found**`, `**500 Internal Server Error**`).
- Headers: inline in backticks, separated by ` · `. Omit `Content-Length`, `Date`, `Connection`, `Keep-Alive`.
  - If no headers: omit the headers line entirely.
- Body: ` ```json ` block, ` ``` ` block, or `—` for empty body.

### 5.3 `invoke` → `InvokeResponse` (large response with FileReference)

```
📄 Response body (2.5 KB) exceeds max size (1 KB). Saved to:
`~/.swag2mcp/responses/petstore-get-pet-findByStatus-abc123.json`
```

**Format:**
- One compact line: `📄 Response body ({size}) exceeds max size ({maxSize}). Saved to:`
- File path in backticks on the next line.
- Do not show status or headers unless the user asks for full details.

### 5.4 Full details on explicit request

If the user asks for details (`show full response`, `show headers`, `what was the status?`, etc.):

```
**200 OK** · 2 headers · 48 B body

`content-type: application/json` · `x-request-id: abc-123-def`

```json
{
  "id": 1,
  "name": "Buddy",
  "status": "available"
}
```
```

**Format:**
- Status line: `**{code} {text}** · {headerCount} headers · {bodySize} body`.
- Headers: inline in backticks, separated by ` · `.
- Body: ` ```json ` block, ` ``` ` block, or `—`.

---

## 6. Info Response (Two-Column Layout)

### 6.1 `info` → `InfoResponse`

```
**System** · 1.2.3 · `~/.swag2mcp` · uptime 2h15m30s

**Specs** · total 3 · active 2 · disabled 1 · collections 7 · endpoints 42

**HTTP Client** · timeout 30s · follow redirects yes · max redirects 5 · max response 1 KB

**MCP** · stdio · auth yes

**Auth** · bearer, oauth2-cc

**Mock** · no
```

**Format:**
- Render each section as one compact line: `**{Section}** · {key value pairs}`.
- Booleans: `yes` / `no`.
- Paths and values in backticks where useful.
- `MaxResponseSize`: convert bytes to human-readable (e.g., `1048` → `1 KB`).
- Omit empty sections entirely.
- If nested objects (Proxy, Headers, Cookies) are present and non-empty, append them inline after the main HTTP Client line or add a second compact line only when necessary.

---

## 7. Auth Response (Key-Value)

### 7.1 `auth` → `AuthResponse`

```
**Token:** `Bearer eyJhbGciOiJIUzI1NiIs...`

**Headers:** `Authorization: Bearer eyJhbGciOiJIUzI1NiIs...` · `X-API-Key: abc123`

**Query:** `api_key: abc123`
```

**Format:**
- Token line: `**Token:** \`{value}\``. Truncate long tokens to 40 chars + `...`.
- Headers line: `**Headers:** \`{name}: {value}\`` separated by ` · `. Omit section if empty.
- Query line: `**Query:** \`{name}: {value}\`` separated by ` · `. Omit section if empty.
- If `disableLLMAuth` is active and response is empty: `**Auth** — disabled`.

---

## 8. Error Responses

### 8.1 `LLMError`

```
**validation_failed** · Invalid endpoint ID — must be 32 hex chars. Use search to find the correct ID.
```

**Format:**
- Single line: `**{code}** · {message}` (or wrap to next line if message is long).
- If `hint` is present, append on a new line in a ` ``` ` block:

```
**invoke_error** · The API request failed — the server may be unreachable.

```
connection refused: dial tcp 127.0.0.1:8080: connect: connection refused
```
```

**Error code prefixes:**

| Code | Tone |
|------|------|
| `validation_failed` | "Fix your input and try again" |
| `not_found` | "Search for the correct ID" |
| `rate_limit` | "Wait before retrying" |
| `invoke_error` | "The server may be down" |

---

## 9. Quick Reference

| Tool | Response Type | Format | Section |
|------|-------------|--------|---------|
| `spec_list` | `SpecsResponse` | Compact table | 2.1 |
| `spec_by_id` | `SpecByIDResponse` | One-line header + table | 3.1 |
| `collection_by_spec` | `CollectionsResponse` | One-line header + table | 2.2 |
| `collection_by_id` | `CollectionByIDResponse` | One-line header + table | 2.3 |
| `tag_by_collection` | `TagsByCollectionResponse` | Same as `collection_by_id` | 2.4 |
| `tag_by_spec` | `TagsBySpecResponse` | One-line header + table | 2.5 |
| `tag_by_id` | `TagByIDResponse` | Single line | 2.6 |
| `endpoint_by_tag` | `EndpointsByTagResponse` | One-line header + table | 2.7 |
| `endpoint_by_collection` | `EndpointsByCollectionResponse` | One-line header + table | 2.8 |
| `endpoint_by_spec` | `EndpointsBySpecResponse` | One-line header + table | 2.9 |
| `endpoint_by_id` | `EndpointByIDResponse` | Single line + single-row table | 3.2 |
| `search` | `SearchResponse` | One-line header + table | 2.10 |
| `inspect` | `InspectResponse` | Compact sections + schema table | 4.1 |
| `invoke` | `InvokeResponse` | Status line + inline headers + body | 5.1 / 5.2 |
| `auth` | `AuthResponse` | Inline key-value lines | 7.1 |
| `info` | `InfoResponse` | Inline compact lines | 6.1 |
| error | `LLMError` | `**code** · message` | 8.1 |
