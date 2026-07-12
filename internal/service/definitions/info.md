---
name: info
---

# info

Returns a comprehensive summary of the swag2mcp runtime: version, configuration, active specs, HTTP client settings, MCP transport, auth methods, and mock mode status.

## When to use

Use this tool when:
- The user asks "what's the current configuration?" or "show me the system status"
- You need to understand how the HTTP client is configured (timeout, proxy, headers, cookies)
- You need to know which specs are active, disabled, and their endpoint counts
- You want to check the MCP transport type and whether auth is enabled
- You need to see which auth methods are configured across all specs
- You want to check if mock mode is enabled

This tool takes **no arguments** — it returns the full runtime summary.

## Parameters

This tool has no parameters.

## Returns

A JSON object with version, latest_version (from GitHub), workspace path, uptime, specs summary, HTTP client configuration, MCP configuration, auth methods, and mock mode status.
