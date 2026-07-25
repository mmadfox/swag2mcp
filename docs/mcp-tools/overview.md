# MCP Tools

## Overview

swag2mcp provides **19 MCP tools** that give an LLM agent full access to your APIs through the Model Context Protocol. These tools cover the complete workflow: discovering what APIs are available, navigating the spec hierarchy, searching and inspecting endpoints, executing API calls, and working with large responses.

### What the tools solve

- **Discovery** — the LLM can find specs, collections, and tags without knowing IDs in advance
- **Navigation** — drill down from spec → collection → tag → endpoint in a structured hierarchy
- **Search** — full-text search across all endpoints when you don't have an ID
- **Inspection** — get the full OpenAPI operation object before making a call
- **Execution** — invoke real API calls with automatic authentication
- **Large response handling** — outline, compress, and slice oversized responses that don't fit inline

### Read-only vs Mutable

| Type | Count | Tools |
|------|-------|-------|
| **Read-only** | 17 | All discovery, endpoint, search, inspect, info, and response tools |
| **Mutable** | 2 | `invoke` (makes real HTTP calls), `auth` (retrieves tokens) |

Read-only tools are marked with `ReadOnlyHint=true` and `IdempotentHint=true` in the MCP protocol, signaling to the LLM that they are safe to call without side effects.

### Error handling

All tools return errors as structured `LLMError` objects with a machine-readable code and a human-readable message that explains what went wrong and what to do next:

| Error code | Meaning |
|------------|---------|
| `validation_failed` | Invalid input (bad ID format, missing required fields) |
| `not_found` | Entity not found in the index or workspace |
| `rate_limit` | Second `invoke` call within 10 seconds on the same endpoint |
| `invoke_error` | HTTP call failure, download failure |
| `auth_error` | Auth token retrieval failure |
| `config_error` | Config file load or save failure |
| `parse_error` | Spec file parse failure |

## Categories

| Category | Tools | Description |
|----------|-------|-------------|
| **Discovery** | `spec_list`, `spec_by_id`, `collection_by_spec`, `collection_by_id`, `tag_by_spec`, `tag_by_collection`, `tag_by_id` | Navigate the spec hierarchy: find specs, collections, and tags |
| **Endpoints** | `endpoint_by_spec`, `endpoint_by_collection`, `endpoint_by_tag`, `endpoint_by_id` | View endpoints at different levels of the hierarchy |
| **Execution** | `search`, `inspect`, `invoke` | Search, inspect the full contract, and call APIs |
| **Utilities** | `auth`, `info`, `response_outline`, `response_compress`, `response_slice` | Auth tokens, runtime info, and large response handling |
| **Skills** | [Formatting guide](/mcp-tools/skills) | Customize how tool responses are displayed |

## Full List

| Tool | Description |
|------|-------------|
| `spec_list` | List all API specifications in the workspace |
| `spec_by_id` | Get detailed spec information with collections |
| `collection_by_spec` | List collections within a spec |
| `collection_by_id` | Get collection details with tags |
| `tag_by_spec` | List all tags across a spec |
| `tag_by_collection` | List tags within a collection |
| `tag_by_id` | Get tag details (ID, title, method count) |
| `endpoint_by_spec` | List all endpoints in a spec |
| `endpoint_by_collection` | List endpoints in a collection |
| `endpoint_by_tag` | List endpoints in a tag |
| `endpoint_by_id` | Quick endpoint summary (method, path, summary) |
| `search` | Full-text search across all endpoints |
| `inspect` | Full OpenAPI operation details (parameters, schemas) |
| `invoke` | Execute a real API call |
| `auth` | Get auth token or headers for a spec |
| `info` | Runtime information (version, specs, config) |
| `response_outline` | Structural summary of a large response file |
| `response_compress` | Compress a large response to fit inline |
| `response_slice` | Extract a fragment of a large response |

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

When you don't have an ID, use `search` to find endpoints by query. When `invoke` returns a `fileRef` (response too large), use `response_outline` → `response_compress` or `response_slice` to explore the data.
