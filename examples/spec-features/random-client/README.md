# Random Client Configuration

Demonstrates how to enable browser-like random headers for all HTTP requests,
including OAuth2 token exchange and Digest challenge requests.

## What it demonstrates

- `random: true` — enables random browser-like headers at startup
- `user_agent` — custom User-Agent (if not set, a random one is generated)
- Per-spec headers and per-collection cookies still work alongside random mode

## Random headers generated

| Header | Example value |
|--------|---------------|
| `User-Agent` | `Mozilla/5.0 (Windows NT 10.0; Win64; x64) ...` |
| `Accept` | `application/json, text/plain, */*` |
| `Accept-Language` | `ru-RU,ru;q=0.9,en;q=0.8` (from system locale) |
| `Accept-Encoding` | `gzip, deflate, br` |
| `Referer` | `https://www.google.com/` |
| `Sec-Ch-Ua` | `"Chromium";v="125", "Google Chrome";v="125"` |
| `Sec-Ch-Ua-Platform` | `"Windows"` |
| `Sec-Fetch-Site` | `same-origin` |
| `Sec-Fetch-Mode` | `cors` |
| `Sec-Fetch-Dest` | `document` |

## Expected behavior

- All API requests include random browser-like headers
- User-configured headers (`X-Store-ID`) are preserved and not overwritten
- User-configured cookies are applied alongside random headers
- OAuth2 token requests also include the same random headers
- Values are generated once at startup and stay consistent for the session
