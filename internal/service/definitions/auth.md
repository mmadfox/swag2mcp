---
name: auth
---
Retrieves an authentication token (bearer) for a given spec/domain.

Use this tool when you need to obtain an auth token to make authenticated API calls on behalf of the user. The token is obtained from the configured auth provider (OAuth2, script, etc.) for the specified domain.

If auth is disabled, an empty token is returned.

Arguments:
- domainId: The 32-character MD5 hash ID of the spec/domain to get an auth token for (required).
