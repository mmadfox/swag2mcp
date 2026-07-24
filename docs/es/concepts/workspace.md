# Workspace

The workspace is the directory where swag2mcp stores all its data — config, cached specs, local spec files, saved responses, and auth scripts.

## Structure

```
~/.swag2mcp/                          # Workspace root (default)
├── swag2mcp.yaml                     # Configuration file
├── cache/                            # Cached remote spec files
│   ├── a1b2c3d4e5f6...spec          # Cached spec content
│   └── a1b2c3d4e5f6...meta          # Cache metadata (JSON)
├── specs/                            # Local spec files
│   └── my-api.yaml
├── responses/                        # Saved API responses (large responses)
│   ├── meteo-get-forecast-abc123.json
│   └── response-fragment-def456.json
└── auth_scripts/                     # Authentication scripts
    ├── meteo.sh                      # Unix shell script
    └── meteo.bat                     # Windows batch script
```

## Default Path

- **Linux/macOS**: `~/.swag2mcp/`
- **Windows**: `%USERPROFILE%\.swag2mcp\`

## Custom Path

```bash
swag2mcp mcp /path/to/workspace
swag2mcp mcp ./my-workspace
```

## Directories

### cache/

Stores downloaded remote spec files. Each file is cached with a SHA-256 hash of its URL as the filename:

- `{hash}.spec` — the cached spec file content
- `{hash}.meta` — JSON metadata (source URL, cache time, TTL)

Each cached file has a random TTL between 1 hour and 48 hours. The cache is automatically checked on every startup — if a valid (non-expired) entry exists, it is reused without downloading.

**Commands:**
- `swag2mcp update` — clears cache and re-downloads all specs
- `swag2mcp clean` — clears cache and responses

### specs/

Stores local spec files that collections point to via `location: specs/{name}`. Files here are used directly without caching.

This directory is populated by:
- `swag2mcp import <source> <name>` — downloads a remote spec and saves it here
- `swag2mcp export` — copies specs here into the export ZIP
- Manual placement — you can copy spec files here yourself

### responses/

Stores API responses that exceed the `max_response_size` limit (default 1 MB). When the LLM invokes an endpoint and the response is too large, swag2mcp saves it here and returns a file reference instead.

Naming convention: `{domain}-{method}-{path_with_underscores}-{6char_hex}.json`

Old responses are cleaned automatically after 48 hours on MCP server start.

### auth_scripts/

Stores authentication scripts for the `script` auth type. Each script is named after the spec's domain.

#### Naming Convention

| Platform | Filename | Example |
|----------|----------|---------|
| Unix (Linux, macOS) | `{domain}.sh` | `meteo.sh` |
| Windows | `{domain}.bat` | `meteo.bat` |

The domain must not contain `/` or `\` characters.

#### How Scripts Work

1. swag2mcp runs the script with a 30-second timeout
2. The script must output valid JSON to stdout
3. swag2mcp parses the JSON and uses the token for API requests

#### Expected Output Format

```json
{
  "token": "your-token-here",
  "expires_in": 3600
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `token` | string | ✅ | The authentication token |
| `access_token` | string | ❌ | Alternative to `token` (checked first) |
| `token_type` | string | ❌ | Token type (e.g., "Bearer") |
| `expires_in` | number | ❌ | Token lifetime in seconds (default: 3600) |

#### Execution

| Platform | Command |
|----------|---------|
| Unix | `sh {domain}.sh` |
| Windows | `cmd /c {domain}.bat` |

#### Token Caching

The token is cached in memory until it expires. On each API call, swag2mcp checks the cache first — the script is only executed when the cached token has expired.

#### Stub Creation

When you configure `auth: { type: script, config: { domain: "myapi" } }`, swag2mcp creates a stub script automatically:

**Unix (`auth_scripts/myapi.sh`):**
```bash
#!/bin/sh
echo '{"token": "your-token-here", "expires_in": 3600}'
```

**Windows (`auth_scripts/myapi.bat`):**
```bat
@echo off
echo {"token": "your-token-here", "expires_in": 3600}
```

Replace the placeholder token with your actual authentication logic.

#### Orphan Cleanup

When you delete a spec, its auth script becomes orphaned. swag2mcp automatically removes orphan scripts on:
- `swag2mcp update`
- `swag2mcp clean`

## Commands

### update

```bash
swag2mcp update [path]
```

Validates the config, clears the cache and responses, then re-downloads all spec files. Also ensures auth scripts exist and removes orphan scripts.

Use this command after:
- Adding or removing collections
- Changing collection locations
- Editing spec files that need re-caching

### clean

```bash
swag2mcp clean [path]
```

Removes all contents of `cache/` and `responses/`, plus orphan auth scripts. Does NOT re-cache specs — use `update` for that.

### validate

```bash
swag2mcp validate [path]
```

Validates the config including all collection locations. See [CLI: validate](../cli/validate.md).

## Export and Import

```bash
# Export workspace to ZIP (default name: swag2mcp-backup-{date}.zip)
swag2mcp export

# Export to a specific path
swag2mcp export /path/to/workspace /path/to/backup.zip

# Export only specific specs
swag2mcp export --spec meteo

# Restore from backup
swag2mcp import --from-zip /path/to/backup.zip
swag2mcp import /path/to/workspace /path/to/backup.zip
```

Export includes: `swag2mcp.yaml`, `specs/`, `auth_scripts/`. Cache and responses are excluded (they are local data).

## .gitignore

If your workspace is inside a Git repository, add these entries to `.gitignore`:

```gitignore
# swag2mcp — local data only
.swag2mcp/cache/
.swag2mcp/responses/
```

The `cache/` and `responses/` directories contain local, machine-specific data that should not be committed. Everything else (`swag2mcp.yaml`, `specs/`, `auth_scripts/`) should be in the repository so the configuration is shared across the team.
