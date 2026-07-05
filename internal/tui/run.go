package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/mmadfox/swag2mcp/internal/service"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

type runState int

const (
	runMenu runState = iota
	runSearchQuery
	runSearchResults
	runBrowseSpecs
	runBrowseCollections
	runBrowseTags
	runBrowseEndpoints
	runEndpointDetail
	runDone
)

const randSuffixLen = 6

type runModel struct {
	state         runState
	svc           *service.Service
	ws            *workspace.Workspace
	input         textinput.Model
	err           error
	msg           string
	width         int
	searchResults []service.EndpointSearchItem
	specs         []service.SpecItem
	collections   []service.CollectionItem
	tags          []service.TagListItem
	endpoints     []service.EndpointTagItem
	selectedSpec  service.SpecItem
	selectedColl  service.CollectionItem
	selectedTag   service.TagListItem
	selectedEp    service.EndpointTagItem
	selectedEpID  string
	epDetail      *service.InspectResponse
}

func newRunModel(svc *service.Service, ws *workspace.Workspace) runModel {
	ti := textinput.New()
	ti.CharLimit = 256
	ti.Width = 60
	ti.Focus()
	return runModel{
		state: runMenu,
		svc:   svc,
		ws:    ws,
		input: ti,
	}
}

func (m runModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m runModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			return m.handleEnter()
		case "b", "B":
			return m.handleBack()
		case "M":
			return m.handleMenu()
		case "s", "S":
			if m.state == runEndpointDetail {
				return m.showEndpoint()
			}
		case "1":
			if m.state == runMenu {
				m.input.SetValue("")
				m.input.Placeholder = "Search endpoints..."
				m.state = runSearchQuery
				m.input.Focus()
				return m, textinput.Blink
			}
		case "2":
			if m.state == runMenu {
				return m.loadSpecs()
			}
		case "3":
			if m.state == runMenu {
				m.state = runDone
				return m, tea.Quit
			}
		}
	}

	if m.state == runEndpointDetail || m.state == runMenu || m.state == runDone || !m.input.Focused() {
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m runModel) handleEnter() (tea.Model, tea.Cmd) {
	val := m.input.Value()

	switch m.state {
	case runSearchQuery:
		if val == "" {
			return m, nil
		}
		return m.doSearch(val)

	case runSearchResults:
		return m.selectSearchResult(val)

	case runBrowseSpecs:
		return m.selectSpec(val)

	case runBrowseCollections:
		return m.selectCollection(val)

	case runBrowseTags:
		return m.selectTag(val)

	case runBrowseEndpoints:
		return m.selectBrowseEndpoint(val)
	}

	return m, nil
}

func (m runModel) handleBack() (tea.Model, tea.Cmd) {
	m.msg = ""
	switch m.state {
	case runSearchQuery:
		return m, nil
	case runSearchResults:
		m.input.SetValue("")
		m.state = runSearchQuery
		return m, nil
	case runBrowseSpecs:
		m.state = runMenu
		return m, nil
	case runBrowseCollections:
		return m.loadSpecs()
	case runBrowseTags:
		return m.loadCollections(m.selectedSpec.ID)
	case runBrowseEndpoints:
		return m.loadTags(m.selectedColl.ID)
	case runEndpointDetail:
		m.state = runBrowseEndpoints
		m.input.Focus()
		return m, textinput.Blink
	}
	return m, nil
}

func (m runModel) handleMenu() (tea.Model, tea.Cmd) {
	m.msg = ""
	m.input.SetValue("")
	m.state = runMenu
	return m, nil
}

func (m runModel) doSearch(query string) (tea.Model, tea.Cmd) {
	results, err := m.svc.Search(context.Background(), service.SearchRequest{
		Query: query,
		Limit: 50, //nolint:mnd // max search results
	})
	if err != nil {
		m.err = err
		return m, nil
	}
	m.searchResults = results.Endpoints
	m.input.SetValue("")
	m.input.Focus()
	m.state = runSearchResults
	return m, nil
}

func (m runModel) selectSearchResult(val string) (tea.Model, tea.Cmd) {
	idx := 0
	if _, err := fmt.Sscanf(val, "%d", &idx); err != nil || idx < 1 || idx > len(m.searchResults) {
		return m, nil
	}
	ep := m.searchResults[idx-1]
	m.selectedEpID = ep.ID
	return m.loadEndpointDetail(ep.ID)
}

func (m runModel) loadSpecs() (tea.Model, tea.Cmd) {
	specs, err := m.svc.Specs(context.Background())
	if err != nil {
		m.err = err
		return m, nil
	}
	m.specs = specs.Specs
	m.input.SetValue("")
	m.input.Focus()
	m.state = runBrowseSpecs
	return m, nil
}

func (m runModel) selectSpec(val string) (tea.Model, tea.Cmd) {
	idx := 0
	if _, err := fmt.Sscanf(val, "%d", &idx); err != nil || idx < 1 || idx > len(m.specs) {
		return m, nil
	}
	m.selectedSpec = m.specs[idx-1]
	return m.loadCollections(m.selectedSpec.ID)
}

func (m runModel) loadCollections(specID string) (tea.Model, tea.Cmd) {
	collections, err := m.svc.CollectionsBySpec(context.Background(), service.CollectionsRequest{SpecID: specID})
	if err != nil {
		m.err = err
		return m, nil
	}
	m.collections = collections.Collections
	m.input.SetValue("")
	m.input.Focus()
	m.state = runBrowseCollections
	return m, nil
}

func (m runModel) selectCollection(val string) (tea.Model, tea.Cmd) {
	idx := 0
	if _, err := fmt.Sscanf(val, "%d", &idx); err != nil || idx < 1 || idx > len(m.collections) {
		return m, nil
	}
	m.selectedColl = m.collections[idx-1]
	return m.loadTags(m.selectedColl.ID)
}

func (m runModel) loadTags(collectionID string) (tea.Model, tea.Cmd) {
	tags, err := m.svc.TagsByCollection(
		context.Background(),
		service.TagsByCollectionRequest{CollectionID: collectionID},
	)
	if err != nil {
		m.err = err
		return m, nil
	}
	m.tags = tags.Tags
	m.input.SetValue("")
	m.input.Focus()
	m.state = runBrowseTags
	return m, nil
}

func (m runModel) selectTag(val string) (tea.Model, tea.Cmd) {
	idx := 0
	if _, err := fmt.Sscanf(val, "%d", &idx); err != nil || idx < 1 || idx > len(m.tags) {
		return m, nil
	}
	m.selectedTag = m.tags[idx-1]
	return m.loadEndpoints(m.selectedTag.ID)
}

func (m runModel) loadEndpoints(tagID string) (tea.Model, tea.Cmd) {
	endpoints, err := m.svc.EndpointsByTag(context.Background(), service.EndpointsByTagRequest{TagID: tagID})
	if err != nil {
		m.err = err
		return m, nil
	}
	m.endpoints = endpoints.Endpoints
	m.input.SetValue("")
	m.input.Focus()
	m.state = runBrowseEndpoints
	return m, nil
}

func (m runModel) selectBrowseEndpoint(val string) (tea.Model, tea.Cmd) {
	idx := 0
	if _, err := fmt.Sscanf(val, "%d", &idx); err != nil || idx < 1 || idx > len(m.endpoints) {
		return m, nil
	}
	ep := m.endpoints[idx-1]
	m.selectedEp = ep
	m.selectedEpID = ep.ID
	return m.loadEndpointDetail(ep.ID)
}

func (m runModel) loadEndpointDetail(endpointID string) (tea.Model, tea.Cmd) {
	detail, err := m.svc.Inspect(context.Background(), service.InspectRequest{EndpointID: endpointID})
	if err != nil {
		m.err = err
		return m, nil
	}
	m.epDetail = &detail
	m.msg = ""
	m.state = runEndpointDetail
	m.input.Blur()
	return m, nil
}

func (m runModel) showEndpoint() (tea.Model, tea.Cmd) {
	if m.epDetail == nil {
		return m, nil
	}
	data, err := json.MarshalIndent(m.epDetail, "", "  ")
	if err != nil {
		m.err = fmt.Errorf("marshal endpoint: %w", err)
		return m, nil
	}

	method := strings.ToLower(m.epDetail.Method)
	path := strings.TrimPrefix(m.epDetail.Path, "/")
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, "{", "")
	path = strings.ReplaceAll(path, "}", "")
	randSuffix := make([]byte, randSuffixLen)
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	for i := range randSuffix {
		randSuffix[i] = letters[rand.IntN(len(letters))]
	}
	filename := fmt.Sprintf("%s-%s-%s-%s.json", m.epDetail.SpecDomain, method, path, string(randSuffix))

	if err := os.WriteFile(filename, data, 0600); err != nil {
		m.err = fmt.Errorf("save endpoint: %w", err)
		return m, nil
	}

	absPath, _ := filepath.Abs(filename)
	m.msg = fmt.Sprintf("✅ Saved to: %s", absPath)
	return m, nil
}

func (m runModel) View() string {
	var s string

	s += "\n  ╭──────────────────────────────────────────────╮\n"
	s += "  │              swag2mcp — Explorer              │\n"
	s += "  ╰──────────────────────────────────────────────╯\n\n"

	if m.err != nil {
		s += fmt.Sprintf("  ❌ Error: %s\n\n", m.err)
	}

	switch m.state {
	case runMenu:
		s += "  What would you like to do?\n"
		s += "  ──────────────────────────\n\n"
		s += "  1. 🔍  Search endpoints\n"
		s += "  2. 📂  Browse by specification\n"
		s += "  3. ❌  Exit\n\n"
		s += "  Press 1, 2, or 3.  (Esc/Ctrl+C to exit)\n"

	case runSearchQuery:
		s += "  Search endpoints\n"
		s += "  ────────────────\n\n"
		s += "  Search fields like method:GET, tag:auth, or path:/users\n"
		s += "  using + for required matches and - to exclude.\n"
		s += "  Use wildcards (*), exact phrases (\"...\"), or fuzzy matching (~)\n"
		s += "  for advanced filtering.\n\n"
		s += "  Examples:\n"
		s += "    +method:POST +tag:user | path:*/v2/* -deprecated\n"
		s += "    summary:\"login\"~\n\n"
		s += "  Enter a search query:\n\n"
		s += "  " + m.input.View() + "\n\n"
		s += "  Press Enter to search, [B]ack  [M]enu.\n"

	case runSearchResults:
		s += fmt.Sprintf("  Search results (%d):\n", len(m.searchResults))
		s += "  ──────────────────────\n\n"
		for i, ep := range m.searchResults {
			s += fmt.Sprintf("  %d. %-6s %s\n", i+1, ep.Method, ep.Path)
			s += fmt.Sprintf("     %s (%s)\n", ep.Summary, ep.SpecDomain)
		}
		s += "\n  " + m.input.View() + "\n\n"
		s += "  Enter number and press Enter.  [B]ack  [M]enu.\n"

	case runBrowseSpecs:
		s += "  Specifications:\n"
		s += "  ────────────────\n\n"
		for i, sp := range m.specs {
			s += fmt.Sprintf("  %d. %s\n", i+1, sp.Domain)
		}
		s += "\n  " + m.input.View() + "\n\n"
		s += "  Enter number and press Enter.  [B]ack  [M]enu.\n"

	case runBrowseCollections:
		s += fmt.Sprintf("  Collections for \"%s\":\n", m.selectedSpec.Domain)
		s += "  ──────────────────────────────\n\n"
		for i, col := range m.collections {
			s += fmt.Sprintf("  %d. %s (%d tags, %d methods)\n", i+1, col.Title, col.CountTags, col.CountMethods)
		}
		s += "\n  " + m.input.View() + "\n\n"
		s += "  Enter number and press Enter.  [B]ack  [M]enu.\n"

	case runBrowseTags:
		s += fmt.Sprintf("  Tags for \"%s\":\n", m.selectedColl.Title)
		s += "  ────────────────────────\n\n"
		for i, tag := range m.tags {
			s += fmt.Sprintf("  %d. %s (%d methods)\n", i+1, tag.Title, tag.CountMethods)
		}
		s += "\n  " + m.input.View() + "\n\n"
		s += "  Enter number and press Enter.  [B]ack  [M]enu.\n"

	case runBrowseEndpoints:
		s += fmt.Sprintf("  Endpoints for tag \"%s\":\n", m.selectedTag.Title)
		s += "  ─────────────────────────────\n\n"
		for i, ep := range m.endpoints {
			s += fmt.Sprintf("  %d. %-6s %s\n", i+1, ep.Method, ep.Path)
			s += fmt.Sprintf("     %s\n", ep.Summary)
		}
		s += "\n  " + m.input.View() + "\n\n"
		s += "  Enter number and press Enter.  [B]ack  [M]enu.\n"

	case runEndpointDetail:
		if m.epDetail != nil {
			s += "  Endpoint Details\n"
			s += "  ────────────────\n\n"
			s += fmt.Sprintf("  Spec:    %s\n", m.epDetail.SpecDomain)
			s += fmt.Sprintf("  Method:  %s\n", m.epDetail.Method)
			s += fmt.Sprintf("  Path:    %s\n", m.epDetail.Path)
			s += fmt.Sprintf("  BaseURL: %s\n\n", m.epDetail.BaseURL)

			if m.epDetail.Operation != nil {
				if m.epDetail.Operation.Summary != "" {
					s += fmt.Sprintf("  Summary: %s\n\n", m.epDetail.Operation.Summary)
				}
				if m.epDetail.Operation.Description != "" {
					s += fmt.Sprintf("  Description: %s\n\n", m.epDetail.Operation.Description)
				}

				if len(m.epDetail.Operation.Parameters) > 0 {
					s += "  Parameters:\n"
					for _, p := range m.epDetail.Operation.Parameters {
						req := ""
						if p.Required {
							req = " (required)"
						}
						s += fmt.Sprintf("    %s (%s%s) — %s\n", p.Name, p.In, req, p.Description)
					}
					s += "\n"
				}

				if m.epDetail.Operation.RequestBody != nil {
					s += "  Request Body:\n"
					if m.epDetail.Operation.RequestBody.Required {
						s += "    Required\n"
					}
					if m.epDetail.Operation.RequestBody.Description != "" {
						s += fmt.Sprintf("    %s\n", m.epDetail.Operation.RequestBody.Description)
					}
					s += "\n"
				}

				if len(m.epDetail.Operation.Responses) > 0 {
					s += "  Responses:\n"
					for code, resp := range m.epDetail.Operation.Responses {
						s += fmt.Sprintf("    %s — %s\n", code, resp.Description)
					}
					s += "\n"
				}
			}

			if m.msg != "" {
				s += fmt.Sprintf("  %s\n\n", m.msg)
			}

			s += "  [S]how JSON  [B]ack  [M]enu\n"
		}
	}

	return s
}

// RunExplorer starts the interactive explorer TUI.
func RunExplorer(svc *service.Service, ws *workspace.Workspace) error {
	p := tea.NewProgram(newRunModel(svc, ws))
	_, err := p.Run()
	return err
}
