---
name: spec_by_id
---

# spec_by_id

Retrieves detailed information about a specific API specification by its ID, including its domain and all associated collections.

## When to use

Use this tool when:
- You need to explore a specific API's structure after discovering it via `spec_list`
- The user asks "show me details about API X" or "what collections are in spec Y?"
- You have a `specId` and need to get its collections before drilling into tags

## Parameters

- `id` (required): The 32-character MD5 hash ID of the specification

## Returns

The specification's ID and domain, plus a list of collections with their IDs, titles, tag counts, and method counts.
