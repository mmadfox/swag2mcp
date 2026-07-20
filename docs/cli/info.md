# info

Show system information.

## Syntax

```bash
swag2mcp info [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `-w, --workspace` | Workspace path |

## Usage

```bash
swag2mcp info
```

## Output

```
Version: dev
Workspace: ~/.swag2mcp
Uptime: 2h 15m
Specs: 3 active, 1 disabled
Endpoints: 42 total
HTTP Client:
  Timeout: 30s
  Max Response: 2 KB
  Proxy: none
MCP:
  Transport: stdio
Auth Methods: bearer, api-key
```
