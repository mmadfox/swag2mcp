# validate

Validate configuration.

## Syntax

```bash
swag2mcp validate [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `-w, --workspace` | Workspace path |

## Usage

```bash
swag2mcp validate
```

## What's Checked

- YAML syntax
- Config structure
- Spec file existence
- Spec URL accessibility
- Spec format validity (OpenAPI/Swagger/Postman)
- Auth settings validity
- HTTP client correctness

## Example Output

```
✓ Configuration is valid
✓ Spec petstore: OK
✓ Spec meteo: OK
✗ Spec old-api: file not found
```
