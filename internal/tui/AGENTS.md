# TUI Package — Architecture Guide for LLM Agents

## Overview

The TUI package (`internal/tui/`) implements an interactive terminal UI using the [Bubbletea](https://github.com/charmbracelet/bubbletea) framework. It follows the Elm Architecture: **Model → Update → View**.

The main entry point is `RunExplorer(svc, ws)` which creates a `tea.Program` with a `runModel`.

---

## Core Concepts

### 1. Model (`runModel`)

The model holds all UI state. Key fields:

```go
type runModel struct {
    state runState      // current screen
    mode  runMode       // current navigation mode
    input textinput.Model  // text input field
    err   error         // error message to display
    msg   string        // success/info message to display
    // ... data fields for each screen
}
```

### 2. States (`runState`)

Each screen is a state. States are grouped by **mode**:

```
modeSearch:
  runSearchQuery → runSearchResults → runEndpointDetail

modeBrowse:
  runBrowseSpecs → runBrowseCollections → runBrowseTags → runBrowseEndpoints → runEndpointDetail
```

`runMenu` and `runDone` are shared across all modes.

### 3. Modes (`runMode`)

Modes isolate navigation chains. `[B]ack` from `runEndpointDetail` uses `m.mode` to determine where to return:

```go
case runEndpointDetail:
    switch m.mode {
    case modeSearch: m.state = runSearchResults
    case modeBrowse: m.state = runBrowseEndpoints
    }
```

`[M]enu` always returns to `runMenu` regardless of mode.

---

## How to Add a New Component

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

---

## Navigation Rules

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

---

## Input Handling

- Input is always focused in list/search states
- Digits are appended to input via `handleDigit`
- `Enter` triggers `handleEnter` which parses the input value
- In `runSearchResults`: `N`/`P` for pagination, numbers for endpoint selection
- In `runEndpointDetail`: input is blurred, `S` for save, `B` for back, `M` for menu

---

## View Pattern

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

---

## Key Methods Reference

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
