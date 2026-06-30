---
name: search
---

# search

Searches endpoints across all specifications using full-text and structured queries.

This is the **primary discovery tool** — use it whenever you don't know the endpoint ID.

## When to use

Use this tool when:
- The user asks a general question about what endpoints exist
- The user describes functionality without knowing specific paths or tags
- You need to find relevant endpoints based on natural language descriptions
- You want to filter by HTTP method (`method:GET`), tag (`tag:auth`), or path (`path:user`)

## Parameters

- `query` (required): Natural language or structured search query. Supports field filters (`method:POST`, `tag:pet`, `path:/api/v1/*`), boolean operators (`+` must, `-` exclude), fuzzy (`term~`), wildcards (`*`, `?`), and phrases (`"exact phrase"`)
- `limit` (required): Maximum number of results to return (min: 1, max: 50)

## Returns

A list of endpoints matching the query with their IDs, methods, paths, and summaries.
