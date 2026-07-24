# Response Size Management

## Overview

API responses can be very large — sometimes too large to fit in the LLM's context window. swag2mcp automatically manages response sizes by saving oversized responses to disk and providing tools to explore them.

## How it works

1. **You call `invoke`** — swag2mcp makes the API request
2. **If the response is small** (within the limit) — it is returned inline to the LLM
3. **If the response is too large** (exceeds the limit) — it is saved to `{workspace}/responses/` as a JSON file. The LLM receives a file reference instead of the full response

### Example: small response (inline)

```json
{
  "statusCode": 200,
  "body": {
    "id": 1,
    "name": "Rex",
    "status": "available"
  }
}
```

### Example: large response (file reference)

```json
{
  "statusCode": 200,
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

## Configuration

```yaml
http_client:
  max_response_size: 1048576  # 1 MB in bytes
```

### max_response_size

- **Type:** `int` (bytes)
- **Default:** `1048576` (1 MB)
- **Range:** 256 to 10,485,760 bytes (10 MB)
- **Effect:** Responses larger than this are saved to disk instead of returned inline
- **When to increase:** APIs that return large datasets (reports, logs, analytics)
- **When to decrease:** Limited LLM context window, or when you prefer file-based access

## Working with large responses

When `invoke` returns a `fileRef`, use these three tools to explore the data:

### 1. response_outline — understand the structure

Get a structural summary of the response: keys, types, array lengths, and navigation hints.

```json
→ response_outline(path: "/path/to/file.json")
← {
    "type": "object",
    "size": 1572864,
    "keys": ["data", "meta"],
    "itemCount": 500,
    "compressionHints": [
      "response_compress(path, 'first_of_array', 'data')",
      "response_compress(path, 'sample_array', 'data', arrayHead=3, arrayTail=2)"
    ]
  }
```

### 2. response_compress — get a smaller version

Compress the data to fit inline. Multiple compression modes let you choose the right trade-off.

| Mode | Description | Best for |
|------|-------------|----------|
| `first_of_array` | Keep only the first element of an array | When all elements have the same structure |
| `sample_array` | Keep head (3) and tail (2) of an array | When you need to see the range of values |
| `truncate_strings` | Shorten every string to N characters | When strings are very long |
| `keys_only` | Replace values with their type names | When you only need the structure |
| `select_keys` | Keep only specified keys | When you need specific fields |

```json
→ response_compress(path: "/path/to/file.json", mode: "first_of_array", jsonPath: "data")
← {
    "body": [{ "id": 1, "name": "Rex" }],
    "hint": "Compressed array from 500 to 1 item using first_of_array mode"
  }
```

### 3. response_slice — extract a specific fragment

Get a specific element or value by JSON path or line range.

```json
→ response_slice(path: "/path/to/file.json", jsonPath: "data.0")
← {
    "slice": {
      "value": { "id": 1, "name": "Rex" },
      "jsonPath": "data.0",
      "nextPath": "data.1",
      "prevPath": null
    }
  }
```

## Complete workflow

```
1. invoke(endpoint) → fileRef (response is 1.5 MB)
2. response_outline(path) → structure: { data: Array(500) }
3. response_compress(path, mode: "first_of_array", jsonPath: "data") → first item
4. response_slice(path, jsonPath: "data.0") → full first item details
5. response_slice(path, jsonPath: "data.1") → second item
```

## Automatic cleanup

When the MCP server starts (`swag2mcp mcp`), response files older than 48 hours are automatically removed. You can also clean them manually:

```bash
swag2mcp clean
```

## Important notes

- **The limit is in bytes** — `1048576` = 1 MB, `2097152` = 2 MB, etc.
- **File references include an open command** — on macOS it's `open`, on Linux it's `xdg-open`
- **Response files are named with random suffixes** — no conflicts between concurrent calls
- **The responses directory is created automatically** — no manual setup needed
