# Execution Tools

Execution tools are the core of swag2mcp: **search** finds endpoints when you don't have an ID, **inspect** reveals the full OpenAPI contract, and **invoke** makes the actual API call. Always use them in this order: search → inspect → invoke.

---

## search

### Purpose

The only tool for finding endpoints when you don't have an endpoint ID. Performs full-text search across all endpoints in all specs using the bluge search engine.

### When to use

- When you don't know the endpoint ID
- When you want to find endpoints by keywords, method, tag, or path
- When you need to discover what endpoints exist for a specific feature

### How it works

Searches the full-text index across all specs. Supports structured queries with field filters, boolean operators, fuzzy matching, wildcards, and more.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `query` | string | Yes | Search query (supports structured syntax) |
| `limit` | int | Yes | Maximum results to return (1-50) |

### Query syntax

| Example | Description |
|---------|-------------|
| `pet` | Simple text search across all fields |
| `method:GET` | Filter by HTTP method |
| `tag:pet` | Filter by tag name |
| `path:"/api/v1/users"` | Exact path search |
| `+method:POST +tag:pet` | Must match both conditions |
| `-method:DELETE` | Exclude DELETE methods |
| `create~` | Fuzzy search (typo-tolerant) |
| `path:/api/v1/*` | Wildcard path search |
| `/pattern/` | Regex search |
| `term^3` | Boost a term's relevance |

**Searchable fields:** `method` (keyword), `tag` (keyword), `path` (text), `summary` (text), `_all` (default text field).

**Not supported:** parentheses for grouping, explicit `AND`/`OR` operators, field grouping.

### Response

```json
{
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "tagName": "forecast",
      "collectionId": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "collectionTitle": "Weather Forecast",
      "specId": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
      "specDomain": "meteo",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Get weather forecast for a location"
    }
  ]
}
```

Each result includes the full ancestry (spec → collection → tag) so the LLM can navigate to related endpoints.

### Nuances

- `limit` must be between 1 and 50 (returns `validation_failed` otherwise)
- `query` is required (returns `validation_failed` if empty)
- Results are returned in relevance order (best match first)
- Use field filters (`method:GET`, `tag:pet`) to narrow results
- For exact path matching, use quotes: `path:"/v1/forecast"`

---

## inspect

### Purpose

Retrieve the full OpenAPI operation object for an endpoint: all parameters, request body schema, response schemas, base URL, and full URL. This is the tool to call **before** `invoke` to understand the endpoint's contract.

### When to use

- Always before `invoke` — you need the full contract to make a correct call
- When you need to explain an API's technical details to the user
- When you need to know required parameters, request body structure, or response format

### How it works

Looks up the endpoint in the index and returns the complete OpenAPI operation object with all schemas resolved.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `endpointId` | string | Yes | 32-character MD5 hash of the endpoint |

### Response

```json
{
  "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
  "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
  "collectionId": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
  "specId": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
  "specDomain": "meteo",
  "method": "POST",
  "path": "/pet",
  "baseUrl": "https://meteo.swagger.io/v2",
  "fullUrl": "https://meteo.swagger.io/v2/pet",
  "operation": {
    "id": "addPet",
    "tags": ["pet"],
    "summary": "Add a new pet",
    "description": "Add a new pet to the store",
    "deprecated": false,
    "parameters": [
      {
        "name": "petId",
        "in": "path",
        "description": "ID of the pet",
        "required": true,
        "schema": {
          "type": "integer",
          "format": "int64"
        }
      }
    ],
    "requestBody": {
      "description": "Pet object to add",
      "required": true,
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "name": { "type": "string" },
              "status": { "type": "string", "enum": ["available", "pending", "sold"] }
            },
            "required": ["name"]
          }
        }
      }
    },
    "responses": {
      "200": {
        "description": "Successful operation",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/Pet"
            }
          }
        }
      },
      "405": {
        "description": "Invalid input"
      }
    }
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `baseUrl` | string | Base URL of the API (from config) |
| `fullUrl` | string | Full URL of the endpoint (base + path) |
| `operation.parameters[]` | array | Parameters with name, location (path/query/header/cookie), description, required flag, and schema |
| `operation.requestBody` | object | Request body with content type and schema |
| `operation.responses` | map | Response codes with descriptions and schemas |
| `operation.deprecated` | bool | Whether the endpoint is deprecated |

### Nuances

- Returns `not_found` if the endpoint does not exist
- This is the **only** tool that returns the full OpenAPI operation — `endpoint_by_id` returns only a summary
- Always call `inspect` before `invoke` to understand required parameters and body structure
- The `operation` object includes `$ref` references that are resolved to their full schema definitions

---

## invoke

### Purpose

Execute a real API call to an endpoint. This is the only tool that makes actual HTTP requests. Auth is applied automatically — you do not need to call `auth` first.

### When to use

- Only after calling `inspect` to understand the endpoint's contract
- Only with explicit user confirmation for destructive operations (POST, PUT, PATCH, DELETE)
- When the user asks to call an API and you have all required parameters

### How it works

1. Looks up the endpoint in the index
2. Substitutes path parameters into the URL
3. Appends query parameters
4. Adds headers and cookies
5. Serializes the request body as JSON
6. Automatically obtains and applies auth (token, headers, query params)
7. Makes the HTTP request
8. Returns the response or saves it to a file if too large

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `endpointId` | string | Yes | 32-character MD5 hash of the endpoint |
| `parameters` | object | No | Path, query, and header parameters as key-value pairs |
| `requestBody` | object | No | Request body for POST/PUT/PATCH requests |
| `headers` | object | No | Additional HTTP headers to send |
| `cookies` | object | No | Additional HTTP cookies to send |

### Response (inline)

```json
{
  "statusCode": 200,
  "headers": {
    "content-type": "application/json"
  },
  "body": {
    "id": 1,
    "name": "Rex",
    "status": "available"
  }
}
```

### Response (file reference — when body exceeds size limit)

```json
{
  "statusCode": 200,
  "headers": {
    "content-type": "application/json"
  },
  "fileRef": {
    "path": "/Users/user/.swag2mcp/responses/response_a1b2c3d4.json",
    "size": 1572864,
    "sizeHint": "1.5 MB",
    "maxSizeHint": "2 KB",
    "message": "Response exceeds the 2 KB limit and has been saved to disk.",
    "openCmd": "open /Users/user/.swag2mcp/responses/response_a1b2c3d4.json"
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `statusCode` | int | HTTP response status code |
| `headers` | object | HTTP response headers |
| `body` | any | Response body (present when within size limit) |
| `fileRef` | object | File reference (present when body exceeds size limit) |

### Working with large responses

When `invoke` returns a `fileRef`, use the response tools to explore the data:

1. **`response_outline(path)`** — get the structural summary (keys, types, array lengths)
2. **`response_compress(path, mode)`** — compress the data to fit inline
3. **`response_slice(path, jsonPath)`** — extract a specific fragment

### Nuances

- **Auth is automatic:** The `invoke` tool automatically obtains and applies authentication from the spec's auth configuration. You do **not** need to call `auth` first.
- **Rate limiting:** Each endpoint has a 10-second cooldown. A second call to the same endpoint within 10 seconds is silently blocked (returns `rate_limit` error).
- **Response size limit:** Default is 2 KB (configurable via `max_response_size`). If the response exceeds this limit, it is saved to `{workspace}/responses/` and a `FileReference` is returned instead of inline `body`.
- **Parameter handling:** Path parameters are substituted into the URL. Query parameters are appended. Parameters from the request override operation spec defaults.
- **Request body:** For POST/PUT/PATCH, the body is serialized as JSON. `Content-Type` is set to `application/json` automatically.
- **Error handling:** HTTP errors (non-2xx) are returned as `invoke_error` with the status code and response body in the hint.
- **Destructive operations:** Never invoke POST/PUT/PATCH/DELETE without explicit user confirmation.
