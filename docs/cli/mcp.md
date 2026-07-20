# mcp

Start the MCP server.

## Syntax

```bash
swag2mcp mcp [workspace] [flags]
```

## Arguments

| Argument | Description |
|----------|-------------|
| `workspace` | Workspace path (optional) |

## Flags

| Flag | Description |
|------|-------------|
| `--transport` | Transport: stdio, sse, streamable-http |
| `--http-addr` | HTTP server address (default 127.0.0.1:8080) |
| `--http-path` | MCP endpoint path (default /mcp) |
| `--auth-token` | Bearer token for HTTP auth |
| `--logfile` | Log file path |
| `--disable-llm-auth` | Disable auth tool |
| `--dump-dir` | Dump directory |
| `--tags` | Tag filter |

## Usage

=== "stdio (default)"
    ```bash
    swag2mcp mcp
    ```

=== "HTTP SSE"
    ```bash
    swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080
    ```

=== "Streamable HTTP"
    ```bash
    swag2mcp mcp --transport streamable-http --http-addr 0.0.0.0:8080
    ```

=== "With auth"
    ```bash
    swag2mcp mcp --transport sse --auth-token "my-secret"
    ```

=== "Custom workspace"
    ```bash
    swag2mcp mcp ./my-workspace
    ```

=== "Disable auth tool"
    ```bash
    swag2mcp mcp --disable-llm-auth
    ```

## Output

On success:

```
MCP server listening on http://127.0.0.1:8080/mcp
```
