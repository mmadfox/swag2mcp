# Endpoints

An endpoint is a specific HTTP method + path that can be invoked.

## Structure

Each endpoint contains:

- **HTTP method**: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
- **Path**: `/api/v1/users/{id}`
- **Parameters**: path, query, header, cookie
- **Request body**: for POST/PUT/PATCH
- **Responses**: status codes and response schemas

## Identification

Each endpoint gets a unique MD5 hash:

```go
id = md5(method + ":" + path)
```

## MCP Tools for Endpoints

| Tool | Description |
|------|-------------|
| `endpoint_by_spec` | All endpoints in a spec |
| `endpoint_by_collection` | Endpoints in a collection |
| `endpoint_by_tag` | Endpoints in a tag |
| `endpoint_by_id` | Quick endpoint summary |
| `inspect` | Full endpoint details (schemas, params) |
| `invoke` | Call the endpoint |
| `search` | Search endpoints |

## Deprecated Endpoints

Endpoints marked as `deprecated` in the spec are shown with a notice.

## Example

```
Query: "Show details for GET /pet/{petId}"
→ inspect(endpointId: "abc123...")
→ Result:
  GET /pet/{petId}
  Parameters:
    - petId (path, integer, required)
  Responses:
    - 200: Pet object
    - 400: Error
    - 404: Not found
```
