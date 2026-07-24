# validate

## Purpose

Check the configuration file and all referenced spec files for errors. This is a **read-only** diagnostic command — it never modifies anything.

## When to use

- After editing `swag2mcp.yaml` manually
- Before running `mcp` or `update` to catch issues early
- When troubleshooting why a spec isn't loading
- In CI/CD pipelines to validate configuration changes

## Syntax

```bash
swag2mcp validate [path] [flags]
```

## Arguments

| Argument | Position | Required | Description |
|----------|----------|----------|-------------|
| `path` | 1 | No | Workspace directory. If omitted, resolves via path resolution rules. |

## Flags

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--tags` | `-t` | `string` | `""` | Validate only specs with matching tags (comma-separated) |

## How it works

```bash
swag2mcp validate
swag2mcp validate ./my-workspace
swag2mcp validate --tags=public
```

## What is checked

| Check | Description |
|-------|-------------|
| YAML syntax | The config file must be valid YAML |
| Config structure | All required fields present, types are correct |
| Domain uniqueness | No duplicate domains |
| Domain format | Lowercase, digits, hyphens only |
| Spec file existence | The `location` file or URL must be reachable |
| Spec format | The file must be valid OpenAPI 3.x, Swagger 2.0, or Postman collection |
| Auth settings | Auth type and config are valid for the selected method |
| HTTP client | HTTP client settings are valid |

## What is NOT checked

| Not checked | Reason |
|-------------|--------|
| Authentication endpoints | `validate` checks auth config syntax but does not test login/token exchange |
| API endpoint availability | Only the spec file URL is checked, not the `base_url` |
| `base_url` correctness | Format is validated, but no test request is made |
| Mock server configuration | `base_mock_url` is not verified for connectivity |

## Example output

```
✅ Configuration is valid.
✓ Spec petstore: OK
✓ Spec meteo: OK
✗ Spec old-api: file not found
```

## Post-command verification

If validation passes, the configuration is ready for `mcp`, `update`, or `run`.

## Nuances

- **No auto-init:** Unlike `add`, `ls`, or `run`, `validate` does **not** auto-initialize if the config is missing. It returns an error: `"configuration not found at <path>"`.
- **Network access:** Remote spec URLs are fetched during validation. The command may take longer if specs are hosted on slow servers.
- **Tag filtering:** When `--tags` is set, only specs matching the specified tags are validated. Other specs are skipped.
