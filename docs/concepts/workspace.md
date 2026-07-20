# Workspace

The workspace is the directory where swag2mcp stores all its data.

## Structure

```
~/.swag2mcp/
├── swag2mcp.yaml          # Configuration file
├── cache/                  # Downloaded spec cache
│   └── api.example.com/
│       └── openapi.json
├── specs/                  # Local spec copies
│   └── my-api.yaml
├── responses/              # Saved API responses
│   └── 2024-01-01/
│       └── get-pets-abc123.json
└── auth_scripts/           # Auth scripts
    └── get-token.sh
```

## Default Path

- **Linux/macOS**: `~/.swag2mcp/`
- **Windows**: `%USERPROFILE%\.swag2mcp\`

## Custom Path

```bash
swag2mcp mcp /path/to/workspace
swag2mcp mcp ./my-workspace
```

## Cleanup

```bash
# Clear cache
swag2mcp clean

# Old responses cleaned automatically (48h TTL)
# Happens on mcp server start
```

## Export and Import

```bash
# Export workspace to ZIP
swag2mcp export --output workspace.zip

# Import from ZIP
swag2mcp import --from-zip workspace.zip
```

## .gitignore

Recommended:

```gitignore
.swag2mcp/
cache/
responses/
```
