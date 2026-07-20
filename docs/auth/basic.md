# Basic Auth

HTTP Basic Authentication.

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
        password: "{{PASSWORD}}"
```

## How It Works

Header `Authorization: Basic base64(username:password)` is added to every request.

## Environment Variables

```yaml
auth:
  type: basic
  config:
    username: "admin"
    password: "$(API_PASSWORD)"
```

## Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `username` | string | Username |
| `password` | string | Password |
