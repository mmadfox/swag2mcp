# Script Auth

Authentication via external script.

## Configuration

```yaml
specs:
  - domain: "api.example.com"
    location: "https://api.example.com/openapi.json"
    auth:
      type: script
      script:
        path: "./auth_scripts/get-token.sh"
        args: ["--env", "production"]
        timeout: 10s
```

## How It Works

1. swag2mcp runs the specified script
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
| `path` | string | Script path |
| `args` | array | Script arguments |
| `timeout` | duration | Execution timeout (default 30s) |

## Script Location

Scripts are stored in `{workspace}/auth_scripts/`.
