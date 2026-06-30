---
name: tag_by_id
---

# tag_by_id

Returns information about a single tag: its ID, title, and how many methods it contains.

This tool tells you **about the tag itself**. To get the actual list of endpoints inside a tag, use `endpoint_by_tag` instead.

## When to use

Use this tool when:
- You have a tag ID and want to verify it exists or see its metadata
- The user asks "what is this tag?" — title and method count
- You need tag statistics before deciding to explore its endpoints

Do **NOT** use this tool to get the list of endpoints — use `endpoint_by_tag` for that.

## Parameters

- `id` (required): The 32-character MD5 hash ID of the tag

## Returns

The tag's ID, human-readable title, and the number of API methods grouped under it.
