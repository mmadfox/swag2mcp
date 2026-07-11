# Mock Server

This example demonstrates how to use the swag2mcp mock server to generate
random API responses without connecting to a real backend.

## What it demonstrates

- `mock_enabled: true` — enables mock server mode
- `base_mock_url` — per-collection mock server address
- Auth mock servers — two global servers (OAuth2 on port 9090, Digest on port 9091); other auth types are handled automatically by the MCP server
- Random data generation — responses match OpenAPI schemas

## Configuration

```yaml
mock_enabled: true

specs:
  - domain: petstore
    base_url: https://petstore.swagger.io/v2
    llm_title: Petstore API
    auth:
      type: bearer
      config:
        token: any-token
    collections:
      - llm_title: Petstore Swagger
        location: specs/petstore.json
        base_mock_url: localhost:8080
```

## Usage

```bash
# 1. Start mock server (terminal 1)
swag2mcp-mock

# 2. Start MCP server (terminal 2)
swag2mcp mcp

# 3. Invoke endpoints — they will hit the mock server
```

## Mock Auth

When `auth` is configured in a spec, the MCP server applies authentication
automatically. Only two auth types need a dedicated mock server:

| Auth Type | Mock Endpoint | Behavior |
|-----------|---------------|----------|
| `oauth2-cc` / `oauth2-pwd` | `POST /token` on port 9090 | Accepts any `client_id`/`username`+`password`, returns `{"access_token":"<random>","token_type":"Bearer","expires_in":3600}` |
| `digest` | `GET /` on port 9091 | Sends a 401 challenge with `algorithm=MD5`, accepts any Digest response, returns `{"status":"authenticated","method":"digest"}` |

Other auth types (`basic`, `bearer`, `api-key`, `script`) do **not** require
a mock server — the MCP server applies the configured credentials to every
request automatically.

## Expected behavior

- Mock server starts on `localhost:8080`
- Auth mock servers start on ports 9090 (OAuth2) and 9091 (Digest)
- Other auth types (basic, bearer, api-key, script) are handled automatically by the MCP server
- All API calls go to the mock server instead of the real API
- Responses contain random data matching the OpenAPI schemas
