# HTTP Client

swag2mcp uses a configurable HTTP client for all API calls. These settings are defined globally and can be overridden at the spec and collection levels.

## Configuration

```yaml
http_client:
  timeout: 30s
  max_response_size: 1048576
  user_agent: "swag2mcp-global/1.0"
  follow_redirects: true
  max_redirects: 10
  random: false
  proxy:
    url: ""
    username: ""
    password: ""
    bypass: []
  headers:
    "Accept": "application/json"
  cookies:
    - name: "session"
      value: "abc123"
      domain: ".example.com"
      path: "/"
```

## Timeout

Controls how long swag2mcp waits for an API response before giving up.

- **Type:** duration (Go format: `30s`, `60s`, `2m`)
- **Default:** `30s`
- **Range:** 1 second to 5 minutes
- **Effect:** If the API does not respond within this time, the request fails with a timeout error.
- **When to increase:** Slow APIs, large payloads, unreliable networks.
- **When to decrease:** Internal APIs, health checks, fast-fail scenarios.

```yaml
http_client:
  timeout: 60s
```

## Max Response Size

Limits how large a response can be before swag2mcp saves it to disk instead of returning it inline to the LLM.

- **Type:** `int` (bytes)
- **Default:** `1048576` (1 MB)
- **Range:** 256 to 10,485,760 bytes (10 MB)
- **Effect:** When a response exceeds this limit, it is saved to `{workspace}/responses/` as a JSON file. The LLM receives a file reference and can explore it with `response_outline`, `response_compress`, and `response_slice` tools.
- **When to increase:** APIs that return large datasets (reports, logs, analytics).
- **When to decrease:** Limited LLM context window, or when you prefer file-based access for all responses.

```yaml
http_client:
  max_response_size: 4194304  # 4 MB
```

## User-Agent

The `User-Agent` header sent with every request. Some APIs require a specific user-agent or block known bot user-agents.

- **Type:** `string`
- **Default:** `"swag2mcp-global/1.0"`
- **Effect:** Identifies your application to the API server.
- **When to change:** The API requires a specific user-agent, or you want to identify your application for analytics.

```yaml
http_client:
  user_agent: "MyApp/1.0"
```

## Follow Redirects

Controls whether swag2mcp automatically follows HTTP redirects (3xx status codes).

- **Type:** `bool`
- **Default:** `true`
- **Effect:** When `true`, swag2mcp follows redirects up to `max_redirects` times. When `false`, the redirect response is returned as-is.
- **When to disable:** APIs that redirect in a loop, security-sensitive endpoints where you want to inspect redirect targets manually.

```yaml
http_client:
  follow_redirects: false
```

## Max Redirects

Limits how many redirects swag2mcp follows before stopping.

- **Type:** `int`
- **Default:** `10`
- **Range:** 0 to 50
- **Effect:** If the API redirects more times than this limit, the request fails.
- **When to change:** APIs with long redirect chains, or reduce for faster failure on redirect loops.

```yaml
http_client:
  max_redirects: 5
```

## Randomizer

Adds random browser-like headers to each request to avoid fingerprinting and blocking.

- **Type:** `bool`
- **Default:** `false`
- **Effect:** When `true`, swag2mcp generates random headers for each request: `User-Agent` (from a pool of real browser strings), `Accept`, `Accept-Language`, `Accept-Encoding`, `Cache-Control`. This overrides the `user_agent` setting.
- **When to enable:** APIs that block requests based on User-Agent or header patterns, scraping scenarios.

```yaml
http_client:
  random: true
```

## Proxy

A proxy server acts as an intermediary between swag2mcp and the target API. All HTTP traffic is routed through it.

**When you might need a proxy:**
- **Corporate network** — all outbound traffic must go through a company proxy
- **Geographic restrictions** — some APIs are region-locked, a proxy in the right region bypasses this
- **Static IP** — APIs that require IP allowlisting
- **Anonymity** — hide the origin IP from the target API

### Proxy URL

- **Type:** `string`
- **Default:** `""` (no proxy)
- **Supported schemes:** `http`, `https`, `socks5`, `socks5h`
- **Supports `$(VAR)`:** ✅ resolved at runtime

| Scheme | Description | Use Case |
|--------|-------------|----------|
| `http` | HTTP proxy for HTTP traffic | Corporate proxies, basic proxying |
| `https` | HTTPS proxy (CONNECT tunnel) | Secure corporate proxies |
| `socks5` | SOCKS5 proxy (DNS resolved locally) | General purpose, any protocol |
| `socks5h` | SOCKS5 proxy (DNS resolved on proxy) | When proxy has better DNS resolution |

### Proxy Authentication

If the proxy requires authentication, provide `username` and `password`:

- **Supports `$(VAR)`:** ✅ resolved at runtime for all three fields (`url`, `username`, `password`)

```yaml
http_client:
  proxy:
    url: "http://proxy.example.com:8080"
    username: "proxyuser"
    password: "$(PROXY_PASSWORD)"
```

### Proxy Bypass

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

Custom HTTP headers added to every request. Headers are merged across cascade levels:

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

Header values support `$(ENV_VAR)` resolution.

## Cookies

Cookies sent with every request. Cookies are merged across cascade levels (lower level overrides global for the same cookie name).

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
| `name` | Yes | Cookie name |
| `value` | Yes | Cookie value (supports `$(ENV_VAR)` resolution) |
| `domain` | No | Domain scope (e.g., `.example.com`) |
| `path` | No | Path scope (e.g., `/`) |
| `secure` | No | Only send over HTTPS |
| `http_only` | No | Not accessible via JavaScript |

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

HTTP client settings cascade from global to spec to collection. All settings can be overridden at every level:

```
Global (http_client)
    ↓ overrides (all settings)
Spec (specs[].http_client)
    ↓ overrides (all settings)
Collection (specs[].collections[].http_client)
```

**All HTTP client settings** (timeout, proxy, user-agent, redirects, response size, randomizer, headers, cookies) can be overridden at both spec and collection levels.

See [Configuration Cascade](./cascade) for details.
