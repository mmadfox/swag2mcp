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
	"github.com/mmadfox/swag2mcp/internal/spec"
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
	runAuthSpecs
	runAuthConfirm
	runAuthResult
	runDone
)

type runMode int

const (
	modeSearch runMode = iota
	modeBrowse
	modeAuth
)

const (
	randSuffixLen = 6
	pageSize      = 10
	actionHint    = "  Enter number and press Enter.  [B]ack  [M]enu.\n"
)

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
	authResult    *service.AuthResponse
	page          int
	totalPages    int
	mode          runMode
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
		page:  1,
	}
}

func (m runModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m runModel) shouldHandleInput() bool {
	return m.input.Focused() && m.state != runEndpointDetail && m.state != runAuthResult && m.state != runMenu && m.state != runDone
}

func (m runModel) transitionTo(state runState) (tea.Model, tea.Cmd) {
	m.input.SetValue("")
	m.input.Focus()
	m.state = state
	return m, textinput.Blink
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
				m.mode = modeSearch
				m.input.SetValue("")
				m.input.Placeholder = "Search endpoints..."
				m.state = runSearchQuery
				m.input.Focus()
				return m, textinput.Blink
			}
			if m.state != runSearchQuery {
				return m.handleDigit("1"), nil
			}
		case "2":
			if m.state == runMenu {
				m.mode = modeBrowse
				return m.loadSpecs()
			}
			if m.state != runSearchQuery {
				return m.handleDigit("2"), nil
			}
		case "3":
			if m.state == runMenu {
				m.mode = modeAuth
				return m.loadAuthSpecs()
			}
			if m.state != runSearchQuery {
				return m.handleDigit("3"), nil
			}
		case "4":
			if m.state == runMenu {
				m.state = runDone
				return m, tea.Quit
			}
			if m.state != runSearchQuery {
				return m.handleDigit("4"), nil
			}
		default:
			if len(msg.String()) == 1 && msg.String()[0] >= '0' && msg.String()[0] <= '9' {
				if m.state != runSearchQuery {
					return m.handleDigit(msg.String()), nil
				}
			}
		}
	}

	if !m.shouldHandleInput() {
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m runModel) handleDigit(digit string) tea.Model {
	switch m.state {
	case runSearchResults, runBrowseSpecs, runBrowseCollections, runBrowseTags, runBrowseEndpoints, runAuthSpecs:
		m.input.SetValue(m.input.Value() + digit)
	}
	return m
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
		switch strings.ToUpper(val) {
		case "N":
			if m.page < m.totalPages {
				m.page++
			}
			m.input.SetValue("")
			return m, nil
		case "P":
			if m.page > 1 {
				m.page--
			}
			m.input.SetValue("")
			return m, nil
		default:
			return m.selectSearchResult(val)
		}

	case runBrowseSpecs:
		return m.selectSpec(val)

	case runBrowseCollections:
		return m.selectCollection(val)

	case runBrowseTags:
		return m.selectTag(val)

	case runBrowseEndpoints:
		return m.selectBrowseEndpoint(val)

	case runAuthSpecs:
		return m.selectAuthSpec(val)

	case runAuthConfirm:
		if strings.ToUpper(val) == "Y" {
			return m.doAuth()
		}
		return m, nil
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
		m.input.Placeholder = "Search endpoints..."
		m.input.Focus()
		m.state = runSearchQuery
		return m, textinput.Blink
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
		switch m.mode {
		case modeSearch:
			m.state = runSearchResults
		case modeBrowse:
			m.state = runBrowseEndpoints
		}
		m.input.SetValue("")
		m.input.Focus()
		return m, textinput.Blink
	case runAuthSpecs:
		m.state = runMenu
		return m, nil
	case runAuthConfirm:
		return m.loadAuthSpecs()
	case runAuthResult:
		m.authResult = nil
		return m.loadAuthSpecs()
	}
	return m, nil
}

func (m runModel) handleMenu() (tea.Model, tea.Cmd) {
	m.msg = ""
	m.authResult = nil
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
	m.page = 1
	m.totalPages = (len(results.Endpoints) + pageSize - 1) / pageSize
	if m.totalPages < 1 {
		m.totalPages = 1
	}
	m.input.SetValue("")
	m.input.Placeholder = "Endpoint #"
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
	return m.loadEndpointDetail(ep.ID)
}

func (m runModel) loadSpecs() (tea.Model, tea.Cmd) {
	specs, err := m.svc.Specs(context.Background())
	if err != nil {
		m.err = err
		return m, nil
	}
	m.specs = specs.Specs
	return m.transitionTo(runBrowseSpecs)
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
	return m.transitionTo(runBrowseCollections)
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
	return m.transitionTo(runBrowseTags)
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
	return m.transitionTo(runBrowseEndpoints)
}

func (m runModel) selectBrowseEndpoint(val string) (tea.Model, tea.Cmd) {
	idx := 0
	if _, err := fmt.Sscanf(val, "%d", &idx); err != nil || idx < 1 || idx > len(m.endpoints) {
		return m, nil
	}
	ep := m.endpoints[idx-1]
	m.selectedEp = ep
	return m.loadEndpointDetail(ep.ID)
}

func (m runModel) loadEndpointDetail(endpointID string) (tea.Model, tea.Cmd) {
	detail, err := m.svc.Inspect(context.Background(), service.InspectRequest{EndpointID: endpointID})
	if err != nil {
		m.err = err
		return m, nil
	}
	m.epDetail = &detail
	m.selectedEpID = endpointID
	m.msg = ""
	m.state = runEndpointDetail
	m.input.Blur()
	return m, nil
}

func (m runModel) loadAuthSpecs() (tea.Model, tea.Cmd) {
	specs, err := m.svc.Specs(context.Background())
	if err != nil {
		m.err = err
		return m, nil
	}
	m.specs = specs.Specs
	return m.transitionTo(runAuthSpecs)
}

func (m runModel) selectAuthSpec(val string) (tea.Model, tea.Cmd) {
	idx := 0
	if _, err := fmt.Sscanf(val, "%d", &idx); err != nil || idx < 1 || idx > len(m.specs) {
		return m, nil
	}
	m.selectedSpec = m.specs[idx-1]
	m.input.SetValue("")
	m.input.Focus()
	m.state = runAuthConfirm
	return m, textinput.Blink
}

func (m runModel) doAuth() (tea.Model, tea.Cmd) {
	result, err := m.svc.Auth(context.Background(), service.AuthRequest{DomainID: m.selectedSpec.ID})
	if err != nil {
		m.err = err
		return m, nil
	}
	m.authResult = &result
	m.state = runAuthResult
	m.input.Blur()
	return m, nil
}

func randomSuffix(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.IntN(len(letters))]
	}
	return string(b)
}

func renderSchema(schema *spec.Schema, indent string) string {
	if schema == nil {
		return ""
	}
	var b strings.Builder

	if schema.Description != "" {
		fmt.Fprintf(&b, "%s%s\n", indent, schema.Description)
	}

	typ := schema.Type
	if typ == "" {
		typ = "any"
	}
	fmt.Fprintf(&b, "%stype: %s", indent, typ)
	if schema.Format != "" {
		fmt.Fprintf(&b, " (%s)", schema.Format)
	}
	b.WriteString("\n")

	if len(schema.Required) > 0 {
		fmt.Fprintf(&b, "%srequired: %s\n", indent, strings.Join(schema.Required, ", "))
	}

	if schema.Default != nil {
		fmt.Fprintf(&b, "%sdefault: %v\n", indent, schema.Default)
	}

	if len(schema.Enum) > 0 {
		vals := make([]string, len(schema.Enum))
		for i, v := range schema.Enum {
			vals[i] = fmt.Sprintf("%v", v)
		}
		fmt.Fprintf(&b, "%senum: [%s]\n", indent, strings.Join(vals, ", "))
	}

	if schema.Example != nil {
		exampleJSON, _ := json.MarshalIndent(schema.Example, indent+"  ", "  ")
		fmt.Fprintf(&b, "%sexample:\n%s%s\n", indent, indent, string(exampleJSON))
	}

	if len(schema.Properties) > 0 {
		fmt.Fprintf(&b, "%sproperties:\n", indent)
		for name, prop := range schema.Properties {
			req := ""
			for _, r := range schema.Required {
				if r == name {
					req = " (required)"
					break
				}
			}
			if req == "" {
				req = " (optional)"
			}
			propType := prop.Type
			if propType == "" {
				propType = "any"
			}
			fmt.Fprintf(&b, "%s  %s (%s)%s", indent, name, propType, req)
			if prop.Description != "" {
				fmt.Fprintf(&b, " — %s", prop.Description)
			}
			b.WriteString("\n")
			if prop.Items != nil {
				fmt.Fprintf(&b, "%s    items:\n", indent)
				b.WriteString(renderSchema(prop.Items, indent+"      "))
			}
			if len(prop.Properties) > 0 {
				b.WriteString(renderSchema(prop, indent+"    "))
			}
		}
	}

	if schema.Items != nil {
		fmt.Fprintf(&b, "%sitems:\n", indent)
		b.WriteString(renderSchema(schema.Items, indent+"  "))
	}

	return b.String()
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
	filename := fmt.Sprintf("%s-%s-%s-%s.json", m.epDetail.SpecDomain, method, path, randomSuffix(randSuffixLen))

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
		s += "  3. 🔑  Get auth token\n"
		s += "  4. ❌  Exit\n\n"
		s += "  Press 1, 2, 3, or 4.  (Esc/Ctrl+C to exit)\n"

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
		start := (m.page - 1) * pageSize
		end := start + pageSize
		if end > len(m.searchResults) {
			end = len(m.searchResults)
		}
		pageItems := m.searchResults[start:end]

		s += fmt.Sprintf("  Search results (%d):\n", len(m.searchResults))
		s += "  ──────────────────────\n\n"
		for i, ep := range pageItems {
			globalIdx := start + i + 1
			s += fmt.Sprintf("  %d. %-6s %s\n", globalIdx, ep.Method, ep.Path)
			s += fmt.Sprintf("     %s (%s)\n", ep.Summary, ep.SpecDomain)
		}
		s += fmt.Sprintf("\n  Page %d/%d\n", m.page, m.totalPages)
		s += "\n  " + m.input.View() + "\n\n"
		s += "  Enter number and press Enter (N for next page, P for previous).  [B]ack  [M]enu.\n"

	case runBrowseSpecs:
		s += "  Specifications:\n"
		s += "  ────────────────\n\n"
		for i, sp := range m.specs {
			s += fmt.Sprintf("  %d. %s\n", i+1, sp.Domain)
		}
		s += "\n  " + m.input.View() + "\n\n"
		s += actionHint

	case runBrowseCollections:
		s += fmt.Sprintf("  Collections for \"%s\":\n", m.selectedSpec.Domain)
		s += "  ──────────────────────────────\n\n"
		for i, col := range m.collections {
			s += fmt.Sprintf("  %d. %s (%d tags, %d methods)\n", i+1, col.Title, col.CountTags, col.CountMethods)
		}
		s += "\n  " + m.input.View() + "\n\n"
		s += actionHint

	case runBrowseTags:
		s += fmt.Sprintf("  Tags for \"%s\":\n", m.selectedColl.Title)
		s += "  ────────────────────────\n\n"
		for i, tag := range m.tags {
			s += fmt.Sprintf("  %d. %s (%d methods)\n", i+1, tag.Title, tag.CountMethods)
		}
		s += "\n  " + m.input.View() + "\n\n"
		s += actionHint

	case runBrowseEndpoints:
		s += fmt.Sprintf("  Endpoints for tag \"%s\":\n", m.selectedTag.Title)
		s += "  ─────────────────────────────\n\n"
		for i, ep := range m.endpoints {
			s += fmt.Sprintf("  %d. %-6s %s\n", i+1, ep.Method, ep.Path)
			s += fmt.Sprintf("     %s\n", ep.Summary)
		}
		s += "\n  " + m.input.View() + "\n\n"
		s += actionHint

	case runEndpointDetail:
		if m.epDetail != nil {
			s += "  Endpoint Details\n"
			s += "  ────────────────\n\n"
			s += fmt.Sprintf("  Spec:    %s\n", m.epDetail.SpecDomain)
			s += fmt.Sprintf("  Method:  %s\n", m.epDetail.Method)
			s += fmt.Sprintf("  Path:    %s\n", m.epDetail.Path)
			s += fmt.Sprintf("  BaseURL: %s\n", m.epDetail.BaseURL)
			s += fmt.Sprintf("  FullURL: %s%s\n\n", m.epDetail.BaseURL, m.epDetail.Path)

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
					for contentType, mt := range m.epDetail.Operation.RequestBody.Content {
						if !strings.Contains(contentType, "json") {
							continue
						}
						s += fmt.Sprintf("\n    %s:\n", contentType)
						if mt.Schema != nil {
							s += renderSchema(mt.Schema, "      ")
						}
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

	case runAuthSpecs:
		s += "  Select spec to get auth token:\n"
		s += "  ───────────────────────────────\n\n"
		for i, sp := range m.specs {
			s += fmt.Sprintf("  %d. %s\n", i+1, sp.Domain)
		}
		s += "\n  " + m.input.View() + "\n\n"
		s += actionHint

	case runAuthConfirm:
		s += fmt.Sprintf("  Get token for \"%s\"?\n", m.selectedSpec.Domain)
		s += "  ───────────────────────────────\n\n"
		s += "  (Y/N)\n\n"
		s += "  " + m.input.View() + "\n\n"
		s += "  [B]ack  [M]enu\n"

	case runAuthResult:
		if m.authResult != nil {
			s += "  Auth Token\n"
			s += "  ──────────\n\n"
			if token := m.authResult.Headers["Authorization"]; token != "" {
				s += fmt.Sprintf("  Token: %s\n\n", token)
			}
			if len(m.authResult.Headers) > 0 {
				s += "  Headers:\n"
				for k, v := range m.authResult.Headers {
					s += fmt.Sprintf("    %s: %s\n", k, v)
				}
				s += "\n"
			}
			if len(m.authResult.QueryParams) > 0 {
				s += "  Query Params:\n"
				for k, v := range m.authResult.QueryParams {
					s += fmt.Sprintf("    %s=%s\n", k, v)
				}
				s += "\n"
			}
			if len(m.authResult.Headers) == 0 && len(m.authResult.QueryParams) == 0 {
				s += "  No auth data available.\n\n"
			}
			s += "  [B]ack  [M]enu\n"
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
