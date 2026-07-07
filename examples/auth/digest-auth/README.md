# HTTP Digest Authentication

Uses HTTP Digest Access Authentication as defined in RFC 2617. The client
automatically handles the challenge-response flow: it sends an unauthenticated
request, receives a 401 with a `WWW-Authenticate: Digest` challenge, computes
the response hash, and retries with the proper `Authorization: Digest` header.

## What it demonstrates

- `auth.type: digest` configuration
- `username` and `password` fields
- Automatic challenge-response handshake
- Nonce caching within TTL (5 minutes)
- Nonce count (`nc`) incrementation on each request

## Expected behavior

- First `invoke` call makes 2 HTTP requests (401 challenge + authenticated)
- Subsequent calls reuse the cached nonce (1 request)
- The `auth` tool returns the `Authorization: Digest ...` header
- After 5 minutes, a new challenge is fetched automatically
