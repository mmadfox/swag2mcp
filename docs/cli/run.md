# run

## Purpose

Launch the interactive **TUI (Terminal User Interface)** API explorer. This is a full-screen application for searching, browsing, inspecting, and invoking API endpoints without leaving the terminal.

## When to use

- You want to explore your APIs interactively
- You need to search for a specific endpoint across all specs
- You want to browse the spec → collection → tag → endpoint hierarchy
- You want to test an API call before configuring the MCP server

## Syntax

```bash
swag2mcp run [path]
```

## Arguments

| Argument | Position | Required | Description |
|----------|----------|----------|-------------|
| `path` | 1 | No | Workspace directory. If omitted, resolves via path resolution rules. |

## Flags

None.

## Modes

### Search mode

Full-text search across all endpoints in all specs. Supports filtering by HTTP method, tag, and path.

- Type a query to search endpoint names, paths, and descriptions
- Filter results by method (GET, POST, PUT, DELETE, etc.)
- View endpoint details with a single keystroke

### Browse mode

Tree navigation through the spec hierarchy:

```
Spec → Collection → Tag → Endpoint
```

- Navigate down the tree to find specific endpoints
- View endpoint details (parameters, request body, responses)
- Invoke the API directly from the TUI

## Navigation

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate up/down |
| `Enter` | Select or open |
| `Esc` | Go back |
| `Tab` | Switch between Search and Browse modes |
| `/` | Focus search input |
| `q` | Quit |

## Post-command verification

The TUI loads all specs from the workspace. If a spec fails to load, an error message is shown in the interface.

## Nuances

- **Auto-init:** If no config file exists, `run` automatically runs the init wizard first.
- **No flags:** The `run` command has no flags — all configuration comes from the workspace.
- **Terminal size:** The TUI requires a terminal with at least 80×24 characters. It may not render correctly in very small terminals.
- **Dependencies:** The TUI uses Bubbletea. It works over SSH and in most terminal emulators.
