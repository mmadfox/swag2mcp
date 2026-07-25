---
name: search
---

# search

**The ONLY tool for finding endpoints when you don't have the endpoint ID.**

Searches endpoints across all specifications using full-text and structured queries.

## When to use

Use this tool when:
- The user asks to find a method/endpoint by description, name, path, tag, or functionality
- The user describes functionality without knowing specific paths or tags
- You need to find relevant endpoints based on natural language descriptions
- You want to filter by HTTP method (`method:GET`), tag (`tag:auth`), or path (`path:user`)

## DO NOT

- ❌ **Do NOT manually traverse** spec → collection → tag → endpoint to find something. Use `search`.
- ❌ **Do NOT guess** endpoint IDs. Use `search`.
- ❌ **Do NOT use** `endpoint_by_tag` / `endpoint_by_collection` / `endpoint_by_spec` for discovery — those are for navigation after you already know what you're looking for.
- ❌ **Do NOT skip** `search` and try to brute-force your way through collections. One `search` call replaces dozens of manual navigation steps.

## User Intent → Search Query Examples

| User says | What to search |
|-----------|---------------|
| "Find the create user endpoint" | `+method:POST +summary:create +summary:user` |
| "Show all GET endpoints" | `method:GET` |
| "What relates to orders?" | `order` |
| "Find endpoint by path /api/v1/users" | `path:"/api/v1/users"` |
| "How do I delete a pet?" | `+method:DELETE +summary:pet` |
| "Show all auth endpoints" | `tag:auth` |
| "Find something about inventory" | `inventory` |
| "Give me all POST requests in the store section" | `+method:POST +tag:store` |

## Parameters

- `query` (required): Natural language or structured search query. Supports field filters (`method:POST`, `tag:pet`, `path:/api/v1/*`), boolean operators (`+` must, `-` exclude), fuzzy (`term~`), wildcards (`*`, `?`), and phrases (`"exact phrase"`)
- `limit` (required): Maximum number of results to return (min: 1, max: 50)

## Returns

A list of endpoints matching the query with their IDs, methods, paths, and summaries.
