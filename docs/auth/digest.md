# Digest Auth

## Purpose

HTTP Digest Access Authentication — a more secure alternative to Basic Auth. The password is not sent in plain text; instead, MD5 hashes are used.

## When to use

- Legacy APIs that only support Digest
- When you need authentication without sending the password in plain text
- Internal enterprise systems

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

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `username` | Yes | Username |
| `password` | Yes | Password |

## Notes

- swag2mcp first sends a request without authentication, receives a challenge from the server (HTTP 401), computes the response, and retries with the `Authorization: Digest ...` header
- The challenge is cached for 5 minutes — subsequent requests don't need an extra round-trip
- Store the password in an environment variable: `password: "$(API_PASSWORD)"`
