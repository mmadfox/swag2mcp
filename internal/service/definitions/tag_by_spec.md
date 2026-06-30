---
name: tag_by_spec
---

# tag_by_spec

Lists all tags across an entire API specification, spanning all collections.

## When to use

Use this tool when:
- You want to see every tag available in an API spec
- The user asks "what categories exist in spec Y?"
- You need a global view of all endpoint categories without drilling into each collection

Use `tag_by_collection` instead if you only need tags within a single, specific collection.

## Parameters

- `specId` (required): The 32-character MD5 hash ID of the specification

## Returns

A list of tags with their IDs, titles, and method counts.
