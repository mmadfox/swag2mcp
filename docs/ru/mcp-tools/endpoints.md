# Endpoint Tools

Endpoint tools let the LLM view API endpoints at different levels of the hierarchy: all endpoints in a spec, in a collection, in a tag, or a single endpoint summary. Use these to discover available operations before inspecting or invoking.

---

## endpoint_by_spec

### Purpose

List all endpoints across an entire spec, spanning all collections and tags. Returns the most comprehensive view — every endpoint in the spec with its full context (tag, collection, spec).

### When to use

- When you want to see every endpoint available in a spec
- When you don't know which collection or tag contains the endpoint you need
- After `spec_by_id` to get the full endpoint list

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `specId` | string | Yes | 32-character MD5 hash of the spec |

### Response

```json
{
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "tagName": "forecast",
      "collectionId": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "collectionTitle": "Weather Forecast",
      "specId": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
      "specDomain": "meteo",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Get weather forecast for a location"
    }
  ]
}
```

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Endpoint identifier |
| `tagId` | string | Parent tag identifier |
| `tagName` | string | Human-readable tag name |
| `collectionId` | string | Parent collection identifier |
| `collectionTitle` | string | Human-readable collection title |
| `specId` | string | Parent spec identifier |
| `specDomain` | string | Spec domain name |
| `method` | string | HTTP method (GET, POST, PUT, DELETE, etc.) |
| `path` | string | API path (e.g. /v1/forecast) |
| `summary` | string | Human-readable summary of what the endpoint does |

### Nuances

- Returns `not_found` if the spec does not exist
- Each endpoint includes its full ancestry (spec → collection → tag) for context
- For a quick summary of a single endpoint, use `endpoint_by_id`

---

## endpoint_by_collection

### Purpose

List all endpoints within a specific collection, regardless of their tag. Returns endpoints grouped by collection with spec and collection metadata.

### When to use

- After `collection_by_id` to see all endpoints in a collection
- When you want to explore a collection's full API surface

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `collectionId` | string | Yes | 32-character MD5 hash of the collection |

### Response

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "tagName": "forecast",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Get weather forecast for a location"
    }
  ]
}
```

### Nuances

- Returns `not_found` if the collection does not exist
- Includes spec and collection metadata for context
- Endpoints from all tags within the collection are returned together

---

## endpoint_by_tag

### Purpose

List all endpoints grouped under a specific tag. This is the most focused view — endpoints in one tag within one collection.

### When to use

- After `tag_by_id` to see the actual endpoints in a tag
- When you know the tag and want to see its operations

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `tagId` | string | Yes | 32-character MD5 hash of the tag |

### Response

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  },
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Get weather forecast for a location"
    }
  ]
}
```

### Nuances

- Returns `not_found` if the tag does not exist
- Includes full context: spec, collection, and tag metadata
- Endpoints are scoped to a single tag within a single collection

---

## endpoint_by_id

### Purpose

Get a quick summary of a single endpoint: method, path, summary, and deprecation status. This is a lightweight tool — for the full OpenAPI operation object (parameters, request body, response schemas), use `inspect`.

### When to use

- When you have an endpoint ID and want a quick reminder of what it does
- Before deciding whether to call `inspect` for full details
- When you need to confirm the method and path before invoking

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | 32-character MD5 hash of the endpoint |

### Response

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  },
  "endpoint": {
    "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
    "method": "GET",
    "path": "/v1/forecast",
    "summary": "Get weather forecast for a location"
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `endpoint.id` | string | Endpoint identifier |
| `endpoint.method` | string | HTTP method |
| `endpoint.path` | string | API path |
| `endpoint.summary` | string | Human-readable summary |

### Nuances

- Returns `not_found` if the endpoint does not exist
- This is a **quick summary** — it does not return parameters, request body, or response schemas
- For full technical details (required before `invoke`), use `inspect`
