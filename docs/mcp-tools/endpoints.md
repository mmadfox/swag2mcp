# Endpoint Tools

Tools for viewing endpoints at different levels.

## endpoint_by_spec

All endpoints in a spec.

**Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `specId` | string | Spec ID |

## endpoint_by_collection

Endpoints in a collection.

**Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `collectionId` | string | Collection ID |

## endpoint_by_tag

Endpoints in a tag.

**Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `tagId` | string | Tag ID |

## endpoint_by_id

Quick endpoint summary.

**Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Endpoint ID |

**Example**:
```
→ endpoint_by_id(id: "def456")
← GET /pet/{petId}
  Summary: Find pet by ID
  Deprecated: false
```

## Output Format

Results are returned as a table:

| Method | Path | Description |
|--------|------|-------------|
| GET | /pet | Find pet |
| POST | /pet | Add pet |
| DELETE | /pet/{id} | Delete pet |
