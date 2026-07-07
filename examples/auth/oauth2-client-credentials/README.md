# OAuth2 Client Credentials

Uses the OAuth2 Client Credentials grant (RFC 6749 section 4.4) for
machine-to-machine authentication. The client automatically fetches a bearer
token from the token endpoint and caches it until expiration.

## What it demonstrates

- `auth.type: oauth2-cc` configuration
- `client_id`, `client_secret`, `token_url` fields
- Optional `scopes` list
- Automatic token fetching and caching
- Token refresh on expiration

## Expected behavior

- First `invoke` call fetches a token from the token endpoint
- Subsequent calls reuse the cached token until it expires
- The `auth` tool returns the bearer token
- Token is automatically refreshed when expired
