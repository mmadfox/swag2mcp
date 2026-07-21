# Digest Auth

HTTP Digest Access Authentication.

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
      type: digest
      config:
        username: "admin"
        password: "$(PASSWORD)"
```

## How It Works

1. swag2mcp sends a request without auth
2. Server responds 401 with `WWW-Authenticate: Digest ...`
3. swag2mcp computes MD5 hashes and retries with `Authorization: Digest ...`

## Environment Variables

```yaml
auth:
  type: digest
  config:
    username: "admin"
    password: "$(API_PASSWORD)"
```

## Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `username` | string | Username |
| `password` | string | Password |
