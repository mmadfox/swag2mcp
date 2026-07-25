# Custom Headers

Shows how to add custom HTTP headers to every API call made by a spec.
Headers can be defined at the spec level (applied to all collections) and
at the collection level (applied only to that collection).

## What it demonstrates

- `headers` at spec level — applied to all collections in the spec
- `headers` at collection level — applied only to that collection
- Headers are merged: collection-level headers override spec-level ones
- Useful for API keys, tracing headers, correlation IDs, etc.

## Expected behavior

- Every `invoke` call includes spec-level headers
- Collection-specific calls also include collection-level headers
- Collection-level headers override spec-level ones with the same key
