# Full-Text Search

## Overview

swag2mcp includes a built-in full-text search engine (bluge) that indexes all endpoints across all specs. The LLM can search for endpoints by method, path, summary, or tag — even without knowing the endpoint ID.

## How indexing works

When a spec is added or updated, every endpoint is indexed. The following fields are searchable:

| Field | Description | Example |
|-------|-------------|---------|
| `method` | HTTP method | `GET`, `POST`, `PUT` |
| `path` | API endpoint path | `/api/v1/users/{id}` |
| `summary` | OpenAPI summary | "Find pet by ID" |
| `tag` | Endpoint category | "pets", "users" |
| `_all` | All fields combined | method + path + tag + summary |

The index is rebuilt on every MCP server start. It is stored in memory for fast searches.

## Query syntax

The search supports a rich query syntax for precise filtering:

| Example | Description |
|---------|-------------|
| `pet` | Simple text search across all fields |
| `method:GET` | Find all GET endpoints |
| `tag:pets` | Find endpoints in the "pets" tag |
| `path:"/api/v1/users"` | Exact path match |
| `+method:POST +tag:pet` | Must match both conditions |
| `-method:DELETE` | Exclude DELETE methods |
| `create~` | Fuzzy search (typo-tolerant) |
| `cr*` | Wildcard search |
| `"find pet"` | Phrase search |
| `+summary:pet -method:DELETE` | Include "pet" in summary, exclude DELETE |

### Field-specific search

You can search within specific fields using the `field:value` syntax:

```
method:GET
tag:pets
path:"/pet/findByStatus"
summary:"find pet by status"
```

### Boolean operators

- `+` — the term must match (AND)
- `-` — the term must not match (NOT)
- Space between terms — OR (any term can match)

### Fuzzy and wildcard

- `term~` — fuzzy search (matches similar words, handles typos)
- `te*` — wildcard (matches any characters)
- `te?t` — single character wildcard

## Examples

```
# Find all GET requests
method:GET

# Find POST requests in the pet tag
+method:POST +tag:pet

# Find endpoints by exact path
path:"/pet/findByStatus"

# Find by description
"find pet by status"

# Find everything except DELETE
+summary:pet -method:DELETE

# Fuzzy search for "create" (handles typos)
create~
```

## MCP tool

The `search` MCP tool exposes the search engine to the LLM:

```
→ search(query: "find pet by status", limit: 5)
← GET /pet/findByStatus — Finds Pets by status
   GET /pet/{petId} — Find pet by ID
```

### Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `query` | Yes | Search query (supports structured syntax) |
| `limit` | Yes | Maximum results (1-50) |

## Important notes

- **The index is in-memory** — it is rebuilt every time the MCP server starts. There is no persistent index file.
- **All fields are lowercased** — searches are case-insensitive
- **Limit is capped at 50** — you cannot request more than 50 results
- **Invalid query syntax** returns a helpful error message with examples
- **The `_all` field** combines method, path, tag, and summary for simple text searches
