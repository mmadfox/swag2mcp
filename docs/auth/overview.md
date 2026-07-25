# Authentication

## Overview

swag2mcp supports **9 authentication methods** for working with APIs that require authorization. You configure it once in the config file — after that, every API call through `invoke` automatically includes the right tokens and headers.

### Where to configure

Authentication is set at the **spec** level in `swag2mcp.yaml`:

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: bearer
      config:
        token: "my-token"
```

### How it works

- You specify the auth type and parameters in the config
- swag2mcp automatically applies them to every request when you call `invoke`
- You **don't need** to request a token before calling an API — it happens automatically
- If a token expires (OAuth2, Script), swag2mcp refreshes it on its own

### Environment variables

Sensitive data (tokens, passwords, keys) can be stored in environment variables using `$(VAR_NAME)` syntax:

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_API_TOKEN)"
```

swag2mcp substitutes the value of `MY_API_TOKEN` at startup.

### MCP auth tool

The LLM agent can retrieve a token or headers through the `auth` MCP tool — for example, to build a curl command or show the user.

In **production**, this tool should be disabled with `--disable-llm-auth` (enabled by default) so the LLM never has access to tokens.

### Methods

| Method | Description | Best for |
|--------|-------------|----------|
| [`none`](/auth/none) | No authentication | Public APIs |
| [`basic`](/auth/basic) | HTTP Basic (username + password) | Legacy APIs, simple auth |
| [`bearer`](/auth/bearer) | Bearer Token (JWT, token) | Modern REST APIs |
| [`api-key`](/auth/api-key) | API key in header or query parameter | Services with API keys |
| [`digest`](/auth/digest) | HTTP Digest (username + password) | Legacy APIs, more secure than Basic |
| [`hmac`](/auth/hmac) | HMAC-SHA256 signature (Binance-style) | Cryptocurrency exchanges |
| [`oauth2-cc`](/auth/oauth2-cc) | OAuth2 Client Credentials | Server-to-server, microservices |
| [`oauth2-pwd`](/auth/oauth2-pwd) | OAuth2 Password Grant | Apps with user login |
| [`script`](/auth/script) | External script to obtain a token | Any custom auth scheme |
