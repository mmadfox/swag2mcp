# Script Auth

Authentication via external script.

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
        domain: "auth.example.com"
```

## How It Works

1. swag2mcp runs the specified script from `{workspace}/auth_scripts/`
2. Script outputs auth headers to stdout
3. Output is parsed and added to the request

## Script Output Format

Script must output headers in `Key: Value` format:

```bash
#!/bin/bash
# get-token.sh
TOKEN=$(curl -s -X POST https://auth.example.com/token \
  -d "grant_type=client_credentials" \
  -d "client_id=$CLIENT_ID" \
  -d "client_secret=$CLIENT_SECRET" | jq -r '.access_token')

echo "Authorization: Bearer $TOKEN"
```

## Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `domain` | string | Domain identifier for the script |

## Script Location

Scripts are stored in `{workspace}/auth_scripts/` and named by domain.
