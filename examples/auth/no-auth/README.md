# No Auth

This example shows a spec with no authentication. The `auth` block is either
omitted entirely or explicitly set to `type: none`.

## What it demonstrates

- Omitting the `auth` block means no authentication
- Explicit `type: none` also works
- All API calls are sent without any auth headers or query parameters

## Expected behavior

- The `auth` tool returns empty headers and query params
- API calls via `invoke` are sent without any auth
