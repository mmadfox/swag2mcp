# Mock Server

This example demonstrates how to use the swag2mcp mock server to generate
random API responses without connecting to a real backend.

## What it demonstrates

- `mock_enabled: true` — enables mock server mode
- `base_mock_url` — per-collection mock server address
- Auth mock servers — simulate all 8 authentication methods
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

## Mock Auth Credentials

When `auth` is configured in a spec, the mock server starts an auth mock
on a random port. Each auth type accepts the following credentials:

| Auth Type | Endpoint | Accepts | Example Request |
|-----------|----------|---------|-----------------|
| `basic` | `GET /` | Any `user:password` in Base64 | `Authorization: Basic YWRtaW46cGFzcw==` |
| `bearer` | `GET /` | Any non-empty token | `Authorization: Bearer any-token` |
| `digest` | `GET /` | Any Digest response | `Authorization: Digest username="test", realm="...", nonce="...", uri="/", response="..."` |
| `oauth2-cc` | `POST /token` | Any `client_id` + `client_secret` | `grant_type=client_credentials&client_id=any&client_secret=any` |
| `oauth2-pwd` | `POST /token` | Any `username` + `password` | `grant_type=password&username=any&password=any` |
| `api-key` | `GET /` | Any `X-Api-Key` header or `api_key` query | `X-Api-Key: any-key` |
| `script` | `GET /token` | No credentials required | `GET /token` |

All auth mocks return `{"status":"authenticated","method":"<type>"}`.
OAuth2 mocks return `{"access_token":"<random>","token_type":"Bearer","expires_in":3600}`.

## Expected behavior

- Mock server starts on `localhost:8080`
- Auth mock starts on a random port
- All API calls go to the mock server instead of the real API
- Responses contain random data matching the OpenAPI schemas
