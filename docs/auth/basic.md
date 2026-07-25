# Basic Auth

## Purpose

HTTP Basic Authentication — the simplest way to authenticate with a username and password.

## When to use

- Legacy APIs that only support Basic Auth
- Simple authentication without complex tokens
- Internal services

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
      type: basic
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

- The password is sent in the `Authorization: Basic ...` header encoded in Base64 — this is **not encryption**. Always use HTTPS.
- Store the password in an environment variable: `password: "$(MY_PASSWORD)"`
