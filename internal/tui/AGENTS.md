# TUI Package — Architecture Guide for LLM Agents

## Overview

The TUI package (`internal/tui/`) implements all user-facing interactive and non-interactive interfaces for swag2mcp. It uses two UI paradigms:

- **Bubbletea TUI** (`github.com/charmbracelet/bubbletea`) — full-screen terminal UI for the Explorer and wizards
- **CLI prompts** (`fmt.Scanln`, `fmt.Scanf`) — simple terminal prompts for add/delete operations

---

## File Structure

| File | Purpose |
|------|---------|
| `run.go` | Interactive API Explorer (Bubbletea TUI) — search, browse, inspect, save |
| `wizard.go` | Initialization wizard (Bubbletea TUI) — `swag2mcp init -i` |
| `collect.go` | Sub-wizards for collecting spec/collection data (Bubbletea TUI) |
| `add.go` | Add spec/collection — YAML import and interactive prompts |
| `delete.go` | Delete spec/collection — interactive prompts |
| `ls.go` | List config — formatted table output |
| `initmcp.go` | Setup, WriteConfig, ExampleConfig — embedded files |
| `atomic.go` | Atomic YAML config writer |
| `config.tmpl` | Go template for config YAML (embedded) |
| `init.swag2mcp.yaml` | Example config file (embedded) |
| `wizard_test.go` | Tests for wizard and config building |
| `AGENTS.md` | This file |

---

## 1. Run Explorer (`run.go`)

Full-screen Bubbletea TUI for interactive API exploration. Entry point: `RunExplorer(svc, ws)`.

### Architecture

Follows Elm Architecture: **Model → Update → View**.

```go
type runModel struct {
    state runState      // current screen
    mode  runMode       // current navigation mode (modeSearch or modeBrowse)
    input textinput.Model
    err   error
    msg   string
    // ... data fields for each screen
}
```

### States and Modes

```
runMenu (shared)
  │
  ├── [1] modeSearch:
  │     runSearchQuery → runSearchResults → runEndpointDetail
  │
  └── [2] modeBrowse:
        runBrowseSpecs → runBrowseCollections → runBrowseTags → runBrowseEndpoints → runEndpointDetail
```

Modes isolate navigation chains. `[B]ack` from `runEndpointDetail` uses `m.mode` to determine where to return.

### Navigation Rules

| Action | Behavior |
|--------|----------|
| `[B]ack` | One level up within current mode |
| `[M]enu` | Always returns to `runMenu` |
| `Esc/Ctrl+C` | Quit from anywhere |
| `Enter` | Confirms input (search query, number selection, N/P for pages) |
| Digits | Appended to input field for number selection |

### `handleBack` Chain

```
modeSearch:
  runEndpointDetail → runSearchResults → runSearchQuery → (stays)

modeBrowse:
  runEndpointDetail → runBrowseEndpoints → runBrowseTags → runBrowseCollections → runBrowseSpecs → runMenu
```

### Input Handling

- Input is always focused in list/search states
- Digits are appended to input via `handleDigit`
- `Enter` triggers `handleEnter` which parses the input value
- In `runSearchResults`: `N`/`P` for pagination, numbers for endpoint selection
- In `runEndpointDetail`: input is blurred, `S` for save, `B` for back, `M` for menu

### View Pattern

Every state in `View()` follows the same pattern:

1. **Header**: `s += "  Screen Title\n"`
2. **Divider**: `s += "  ──────────\n\n"`
3. **Data**: loop and format items
4. **Input**: `s += "\n  " + m.input.View() + "\n\n"`
5. **Actions**: `s += "  [A]ction  [B]ack  [M]enu\n"`

### Error/Message Display

```go
if m.err != nil {
    s += fmt.Sprintf("  ❌ Error: %s\n\n", m.err)
}
if m.msg != "" {
    s += fmt.Sprintf("  %s\n\n", m.msg)
}
```

### Run Explorer Methods Reference

| Method | Purpose |
|--------|---------|
| `handleEnter()` | Process Enter key — search, select, paginate |
| `handleBack()` | Navigate one level up within current mode |
| `handleMenu()` | Return to main menu |
| `handleDigit(digit)` | Append digit to input in list states |
| `loadSpecs()` | Fetch and display specs list |
| `loadCollections(specID)` | Fetch collections for a spec |
| `loadTags(collectionID)` | Fetch tags for a collection |
| `loadEndpoints(tagID)` | Fetch endpoints for a tag |
| `loadEndpointDetail(id)` | Fetch and display endpoint details |
| `doSearch(query)` | Execute search and show results |
| `showEndpoint()` | Save endpoint detail as JSON file |
| `selectSearchResult(val)` | Select endpoint from search results |
| `selectSpec(val)` | Select spec from list |
| `selectCollection(val)` | Select collection from list |
| `selectTag(val)` | Select tag from list |
| `selectBrowseEndpoint(val)` | Select endpoint from browse list |

---

## 2. Initialization Wizard (`wizard.go`)

Full-screen Bubbletea TUI for `swag2mcp init -i`. Entry point: `RunTUI()`.

### States (18 total)

```
stateWorkspaceDir → stateConfigPath → stateAskAddSpec
  → stateSpecDomain → stateSpecTitle → stateSpecInstruction → stateSpecBaseURL → stateSpecTags
  → stateAuthType → stateAuthField → stateAskAddCollection
  → stateCollTitle → stateCollLocation → stateAskAddAnotherSpec
  → stateConfirm → stateDone
```

### Key Functions

| Function | Purpose |
|----------|---------|
| `RunTUI()` | Start the wizard, return config path, workspace dir, specs |
| `BuildConfigYAML(specs)` | Render config YAML from collected data |
| `WriteResult(configPath, workspaceDir, specs)` | Write config file and init workspace |
| `authMethodsList()` | Format available auth methods for display |
| `authFieldsFor(authType)` | Get field definitions for an auth type |

### Auth Methods

8 methods supported: `none`, `basic`, `bearer`, `digest`, `api-key`, `oauth2-cc`, `oauth2-pwd`, `script`. Each has configurable fields, some support `$(ENV_VAR)` syntax.

---

## 3. Collection Sub-Wizards (`collect.go`)

Bubbletea TUI for collecting a single spec or collection. Used by `add.go` for interactive mode.

### Models

- `collectModel` — collects a full spec (domain, title, instruction, baseURL, tags, auth, collections)
- `collectCollectionModel` — collects a single collection (title, location)

### Key Functions

| Function | Purpose |
|----------|---------|
| `collectSpec(specNum)` | Run TUI to collect a spec, return `SpecInput` |
| `collectCollection(specNum, collNum, domain)` | Run TUI to collect a collection, return `CollectionInput` |

### States

```
colDomain → colTitle → colInstruction → colBaseURL → colTags
  → colAuthType → colAuthField → colAskAddCollection
  → colCollTitle → colCollLocation → colDone
```

---

## 4. Add Operations (`add.go`)

Non-interactive and interactive functions for adding specs/collections to config.

### Key Functions

| Function | Purpose |
|----------|---------|
| `AddSpecFromYAML(configPath, data)` | Add spec from YAML string (non-interactive) |
| `AddCollectionFromYAML(configPath, data)` | Add collection from YAML string (non-interactive) |
| `AddSpecTUI(configPath)` | Interactive wizard to add a spec |
| `AddCollectionTUI(configPath)` | Interactive wizard to add a collection |

### YAML Input Format

```yaml
# add spec
domain: petstore
llm_title: Petstore API
llm_instruction: Use this API to manage pets.
base_url: https://petstore.swagger.io/v2
tags: [public, demo]
auth:
  type: bearer
  config:
    token: $(TOKEN)
collections:
  - llm_title: Petstore Swagger
    location: https://petstore.swagger.io/v2/swagger.json

# add collection
spec_domain: petstore
llm_title: Orders Collection
location: https://petstore.example.com/orders.json
```

---

## 5. Delete Operations (`delete.go`)

Interactive functions for deleting specs/collections from config.

### Key Functions

| Function | Purpose |
|----------|---------|
| `DeleteSpecTUI(configPath)` | Interactive wizard to delete a spec |
| `DeleteCollectionTUI(configPath)` | Interactive wizard to delete a collection |

Both functions:
1. Load config
2. List items with numbers
3. Prompt for selection
4. Ask for confirmation
5. Use `AtomicWriteConfig` to safely update

---

## 6. List Config (`ls.go`)

Formatted table output for `swag2mcp ls`.

### Key Function

| Function | Purpose |
|----------|---------|
| `ListConfig(configPath, tags)` | Load config, filter by tags, return formatted string |

Uses `text/tabwriter` for aligned columns. Shows domain, title, baseURL, tags, auth type, and collections.

---

## 7. Init MCP (`initmcp.go`)

Non-interactive workspace initialization and embedded file access.

### Key Functions

| Function | Purpose |
|----------|---------|
| `Setup(configPath, workspaceDir)` | Create workspace dirs and write example config |
| `ExampleConfig()` | Return embedded example config YAML |
| `WriteConfig(configPath, data)` | Write YAML data to config file |

### Embedded Files

| File | Content |
|------|---------|
| `init.swag2mcp.yaml` | Example config with 4 sample specs (train-booking, petfood-shop, music-stream, cinema-tickets) |
| `config.tmpl` | Go template for config YAML generation |

---

## 8. Atomic Config Writer (`atomic.go`)

Safe YAML config updater that prevents file corruption.

### Key Function

| Function | Purpose |
|----------|---------|
| `AtomicWriteConfig(configPath, fn)` | Read config, apply mutation, validate, write to `.tmp`, atomically rename |

**Flow:**
1. Load config from `configPath`
2. Call `fn(cfg)` to mutate the config
3. Validate the result
4. Marshal to YAML
5. Write to `configPath.tmp`
6. `os.Rename(configPath.tmp, configPath)` — atomic on Unix

---

## How to Add a New Component to the Explorer

### Step 1: Add a new mode constant

```go
const (
    modeSearch runMode = iota
    modeBrowse
    modeFavorites  // new
)
```

### Step 2: Add new states (if needed)

```go
const (
    runMenu runState = iota
    // ... existing states
    runFavoritesList
    runFavoriteDetail
)
```

### Step 3: Add data fields to `runModel`

```go
type runModel struct {
    // ... existing fields
    favorites []service.FavoriteItem
}
```

### Step 4: Add navigation in `Update`

```go
case "4":  // new menu option
    if m.state == runMenu {
        m.mode = modeFavorites
        return m.loadFavorites()
    }
```

### Step 5: Add `load*` and `select*` methods

```go
func (m runModel) loadFavorites() (tea.Model, tea.Cmd) {
    data, err := m.svc.Favorites(...)
    m.favorites = data
    m.input.SetValue("")
    m.input.Focus()
    m.state = runFavoritesList
    return m, nil
}
```

### Step 6: Add `handleBack` case for the new mode

```go
case runFavoriteDetail:
    switch m.mode {
    case modeFavorites:
        m.state = runFavoritesList
    }
    m.input.SetValue("")
    m.input.Focus()
    return m, textinput.Blink
```

### Step 7: Add `View` rendering

```go
case runFavoritesList:
    for i, f := range m.favorites {
        s += fmt.Sprintf("  %d. %s\n", i+1, f.Name)
    }
    s += "\n  " + m.input.View() + "\n\n"
    s += "  Enter number and press Enter.  [B]ack  [M]enu.\n"
```
