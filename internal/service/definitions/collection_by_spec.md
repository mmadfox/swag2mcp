---
name: collection_by_spec
---

# collection_by_spec

Lists all collections (logical groups of endpoints) within a specific API specification.

## When to use

Use this tool when:
- You have a `specId` from `spec_list` or `spec_by_id` and want to see how endpoints are organized
- The user says "show me collections in the meteo API" or "what groups exist in spec X?"
- You need to navigate from spec → collection → tag → endpoint

After finding a collection of interest, use `collection_by_id` for its tags or `endpoint_by_collection` for its endpoints.

## Parameters

- `specId` (required): The 32-character MD5 hash ID of the specification

## Returns

A list of collections with their IDs, titles, and statistics (tag count, method count).
