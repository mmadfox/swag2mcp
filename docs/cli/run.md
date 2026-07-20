# run

Start the interactive TUI (Terminal User Interface).

## Syntax

```bash
swag2mcp run [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `-w, --workspace` | Workspace path |

## Modes

### Search

- Full-text search across all endpoints
- Filter by method, tag, path
- Quick detail view

### Browse

- Tree navigation: Spec → Collection → Tag → Endpoint
- Endpoint detail view
- API invocation

## Navigation

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate |
| `Enter` | Select / open |
| `Esc` | Back |
| `Tab` | Switch modes |
| `/` | Search |
| `q` | Quit |
