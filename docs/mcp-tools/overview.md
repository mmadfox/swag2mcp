# MCP Tools

swag2mcp provides 19 MCP tools for LLM agents.

## Categories

| Category | Tools | Description |
|----------|-------|-------------|
| **Discovery** | `spec_list`, `spec_by_id`, `collection_by_spec`, `collection_by_id` | Find specs and collections |
| **Tags** | `tag_by_spec`, `tag_by_collection`, `tag_by_id` | Navigate tags |
| **Endpoints** | `endpoint_by_spec`, `endpoint_by_collection`, `endpoint_by_tag`, `endpoint_by_id` | View endpoints |
| **Execution** | `search`, `inspect`, `invoke` | Search, inspect, and call APIs |
| **Utilities** | `auth`, `info`, `response_outline`, `response_compress`, `response_slice` | Helper tools |

## Full List

| Tool | Description |
|------|-------------|
| `spec_list` | List all specs |
| `spec_by_id` | Spec details |
| `collection_by_spec` | Collections in a spec |
| `collection_by_id` | Collection details |
| `tag_by_spec` | Tags in a spec |
| `tag_by_collection` | Tags in a collection |
| `tag_by_id` | Tag details |
| `endpoint_by_spec` | Endpoints in a spec |
| `endpoint_by_collection` | Endpoints in a collection |
| `endpoint_by_tag` | Endpoints in a tag |
| `endpoint_by_id` | Quick endpoint summary |
| `search` | Full-text search |
| `inspect` | Full endpoint details |
| `invoke` | Call an API |
| `auth` | Get auth token/headers |
| `info` | System info |
| `response_outline` | Large response structure |
| `response_compress` | Compress large response |
| `response_slice` | Fragment of large response |

## Navigation Hierarchy

```
spec_list
  └── spec_by_id(id)
        └── collection_by_spec(specId)
              └── collection_by_id(id)
                    └── tag_by_collection(collectionId)
                          └── tag_by_id(id)
                                └── endpoint_by_tag(tagId)
                                      └── endpoint_by_id(id)
                                            └── inspect(endpointId)
                                                  └── invoke(endpointId)
```
