# Rate Limiting

swag2mcp has a built-in rate limiter to prevent duplicate calls.

## How It Works

- **Limit**: 10 seconds per endpoint
- **Reset**: after 10 seconds of inactivity
- **Block**: second call within 10 seconds is silently blocked

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

The rate limiter is built-in and not configurable.
