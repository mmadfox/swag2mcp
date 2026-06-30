---
name: collection_by_id
---

# collection_by_id

Retrieves detailed information about a specific collection including its tags, spec metadata, and method statistics.

## When to use

Use this tool when:
- You have a collection ID and want to see its tags before drilling into endpoints
- The user asks for details about a specific collection like "show me what's in the pets collection"
- You need to see the spec that owns this collection

To list all collections in a spec first, use `collection_by_spec` instead.

## Parameters

- `id` (required): The 32-character MD5 hash ID of the collection

## Returns

The collection's spec (ID, domain), collection details (ID, title, method count), and a list of tags with their IDs, titles, and method counts.
