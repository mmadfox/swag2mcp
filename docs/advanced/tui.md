# TUI Explorer

## Overview

swag2mcp includes a built-in TUI (Terminal User Interface) for interactive API exploration. It is a full-screen terminal application that lets you search, browse, inspect, and invoke API endpoints without leaving the terminal.

## Launch

```bash
swag2mcp run
```

If no config file exists, the TUI will automatically start the initialization wizard first.

## Modes

The TUI has three modes, switchable with the `Tab` key:

### Search mode

Full-text search across all endpoints in all specs. Supports the same query syntax as the `search` MCP tool.

- Type a query to search endpoint names, paths, and descriptions
- Filter results by method, tag, or path
- View endpoint details with a single keystroke
- Navigate through results with pagination (10 items per page)

### Browse mode

Tree navigation through the spec hierarchy:

```
Spec → Collection → Tag → Endpoint
```

- Navigate down the tree to find specific endpoints
- View endpoint details (parameters, request body, responses)
- Invoke the API directly from the TUI
- Save endpoint details as a JSON file

### Auth mode

View authentication tokens and headers for any spec. Useful for debugging or generating curl commands.

## Controls

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate up/down |
| `Enter` | Select or open |
| `Esc` | Go back one level |
| `Tab` | Switch between Search, Browse, and Auth modes |
| `/` | Focus search input |
| `N` / `P` | Next / previous page |
| `B` | Back to previous screen |
| `M` | Return to main menu |
| `S` | Save endpoint detail as JSON file |
| `q` / `Ctrl+C` | Quit |

## States

The TUI goes through these states as you navigate:

1. **Loading** — loading data from the workspace
2. **Search** — search mode with query input
3. **Browse** — browse mode with spec list
4. **Spec List** — list of all specs
5. **Collection List** — collections within a spec
6. **Tag List** — tags within a collection
7. **Endpoint List** — endpoints within a tag
8. **Endpoint Detail** — full endpoint information
9. **Invoke Result** — API call result
10. **Error** — error state with message

## Endpoint detail view

When you select an endpoint, the TUI shows:

- HTTP method and path
- Base URL and full URL
- Summary and description
- All parameters (name, location, type, required)
- Request body schema (if applicable)
- Response codes and schemas
- Deprecation status

## Requirements

- **Terminal size:** At least 80×24 characters
- **Terminal emulator:** Works in most modern terminals (iTerm2, Terminal.app, GNOME Terminal, Windows Terminal, etc.)
- **SSH:** Works over SSH connections

## Important notes

- **Auto-init** — if no config file exists, the TUI automatically starts the initialization wizard
- **Pagination** — lists are paginated at 10 items per page. Use `N` and `P` to navigate
- **Save endpoint details** — press `S` in the endpoint detail view to save the full detail as a JSON file in the current directory
- **Auth mode** — shows tokens and headers for debugging. In production, the auth tool can be disabled with `--disable-llm-auth`
