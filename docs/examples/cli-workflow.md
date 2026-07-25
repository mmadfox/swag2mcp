# CLI Workflow

This page shows real examples of using swag2mcp from the terminal — from initialization to daily operations.

## Quick start

```bash
# 1. Initialize a workspace
mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp

# 2. List your specs
swag2mcp ls
```

## Adding a spec with YAML

### Simple spec (public API)

```bash
swag2mcp add spec --yaml - <<EOF
domain: meteo
llm_title: Open-Meteo Weather API
base_url: https://api.open-meteo.com
collections:
  - llm_title: Weather Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
EOF
```

### Spec with auth (bearer token from env)

```bash
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My Protected API
base_url: https://api.example.com/v1
auth:
  type: bearer
  config:
    token: \$(MY_TOKEN)
collections:
  - llm_title: Users
    location: https://raw.githubusercontent.com/my-org/my-api/main/users.yaml
EOF
```

### Spec with multiple collections

```bash
swag2mcp add spec --yaml - <<EOF
domain: meteo
llm_title: Open-Meteo APIs
base_url: https://api.open-meteo.com
collections:
  - llm_title: Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
  - llm_title: Air Quality
    base_url: https://air-quality-api.open-meteo.com
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
EOF
```

## Adding a collection to an existing spec

```bash
swag2mcp add collection --yaml - <<EOF
spec_domain: meteo
llm_title: Marine Weather
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
EOF
```

## Listing specs

```bash
$ swag2mcp ls
Specifications:
  dadjoke (https://icanhazdadjoke.com)
    jokes (3 endpoints)
  meteo (https://api.open-meteo.com)
    forecast (5 endpoints)
    air-quality (8 endpoints)
    marine (4 endpoints)
```

### Filter by tags

```bash
swag2mcp ls --tags=public
```

## Viewing runtime info

```bash
$ swag2mcp info
{
  "version": "v1.2.0",
  "workspace": "/home/user/.swag2mcp",
  "specs": {
    "total": 2,
    "active": 2,
    "disabled": 0,
    "collections": 4,
    "endpoints": 20
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "follow_redirects": true
  },
  "mcp": {
    "transport": "stdio"
  },
  "auth": {
    "methods": ["bearer"]
  }
}
```

## Validating configuration

```bash
$ swag2mcp validate
✅ Configuration is valid.
✓ Spec dadjoke: OK
✓ Spec meteo: OK
```

## Starting the MCP server

### stdio (for IDE integration)

```bash
swag2mcp mcp
```

### HTTP (for remote access)

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

### With tag filter

```bash
swag2mcp mcp --tags=public
```

## Updating specs

Refresh all cached spec files:

```bash
swag2mcp update
```

## Cleaning cache

```bash
swag2mcp clean
```

## Export and import

### Backup your workspace

```bash
swag2mcp export --output ~/backups/swag2mcp-2026-07-24.zip
```

### Restore on another machine

```bash
# On the new machine
swag2mcp import --from-zip swag2mcp-2026-07-24.zip
```

## Interactive TUI explorer

```bash
swag2mcp run
```

Opens a full-screen terminal UI for searching, browsing, and invoking APIs.

## Mock server

```bash
# Install the mock binary
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest

# Start mock servers
swag2mcp-mock mockserver
```
