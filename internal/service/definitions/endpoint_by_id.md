---
name: endpoint_by_id
---

# endpoint_by_id

Returns a quick summary of a single endpoint: method, path, summary, and deprecation status.

Use **`inspect`** instead when you need the full OpenAPI operation object with parameters, request body, and response schemas.

## When to use

Use this tool when:
- You already have an endpoint ID and want a quick overview
- The user asks "what is this endpoint?" at a high level
- You need the method, path, or summary to present to the user

Do **NOT** use this tool when you need technical details (schemas, parameters, request body) — use `inspect` instead.

## Parameters

- `id` (required): The 32-character MD5 hash ID of the endpoint

## Returns

The endpoint's method (GET/POST/etc.), path, summary, and whether it's deprecated.
