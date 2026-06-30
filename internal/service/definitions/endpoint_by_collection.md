---
name: endpoint_by_collection
---

# endpoint_by_collection

Lists all endpoints within a specific collection, regardless of their tag.

## When to use

Use this tool when:
- You have a collection ID and want to see every endpoint it contains
- The user asks "show me all endpoints in collection X"
- You need a complete inventory of a collection's API surface

For a filtered view by tag within a collection, use `endpoint_by_tag` instead.

## Parameters

- `collectionId` (required): The 32-character MD5 hash ID of the collection

## Returns

A list of endpoints with their IDs, HTTP methods, paths, summaries, and deprecation status.
