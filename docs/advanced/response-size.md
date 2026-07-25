# Response Size Management

swag2mcp automatically manages API response sizes.

## How It Works

1. **Limit**: default 2 KB (2048 bytes)
2. **Exceeded**: response is saved to disk
3. **FileReference**: LLM gets a file reference instead of full response

## Configuration

```yaml
global:
  http_client:
    max_response_size: 2048  # in bytes
```

## FileReference

When response exceeds the limit:

```json
{
  "fileRef": {
    "path": "/Users/user/.swag2mcp/responses/2024-01-01/get-pets-abc123.json",
    "size": 15000,
    "mimeType": "application/json"
  }
}
```

## Tools for Large Responses

| Tool | Description |
|------|-------------|
| `response_outline` | Show response structure |
| `response_compress` | Compress response (first item, sample, keys only) |
| `response_slice` | Get response fragment |

## Example Workflow

```
1. invoke(endpoint) → fileRef (15 KB response)
2. response_outline(path) → structure: { pets: Array(100) }
3. response_compress(path, mode: "first_of_array") → first pet
4. response_slice(path, jsonPath: "pets.0") → first pet details
```

## Compression Modes

| Mode | Description |
|------|-------------|
| `first_of_array` | First array element only |
| `sample_array` | Head (3) and tail (2) of array |
| `truncate_strings` | Shorten strings to N chars |
| `keys_only` | Object keys only |
| `select_keys` | Selected keys only |
