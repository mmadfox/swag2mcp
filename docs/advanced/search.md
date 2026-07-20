# Full-Text Search

swag2mcp uses bluge — a full-text search engine for Go.

## Indexing

When a spec is added, all endpoints are indexed:

- **Method**: GET, POST, PUT, DELETE, PATCH
- **Path**: `/api/v1/users/{id}`
- **Summary**: OpenAPI summary
- **Tags**: endpoint categories
- **ID**: unique identifier

## Query Syntax

| Example | Description |
|---------|-------------|
| `pet` | Search all fields |
| `+method:POST +tag:pet` | POST in pet tag |
| `path:"/api/v1/users"` | Exact path match |
| `create~` | Fuzzy search |
| `cr*` | Wildcard search |
| `"find pet"` | Phrase search |
| `+summary:pet -method:DELETE` | Exclude DELETE |

## Searchable Fields

| Field | Description |
|-------|-------------|
| `method` | HTTP method |
| `path` | Endpoint path |
| `summary` | Endpoint description |
| `tag` | Tag |
| `id` | Endpoint ID |

## Examples

```
# Find all GET requests
method:GET

# Find POST requests in pet tag
+method:POST +tag:pet

# Find endpoints by path
path:"/pet/findByStatus"

# Find by description
"find pet by status"

# Find everything except DELETE
+summary:pet -method:DELETE
```

## MCP Tool

```
→ search(query: "find pet by status", limit: 5)
← GET /pet/findByStatus — Finds Pets by status
   GET /pet/{petId} — Find pet by ID
```
