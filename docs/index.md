# swag2mcp

<a href="https://www.youtube.com/watch?v=9CcvwmfTkds" target="_blank">
  <img src="https://raw.githubusercontent.com/mmadfox/swag2mcp/main/assets/cover.png" width="600" alt="Watch video — swag2mcp in 2 minutes">
</a>

## Your API speaks LLM

One line of config turns any OpenAPI/Swagger/Postman file into an MCP server. LLM agents discover, inspect, and invoke your APIs — zero integration code.

<img src="architecture.svg" width="800" alt="swag2mcp architecture">

## Stop writing wrappers

Every time you connect a new API to an LLM, you write the same boilerplate: spec parsing, authentication, error handling, rate limiting. swag2mcp does it for you — 19 ready-made MCP tools.

## Who needs this

| Role | Why |
|------|-----|
| **AI Agent Developer** | Connect any API in 2 minutes, not 2 days |
| **MCP Engineer** | No handler code — just point to a spec |
| **Architect** | Single API integration layer for all LLMs in your company |
| **Data Analyst** | Access APIs via natural language, no coding |
| **DevOps / SRE** | Monitoring and automation through LLM without extra services |
| **Integrator** | 9 auth methods out of the box — Basic to OAuth2 to HMAC |
| **QA Engineer** | Mock server for isolated testing without real APIs |
| **Product Manager** | Rapid AI feature prototypes without backend work |
| **and many others** | |

## Quick Start

```bash
go install github.com/mmadfox/swag2mcp@latest
swag2mcp init
swag2mcp add https://petstore.swagger.io/v2/swagger.json
swag2mcp mcp
```

## Integrations

OpenCode · Cursor · Claude Desktop · VS Code · Crush
