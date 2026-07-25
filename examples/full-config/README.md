# Full Configuration

This example demonstrates every feature swag2mcp has to offer in a single
configuration file. Use it as a reference for building your own configs.

## What it demonstrates

- Multiple specs with different auth types
- All 8 authentication methods
- Spec-level and collection-level `llm_title` / `llm_instruction`
- `disable` flag on specs and collections
- `tags` for filtering with `--tags`
- Custom headers at spec and collection level
- Collection-level `base_url` override
- Environment variable resolution with `$(VAR_NAME)` syntax
- Multiple collections per spec
- Remote (HTTP) and local (file) collection locations
- Mock server mode with `mock_enabled: true` and `base_mock_url`

## Expected behavior

- 4 specs are registered (2 enabled, 2 disabled)
- Each spec uses a different auth method
- Tags allow filtering: `swag2mcp mcp --tags=production`
- Disabled specs/collections are ignored
- Environment variables are resolved at startup
- Mock server can be started with `swag2mcp-mock`
