# Troubleshooting

## Installation

### command not found

```bash
# Check Go installation
go version

# Check GOPATH
go env GOPATH

# Add to PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### permission denied

```bash
# For go install
sudo chown -R $(whoami) $(go env GOPATH)

# For binary
sudo chmod +x /usr/local/bin/swag2mcp
```

## Configuration

### config not found

```bash
# Create config
swag2mcp init

# Or specify path
swag2mcp mcp /path/to/workspace
```

### YAML parsing error

```bash
# Check syntax
swag2mcp validate

# Ensure correct indentation (2 spaces)
```

### invalid spec location

```bash
# Check URL
curl -I https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml

# Check path
ls -la ./specs/my-api.yaml
```

## MCP Server

### port already in use

```bash
# Find process
lsof -i :8080

# Kill process
kill <PID>

# Or change port
swag2mcp mcp --http-addr 127.0.0.1:9090
```

### connection refused

```bash
# Ensure server is running
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080

# Check in another terminal
curl http://127.0.0.1:8080/health
```

### MCP tools not showing

```bash
# Check specs are loaded
swag2mcp ls

# Check logs
swag2mcp mcp --logfile /tmp/swag2mcp.log
```

## Authentication

### 401 Unauthorized

```bash
# Check token
swag2mcp info

# Check config
swag2mcp validate
```

### 403 Forbidden

Check API permissions. Additional scopes may be needed.

### Token expired

OAuth2 tokens refresh automatically. Bearer tokens need manual update.

## Search

### No results

```bash
# Check spec is loaded
swag2mcp ls

# Try a different query
swag2mcp search "pet"
```

## Performance

### Slow responses

```yaml
http_client:
  timeout: 60s
```

### Large responses

```yaml
http_client:
  max_response_size: 8388608
```
