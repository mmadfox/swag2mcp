# Endpoints

An endpoint is a specific HTTP method + path that can be invoked (e.g., `GET /api/users/{id}`). Endpoints are the actual API operations that the LLM discovers, inspects, and calls.

## Structure

Each endpoint contains:

- **HTTP method**: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
- **Path**: `/api/v1/users/{id}`
- **Summary**: a short description of what the endpoint does — very useful for the LLM to understand its purpose at a glance
- **Description**: a detailed explanation of the endpoint's behavior, parameters, and use cases
- **Parameters**: path, query, header, cookie
- **Request body**: for POST/PUT/PATCH
- **Responses**: status codes and response schemas

The `summary` and `description` fields come from the OpenAPI/Swagger/Postman file. They are the primary way the LLM understands what an endpoint does. Well-written summaries make endpoint discovery much more effective.

## MCP Tools for Endpoints

| Tool | Description |
|------|-------------|
| `endpoint_by_spec` | All endpoints in a spec |
| `endpoint_by_collection` | Endpoints in a collection |
| `endpoint_by_tag` | Endpoints in a tag |
| `endpoint_by_id` | Quick endpoint summary |
| `inspect` | Full endpoint details (schemas, params) |
| `invoke` | Call the endpoint |
| `search` | Search endpoints by text |

## Deprecated Endpoints

Endpoints marked as `deprecated` in the spec are shown with a notice when inspected.

## Configuration

Endpoints are **read-only** from swag2mcp's perspective. There are no YAML config settings for endpoints — you cannot add, remove, rename, or modify them in `swag2mcp.yaml`.

To change endpoints (add new ones, update summaries, modify parameters, mark as deprecated), edit the original OpenAPI/Swagger/Postman file and run `swag2mcp update` to re-parse and re-index.

## Example

```
Query: "Show details for GET /pet/{petId}"
→ inspect(endpointId: "abc123...")
→ Result:
  GET /pet/{petId}
  Summary: Find pet by ID
  Description: Returns a single pet by its ID
  Parameters:
    - petId (path, integer, required)
  Responses:
    - 200: Pet object
    - 400: Error
    - 404: Not found
```
