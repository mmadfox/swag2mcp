# Discovery Tools

Discovery tools let the LLM navigate the spec hierarchy: find all specs, drill into a spec to see its collections, and explore tags within a collection. Start with `spec_list` to see what APIs are available, then use IDs to drill deeper.

---

## spec_list

### Purpose

List all API specifications registered in the workspace. This is the starting point for any session — the LLM calls it first to discover what APIs are available.

### When to use

- At the start of a session to see what APIs are configured
- After adding or removing specs to refresh the list
- When you need a spec ID for other tools

### How it works

Returns a list of all specs with their unique ID and domain name. No parameters needed.

### Parameters

None.

### Response

```json
{
  "specs": [
    {
      "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
      "domain": "meteo"
    },
    {
      "id": "b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7",
      "domain": "dadjoke"
    }
  ]
}
```

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | 32-character MD5 hash, unique identifier for the spec |
| `domain` | string | Domain name of the spec (e.g. "meteo", "dadjoke") |

### Nuances

- Returns only `id` and `domain` — for full details (collections, tags), use `spec_by_id`
- All IDs are 32-character MD5 hex strings (`^[0-9a-f]{32}$`)
- If no specs are configured, returns an empty array

---

## spec_by_id

### Purpose

Get detailed information about a specific spec: its domain, all collections, and their statistics (tag count, method count).

### When to use

- After `spec_list` to see the collections inside a spec
- When you need collection IDs for further navigation

### How it works

Takes a spec ID and returns the spec metadata plus all its collections with counts.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | 32-character MD5 hash of the spec |

### Response

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collections": [
    {
      "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "title": "Weather Forecast",
      "llmTitle": "Forecast API",
      "countTags": 3,
      "countMethods": 12
    }
  ]
}
```

| Field | Type | Description |
|-------|------|-------------|
| `spec.id` | string | Spec identifier |
| `spec.domain` | string | Spec domain name |
| `collections[].id` | string | Collection identifier |
| `collections[].title` | string | Human-readable title |
| `collections[].llmTitle` | string | LLM-friendly title (optional) |
| `collections[].countTags` | int | Number of tags in the collection |
| `collections[].countMethods` | int | Number of HTTP methods in the collection |

### Nuances

- Returns `not_found` error if the spec ID does not exist
- The `id` must be a valid 32-character MD5 hex string

---

## collection_by_spec

### Purpose

List all collections within a specific spec. Similar to `spec_by_id` but returns only the collection list without extra spec metadata.

### When to use

- When you already have the spec ID and just need the collection list
- As a lighter alternative to `spec_by_id`

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `specId` | string | Yes | 32-character MD5 hash of the spec |

### Response

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collections": [
    {
      "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "title": "Weather Forecast",
      "llmTitle": "Forecast API",
      "countTags": 3,
      "countMethods": 12
    }
  ]
}
```

### Nuances

- Returns `not_found` if the spec does not exist
- Same data as `spec_by_id` but without the extra spec wrapper

---

## collection_by_id

### Purpose

Get detailed information about a specific collection: its metadata, the parent spec, and all tags within the collection.

### When to use

- After `collection_by_spec` to see the tags inside a collection
- When you need tag IDs for `tag_by_id` or `endpoint_by_tag`

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | 32-character MD5 hash of the collection |

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
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    },
    {
      "id": "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7",
      "title": "current",
      "countMethods": 7
    }
  ]
}
```

| Field | Type | Description |
|-------|------|-------------|
| `spec` | object | Parent spec (id, domain) |
| `collection` | object | Collection metadata (id, title, countMethods) |
| `tags[]` | array | List of tags with id, title, countMethods |

### Nuances

- Returns `not_found` if the collection ID does not exist
- Tags are returned with their IDs — use `endpoint_by_tag(tagId)` to see the actual endpoints

---

## tag_by_spec

### Purpose

List all tags across an entire spec, spanning all collections. Useful for getting a bird's-eye view of all available tags.

### When to use

- When you want to see all tags in a spec without drilling into each collection
- When you don't know which collection contains the tag you need

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `specId` | string | Yes | 32-character MD5 hash of the spec |

### Response

```json
{
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    },
    {
      "id": "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7",
      "title": "current",
      "countMethods": 7
    }
  ]
}
```

### Nuances

- Returns `not_found` if the spec does not exist
- Tags are aggregated from all collections in the spec

---

## tag_by_collection

### Purpose

List all tags within a specific collection. Unlike `tag_by_spec`, this also returns the parent spec and collection metadata.

### When to use

- After `collection_by_id` to confirm the tag list
- When you need the full context (spec + collection + tags)

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
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    }
  ]
}
```

### Nuances

- Returns `not_found` if the collection does not exist
- Same tag data as `tag_by_spec` but scoped to one collection

---

## tag_by_id

### Purpose

Get information about a single tag: its ID, title, and how many methods it contains. This tells you about the tag itself — to see the actual endpoints, use `endpoint_by_tag`.

### When to use

- When you have a tag ID and want to confirm its name and size
- Before calling `endpoint_by_tag` to understand how many endpoints to expect

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | 32-character MD5 hash of the tag |

### Response

```json
{
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `tag.id` | string | Tag identifier |
| `tag.title` | string | Human-readable tag name |
| `tag.countMethods` | int | Number of HTTP methods in this tag |

### Nuances

- Returns `not_found` if the tag does not exist
- This tool returns tag metadata only — use `endpoint_by_tag` to get the actual list of endpoints
