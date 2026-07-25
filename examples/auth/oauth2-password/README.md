# OAuth2 Password Grant

Uses the OAuth2 Resource Owner Password Credentials grant (RFC 6749 section 4.3)
for first-party applications. The client authenticates with username/password
and receives a bearer token.

## What it demonstrates

- `auth.type: oauth2-pwd` configuration
- `username`, `password`, `client_id`, `client_secret`, `token_url` fields
- Optional `scopes` list
- Automatic token fetching and caching
- Token refresh on expiration

## Expected behavior

- First `invoke` call fetches a token using password grant
- Subsequent calls reuse the cached token
- The `auth` tool returns the bearer token
- Token is automatically refreshed when expired
