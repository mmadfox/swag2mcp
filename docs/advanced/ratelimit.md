# Rate Limiting

swag2mcp has a built-in rate limiter to prevent duplicate calls.

## How It Works

- **Default limit**: 10 seconds per endpoint
- **Reset**: after the interval of inactivity
- **Block**: second call within the interval is silently blocked

## Behavior

```
t=0s  → invoke(endpoint) → executes
t=2s  → invoke(endpoint) → blocked (silent, no error)
t=12s → invoke(endpoint) → executes
```

## Why

- Prevents accidental duplicate calls
- Protects against API rate limits
- Saves resources

## Configuration

```yaml
# Disable the rate limiter entirely
disable_ratelimiter: true

# Custom rate limit interval (Go duration format)
rate_limit_interval: 30s
```

- `disable_ratelimiter` — set to `true` to disable the per-endpoint rate limiter. Useful when testing or when you need to call the same endpoint repeatedly.
- `rate_limit_interval` — custom interval between calls to the same endpoint. Default: `10s`. Range: any valid Go duration (e.g., `5s`, `30s`, `1m`).
