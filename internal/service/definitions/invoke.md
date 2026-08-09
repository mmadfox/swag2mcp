---
name: invoke
---

# invoke

Executes a real API call to an endpoint using the provided parameters and returns the response data, status code, and headers.

## When to use

Use this tool **only** when the user explicitly asks to perform an action, such as:
- "Get all pets" or "Create a user"
- "Call the API" or "Make a request"
- "Test the endpoint" or "Try it out"

Always use `inspect` first to understand the required parameters, headers, and request body before invoking.

**Never invoke a destructive operation (POST/PUT/PATCH/DELETE) without explicit user confirmation.**

## Parameters

- `endpointId` (required): The 32-character MD5 hash ID of the endpoint to invoke
- `parameters` (optional): Object containing path, query, and header parameters as key-value pairs
- `requestBody` (optional): The request body for POST/PUT/PATCH requests. Provide as a JSON object matching the schema from `inspect`

## Returns

The API response data, HTTP status code, and response headers.

**Large responses** (over ~50 KB) are automatically saved to disk. You will receive a file path — use it to reference the result rather than displaying the full content inline.

## Important

- **Auth is automatic**: `invoke` handles authentication automatically. Do NOT pass `headers` or `cookies` — swag2mcp applies auth under the hood.
- **Auth-injected parameters are automatic**: If `inspect` shows a required parameter that is an auth credential (e.g. `api_key`, `timestamp`, `signature`, `recvWindow`), do NOT pass it in `parameters` — swag2mcp injects it automatically. Only pass genuine business parameters.
- **One at a time, bounded retries**: Make at most **one outstanding invoke** at a time — never launch a batch. If a call is `rate_limit`-ed, do not retry immediately or in a batch; back off (~15s, then ~30s) and retry that endpoint at most twice more, then mark it as **"rate limited"** and move on. Non-rate-limit errors are final — do not retry.
- **Rate limits**: Before calling `invoke` in a loop, call `info` to check `rate_limiting` (per_endpoint_interval, global_limit). Respect these limits to avoid throttling.
