# Script Auth

## Purpose

Authentication via an external script — the most flexible method. You can write a script in any language (bash, Python, etc.) that obtains a token however you like and returns it to swag2mcp.

## When to use

- Custom or non-standard authentication schemes
- Complex token acquisition logic (multi-step, with additional checks)
- When none of the standard methods fit your needs

## Configuration

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: script
      config:
        domain: "my-auth"
```

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `domain` | Yes | Script file name (without extension) |

## Script location

The script must be placed in the `auth_scripts` directory of your workspace:

- **Linux / macOS:** `{workspace}/auth_scripts/{domain}.sh`
- **Windows:** `{workspace}/auth_scripts/{domain}.bat`

## Script output format

The script must output JSON to stdout with the token and its expiry time:

```bash
#!/bin/bash
# auth_scripts/my-auth.sh

TOKEN=$(curl -s -X POST https://auth.example.com/token \
  -d "grant_type=client_credentials" \
  -d "client_id=$CLIENT_ID" \
  -d "client_secret=$CLIENT_SECRET" | jq -r '.access_token')

echo "{\"token\": \"$TOKEN\", \"expires_in\": 3600}"
```

### JSON fields

| Field | Required | Description |
|-------|----------|-------------|
| `token` | Yes | Authentication token |
| `expires_in` | No | Token lifetime in seconds (default: 3600) |

## Notes

- swag2mcp runs the script on every request if the cached token has expired
- The script must complete within 30 seconds
- The token is cached until its expiry time
- Script filename = `{domain}.sh` (Unix) or `{domain}.bat` (Windows)
- `domain` must not contain `/` or `\`
