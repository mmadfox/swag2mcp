# init

Initialize a workspace.

## Syntax

```bash
swag2mcp init [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `-i, --interactive` | Interactive mode (TUI) |
| `-f, --force` | Overwrite existing config |

## Usage

::: code-group

```bash [Default]
swag2mcp init
```
Creates `~/.swag2mcp/swag2mcp.yaml` with minimal config.

```bash [Interactive]
swag2mcp init -i
```
18-step TUI wizard:
1. Choose directory
2. Add specs
3. Configure collections
4. Configure auth
5. Configure HTTP client

```bash [Force]
swag2mcp init -f
```
Overwrites existing config.

:::

## Result

```
~/.swag2mcp/
├── swag2mcp.yaml
├── cache/
├── specs/
├── responses/
└── auth_scripts/
```
