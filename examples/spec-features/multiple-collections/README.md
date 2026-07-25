# Multiple Collections

Demonstrates a spec with multiple collections. Each collection points to a
different OpenAPI/Swagger specification file. This is useful when an API is
split across multiple spec files.

## What it demonstrates

- Multiple `collections` per spec (up to 30)
- Each collection can have its own `llm_title` and `llm_instruction`
- Collections are independently indexed and searchable
- The LLM can discover endpoints across all collections

## Expected behavior

- All 4 collections are indexed and searchable
- `collection_by_spec` returns all 4 collections
- Endpoints from all collections appear in search results
- Each collection is independently discoverable
