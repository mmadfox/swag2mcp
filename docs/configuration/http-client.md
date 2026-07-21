# HTTP Client

swag2mcp uses a configurable HTTP client for all API calls. These settings are defined globally and apply to every request unless overridden at the spec or collection level (headers and cookies only).

## Configuration

```yaml
http_client:
  timeout: 30s
  max_response_size: 1048576
  user_agent: "swag2mcp/1.0"
  follow_redirects: true
  max_redirects: 10
  random: false
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

## Timeout

Controls how long swag2mcp waits for an API response before giving up.

- **Default**: 30 seconds
- **Range**: 1 second to 5 minutes
- **When to increase**: Slow APIs, large payloads, unreliable networks
- **When to decrease**: Internal APIs, health checks, fast-fail scenarios

```yaml
http_client:
  timeout: 60s
```

## Max Response Size

Limits how large a response can be before swag2mcp saves it to disk instead of returning it inline.

- **Default**: 1,048,576 bytes (1 MB)
- **Range**: 256 to 10,485,760 bytes (10 MB)
- **What happens when exceeded**: The response is saved to `~/.swag2mcp/responses/` as a JSON file. The LLM receives a file reference and can explore it with `response_outline`, `response_compress`, and `response_slice` tools.

```yaml
http_client:
  max_response_size: 4194304  # 4 MB
```

## User-Agent

The `User-Agent` header sent with every request. Some APIs require a specific user-agent or block known bot user-agents.

- **Default**: `swag2mcp-global/1.0`
- **When to change**: API requires a specific user-agent, or you want to identify your application

```yaml
http_client:
  user_agent: "MyApp/1.0"
```

## Follow Redirects

Controls whether swag2mcp automatically follows HTTP redirects (3xx status codes).

- **Default**: `true`
- **When to disable**: APIs that redirect in a loop, security-sensitive endpoints where you want to inspect redirect targets manually

```yaml
http_client:
  follow_redirects: false
```

## Max Redirects

Limits how many redirects swag2mcp follows before stopping.

- **Default**: 10
- **Range**: 0 to 50
- **When to change**: APIs with long redirect chains, or reduce for faster failure on redirect loops

```yaml
http_client:
  max_redirects: 5
```

## Randomizer

Adds random browser-like headers to each request to avoid fingerprinting and blocking.

- **Default**: `false`
- **What it adds**: Random `User-Agent` (from a pool of real browser strings), `Accept`, `Accept-Language`, `Accept-Encoding`, `Cache-Control`
- **When to enable**: APIs that block requests based on User-Agent or header patterns, scraping scenarios

```yaml
http_client:
  random: true
```

## Proxy

A proxy server acts as an intermediary between swag2mcp and the target API. All HTTP traffic is routed through it.

**When you might need a proxy:**
- Corporate network — all outbound traffic must go through a company proxy
- Geographic restrictions — some APIs are region-locked, a proxy in the right region bypasses this
- Static IP — APIs that require IP allowlisting
- Anonymity — hide the origin IP from the target API

### Supported Proxy Schemes

| Scheme | Description | Use Case |
|--------|-------------|----------|
| `http` | HTTP proxy for HTTP traffic | Corporate proxies, basic proxying |
| `https` | HTTPS proxy (CONNECT tunnel) | Secure corporate proxies |
| `socks5` | SOCKS5 proxy (DNS resolved locally) | General purpose, any protocol |
| `socks5h` | SOCKS5 proxy (DNS resolved on proxy) | When proxy has better DNS resolution |

### Authentication

If the proxy requires authentication, provide `username` and `password`:

```yaml
http_client:
  proxy:
    url: "http://proxy.example.com:8080"
    username: "proxyuser"
    password: "$(PROXY_PASSWORD)"
```

### Bypass

A list of domains that should **not** go through the proxy. Useful for internal services, localhost, or APIs that are only accessible directly.

```yaml
http_client:
  proxy:
    url: "http://proxy.example.com:8080"
    bypass:
      - "localhost"
      - "127.0.0.1"
      - "*.internal.company.com"
      - "api.local"
```

Bypass supports wildcard patterns (`*.example.com` matches any subdomain).

## Headers

Headers added to every request. They are merged with spec-level and collection-level headers:

```
Global headers → Spec headers (merged) → Collection headers (merged)
```

Collection headers override spec headers, which override global headers for the same key.

```yaml
http_client:
  headers:
    "Accept": "application/json"
    "Accept-Language": "en-US"
```

## Cookies

Cookies sent with every request. They are merged with spec-level cookies.

- **Global → Spec**: cookies are merged (spec overrides global for the same name)
- **Collection**: cookies are **not** supported at the collection level

```yaml
http_client:
  cookies:
    - name: "session"
      value: "abc123"
      domain: ".example.com"
      path: "/"
      secure: false
      http_only: false
```

### Cookie Fields

| Field | Required | Description |
|-------|----------|-------------|
| `name` | ✅ | Cookie name |
| `value` | ✅ | Cookie value |
| `domain` | ❌ | Domain scope (e.g., `.example.com`) |
| `path` | ❌ | Path scope (e.g., `/`) |
| `secure` | ❌ | Only send over HTTPS |
| `http_only` | ❌ | Not accessible via JavaScript |

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
          value: "$(CSRF_TOKEN)"
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Cascade

HTTP client settings cascade from global to spec to collection:

```
Global (http_client)
    ↓ overrides (headers, cookies only)
Spec (specs[].http_client)
    ↓ overrides (headers only)
Collection (specs[].collections[].http_client)
```

**Only `headers` and `cookies` can be overridden at the spec and collection levels.** All other settings (timeout, proxy, user-agent, redirects, response size, randomizer) are global only.

See [Configuration Cascade](./cascade) for details.
