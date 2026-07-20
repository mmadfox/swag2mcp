# HTTP Client

swag2mcp uses a configurable HTTP client for API calls.

## Configuration

```yaml
http_client:
  timeout: 30s
  max_response_size: 1048576
  user_agent: "swag2mcp/1.0"
  follow_redirects: true
  max_redirects: 10
  proxy:
    url: ""
    username: ""
    password: ""
    bypass: []
  headers:
    "User-Agent": "swag2mcp/1.0"
  cookies:
    - name: "session"
      value: "abc123"
      domain: ".example.com"
      path: "/"
```

## Parameters

| Parameter | Description |
|-----------|-------------|
| `timeout` | HTTP request timeout (1s-5m) |
| `max_response_size` | Max response size in bytes (256-10485760) |
| `user_agent` | User-Agent header value |
| `follow_redirects` | Follow HTTP redirects |
| `max_redirects` | Max redirects to follow (0-50) |
| `proxy.url` | HTTP proxy URL (http, https, socks5, socks5h) |
| `proxy.username` | Proxy username |
| `proxy.password` | Proxy password |
| `proxy.bypass` | Domains to bypass proxy |
| `headers` | Headers added to every request |
| `cookies` | Cookies added to every request |

## Randomizer

Add random browser-like headers:

```yaml
http_client:
  random: true
```

Adds random `User-Agent`, `Accept`, `Accept-Language` headers.

## Proxy

```yaml
http_client:
  proxy:
    url: "http://user:pass@proxy.example.com:8080"
    bypass:
      - "localhost"
      - "*.internal.com"
```

## Custom Headers at Spec Level

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    http_client:
      headers:
        "Accept": "application/json"
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Cookies at Spec Level

```yaml
specs:
  - domain: example
    llm_title: Example API
    base_url: https://api.example.com
    http_client:
      cookies:
        - name: "session"
          value: "abc123"
        - name: "csrf"
          value: "{{CSRF_TOKEN}}"
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Cascade

HTTP client settings cascade:

1. Global (`http_client`)
2. Spec (`specs[].http_client`)
3. Collection (`specs[].collections[].http_client`)

Each level overrides the previous.
