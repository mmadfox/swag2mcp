# Rate Limiting

## Overview

swag2mcp has a built-in rate limiter that prevents the LLM from calling the same API endpoint too frequently. This protects against accidental duplicate calls and respects API rate limits.

## How it works

Each endpoint has a cooldown period. If the LLM tries to call the same endpoint again within the cooldown, the call is rejected with a structured error.

```
t=0s  → invoke(endpoint) → executes
t=2s  → invoke(endpoint) → rejected with rate_limit error
t=12s → invoke(endpoint) → executes (cooldown has passed)
```

### Default behavior

- **Cooldown:** 10 seconds per endpoint
- **Scope:** Per-endpoint — calling endpoint A does not affect endpoint B
- **Error response:** The LLM receives an `LLMError` with code `rate_limit` and a message indicating how long to wait
- **Reset:** After 10 seconds of inactivity on that endpoint

### Error format

When rate limited, the LLM receives:

```json
{
  "code": "rate_limit",
  "message": "rate limit exceeded for endpoint \"abc123\": try again in 8 seconds",
  "hint": "Wait for the cooldown period to expire, then try invoking the endpoint again. Use the search tool to find other endpoints you can call in the meantime."
}
```

The LLM can use this information to wait and retry, or switch to a different endpoint.

### Why it exists

- **Prevents accidental duplicate calls** — the LLM might call the same endpoint multiple times in quick succession
- **Protects against API rate limits** — many APIs have their own rate limits, and hitting them would cause errors
- **Saves resources** — reduces unnecessary network traffic

## Configuration

You can disable the rate limiter or change the cooldown interval:

```yaml
# Disable the rate limiter entirely
disable_ratelimiter: true

# Custom cooldown interval
rate_limit_interval: 30s
```

### disable_ratelimiter

- **Type:** `bool`
- **Default:** `false`
- **Effect:** When `true`, the per-endpoint rate limiter is disabled. The LLM can call the same endpoint repeatedly without waiting.
- **When to enable:** Testing, debugging, or when you need to call the same endpoint multiple times in quick succession.
- **When to keep disabled (recommended):** Production. The rate limiter prevents accidental abuse.

### rate_limit_interval

- **Type:** duration (Go format: `10s`, `30s`, `1m`)
- **Default:** `10s`
- **Effect:** Sets the cooldown period between calls to the same endpoint.
- **When to increase:** APIs with strict rate limits (e.g., 10 requests per minute).
- **When to decrease:** Internal APIs where you control the load.
- **Examples:** `5s`, `30s`, `1m`, `2m`

## Important notes

- **Per-endpoint tracking** — each endpoint is tracked independently. Calling one endpoint does not affect others.
- **Error returned to LLM** — the second call within the cooldown is rejected with a `rate_limit` error. The LLM receives the cooldown duration and can retry after waiting.
- **No cleanup needed** — the rate limiter tracks endpoints automatically and does not require maintenance.
