# Bearer Auth

Bearer Token authentication (JWT, OAuth2 tokens).

## Configuration

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: bearer
      config:
        token: "eyJhbGciOiJIUzI1NiIs..."
```

## How It Works

Header `Authorization: Bearer <token>` is added to every request.

## Environment Variables

```yaml
auth:
  type: bearer
  config:
    token: "$(API_TOKEN)"
```

## Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `token` | string | Bearer token value |
