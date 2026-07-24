# add

## Purpose

Add a new **spec** (API service) or **collection** (OpenAPI/Swagger/Postman file) to an existing configuration. This is the primary way to grow your workspace with new APIs.

## When to use

- You have a new API to connect to your LLM agent
- You found an OpenAPI spec URL and want to add it
- You want to add an additional spec file (collection) to an existing spec
- You prefer writing YAML directly instead of using the interactive wizard

## Syntax

```bash
swag2mcp add spec [path] [flags]
swag2mcp add collection [path] [flags]
```

## Arguments

| Argument | Position | Required | Description |
|----------|----------|----------|-------------|
| `path` | 1 | No | Workspace directory. If omitted, resolves via path resolution rules. |

## Flags

### `add spec`

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--yaml` | `-y` | `string` | `""` | YAML input inline or `-` for stdin |
| `--example` | `-e` | `bool` | `false` | Print a YAML template and exit |

### `add collection`

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--yaml` | `-y` | `string` | `""` | YAML input inline or `-` for stdin |
| `--example` | `-e` | `bool` | `false` | Print a YAML template and exit |

## How it works

### Interactive mode (default)

Launches a TUI wizard that lets you fill in the spec or collection fields step by step.

```bash
swag2mcp add spec
swag2mcp add collection
```

### YAML inline mode

Pass the YAML directly as a string. **Be careful with shell quoting** — special characters like `:`, `#`, `&`, `{` can break the command.

```bash
swag2mcp add spec --yaml 'domain: meteo
llm_title: Open-Meteo API
base_url: https://meteo.swagger.io/v2
collections:
  - llm_title: Main
    location: https://example.com/spec.json'
```

### YAML from stdin (recommended for complex YAML)

Pipe from a file or use a heredoc to avoid shell quoting issues entirely:

```bash
# Pipe from file
cat spec.yaml | swag2mcp add spec --yaml -

# Heredoc
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My API
llm_instruction: "Use this API for X & Y # important"
base_url: https://api.example.com/v1
collections:
  - llm_title: Main
    location: https://raw.githubusercontent.com/org/repo/main/spec.yaml
EOF
```

### YAML template

Print the expected YAML structure and exit:

```bash
swag2mcp add spec --example
swag2mcp add collection --example
```

## YAML format

### Spec

```yaml
domain: meteo
llm_title: Open-Meteo API
llm_instruction: Use this API to manage pets.
base_url: https://meteo.swagger.io/v2
tags: [public, demo]
auth:
  type: bearer
  config:
    token: $(TOKEN)
collections:
  - llm_title: Open-Meteo Swagger
    location: https://example.com/spec.json
```

### Collection

```yaml
spec_domain: meteo
llm_title: Orders Collection
location: https://example.com/orders.json
```

## Post-command verification

```bash
swag2mcp ls [path]
# The new spec or collection should appear in the list
```

## Nuances

- **Auto-init:** If no config file exists, `add` automatically runs the init wizard first. You don't need to run `init` separately.
- **Shell quoting:** Inline YAML (`--yaml '...'`) is fragile with special characters. Prefer `--yaml -` with a heredoc or pipe for anything beyond simple values.
- **`--example` exits immediately** without checking for an existing config or modifying anything.
- **`add spec` vs `add collection`:** Use `add spec` for a new API service (new domain). Use `add collection` to add another spec file to an existing spec.
