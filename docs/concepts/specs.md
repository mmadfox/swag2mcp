# Specs

A spec is a logical container representing an API domain or service (e.g., YouTube, Binance, Open-Meteo). Each spec has a unique `domain`, a `base_url`, optional `auth`, and contains one or more collections.

[Collections](./collections) point to OpenAPI/Swagger/Postman files — the spec itself is not a file, it's the grouping around them.

## Domain — Naming Rules

The `domain` is the unique identifier of a spec. It is used as the primary key throughout the system.

| Rule | Constraint |
|------|------------|
| Characters | `a-z`, `0-9`, `_`, `-` only |
| Length | 1–60 characters |
| Uniqueness | **No duplicates allowed** — two active specs cannot share the same domain |

**Valid examples:** `meteo`, `binance`, `github-api`, `my_service`, `openai-v1`

**Invalid examples:** `Meteo` (uppercase), `my api` (space), `my.api` (dot), `a-very-long-domain-name-that-exceeds-sixty-characters` (too long)

## Spec Fields

| Field | YAML key | Required | Description |
|-------|----------|----------|-------------|
| [Domain](#domain--naming-rules) | `domain` | ✅ | Unique API identifier (1–60 chars, `a-z0-9_-`) |
| LLM Title | `llm_title` | ✅ | Human-readable name the LLM uses to reference this API (5–120 chars) |
| [LLM Instruction](#llm-instruction) | `llm_instruction` | ❌ | Short hint injected into the swag2mcp system prompt (max 500 chars) |
| Base URL | `base_url` | ✅ | Base URL for all API requests (valid URL) |
| [Disable](#disable) | `disable` | ❌ | Skip this spec during loading and indexing |
| [Tags](#tags) | `tags` | ❌ | Tags for filtering (e.g., `["public", "demo"]`) |
| [Auth](#auth) | `auth` | ❌ | Authentication configuration |
| [HTTP Client](#http-client) | `http_client` | ❌ | Per-spec HTTP settings (headers, cookies) |
| [Collections](./collections) | `collections` | ✅ | List of 1–30 collections |

## Validation

When swag2mcp validates the config, these rules are checked for every spec:

| Check | Rule |
|-------|------|
| **Duplicate domains** | No two active specs may share the same `domain` |
| **Domain format** | Must match `^[a-z0-9_-]{1,60}$` |
| **LLM Title** | Required, 5–120 characters, letters/digits/spaces/basic punctuation |
| **LLM Instruction** | Max 500 characters, same character set as title |
| **Base URL** | Required, must be a valid URL |
| **Collections** | Required, 1–30 items |
| **Auth** | Validated per auth type (e.g., bearer requires `token`, basic requires `username` + `password`) |
| **Location** | Each collection's `location` must be a valid URL or file path (5–250 chars) |

Validation runs on every `swag2mcp mcp` startup. If it fails, the MCP server will not start — in some IDEs this means the server simply won't connect, and the LLM receives a clear error message explaining what to fix.

To diagnose issues before starting the server, use the [`validate`](../cli/validate.md) command:

```bash
# Validate default workspace (~/.swag2mcp)
swag2mcp validate

# Validate a custom project workspace
swag2mcp validate ./my-project
```

## LLM Instruction

It is recommended to set `llm_instruction` on each spec — a short hint (up to 500 chars) that tells the LLM what this API is for and when to use it. This instruction is injected into the swag2mcp system prompt, helping the LLM understand the spec's purpose without extra context.

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    llm_instruction: "Use this API to get random dad jokes or search for specific jokes by keyword."
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

Collections can also have their own `llm_instruction` (up to 360 chars) for more specific guidance.

## Auth

Authentication is configured at the spec level and applies to all its collections. swag2mcp supports 9 auth methods:

| Method | YAML type | Key fields |
|--------|-----------|------------|
| [None](../auth/none.md) | `none` | — |
| [Basic](../auth/basic.md) | `basic` | `username`, `password` |
| [Bearer](../auth/bearer.md) | `bearer` | `token` |
| [Digest](../auth/digest.md) | `digest` | `username`, `password` |
| [OAuth2 Client Credentials](../auth/oauth2-cc.md) | `oauth2-cc` | `client_id`, `client_secret`, `token_url` |
| [OAuth2 Password](../auth/oauth2-pwd.md) | `oauth2-pwd` | `username`, `password`, `client_id`, `token_url` |
| [API Key](../auth/api-key.md) | `api-key` | `key`, `value`, `in` (`header` or `query`) |
| [HMAC](../auth/hmac.md) | `hmac` | `api_key`, `secret_key` |
| [Script](../auth/script.md) | `script` | `domain` |

See [Auth Overview](../auth/overview.md) for full details on each method.

## HTTP Client

You can override HTTP settings at the spec level. These apply to all requests made by this spec's collections.

```yaml
specs:
  - domain: slow-api
    llm_title: Slow API
    base_url: https://slow-api.example.com
    http_client:
      headers:
        X-API-Version: "2"
      cookies:
        - name: session
          value: abc123
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

Settings cascade: global → spec → collection. See [Configuration Cascade](../configuration/cascade.md) for details.

## Tags

Tags let you filter specs by category. Use them with the `--tags` flag on `swag2mcp ls` or during bootstrap.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    tags: ["weather", "public"]
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

```bash
# List only specs tagged "weather"
swag2mcp ls --tags weather
```

## Disable

Set `disable: true` to skip a spec entirely. It won't be loaded, indexed, or available to the LLM.

```yaml
specs:
  - domain: old-api
    llm_title: Old API (Deprecated)
    base_url: https://old-api.example.com
    disable: true
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Examples

### Minimal Spec

```yaml
specs:
  - domain: dadjokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

### Spec with Auth

```yaml
specs:
  - domain: binance
    llm_title: Binance Market Data API
    base_url: https://api.binance.com
    auth:
      type: hmac
      config:
        api_key: $(BINANCE_API_KEY)
        secret_key: $(BINANCE_SECRET_KEY)
    collections:
      - llm_title: Market Data
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml
```

### Spec with Multiple Collections

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

### Spec with LLM Instruction and Tags

```yaml
specs:
  - domain: rickandmorty
    llm_title: Rick and Morty API
    llm_instruction: "Use this API to get information about characters, episodes, and locations from the Rick and Morty show."
    base_url: https://rickandmortyapi.com/api
    tags: ["entertainment", "public"]
    collections:
      - llm_title: Characters
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/rick-and-morty.json
```

## Related

- [Spec Settings (config)](../configuration/spec-settings.md) — full YAML reference
- [Configuration Cascade](../configuration/cascade.md) — how settings override each other
- [Auth Overview](../auth/overview.md) — all 9 auth methods
- [HTTP Client](../configuration/http-client.md) — HTTP client configuration
- [Collections](./collections) — spec files within a spec
