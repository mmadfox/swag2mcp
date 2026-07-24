# Troubleshooting

## Installation Problems

### swag2mcp: command not found

The binary is not in your PATH.

```bash
# Check if Go is installed
go version

# Find where Go installs binaries
go env GOPATH
# Usually ~/go or ~/go/bin

# Add to PATH (add this to ~/.zshrc or ~/.bashrc)
export PATH=$PATH:$(go env GOPATH)/bin

# Or use the full path
~/go/bin/swag2mcp --version
```

If you downloaded a binary from GitHub Releases, make sure it's in a directory that's in your PATH:

```bash
# Move to /usr/local/bin (macOS/Linux)
sudo mv swag2mcp /usr/local/bin/
```

### permission denied

The binary does not have execute permissions.

```bash
# For go install (fix ownership)
sudo chown -R $(whoami) $(go env GOPATH)

# For downloaded binary
chmod +x /path/to/swag2mcp
```

### Go version too old

swag2mcp requires Go 1.23+.

```bash
go version
# If version < 1.23, update Go:
# https://go.dev/dl/
```

### Mock server not found

The mock server is a separate binary. Install it explicitly:

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

## Configuration Problems

### Configuration file not found

swag2mcp cannot find `swag2mcp.yaml`.

```bash
# Create a new config
swag2mcp init

# Or specify the path explicitly
swag2mcp mcp /path/to/workspace
swag2mcp ls /path/to/workspace
```

**Common cause:** You ran `swag2mcp mcp` from a random directory and it looked for `~/.swag2mcp/` instead of your project's workspace. Always pass the path explicitly.

### Wrong workspace loaded

swag2mcp loaded a different workspace than expected.

**Resolution order:** Explicit `[path]` → current directory (`./`) → `~/.swag2mcp/`. If you run `swag2mcp mcp` without a path from a directory that doesn't have `swag2mcp.yaml`, it falls back to `~/.swag2mcp/`.

**Fix:** Always pass the workspace path: `swag2mcp mcp /path/to/your/workspace`

### YAML parsing error

The config file has invalid YAML syntax.

```bash
# Validate the config
swag2mcp validate

# Common mistakes:
# - Tabs instead of spaces (YAML requires spaces)
# - Missing indentation for nested fields
# - Unquoted strings with special characters (: # & {)
```

**Tip:** Use a YAML linter or editor with YAML support to catch syntax errors.

### Validation fails: "no specifications defined"

The config file exists but has no specs.

```bash
# Add a spec
swag2mcp add spec

# Or edit swag2mcp.yaml and add at least one spec
```

### Validation fails: "duplicate domain"

Two specs have the same `domain` value. Domains must be unique.

```bash
# List current specs
swag2mcp ls

# Check for duplicate domains in swag2mcp.yaml
```

### Validation fails: "invalid spec location"

The `location` URL or file path is not accessible or not a valid spec file.

```bash
# Check if the URL is reachable
curl -I https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml

# Check if the local file exists
ls -la ./specs/my-api.yaml

# Verify the file is valid OpenAPI/Swagger/Postman
# (not just any JSON or HTML page)
```

**Common cause:** The `location` field points to the API endpoint itself (e.g., `https://api.example.com/v1/users`) instead of the spec file URL. The location must point to an OpenAPI/Swagger/Postman file.

## MCP Server Problems

### Port already in use

Another process is using the port.

```bash
# Find the process
lsof -i :8080

# Kill it
kill <PID>

# Or use a different port
swag2mcp mcp --transport sse --http-addr :9090
```

### Connection refused

The MCP server is not running or not reachable.

```bash
# Make sure the server is running
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080

# In another terminal, check the health endpoint
curl http://127.0.0.1:8080/health

# If using a custom path
curl http://127.0.0.1:8080/custom-path/health
```

### MCP tools not showing up in the LLM client

The LLM client cannot see any tools.

```bash
# Check that specs are loaded
swag2mcp ls

# Check that specs are not disabled
swag2mcp validate

# Check the server logs
swag2mcp mcp --logfile /tmp/swag2mcp.log
cat /tmp/swag2mcp.log

# Verify the workspace path in your IDE config is correct
# (must be an absolute path)
```

**Common causes:**
- Wrong workspace path in IDE config
- All specs have `disable: true`
- Specs are filtered out by `--tags`
- Config file doesn't exist at the specified path

### MCP handshake fails (HTTP transport)

For SSE and Streamable HTTP transports, the MCP protocol requires initialization before tool calls work.

```
Step 1: POST /mcp → {"method":"initialize", ...}
Step 2: POST /mcp → {"method":"notifications/initialized"}
Step 3: POST /mcp → {"method":"tools/list", ...}  ← now works
```

Make sure your LLM client completes the handshake before calling tools.

### Health check returns 404

The health endpoint path may differ from the MCP path.

```bash
# Default health endpoint
curl http://127.0.0.1:8080/health

# If you changed the MCP path, health is still at /health
# (not affected by --http-path)
```

### Auth tool not available

The `auth` MCP tool is not showing up.

The `auth` tool is **disabled by default** (`--disable-llm-auth=true`). This is intentional for production security.

```bash
# Enable the auth tool
swag2mcp mcp --disable-llm-auth=false
```

## Authentication Problems

### 401 Unauthorized

The API rejected the request due to missing or invalid credentials.

```bash
# Check that auth is configured
swag2mcp info

# Validate the config
swag2mcp validate

# Check that environment variables are set
echo $MY_TOKEN

# Verify the token is not expired (bearer tokens are static)
```

**Common causes:**
- Token is missing or empty
- Environment variable not set
- Token has expired (bearer tokens don't auto-refresh)
- Wrong auth type configured

### 403 Forbidden

The API rejected the request due to insufficient permissions.

- The token may not have the required scopes
- The API key may not have access to this resource
- Check the API documentation for required permissions

### OAuth2 token endpoint unreachable

swag2mcp cannot reach the OAuth2 token URL.

```bash
# Check the token_url in your config
# Verify the URL is correct and reachable
curl -X POST https://auth.example.com/oauth/token \
  -d "grant_type=client_credentials" \
  -d "client_id=test" \
  -d "client_secret=test"

# Check network connectivity
# Check proxy settings if behind a corporate proxy
```

### Digest auth fails

swag2mcp cannot complete the Digest authentication handshake.

- The server must return a `WWW-Authenticate: Digest ...` header with a 401 response
- The challenge is cached for 5 minutes — if the server changes its nonce, wait for the cache to expire
- Check that username and password are correct

### HMAC signature mismatch

The API rejected the HMAC-signed request.

- Verify that `api_key` and `secret_key` are correct
- Check that the API uses Binance-style HMAC-SHA256 signing
- Some exchanges use different signing methods — HMAC auth is specifically for Binance-compatible APIs

### Script auth fails

The external auth script failed.

```bash
# Check that the script exists
ls -la ~/.swag2mcp/auth_scripts/my-domain.sh

# Run the script manually to test
sh ~/.swag2mcp/auth_scripts/my-domain.sh

# Check the script output format (must be JSON: {"token": "...", "expires_in": 3600})
# Check that the script completes within 30 seconds
# Check that the script has execute permissions
chmod +x ~/.swag2mcp/auth_scripts/my-domain.sh
```

## Search Problems

### No search results

The search returned no endpoints.

```bash
# Check that specs are loaded
swag2mcp ls

# Check that specs are not disabled
swag2mcp validate

# Try a simpler query
# Try searching by method: method:GET
# Try searching by tag: tag:pets

# The index is rebuilt on every MCP server start
# If you just added a spec, restart the server
```

### Search returns irrelevant results

The query is too broad or ambiguous.

- Use field filters to narrow: `method:GET +tag:pets`
- Use exact phrases: `"find pet by status"`
- Use the `limit` parameter to get more focused results

## API Call Problems

### invoke returns an error

The API call failed.

```bash
# Check the error message — it includes the HTTP status code
# 4xx errors: check parameters, auth, or permissions
# 5xx errors: the API server has a problem

# Always inspect the endpoint before invoking
inspect(endpointId: "...")

# Check that all required parameters are provided
# Check parameter types (string, number, boolean)
```

### Rate limit error

The LLM called the same endpoint too quickly.

Each endpoint has a 10-second cooldown. Wait before calling again, or disable the rate limiter:

```yaml
disable_ratelimiter: true
```

### Response too large (fileRef returned)

The response exceeded `max_response_size`.

This is normal. Use the response tools to explore the data:

```
1. response_outline(path) → understand the structure
2. response_compress(path, mode: "first_of_array") → get a sample
3. response_slice(path, jsonPath: "data.0") → get specific data
```

Or increase the limit:

```yaml
http_client:
  max_response_size: 4194304  # 4 MB
```

### Slow API responses

The API is taking too long to respond.

```yaml
http_client:
  timeout: 120s  # Increase from default 30s
```

## Workspace Problems

### swag2mcp init fails: "directory is not empty"

The target directory already has files.

```bash
# Use --force to overwrite
swag2mcp init --force

# Or use a different directory
swag2mcp init ./new-workspace
```

### swag2mcp update fails

One or more spec files could not be downloaded.

```bash
# Check the error message for which URL failed
# Verify the URL is accessible
curl -I <failed-url>

# Check network connectivity
# Check proxy settings
```

### Export creates no ZIP

The `[output]` argument must be a file path ending in `.zip`, not a directory.

```bash
# Correct
swag2mcp export /path/to/workspace /path/to/backup.zip

# Wrong (no ZIP will be created)
swag2mcp export /path/to/workspace /some/directory
```

### Import fails: "not a valid swag2mcp backup"

The ZIP file was not created by `swag2mcp export`.

Only ZIP archives created by `swag2mcp export` can be imported. The archive has a specific internal structure (`swag2mcp.yaml`, `specs/`, `auth_scripts/`).

## TUI Problems

### TUI doesn't render correctly

The terminal is too small or doesn't support the required features.

- Minimum terminal size: 80×24 characters
- The TUI uses Bubbletea and works in most modern terminals
- Try resizing your terminal window
- Try a different terminal emulator

### TUI shows "no specs found"

The workspace has no configured specs.

```bash
# Check specs
swag2mcp ls

# Add a spec
swag2mcp add spec
```

## Mock Server Problems

### Mock server doesn't start

```bash
# Check that mock_enabled: true in config
# Check that every collection has base_mock_url set
# Check that ports are not in use
lsof -i :9090

# Check the mock server logs
swag2mcp-mock mockserver
```

### Mock server returns empty responses

The spec file may not have response schemas defined.

- Mock server generates data from response schemas
- If no schema is found, it returns `{}`
- Check that your OpenAPI spec has `responses` with `schema` defined

## Network Problems

### Proxy connection failed

swag2mcp cannot connect through the configured proxy.

```bash
# Check proxy URL format (must include scheme: http://, https://, socks5://)
# Check proxy credentials
# Check bypass list — the target may be in the bypass list
# Test the proxy with curl
curl -x http://proxy.company.com:8080 https://api.example.com
```

### TLS/SSL errors

Certificate verification failed.

- If using a self-signed certificate for the MCP server, the client must trust it
- For the mock server with `--tls`, a self-signed certificate is generated automatically
- For API calls, swag2mcp uses the system's certificate store

## Other Problems

### High disk usage

The cache and responses directories can grow over time.

```bash
# Clean everything
swag2mcp clean

# Old responses (>48h) are cleaned automatically on MCP server start
# Cache files expire randomly between 1-48 hours
```

### "command not found" after go install

The `go install` directory is not in your PATH.

```bash
# Find where Go installs binaries
go env GOPATH
# Add to PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### LLM doesn't use the tools correctly

The LLM may need better instructions or a formatting skill.

- Use `llm_instruction` in your spec config to describe what the API does
- Consider using the [swag2mcp-format skill](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md) for consistent output formatting
- The quality of LLM responses depends on the model and the instructions it receives

### How do I report a bug?

Open an issue on [GitHub](https://github.com/mmadfox/swag2mcp/issues) with:
- swag2mcp version (`swag2mcp --version`)
- Your operating system and architecture
- The exact command you ran
- The full error message
- Your config file (with secrets removed)
