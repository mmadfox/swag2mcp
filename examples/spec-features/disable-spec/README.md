# Disable Spec

Shows how to temporarily disable a spec or collection without removing it
from the configuration. Disabled items are completely ignored by swag2mcp.

## What it demonstrates

- `disable: true` at spec level — the entire spec is skipped
- `disable: true` at collection level — only that collection is skipped
- Disabled items do not appear in any tool responses
- Useful for maintenance, deprecation, or A/B testing

## Expected behavior

- "active-api" appears in `spec_list` with all its collections
- "deprecated-api" does NOT appear anywhere
- "active-api" → "Maintenance Mode" collection is visible
- "active-api" → "Deprecated Endpoints" collection is NOT visible
