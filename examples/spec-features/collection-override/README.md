# Collection Override

Demonstrates how a collection can override spec-level settings. Each collection
can have its own `base_url` and `headers`, which take precedence over the
spec-level values. This is useful when different parts of an API live on
different servers.

## What it demonstrates

- `base_url` at collection level — overrides the spec-level `base_url`
- `headers` at collection level — merged with spec-level headers
- `disable` at collection level — skip a specific collection
- Collection-level settings only affect that collection

## Expected behavior

- "Users" collection uses the spec-level `base_url`
- "Billing" collection uses its own `base_url` (different server)
- "Billing" collection adds an extra `X-Billing-Version` header
- "Legacy" collection is disabled and not indexed
