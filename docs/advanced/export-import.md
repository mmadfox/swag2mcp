# Export and Import

swag2mcp supports full workspace round-trip via ZIP.

## Export

```bash
# Export to default file
swag2mcp export

# Export with custom path
swag2mcp export --output ~/backups/swag2mcp-backup.zip
```

### What's Exported

- `swag2mcp.yaml` — configuration file
- `specs/` — all specs
- `cache/` — downloaded spec cache
- `auth_scripts/` — auth scripts

## Import

```bash
# Import from ZIP
swag2mcp import --from-zip backup.zip

# Import with overwrite
swag2mcp import --from-zip backup.zip -f
```

### What's Imported

- Configuration
- Specs
- Cache
- Auth scripts

## Use Cases

::: code-group

```bash [Backup]
swag2mcp export --output swag2mcp-$(date +%Y-%m-%d).zip
```

```bash [Transfer to another machine]
# On old machine
swag2mcp export --output swag2mcp.zip

# On new machine
swag2mcp import --from-zip swag2mcp.zip
```

```bash [Config template]
swag2mcp init
swag2mcp export --output template.zip
```

:::
