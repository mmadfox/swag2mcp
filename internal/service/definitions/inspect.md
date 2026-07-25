---
name: inspect
---

# inspect

Retrieves the **full OpenAPI operation object** for an endpoint — parameters, request body, response schemas, and all definitions.

Use this when you need complete technical details before invoking an endpoint or explaining its contract to the user.

For a **quick summary** (method, path, summary only), use `endpoint_by_id` instead.

## When to use

Use this tool when:
- The user asks "show me the full spec for this endpoint"
- You need to understand the exact schema for request/response bodies
- You need to know which parameters (path, query, header) an endpoint accepts
- You are preparing to call an endpoint via `invoke` and need to build the correct request
- The user asks for examples, response codes, or technical contract details

## Parameters

- `endpointId` (required): The 32-character MD5 hash ID of the endpoint to inspect

## Returns

The full OpenAPI operation object including parameters (with schemas), request body, responses, and all referenced schema definitions.
