---
name: endpoint_by_tag
---

# endpoint_by_tag

Lists all endpoints grouped under a specific tag within a collection.

## When to use

Use this tool when:
- You have a tag ID and want to see all endpoints in that category
- The user asks "show me all pets endpoints" or "what's in the orders tag?"
- You've identified a tag via `tag_by_id` or `tag_by_collection` and want its endpoints

To see the tag's metadata (title, method count) without listing endpoints, use `tag_by_id` instead.

## Parameters

- `tagId` (required): The 32-character MD5 hash ID of the tag

## Returns

A list of endpoints with their IDs, HTTP methods, paths, summaries, and deprecation status.
