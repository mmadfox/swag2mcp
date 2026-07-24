# version

## Purpose

Print the swag2mcp version. Useful for verifying the installed version, reporting bugs, or checking compatibility.

## When to use

- You want to check which version of swag2mcp is installed
- You are reporting a bug and need to include the version
- You want to verify a successful installation

## Syntax

```bash
swag2mcp version
swag2mcp --version
```

## Arguments

None.

## Flags

None.

## How it works

```bash
swag2mcp version
# swag2mcp v1.2.0

swag2mcp --version
# swag2mcp v1.2.0
```

## Output format

```
swag2mcp <version>
```

The version is set at build time via `ldflags`. If not set, it defaults to `"dev"`.

## Nuances

- **Two forms:** Both `swag2mcp version` (subcommand) and `swag2mcp --version` (global flag) produce the same output.
- **No config required:** This command works without a workspace or config file.
