---
name: tag_by_collection
---

# tag_by_collection

Lists all tags within a specific collection.

## When to use

Use this tool when:
- You have a collection ID and want to see how its endpoints are categorized
- The user asks "what tags are in collection X?" or "show me the categories"
- You are navigating hierarchically: spec → collection → tag → endpoint

Use `tag_by_spec` instead if you want all tags across an entire specification (not just one collection).

## Parameters

- `collectionId` (required): The 32-character MD5 hash ID of the collection

## Returns

A list of tags with their IDs, titles, and method counts.
