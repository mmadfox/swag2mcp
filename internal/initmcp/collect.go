package initmcp

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// collectState represents a step in a collection sub-wizard.
type collectState int

const (
	colDomain collectState = iota
	colTitle
	colInstruction
	colBaseURL
	colTags
	colAuthType
	colAuthField
	colAskAddCollection
	colCollTitle
	colCollLocation
	colDone
)

// collectModel is a Bubbletea model for collecting a single spec or collection.
type collectModel struct {
	state          collectState
	specNum        int
	domain         string
	result         SpecInput
	curColl        CollectionInput
	input          textinput.Model
	authFieldIndex int
	err            error
}

func newCollectModel(specNum int) collectModel {
	ti := textinput.New()
	ti.CharLimit = 256
	ti.Width = 60
	ti.Focus()
	return collectModel{
		state:   colDomain,
		specNum: specNum,
		input:   ti,
	}
}

func (m collectModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m collectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			return m.handleEnter()
		case "y", "Y":
			if m.state == colAskAddCollection {
				return m.handleYes()
			}
		case "n", "N":
			if m.state == colAskAddCollection {
				return m.handleNo()
			}
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m collectModel) transitionTo(s collectState) (tea.Model, tea.Cmd) {
	m.state = s
	m.input.Focus()
	return m, textinput.Blink
}

func (m collectModel) handleEnter() (tea.Model, tea.Cmd) {
	val := m.input.Value()
	switch m.state {
	case colDomain:
		if val == "" {
			val = "my-api"
		}
		m.domain = val
		m.result.Domain = val
		m.input.SetValue("")
		m.input.Placeholder = "My API"
		return m.transitionTo(colTitle)
	case colTitle:
		if val == "" {
			val = m.domain + " API"
		}
		m.result.LLMTitle = val
		m.input.SetValue("")
		m.input.Placeholder = "Optional — describe what this API does"
		return m.transitionTo(colInstruction)
	case colInstruction:
		m.result.Instruction = val
		m.input.SetValue("")
		m.input.Placeholder = "https://api.example.com/v1"
		return m.transitionTo(colBaseURL)
	case colBaseURL:
		if val == "" {
			val = "https://api.example.com/v1"
		}
		m.result.BaseURL = val
		m.input.SetValue("")
		m.input.Placeholder = "public, internal (comma-separated)"
		return m.transitionTo(colTags)
	case colTags:
		if val != "" {
			parts := strings.Split(val, ",")
			for _, p := range parts {
				t := strings.TrimSpace(p)
				if t != "" {
					m.result.Tags = append(m.result.Tags, t)
				}
			}
		}
		m.input.SetValue("")
		m.input.Placeholder = "0 (none)"
		return m.transitionTo(colAuthType)
	case colAuthType:
		if val == "" {
			val = "none"
		}
		for idx, am := range availableAuthMethods {
			if fmt.Sprintf("%d", idx) == val || am.Type == val {
				val = am.Type
				break
			}
		}
		m.result.AuthType = val
		m.result.AuthConfig = make(map[string]string)
		m.authFieldIndex = 0
		fields := authFieldsFor(val)
		if len(fields) == 0 {
			m.input.SetValue("")
			return m.transitionTo(colAskAddCollection)
		}
		m.input.SetValue("")
		m.input.Placeholder = fields[0].Placeholder
		return m.transitionTo(colAuthField)
	case colAuthField:
		fields := authFieldsFor(m.result.AuthType)
		if m.authFieldIndex < len(fields) {
			f := fields[m.authFieldIndex]
			if val == "" && !f.Optional {
				val = f.Placeholder
			}
			m.result.AuthConfig[f.Name] = val
		}
		m.authFieldIndex++
		m.input.SetValue("")
		if m.authFieldIndex >= len(fields) {
			return m.transitionTo(colAskAddCollection)
		}
		m.input.Placeholder = fields[m.authFieldIndex].Placeholder
		return m, nil
	case colCollTitle:
		if val == "" {
			val = m.domain + " Collection"
		}
		m.curColl.Title = val
		m.input.SetValue("")
		m.input.Placeholder = "./specs/swagger.json"
		return m.transitionTo(colCollLocation)
	case colCollLocation:
		if val == "" {
			val = "./specs/" + m.domain + ".json"
		}
		m.curColl.Location = val
		m.result.Collections = append(m.result.Collections, m.curColl)
		m.curColl = CollectionInput{}
		m.input.SetValue("")
		return m.transitionTo(colAskAddCollection)
	}
	return m, nil
}

func (m collectModel) handleYes() (tea.Model, tea.Cmd) {
	m.curColl = CollectionInput{}
	m.input.SetValue("")
	m.input.Placeholder = m.domain + " Collection"
	return m.transitionTo(colCollTitle)
}

func (m collectModel) handleNo() (tea.Model, tea.Cmd) {
	m.state = colDone
	return m, tea.Quit
}

func (m collectModel) View() string {
	var s string
	s += "\n  ╭──────────────────────────────────────────────╮\n"
	s += "  │           swag2mcp — Add Specification        │\n"
	s += "  ╰──────────────────────────────────────────────╯\n\n"

	switch m.state {
	case colDomain:
		s += fmt.Sprintf("  Spec #%d — Domain\n", m.specNum)
		s += "  ──────────────────\n\n"
		s += "  A unique identifier for this API.\n  Examples: petstore, github-api, stripe.\n\n"
		s += "  Rules: 1–60 characters. Letters, digits, hyphens, and underscores only.\n\n"
		s += "  ────\n\n"
		s += "  [my-api]\n"
		s += "  " + m.input.View() + "\n\n"
		s += "  Press Enter to confirm. Leave empty for default.\n"
	case colTitle:
		s += fmt.Sprintf("  Spec #%d (%s) — LLM Title\n", m.specNum, m.domain)
		s += "  " + headerLine(fmt.Sprintf("Spec #%d (%s) — LLM Title", m.specNum, m.domain)) + "\n\n"
		s += "  A human-readable name the LLM will see when referencing this API.\n\n"
		s += "  Rules: 20–120 characters. Letters, digits, spaces, and basic punctuation allowed.\n\n"
		s += "  ────\n\n"
		s += "  [" + m.domain + " API]\n"
		s += "  " + m.input.View() + "\n\n"
		s += "  Press Enter to confirm. Leave empty for default.\n"
	case colInstruction:
		s += fmt.Sprintf("  Spec #%d (%s) — LLM Instruction\n", m.specNum, m.domain)
		s += "  " + headerLine(fmt.Sprintf("Spec #%d (%s) — LLM Instruction", m.specNum, m.domain)) + "\n\n"
		s += "  Optional. Tell the LLM how to use this API.\n  Example: Use this API to manage users, roles, and permissions.\n\n"
		s += "  Rules: up to 500 characters. Letters, digits, spaces, and basic punctuation allowed.\n\n"
		s += "  ────\n\n"
		s += "  [Optional — describe what this API does]\n"
		s += "  " + m.input.View() + "\n\n"
		s += "  Press Enter to confirm. Leave empty for default.\n"
	case colBaseURL:
		s += fmt.Sprintf("  Spec #%d (%s) — Base URL\n", m.specNum, m.domain)
		s += "  " + headerLine(fmt.Sprintf("Spec #%d (%s) — Base URL", m.specNum, m.domain)) + "\n\n"
		s += "  The base URL for all API requests.\n  Example: https://api.example.com/v1\n\n"
		s += "  Rules: must be a valid URL (https://... or http://...).\n\n"
		s += "  ────\n\n"
		s += "  [https://api.example.com/v1]\n"
		s += "  " + m.input.View() + "\n\n"
		s += "  Press Enter to confirm. Leave empty for default.\n"
	case colTags:
		s += fmt.Sprintf("  Spec #%d (%s) — Tags\n", m.specNum, m.domain)
		s += "  " + headerLine(fmt.Sprintf("Spec #%d (%s) — Tags", m.specNum, m.domain)) + "\n\n"
		s += "  Optional. Tags let you filter which specifications are loaded when starting the server.\n\n"
		s += "  When you run:\n    swag2mcp mcp --tags=public,internal\n  only specs with matching tags will be activated.\n\n"
		s += "  Enter comma-separated tags, or leave empty to skip.\n\n"
		s += "  ────\n\n"
		s += "  [public, internal (comma-separated)]\n"
		s += "  " + m.input.View() + "\n\n"
		s += "  Press Enter to confirm. Leave empty for default.\n"
	case colAuthType:
		s += fmt.Sprintf("  Spec #%d (%s) — Auth method\n", m.specNum, m.domain)
		s += "  " + headerLine(fmt.Sprintf("Spec #%d (%s) — Auth method", m.specNum, m.domain)) + "\n\n"
		s += "  Choose an authentication method.\n  Enter the number or type name.\n\n" + authMethodsList()
		s += "\n  Tip: You can use environment variables for any field below.\n  Use $(MY_VAR) syntax — swag2mcp resolves it at runtime.\n  Example: $(API_TOKEN) instead of a hardcoded token.\n\n"
		s += "  ────\n\n"
		s += "  [0 (none)]\n"
		s += "  " + m.input.View() + "\n\n"
		s += "  Press Enter to confirm. Leave empty for default.\n"
	case colAuthField:
		fields := authFieldsFor(m.result.AuthType)
		if m.authFieldIndex < len(fields) {
			f := fields[m.authFieldIndex]
			label := ""
			if f.Optional {
				label = "Optional"
			}
			if f.SupportsEnv {
				if label != "" {
					label += ". "
				}
				label += "Supports $(ENV_VAR) syntax"
			}
			s += fmt.Sprintf("  Spec #%d (%s) — Auth: %s\n", m.specNum, m.domain, f.Name)
			s += "  " + headerLine(fmt.Sprintf("Spec #%d (%s) — Auth: %s", m.specNum, m.domain, f.Name)) + "\n\n"
			s += fmt.Sprintf("  Enter the %s for %s authentication.\n\n", f.Name, m.result.AuthType)
			if label != "" {
				s += "  " + label + "\n\n"
			}
			s += "  ────\n\n"
			s += "  [" + f.Placeholder + "]\n"
			s += "  " + m.input.View() + "\n\n"
			s += "  Press Enter to confirm. Leave empty for default.\n"
		}
	case colAskAddCollection:
		s += fmt.Sprintf("  Spec #%d (%s) — Add a collection?\n", m.specNum, m.domain)
		s += "  ───────────────────────────────\n\n"
		s += "  A collection points to a single Swagger/OpenAPI spec file.\n"
		s += "  Each specification can have multiple collections.\n\n"
		s += "  Type y (yes) or n (no), then press Enter.\n"
	case colCollTitle:
		s += fmt.Sprintf("  Spec #%d (%s), Collection #%d — Title\n", m.specNum, m.domain, len(m.result.Collections)+1)
		s += "  " + headerLine(fmt.Sprintf("Spec #%d (%s), Collection #%d — Title", m.specNum, m.domain, len(m.result.Collections)+1)) + "\n\n"
		s += "  A name for this collection (a single Swagger/OpenAPI spec file).\n\n"
		s += "  Rules: up to 120 characters. Letters, digits, spaces, and basic punctuation allowed.\n\n"
		s += "  ────\n\n"
		s += "  [" + m.domain + " Collection]\n"
		s += "  " + m.input.View() + "\n\n"
		s += "  Press Enter to confirm. Leave empty for default.\n"
	case colCollLocation:
		s += fmt.Sprintf("  Spec #%d (%s), Collection #%d — Location\n", m.specNum, m.domain, len(m.result.Collections)+1)
		s += "  " + headerLine(fmt.Sprintf("Spec #%d (%s), Collection #%d — Location", m.specNum, m.domain, len(m.result.Collections)+1)) + "\n\n"
		s += "  Path or URL to the Swagger/OpenAPI spec file.\n"
		s += "  Supports:\n"
		s += "    • Local path:  ./specs/api.json\n"
		s += "    • File URL:    file:///home/user/specs/api.json\n"
		s += "    • HTTP(S):     https://example.com/api/swagger.json\n\n"
		s += "  Rules: 5–250 characters. Must point to a valid spec file or URL.\n\n"
		s += "  ────\n\n"
		s += "  [./specs/" + m.domain + ".json]\n"
		s += "  " + m.input.View() + "\n\n"
		s += "  Press Enter to confirm. Leave empty for default.\n"
	}
	return s
}

// collectSpec runs a TUI to collect a single specification.
func collectSpec(specNum int) (SpecInput, error) {
	p := tea.NewProgram(newCollectModel(specNum))
	final, err := p.Run()
	if err != nil {
		return SpecInput{}, err
	}
	m, ok := final.(collectModel)
	if !ok {
		return SpecInput{}, fmt.Errorf("unexpected model type")
	}
	return m.result, m.err
}

// collectCollectionModel is a minimal TUI for collecting a single collection.
type collectCollectionModel struct {
	state   collectState
	domain  string
	specNum int
	collNum int
	result  CollectionInput
	input   textinput.Model
	err     error
}

func newCollectCollectionModel(specNum, collNum int, domain string) collectCollectionModel {
	ti := textinput.New()
	ti.CharLimit = 256
	ti.Width = 60
	ti.Focus()
	return collectCollectionModel{
		state:   colCollTitle,
		domain:  domain,
		specNum: specNum,
		collNum: collNum,
		input:   ti,
	}
}

func (m collectCollectionModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m collectCollectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			return m.handleEnter()
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m collectCollectionModel) handleEnter() (tea.Model, tea.Cmd) {
	val := m.input.Value()
	switch m.state {
	case colCollTitle:
		if val == "" {
			val = m.domain + " Collection"
		}
		m.result.Title = val
		m.input.SetValue("")
		m.input.Placeholder = "./specs/swagger.json"
		m.state = colCollLocation
		m.input.Focus()
		return m, textinput.Blink
	case colCollLocation:
		if val == "" {
			val = "./specs/" + m.domain + ".json"
		}
		m.result.Location = val
		m.state = colDone
		return m, tea.Quit
	}
	return m, nil
}

func (m collectCollectionModel) View() string {
	var s string
	s += "\n  ╭──────────────────────────────────────────────╮\n"
	s += "  │           swag2mcp — Add Collection          │\n"
	s += "  ╰──────────────────────────────────────────────╯\n\n"

	switch m.state {
	case colCollTitle:
		s += fmt.Sprintf("  Spec #%d (%s), Collection #%d — Title\n", m.specNum, m.domain, m.collNum)
		s += "  " + headerLine(fmt.Sprintf("Spec #%d (%s), Collection #%d — Title", m.specNum, m.domain, m.collNum)) + "\n\n"
		s += "  A name for this collection (a single Swagger/OpenAPI spec file).\n\n"
		s += "  Rules: up to 120 characters. Letters, digits, spaces, and basic punctuation allowed.\n\n"
		s += "  ────\n\n"
		s += "  [" + m.domain + " Collection]\n"
		s += "  " + m.input.View() + "\n\n"
		s += "  Press Enter to confirm. Leave empty for default.\n"
	case colCollLocation:
		s += fmt.Sprintf("  Spec #%d (%s), Collection #%d — Location\n", m.specNum, m.domain, m.collNum)
		s += "  " + headerLine(fmt.Sprintf("Spec #%d (%s), Collection #%d — Location", m.specNum, m.domain, m.collNum)) + "\n\n"
		s += "  Path or URL to the Swagger/OpenAPI spec file.\n"
		s += "  Supports:\n"
		s += "    • Local path:  ./specs/api.json\n"
		s += "    • File URL:    file:///home/user/specs/api.json\n"
		s += "    • HTTP(S):     https://example.com/api/swagger.json\n\n"
		s += "  Rules: 5–250 characters. Must point to a valid spec file or URL.\n\n"
		if strings.HasPrefix(m.result.Location, "http://") || strings.HasPrefix(m.result.Location, "https://") {
			s += "  ℹ️ URL detected — will be cached on first use.\n\n"
		}
		s += "  ────\n\n"
		s += "  [./specs/" + m.domain + ".json]\n"
		s += "  " + m.input.View() + "\n\n"
		s += "  Press Enter to confirm. Leave empty for default.\n"
	}
	return s
}

// collectCollection runs a TUI to collect a single collection.
func collectCollection(specNum, collNum int, domain string) (CollectionInput, error) {
	p := tea.NewProgram(newCollectCollectionModel(specNum, collNum, domain))
	final, runErr := p.Run()
	if runErr != nil {
		return CollectionInput{}, runErr
	}
	m, ok := final.(collectCollectionModel)
	if !ok {
		return CollectionInput{}, fmt.Errorf("unexpected model type")
	}
	return m.result, m.err
}
