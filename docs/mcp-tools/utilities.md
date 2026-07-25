# Utility Tools

Helper tools.

## auth

Get auth token or headers.

**Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `specId` | string | Spec ID |

**Example**:
```
→ auth(specId: "abc123")
← Authorization: Bearer eyJhbGci...
   X-API-Key: my-api-key
```

!!! note
    Disabled with `--disable-llm-auth` flag.

## info

System information.

**Parameters**: none

**Example**:
```
→ info()
← Version: dev
  Uptime: 1h 30m
  Specs: 3 active
  Endpoints: 42
  Transport: stdio
```

## response_outline

Structure of a large response (when invoke returns fileRef).

**Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `path` | string | Response file path |
| `maxDepth` | int | Max depth (default 3) |

## response_compress

Compress a large response.

**Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `path` | string | Response file path |
| `jsonPath` | string | Value path |
| `mode` | string | Compression mode |

**Compression Modes**:
| Mode | Description |
|------|-------------|
| `first_of_array` | First array element |
| `sample_array` | Head and tail of array |
| `truncate_strings` | Shorten strings |
| `keys_only` | Object keys only |
| `select_keys` | Selected keys only |

## response_slice

Fragment of a large response.

**Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `path` | string | Response file path |
| `jsonPath` | string | Value path |
| `line` | int | Line number |
| `range` | string | Line range |
