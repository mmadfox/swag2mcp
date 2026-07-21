# Collections

A collection is a single OpenAPI/Swagger/Postman file that describes a specific API. It points to a `location` (URL or local file path) and belongs to a spec (domain).

One spec can have multiple collections — for example, the "meteo" spec might have "Forecast", "Air Quality", and "Marine" collections, each pointing to a different spec file.

## Collection Fields

| Field | YAML key | Required | Description |
|-------|----------|----------|-------------|
| [LLM Title](#llm-instruction) | `llm_title` | ❌ | Collection display name for the LLM (max 120 chars) |
| [LLM Instruction](#llm-instruction) | `llm_instruction` | ❌ | Short hint for the LLM (max 360 chars) |
| Title | `title` | ❌ | Original spec title override (auto-populated from parsed document) |
| [Location](#location--how-spec-files-are-resolved) | `location` | ✅ | URL or path to the spec file (5–250 chars) |
| [Disable](#disable) | `disable` | ❌ | Skip this collection during loading |
| [HTTP Client](#http-client-override) | `http_client` | ❌ | Per-collection HTTP settings (headers, cookies) |
| [Base URL](#base-url-override) | `base_url` | ❌ | Override the spec's base URL for this collection |
| [Mock Server](#mock-server) | `base_mock_url` | ❌ | Mock server address in `host:port` format |

## Location — How Spec Files Are Resolved

The `location` field tells swag2mcp where to find the OpenAPI/Swagger/Postman file. It supports several source types:

| Source | Example | Description |
|--------|---------|-------------|
| **Remote URL** | `https://raw.githubusercontent.com/.../spec.yaml` | Downloaded and cached |
| **Local file (absolute)** | `/home/user/specs/my-api.yaml` | Read from filesystem |
| **Local file (relative)** | `./specs/my-api.yaml` | Resolved from current directory |
| **Workspace specs/** | `specs/my-api.yaml` | Stored in `~/.swag2mcp/specs/` |
| **file:// URI** | `file:///home/user/spec.yaml` | Converted to local path |

swag2mcp automatically detects the source type:

- `https://` or `http://` → remote URL (cached)
- `file://` → local file (converted to filesystem path)
- Everything else → local file (with `~` expansion for home directory)

### Remote URLs

When you use a remote URL, swag2mcp downloads the file and caches it locally. The cache is reused on subsequent starts to avoid repeated downloads.

### Local Files

Local files are read directly from the filesystem. If the file is outside the `specs/` directory, it is copied to the cache for consistency.

### Workspace specs/ Directory

The `specs/` directory inside the workspace (`~/.swag2mcp/specs/`) is the recommended place for local spec files. Files stored here are used directly without caching.

```bash
# Import a spec file into the workspace
swag2mcp import https://example.com/api.yaml

# After import, the location becomes:
# specs/api.yaml
```

## Cache System

swag2mcp caches remote spec files to avoid downloading them on every startup.

### How It Works

1. When a collection with a remote URL is loaded, swag2mcp checks the cache
2. If a valid (non-expired) cache entry exists, it is used directly
3. If not, the file is downloaded, parsed, and stored in the cache

### Cache Structure

```
~/.swag2mcp/
  cache/
    {sha256_hash}.spec    # Cached spec file content
    {sha256_hash}.meta    # Cache metadata (JSON)
```

Each cached file has a metadata file containing:

```json
{
  "source": "https://example.com/api.yaml",
  "source_type": "url",
  "cached_at": "2024-01-01T00:00:00Z",
  "mod_time": "2024-01-01T00:00:00Z",
  "ttl_sec": 3600
}
```

### Cache TTL

Each cached file gets a **random TTL** between 1 hour and 48 hours. This prevents all cached files from expiring at the same time (thundering herd problem).

### Cache Key

The cache key is a SHA-256 hash of the raw location string (first 16 bytes = 32 hex chars).

### Managing the Cache

```bash
# Clear cache and responses, re-download all spec files
swag2mcp update

# Clear cache and responses only
swag2mcp clean
```

- `swag2mcp update` — validates config, clears `cache/` and `responses/`, then re-caches all collection locations
- `swag2mcp clean` — removes all contents of `cache/` and `responses/`, plus orphan auth scripts
- Old responses are cleaned automatically after 48 hours on MCP server start

## Validation

Every collection is validated when the config is loaded. Validation runs on every `swag2mcp mcp` startup. If it fails, the MCP server will not start — in some IDEs this means the server simply won't connect, and the LLM receives a clear error message explaining what to fix.

| Check | Rule |
|-------|------|
| **Location** | Required, 5–250 characters |
| **Location accessibility** | Must be a reachable URL or existing file |
| **Location validity** | Must be a valid OpenAPI 3.x, Swagger 2.0, or Postman file |
| **LLM Title** | Max 120 characters, letters/digits/basic punctuation |
| **LLM Instruction** | Max 360 characters, same character set as title |
| **Base URL** | Must be a valid URL if set |
| **Base Mock URL** | Must be `host:port` or `host:port/path` where host is `localhost`, `127.0.0.1`, or `0.0.0.0` |
| **Mock required** | If `mock_enabled: true`, every collection must have `base_mock_url` |
| **Duplicate mock ports** | No two collections may share the same mock port |

To diagnose issues before starting the server, use the [`validate`](../cli/validate.md) command:

```bash
# Validate default workspace (~/.swag2mcp)
swag2mcp validate

# Validate a custom project workspace
swag2mcp validate ./my-project
```

## Adding Collections

### Via YAML Config

Edit `~/.swag2mcp/swag2mcp.yaml` directly:

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

After editing, restart the MCP server (`swag2mcp mcp`) for changes to take effect.

### Via CLI

```bash
# Interactive mode
swag2mcp add collection

# Non-interactive with YAML
swag2mcp add collection --yaml 'spec_domain: meteo
llm_title: Forecast
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml'

# Pipe from stdin
cat collection.yaml | swag2mcp add collection --yaml -

# Show YAML example
swag2mcp add collection --example
```

### Via Import

```bash
# Import a spec file into the workspace
swag2mcp import https://example.com/api.yaml
```

## LLM Instruction

Collections can have their own `llm_instruction` (up to 360 chars) for more specific guidance. This is injected into the swag2mcp system prompt alongside the spec-level instruction.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        llm_instruction: "Use this collection for current weather and daily forecasts."
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        llm_instruction: "Use this collection for air quality index and pollution data."
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
```

If `llm_title` is not set, it is automatically populated from the spec document's `title` field. If `llm_instruction` is not set, it is populated from the spec document's `description` field.

## Disable

Set `disable: true` to skip a collection. It won't be loaded, indexed, or available to the LLM.

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
        disable: true
```

## Base URL Override

Each collection can override the spec's `base_url`. This is useful when different collections within the same spec use different API endpoints.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

## HTTP Client Override

Collections can override HTTP settings (headers, cookies) from the spec and global levels.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          headers:
            X-API-Version: "2"
          cookies:
            - name: session
              value: abc123
```

Settings cascade: global → spec → collection. See [Configuration Cascade](../configuration/cascade.md) for details.

## Mock Server

When `mock_enabled: true` is set at the config level, every collection must have `base_mock_url` set. This tells swag2mcp where the mock server is running for this collection.

```yaml
mock_enabled: true
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        base_mock_url: localhost:8080
```

See [Mock Server](../advanced/mock-server.md) for full details.

## Examples

### Minimal Collection

```yaml
specs:
  - domain: dadjokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

### Full Collection with All Fields

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        llm_instruction: "Use for current weather and daily forecasts."
        title: "Custom Title"
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        disable: false
        base_url: https://forecast-api.open-meteo.com
        base_mock_url: localhost:8080
        http_client:
          headers:
            X-Custom: value
```

### Multiple Collections per Spec

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

### Local File in specs/ Directory

```yaml
specs:
  - domain: myapi
    llm_title: My Internal API
    base_url: https://api.mycompany.com
    collections:
      - llm_title: Users
        location: specs/users.openapi.json
      - llm_title: Orders
        location: specs/orders.openapi.json
```

### Disabled Collection

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
        disable: true
```

## Related

- [Collection Settings (config)](../configuration/collection-settings.md) — full YAML reference
- [Configuration Cascade](../configuration/cascade.md) — how settings override each other
- [Specs](./specs) — logical containers for collections
- [HTTP Client](../configuration/http-client.md) — HTTP client configuration
- [Mock Server](../advanced/mock-server.md) — mock server setup
- [CLI: validate](../cli/validate.md) — validate command reference
- [CLI: update](../cli/update.md) — update command reference
- [CLI: clean](../cli/clean.md) — clean command reference
