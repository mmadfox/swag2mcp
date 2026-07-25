# Minimal Configuration

This is the simplest possible swag2mcp configuration. It defines a single API
specification with one collection and no authentication.

## What it demonstrates

- Minimal required fields: `domain`, `base_url`, `llm_title`, and at least one
  `collection` with a `location`
- No authentication (auth block is omitted entirely)
- Single spec, single collection

## Expected behavior

- The server starts and registers all 16 MCP tools
- The spec "weather-api" appears in `spec_list` and `spec_by_id`
- The collection is discoverable via `collection_by_spec`
- All endpoints from the OpenAPI spec are indexed and searchable
- No authentication is applied to any API call
