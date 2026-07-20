# TUI Explorer

swag2mcp has a built-in TUI (Terminal User Interface) for interactive API exploration.

## Launch

```bash
swag2mcp run
```

## Modes

### Search

- Full-text search across all endpoints
- Filter by method, tag, path
- Quick detail view

### Browse

- Tree navigation
- Spec → Collection → Tag → Endpoint
- Endpoint detail view
- API invocation

## Controls

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate list |
| `Enter` | Select / open |
| `Esc` | Back |
| `Tab` | Switch modes |
| `/` | Focus search |
| `q` / `Ctrl+C` | Quit |

## States

The TUI goes through these states:

1. **Loading** — loading data
2. **Search** — search mode
3. **Browse** — browse mode
4. **SpecList** — spec list
5. **CollectionList** — collection list
6. **TagList** — tag list
7. **EndpointList** — endpoint list
8. **EndpointDetail** — endpoint details
9. **InvokeResult** — API call result
10. **Error** — error state
