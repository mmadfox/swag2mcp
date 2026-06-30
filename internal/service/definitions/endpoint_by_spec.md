---
name: endpoint_by_spec
---

# endpoint_by_spec

Lists all endpoints across an entire API specification, spanning all collections and tags.

## When to use

Use this tool when:
- You want a comprehensive view of every endpoint in a spec
- The user asks "show me all endpoints in API X" or "what does spec Y expose?"
- You need to search across all collections within a single spec

For a narrower scope, use `endpoint_by_collection` (single collection) or `endpoint_by_tag` (single tag).

## Parameters

- `specId` (required): The 32-character MD5 hash ID of the specification

## Returns

A list of endpoints with their IDs, HTTP methods, paths, summaries, and deprecation status.
