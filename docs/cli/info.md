# info

## Purpose

Show a comprehensive summary of the swag2mcp runtime as **JSON**. This includes version, workspace path, specs summary, HTTP client settings, MCP transport configuration, auth methods, and mock mode status.

## When to use

- You want a machine-readable overview of the workspace
- You need to check the runtime configuration for debugging
- You want to see how many specs and endpoints are active
- You need to verify HTTP client or MCP transport settings

## Syntax

```bash
swag2mcp info [path]
```

## Arguments

| Argument | Position | Required | Description |
|----------|----------|----------|-------------|
| `path` | 1 | No | Workspace directory. If omitted, resolves via path resolution rules. |

## Flags

None.

## How it works

```bash
swag2mcp info
swag2mcp info ./my-workspace
```

## Output

The output is a JSON object with the following structure:

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
    "proxy": "none",
    "follow_redirects": true,
    "max_redirects": 10,
    "randomize": false
  },
  "mcp": {
    "transport": "stdio",
    "addr": ":8080",
    "path": "/mcp"
  },
  "auth_methods": ["bearer", "api-key"],
  "mock_enabled": false
}
```

## Post-command verification

Use `info` to confirm that the workspace loaded correctly and all specs are active before starting the MCP server.

## Nuances

- **Auto-init:** If no config file exists, `info` automatically runs the init wizard first.
- **JSON only:** The output is always JSON. For human-readable output, use `ls`.
- **`max_response_size`:** Shown in human-readable format (e.g., `"1 KB"`, `"2 MB"`).
- **No full-text index:** `info` disables full-text indexing since it only needs config and spec metadata.
