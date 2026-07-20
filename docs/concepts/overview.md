# Concepts

## Architecture

swag2mcp acts as a bridge between API specifications and LLM agents:

<img src="/architecture.svg" width="800" alt="swag2mcp architecture">

## Core Concepts

**Spec** — an OpenAPI/Swagger/Postman file describing an API. You add a spec to swag2mcp, and it automatically parses it into its parts. Learn more: [Specs](./specs).

**Collection** — a logical group of endpoints within a spec. One spec can have multiple collections. For example, a weather API spec might have "Forecast", "Air Quality", and "Marine" collections. Learn more: [Collections](./collections).

**Tag** — a category of endpoints inside a collection. Helps the LLM find the right operations more precisely. Learn more: [Tags](./tags).

**Endpoint** — a specific HTTP method + path (e.g., `GET /api/users`). The LLM can find an endpoint by description, inspect its parameters and schemas, and then invoke it. Learn more: [Endpoints](./endpoints).

**Workspace** — the directory where swag2mcp stores config, spec cache, saved responses, and auth scripts. Learn more: [Workspace](./workspace).

## How It Works

1. **Add a spec** — via `swag2mcp add <url>` or in the YAML config, point to an OpenAPI/Swagger/Postman file.
2. **swag2mcp parses the spec** — creates Collections, Tags, and Endpoints, indexes them for search.
3. **LLM finds the right endpoint** — through MCP tools (`search`, `endpoint_by_tag`, `inspect`) the LLM searches for a matching endpoint by description, reviews its parameters and request schema.
4. **LLM invokes the endpoint** — via the MCP tool `invoke`, the LLM sends the request. swag2mcp validates every input parameter against the endpoint's OpenAPI schema (path params, query params, headers, request body) before making the call. If something doesn't match the schema, the LLM gets a clear error explaining what's wrong. Once validated, swag2mcp executes the real HTTP call and returns the result.
5. **Result goes back to the LLM** — the API response is passed back to the agent. Large responses are saved to the workspace and accessible via a file reference.

swag2mcp is a bridge between LLMs and the world of APIs. You add API specifications, and the LLM — through the MCP protocol — finds the right endpoints, inspects their documentation, and calls them. All you need to do is add a spec and start the MCP server.

> **Config is editable at any time.** The YAML config file (`~/.swag2mcp/swag2mcp.yaml`) can be edited by hand — add specs, change auth, tweak settings. After every edit, restart the MCP server (`swag2mcp mcp`) for changes to take effect.

## Hierarchy

```
Spec (OpenAPI file)
  └── Collection 1 (logical group)
        └── Tag 1 (category)
              └── Endpoint (GET /api/users)
              └── Endpoint (POST /api/users)
        └── Tag 2
              └── Endpoint (GET /api/users/{id})
  └── Collection 2
        └── Tag 3
              └── Endpoint (DELETE /api/users/{id})
```
