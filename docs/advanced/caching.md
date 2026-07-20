# Caching

swag2mcp caches downloaded specs for faster startup.

## Cache Rules

| Source | Behavior |
|--------|----------|
| HTTP/HTTPS URL | Always cached |
| Local file in `specs/` | Used directly |
| Local file outside `specs/` | Copied to cache |

## TTL (Time To Live)

- **Random**: 1 to 48 hours
- **Reset**: on each MCP server start
- **Force**: `swag2mcp update <id>`

## Cache Structure

```
~/.swag2mcp/cache/
├── api.example.com/
│   └── openapi.json
├── petstore.swagger.io/
│   └── swagger.json
└── ...
```

## Clean Cache

```bash
# Full clean
swag2mcp clean

# Automatic (on mcp start)
swag2mcp mcp  # old responses > 48h removed
```

## Cache Key

Key is derived from the spec URL:

```go
cacheKey = md5(location)
```
