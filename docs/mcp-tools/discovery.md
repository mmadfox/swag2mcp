# Discovery Tools

Tools for navigating the spec hierarchy.

## spec_list

List all registered specs.

**Parameters**: none

**Example**:
```
→ spec_list()
← petstore (15 endpoints)
  meteo (9 endpoints)
  binance (24 endpoints)
```

## spec_by_id

Detailed spec information.

**Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Spec ID |

**Example**:
```
→ spec_by_id(id: "abc123")
← Petstore API
  Domain: petstore.swagger.io
  Collections: pet (3), store (4), user (8)
```

## collection_by_spec

List collections in a spec.

**Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `specId` | string | Spec ID |

## collection_by_id

Detailed collection information.

**Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Collection ID |

## tag_by_spec

All tags in a spec.

**Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `specId` | string | Spec ID |

## tag_by_collection

Tags in a collection.

**Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `collectionId` | string | Collection ID |

## tag_by_id

Detailed tag information.

**Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Tag ID |
