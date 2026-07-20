# Execution Tools

Tools for searching, inspecting, and invoking APIs.

## search

Full-text search across all endpoints.

**Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `query` | string | Search query |
| `limit` | int | Max results (1-50) |

**Query Syntax**:
| Example | Description |
|---------|-------------|
| `pet` | Simple search |
| `+method:POST +tag:pet` | POST in pet tag |
| `path:"/api/v1/users"` | Path search |
| `inventory` | Description search |
| `create~` | Fuzzy search |

**Example**:
```
→ search(query: "find pet by status", limit: 5)
← GET /pet/findByStatus — Finds Pets by status
   GET /pet/{petId} — Find pet by ID
```

## inspect

Full endpoint details: parameters, request body, response schemas.

**Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `endpointId` | string | Endpoint ID |

**Example**:
```
→ inspect(endpointId: "def456")
← POST /pet
  Parameters:
    - name (query, string, required)
    - status (query, string, enum: available/pending/sold)
  Request Body:
    type: object
    properties:
      name: string
      photoUrls: string[]
  Responses:
    200: Pet object
    405: Invalid input
```

## invoke

Call an API endpoint.

**Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `endpointId` | string | Endpoint ID |
| `parameters` | object | Path, query, header params |
| `requestBody` | object | Request body |
| `headers` | object | Additional headers |
| `cookies` | object | Additional cookies |

**Example**:
```
→ invoke(
    endpointId: "def456",
    parameters: { "petId": 1 }
  )
← 200 OK
  {
    "id": 1,
    "name": "Rex",
    "status": "available"
  }
```

!!! warning
    Never invoke destructive operations (POST/PUT/PATCH/DELETE) without explicit user confirmation.
