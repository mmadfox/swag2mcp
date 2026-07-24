# Utility Tools

Utility tools provide supporting functionality: retrieving auth tokens, getting runtime information, and working with large API responses that don't fit inline.

---

## auth

### Purpose

Retrieve an authentication token, headers, or query parameters for a specific spec. This gives the LLM access to credentials that can be used outside of swag2mcp (e.g., generating a curl command).

### When to use

- Only when the user explicitly asks for the raw token or credentials
- When generating a curl command or code snippet that needs auth
- When the user wants to see what auth method is configured

### When NOT to use

- **Do not** call `auth` before `inspect` or `invoke` — `invoke` automatically obtains and applies authentication
- **Do not** call `auth` just to check if auth is configured — use `info` instead

### How it works

Looks up the spec's auth configuration and executes the auth flow (token exchange, script execution, etc.) to obtain the current credentials.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `specId` | string | Yes | 32-character MD5 hash of the spec |

### Response

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "headers": {
    "Authorization": "Bearer eyJhbGciOiJIUzI1NiIs...",
    "X-API-Key": "my-api-key"
  },
  "queryParams": {
    "api_key": "my-api-key"
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `token` | string | Raw token value (bearer token, API key, etc.) |
| `headers` | object | HTTP headers to include in requests |
| `queryParams` | object | Query parameters to include in requests |

### Nuances

- **Disabled by default in production:** The `--disable-llm-auth` flag (default: `true`) removes the `auth` tool from the MCP tool list entirely. The LLM cannot see or request tokens. Set `--disable-llm-auth=false` to enable it for debugging or short-lived tokens.
- **`invoke` handles auth automatically:** You do not need to call `auth` before `invoke`. The invoke service automatically obtains and applies the correct authentication.
- **Supports 9 auth methods:** `none`, `basic`, `bearer`, `digest`, `hmac`, `oauth2-cc` (client credentials), `oauth2-pwd` (password), `api-key`, `script`.
- Returns `auth_error` if the auth method fails (e.g., OAuth2 token endpoint unreachable, script execution failure).

---

## info

### Purpose

Return a comprehensive summary of the swag2mcp runtime: version, workspace path, active specs, HTTP client settings, MCP transport configuration, auth methods, and mock mode status.

### When to use

- When the user asks about the system configuration
- When you need to check runtime settings (timeout, response size limit, transport)
- When you need to know which auth methods are available
- When troubleshooting configuration issues

### How it works

Returns a pre-computed snapshot of the runtime state. No parameters needed.

### Parameters

None.

### Response

```json
{
  "version": "v1.2.0",
  "workspace": "~/.swag2mcp",
  "uptime": "2h 15m",
  "specs": {
    "total": 4,
    "active": 3,
    "disabled": 1,
    "collections": 6,
    "endpoints": 42
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "follow_redirects": true,
    "max_redirects": 10,
    "randomize": false,
    "proxy": null,
    "headers": {},
    "cookies": []
  },
  "mcp": {
    "transport": "stdio",
    "addr": ":8080",
    "path": "/mcp",
    "auth_enabled": false
  },
  "auth": {
    "methods": ["bearer", "api-key"]
  },
  "mock": {
    "enabled": false
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `version` | string | swag2mcp version |
| `workspace` | string | Workspace directory path |
| `uptime` | string | Server uptime (human-readable) |
| `specs` | object | Spec summary: total, active, disabled, collections, endpoints |
| `http_client` | object | HTTP client configuration |
| `http_client.max_response_size` | string | Max response size in human-readable format (e.g. "2 KB") |
| `mcp` | object | MCP server configuration |
| `auth` | object | Available auth methods |
| `mock` | object | Mock server status |

### Nuances

- `max_response_size` is shown in human-readable format (e.g., `"1 KB"`, `"2 MB"`)
- `uptime` is computed from the server start time
- The data is a snapshot taken at bootstrap time — it reflects the state when the MCP server started

---

## response_outline

### Purpose

Get a high-level structural summary of a large JSON response file that was saved to disk by `invoke`. It returns the shape of the data — keys, types, array lengths, and navigation hints — without returning the actual values.

### When to use

- Immediately after `invoke` returns a `fileRef` (response too large for inline)
- This is the **mandatory first step** in the large-response workflow

### How it works

Reads the saved response file and analyzes its structure: top-level type, keys, array lengths, nesting depth, and compression hints.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `path` | string | Yes | Absolute path from `fileRef.path` |
| `maxDepth` | int | No | Maximum recursion depth (default: 3) |
| `maxArrayItems` | int | No | How many array items to inspect (default: 5) |

### Response

```json
{
  "outline": {
    "type": "object",
    "size": 1572864,
    "lineCount": 12500,
    "depth": 3,
    "structure": {
      "type": "object",
      "keys": ["data", "meta", "error"],
      "data": {
        "type": "array",
        "length": 500,
        "items": {
          "type": "object",
          "keys": ["id", "name", "status", "createdAt"]
        }
      }
    },
    "schemaHint": "object with 3 keys: data (array[500]), meta (object), error (null)",
    "keys": ["data", "meta", "error"],
    "itemCount": 500,
    "itemType": "object",
    "compressionHints": [
      "response_compress(path, 'first_of_array', 'data')",
      "response_compress(path, 'sample_array', 'data', arrayHead=3, arrayTail=2)",
      "response_compress(path, 'keys_only', 'data')",
      "response_compress(path, 'select_keys', 'data', selectKeys=[id, name])"
    ],
    "navigationHints": {
      "paths": ["data", "meta", "error"],
      "arrays": [
        {"path": "data", "length": 500}
      ]
    }
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `type` | string | Top-level type: "object" or "array" |
| `size` | int | File size in bytes |
| `lineCount` | int | Number of lines in the file |
| `depth` | int | Maximum nesting depth inspected |
| `structure` | object | Recursive structure with keys, types, array lengths |
| `schemaHint` | string | One-line summary of the top-level shape |
| `keys` | array | Top-level keys (for objects) |
| `itemCount` | int | Array length (for arrays) |
| `compressionHints` | array | Suggested `response_compress` calls with parameters |
| `navigationHints` | object | Top-level paths and arrays with lengths |

### Nuances

- Returns `validation_failed` if the path is invalid or not inside the responses directory
- Returns `not_found` if the file does not exist
- Returns `validation_failed` if the file is not valid JSON
- The `compressionHints` field provides ready-to-use suggestions for `response_compress` calls

---

## response_compress

### Purpose

Reduce a JSON value inside a saved response file so it fits within the response size limit and can be returned to the LLM inline. Multiple compression modes let you choose the right trade-off between size and information.

### When to use

- After `response_outline` to understand the structure
- When you need to get data from a large response inline
- When `response_slice` is too narrow and you need a broader view

### How it works

Reads the saved response file, navigates to the specified JSON path, applies the compression mode, and returns the compressed result. If the result still exceeds the size limit, it is saved to a new file.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `path` | string | Yes | Absolute path from `fileRef.path` |
| `jsonPath` | string | No | Path to the value to compress (e.g. `data` or `data.0`) |
| `mode` | string | Yes | Compression mode (see table below) |
| `arrayHead` | int | No | Leading items to keep in `sample_array` mode (default: 3) |
| `arrayTail` | int | No | Trailing items to keep in `sample_array` mode (default: 2) |
| `stringLen` | int | No | Max string length in `truncate_strings` mode (default: 80) |
| `selectKeys` | array | No | Keys to keep in `select_keys` mode |

### Compression modes

| Mode | Description | Best for |
|------|-------------|----------|
| `first_of_array` | Keep only the first element of an array | When all elements have the same structure |
| `sample_array` | Keep head and tail of an array | When you need to see the range of values |
| `truncate_strings` | Shorten every string to `stringLen` characters | When strings are very long but structure matters |
| `keys_only` | Replace object values with their type names | When you only need the structure |
| `select_keys` | Keep only specified keys in every object | When you need specific fields from many objects |

### Response

```json
{
  "body": [
    { "id": 1, "name": "Rex", "status": "available" },
    { "id": 2, "name": "Max", "status": "pending" }
  ],
  "hint": "Compressed array from 500 to 2 items using first_of_array mode"
}
```

| Field | Type | Description |
|-------|------|-------------|
| `body` | any | Compressed JSON value (present when within size limit) |
| `fileRef` | object | File reference (present when still too large) |
| `hint` | string | Explanation of what was compressed |

### Nuances

- If the compressed result still exceeds `max_response_size`, it is saved to a new file and a `FileReference` is returned
- Default values: `arrayHead=3`, `arrayTail=2`, `stringLen=80`
- Returns `validation_failed` for invalid path, invalid JSONPath, or non-JSON file
- Returns `not_found` if the file does not exist or JSONPath does not match

---

## response_slice

### Purpose

Extract a specific fragment of a saved JSON response file by logical JSON path or by line range. Unlike `response_compress`, this gives you the raw, unmodified data.

### When to use

- When you need a specific element or value from a large response
- When `response_compress` doesn't give you enough detail
- When you want to navigate through a response step by step

### How it works

Reads the saved response file and extracts a fragment by JSON path (e.g., `data.3.name`) or by line range (e.g., `120-240`). Returns navigation hints for stepping through arrays and objects.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `path` | string | Yes | Absolute path from `fileRef.path` |
| `jsonPath` | string | No | Logical path to the value (e.g. `data.3.name`) |
| `line` | int | No | 1-based line number to center the fragment on |
| `range` | string | No | Line range as `start-end` (e.g. `120-240`) |
| `around` | int | No | Lines to include around `line` (default: 20) |

### Response

```json
{
  "slice": {
    "lines": [120, 130],
    "fragment": "{\n  \"id\": 1,\n  \"name\": \"Rex\"\n}",
    "value": {
      "id": 1,
      "name": "Rex"
    },
    "jsonPath": "data.0",
    "context": "object",
    "isComplete": true,
    "nextLine": 131,
    "prevLine": 119,
    "nextPath": "data.1",
    "prevPath": null
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `lines` | array | 1-based line range [start, end] |
| `fragment` | string | Raw JSON text (when small enough) |
| `value` | any | Extracted JSON value |
| `jsonPath` | string | The JSON path used |
| `context` | string | "object", "array", or "value" |
| `isComplete` | bool | True when the value is a valid JSON fragment |
| `nextLine` | int | Suggested next line for line-based navigation |
| `prevLine` | int | Suggested previous line |
| `nextPath` | string | Suggested next JSON path for array navigation |
| `prevPath` | string | Suggested previous JSON path |

### Nuances

- **Prefer `jsonPath` over line numbers** — JSON paths are stable and descriptive, line numbers change if the file is regenerated
- If the extracted fragment exceeds `max_response_size`, it is saved to a new file and a `FileReference` is returned
- Default `around` is 20 lines
- The response includes `nextPath`/`prevPath` for stepping through arrays and `nextLine`/`prevLine` for line-based navigation
- Returns `validation_failed` for invalid path, invalid JSONPath, invalid line/range, or non-JSON file
- Returns `not_found` if the file does not exist or JSONPath does not match
