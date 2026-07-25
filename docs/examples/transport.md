# Transport Examples

## stdio

```bash
swag2mcp mcp
```

Default. LLM client runs swag2mcp as a child process.

## SSE

```bash
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080
```

HTTP server with Server-Sent Events.

## Streamable HTTP

```bash
swag2mcp mcp --transport streamable-http --http-addr 0.0.0.0:8080
```

HTTP server with streaming.

## With Auth

```bash
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080 --auth-token "my-secret"
```

Full examples in `examples/mcp-transport/`.
