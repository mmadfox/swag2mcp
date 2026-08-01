# swag2mcp

Bridges OpenAPI/Swagger/Postman API specifications with LLM agents via the Model Context Protocol (MCP).

Not every API speaks MCP — private endpoints, internal services, legacy systems, and third-party APIs rarely do. swag2mcp wraps any REST API in an MCP interface, giving LLM agents instant access to your entire API surface without modifying a single line of server code. Through live API calls, the LLM gains real-world knowledge to make informed decisions, automate workflows, and act on your data — not just guess.

<a href="https://www.youtube.com/watch?v=1Da4UmE2f9U" target="_blank">
  <img src="https://raw.githubusercontent.com/mmadfox/swag2mcp/main/docs/cover.jpg" alt="Preview">
</a>

## Your API speaks LLM

One line of config turns any OpenAPI/Swagger/Postman file into an MCP server. LLM agents discover, inspect, and invoke your APIs — zero integration code.

<img src="/architecture.svg" width="700" alt="swag2mcp architecture">

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

---

## License

Licensed under the **Apache License, Version 2.0**.

See [LICENSE](https://github.com/mmadfox/swag2mcp/blob/main/LICENSE) for the full license text.
