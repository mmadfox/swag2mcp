---
name: auth
---

# auth

Retrieves an authentication token (bearer/headers/query params) for a given spec/domain.

## When to use

Use this tool **only** when the user explicitly asks for the raw token or credentials, such as:
- "Generate a curl command for this endpoint"
- "Show me the auth token"
- "Give me the headers to use in Postman"
- "I need the token for an external script"

## When NOT to use

- **Do NOT** call `auth` before `inspect` — `inspect` is read-only and does not make API calls.
- **Do NOT** call `auth` before `invoke` — `invoke` automatically obtains and applies authentication under the hood. You do not need to pass the token manually.

## Parameters

- `specId` (required): The 32-character MD5 hash ID of the spec/domain to get an auth token for.

## Returns

The token string, any additional auth headers, and query parameters.
