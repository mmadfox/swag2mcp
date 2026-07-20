# clean

Clear cache.

## Syntax

```bash
swag2mcp clean [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `-w, --workspace` | Workspace path |

## Usage

```bash
swag2mcp clean
```

## What's Cleaned

- `cache/` directory — downloaded spec cache
- `responses/` directory — saved API responses

## Automatic Cleanup

On `swag2mcp mcp` start, responses older than 48 hours are removed automatically.
