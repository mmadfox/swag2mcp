package initpkg

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/mmadfox/swag2mcp/internal/workspace"
)

// state represents the current step in the initialization wizard.
type state int

const (
	stateConfigPath state = iota
	stateWorkspaceDir
	stateAskAddSpec
	stateSpecDomain
	stateSpecTitle
	stateSpecInstruction
	stateSpecBaseURL
	stateSpecTags
	stateAuthType
	stateAuthField
	stateAskAddCollection
	stateCollTitle
	stateCollLocation
	stateAskAddAnotherSpec
	stateConfirm
	stateDone
)

// authFieldDef describes a single field for an auth method.
type authFieldDef struct {
	Name        string
	Placeholder string
	Optional    bool
	SupportsEnv bool
}

// authMethod describes an authentication method and its fields.
type authMethod struct {
	Type   string
	Label  string
	Fields []authFieldDef
}

// availableAuthMethods is the list of supported authentication methods.
var availableAuthMethods = []authMethod{
	{Type: "none", Label: "No authentication", Fields: nil},
	{Type: "basic", Label: "HTTP Basic (username + password)", Fields: []authFieldDef{
		{Name: "username", Placeholder: "admin", SupportsEnv: true},
		{Name: "password", Placeholder: "secret", SupportsEnv: true},
	}},
	{Type: "bearer", Label: "Bearer token", Fields: []authFieldDef{
		{Name: "token", Placeholder: "eyJhbGci...", SupportsEnv: true},
	}},
	{Type: "digest", Label: "HTTP Digest (username + password)", Fields: []authFieldDef{
		{Name: "username", Placeholder: "admin", SupportsEnv: true},
		{Name: "password", Placeholder: "secret", SupportsEnv: true},
	}},
	{Type: "api-key", Label: "API Key (key + value + location)", Fields: []authFieldDef{
		{Name: "key", Placeholder: "X-API-Key"},
		{Name: "value", Placeholder: "your-api-key"},
		{Name: "in", Placeholder: "header or query"},
	}},
	{Type: "oauth2-cc", Label: "OAuth2 Client Credentials", Fields: []authFieldDef{
		{Name: "client_id", Placeholder: "your-client-id", SupportsEnv: true},
		{Name: "client_secret", Placeholder: "your-client-secret", SupportsEnv: true},
		{Name: "token_url", Placeholder: "https://auth.example.com/token"},
		{Name: "scopes", Placeholder: "read write (space-separated)", Optional: true},
	}},
	{Type: "oauth2-pwd", Label: "OAuth2 Password Grant", Fields: []authFieldDef{
		{Name: "username", Placeholder: "user", SupportsEnv: true},
		{Name: "password", Placeholder: "pass", SupportsEnv: true},
		{Name: "client_id", Placeholder: "your-client-id", SupportsEnv: true},
		{Name: "client_secret", Placeholder: "your-client-secret", SupportsEnv: true},
		{Name: "token_url", Placeholder: "https://auth.example.com/token"},
		{Name: "scopes", Placeholder: "read write (space-separated)", Optional: true},
	}},
	{Type: "script", Label: "Script (custom auth logic)", Fields: []authFieldDef{
		{Name: "source", Placeholder: "path/to/script.sh or inline code"},
	}},
}

// SpecInput holds the user-provided data for a single API specification.
type SpecInput struct {
	Domain      string
	LLMTitle    string
	Instruction string
	BaseURL     string
	Tags        []string
	AuthType    string
	AuthConfig  map[string]string
	Collections []CollectionInput
}

// CollectionInput holds the user-provided data for a single collection.
type CollectionInput struct {
	Title    string
	Location string
}

// templateData is passed to the config YAML template.
type templateData struct {
	WorkspaceDir string
	Specs        []SpecInput
}

// hint holds the title, description, placeholder, and validation rules for a wizard step.
type hint struct {
	title       string
	description string
	placeholder string
	rules       string
}

// model is the Bubbletea model for the initialization wizard.
type model struct {
	state          state
	configPath     string
	workspaceDir   string
	specs          []SpecInput
	curSpec        SpecInput
	curColl        CollectionInput
	input          textinput.Model
	err            error
	width          int
	authFieldIndex int
}

// initialModel returns the starting model for the wizard.
func initialModel() model {
	ti := textinput.New()
	ti.CharLimit = 256
	ti.Width = 60
	ti.Focus()

	return model{
		state: stateConfigPath,
		input: ti,
	}
}

// Init implements tea.Model.
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		case "y", "Y":
			if m.state == stateAskAddSpec || m.state == stateAskAddCollection || m.state == stateAskAddAnotherSpec {
				return m.handleYes()
			}
			if m.state == stateConfirm {
				return m.handleConfirm()
			}

		case "n", "N":
			if m.state == stateAskAddSpec || m.state == stateAskAddCollection || m.state == stateAskAddAnotherSpec {
				return m.handleNo()
			}
			if m.state == stateConfirm {
				m.state = stateAskAddAnotherSpec
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// transitionTo sets the state and refocuses the text input.
func (m model) transitionTo(s state) (tea.Model, tea.Cmd) {
	m.state = s
	m.input.Focus()
	return m, textinput.Blink
}

// handleEnter processes the Enter key press for the current state.
func (m model) handleEnter() (tea.Model, tea.Cmd) {
	val := m.input.Value()

	switch m.state {
	case stateConfigPath:
		if val == "" {
			val = "swag2mcp.yaml"
		}
		if info, statErr := os.Stat(val); statErr == nil && info.IsDir() {
			val = filepath.Join(val, "swag2mcp.yaml")
		}
		m.configPath = val
		m.input.SetValue("")
		m.input.Placeholder = "./.swag2mcp"
		return m.transitionTo(stateWorkspaceDir)

	case stateWorkspaceDir:
		if val == "" {
			val = filepath.Join(filepath.Dir(m.configPath), ".swag2mcp")
		}
		m.workspaceDir = val
		m.input.SetValue("")
		return m.transitionTo(stateAskAddSpec)

	case stateSpecDomain:
		if val == "" {
			val = "my-api"
		}
		m.curSpec.Domain = val
		m.input.SetValue("")
		m.input.Placeholder = "My API"
		return m.transitionTo(stateSpecTitle)

	case stateSpecTitle:
		if val == "" {
			val = m.curSpec.Domain + " API"
		}
		m.curSpec.LLMTitle = val
		m.input.SetValue("")
		m.input.Placeholder = "Optional — describe what this API does"
		return m.transitionTo(stateSpecInstruction)

	case stateSpecInstruction:
		m.curSpec.Instruction = val
		m.input.SetValue("")
		m.input.Placeholder = "https://api.example.com/v1"
		return m.transitionTo(stateSpecBaseURL)

	case stateSpecBaseURL:
		if val == "" {
			val = "https://api.example.com/v1"
		}
		m.curSpec.BaseURL = val
		m.input.SetValue("")
		m.input.Placeholder = "petstore, public (comma-separated)"
		return m.transitionTo(stateSpecTags)

	case stateSpecTags:
		if val != "" {
			parts := strings.Split(val, ",")
			for _, p := range parts {
				t := strings.TrimSpace(p)
				if t != "" {
					m.curSpec.Tags = append(m.curSpec.Tags, t)
				}
			}
		}
		m.input.SetValue("")
		m.input.Placeholder = "none"
		return m.transitionTo(stateAuthType)

	case stateAuthType:
		if val == "" {
			val = "none"
		}
		// Resolve by index (e.g. "1", "2") or by type name.
		for idx, am := range availableAuthMethods {
			if fmt.Sprintf("%d", idx) == val || am.Type == val {
				val = am.Type
				break
			}
		}
		m.curSpec.AuthType = val
		m.curSpec.AuthConfig = make(map[string]string)
		m.authFieldIndex = 0

		fields := authFieldsFor(val)
		if len(fields) == 0 {
			m.input.SetValue("")
			return m.transitionTo(stateAskAddCollection)
		}
		m.input.SetValue("")
		m.input.Placeholder = fields[0].Placeholder
		return m.transitionTo(stateAuthField)

	case stateAuthField:
		fields := authFieldsFor(m.curSpec.AuthType)
		if m.authFieldIndex < len(fields) {
			f := fields[m.authFieldIndex]
			if val == "" && !f.Optional {
				val = f.Placeholder
			}
			m.curSpec.AuthConfig[f.Name] = val
		}
		m.authFieldIndex++
		m.input.SetValue("")
		if m.authFieldIndex >= len(fields) {
			return m.transitionTo(stateAskAddCollection)
		}
		m.input.Placeholder = fields[m.authFieldIndex].Placeholder
		return m, nil

	case stateCollTitle:
		if val == "" {
			val = m.curSpec.Domain + " Collection"
		}
		m.curColl.Title = val
		m.input.SetValue("")
		m.input.Placeholder = "./specs/swagger.json"
		return m.transitionTo(stateCollLocation)

	case stateCollLocation:
		if val == "" {
			val = "./specs/" + m.curSpec.Domain + ".json"
		}
		m.curColl.Location = val
		m.curSpec.Collections = append(m.curSpec.Collections, m.curColl)
		m.curColl = CollectionInput{}
		m.input.SetValue("")
		return m.transitionTo(stateAskAddCollection)
	}

	return m, nil
}

// handleYes processes a "yes" answer for yes/no states.
func (m model) handleYes() (tea.Model, tea.Cmd) {
	switch m.state {
	case stateAskAddSpec:
		m.curSpec = SpecInput{}
		m.input.SetValue("")
		m.input.Placeholder = "my-api"
		return m.transitionTo(stateSpecDomain)

	case stateAskAddCollection:
		m.curColl = CollectionInput{}
		m.input.SetValue("")
		m.input.Placeholder = m.curSpec.Domain + " Collection"
		return m.transitionTo(stateCollTitle)

	case stateAskAddAnotherSpec:
		m.curSpec = SpecInput{}
		m.input.SetValue("")
		m.input.Placeholder = "my-api"
		return m.transitionTo(stateSpecDomain)
	}
	return m, nil
}

// handleNo processes a "no" answer for yes/no states.
func (m model) handleNo() (tea.Model, tea.Cmd) {
	switch m.state {
	case stateAskAddSpec:
		return m.transitionTo(stateConfirm)

	case stateAskAddCollection:
		m.specs = append(m.specs, m.curSpec)
		m.curSpec = SpecInput{}
		m.input.SetValue("")
		return m.transitionTo(stateAskAddAnotherSpec)

	case stateAskAddAnotherSpec:
		return m.transitionTo(stateConfirm)
	}
	return m, nil
}

// handleConfirm writes the result and transitions to done.
func (m model) handleConfirm() (tea.Model, tea.Cmd) {
	if err := WriteResult(m.configPath, m.workspaceDir, m.specs); err != nil {
		m.err = err
	}
	m.state = stateDone
	return m, tea.Quit
}

// authFieldsFor returns the field definitions for the given auth type.
func authFieldsFor(authType string) []authFieldDef {
	for _, m := range availableAuthMethods {
		if m.Type == authType {
			return m.Fields
		}
	}
	return nil
}

// currentHint returns the hint for the current wizard step.
func (m model) currentHint() hint {
	switch m.state {
	case stateConfigPath:
		return hint{
			title:       "Config file path",
			description: "Where should the swag2mcp.yaml configuration file be created?",
			placeholder: ".",
		}
	case stateWorkspaceDir:
		return hint{
			title:       "Workspace directory",
			description: "Where should the workspace directory be created?\n  This will store cached specs, local spec files, and API responses.",
			placeholder: "./.swag2mcp",
		}
	case stateSpecDomain:
		return hint{
			title:       fmt.Sprintf("Spec #%d — Domain", len(m.specs)+1),
			description: "A unique identifier for this API.\n  Examples: petstore, github-api, stripe.",
			placeholder: "my-api",
			rules:       "Rules: 1–60 characters. Letters, digits, hyphens, and underscores only.",
		}
	case stateSpecTitle:
		return hint{
			title:       fmt.Sprintf("Spec #%d (%s) — LLM Title", len(m.specs)+1, m.curSpec.Domain),
			description: "A human-readable name the LLM will see when referencing this API.",
			placeholder: m.curSpec.Domain + " API",
			rules:       "Rules: 20–120 characters. Letters, digits, spaces, and basic punctuation allowed.",
		}
	case stateSpecInstruction:
		return hint{
			title:       fmt.Sprintf("Spec #%d (%s) — LLM Instruction", len(m.specs)+1, m.curSpec.Domain),
			description: "Optional. Tell the LLM how to use this API.\n  Example: Use this API to manage users, roles, and permissions.",
			placeholder: "Optional — describe what this API does",
			rules:       "Rules: up to 500 characters. Letters, digits, spaces, and basic punctuation allowed.",
		}
	case stateSpecBaseURL:
		return hint{
			title:       fmt.Sprintf("Spec #%d (%s) — Base URL", len(m.specs)+1, m.curSpec.Domain),
			description: "The base URL for all API requests.\n  Example: https://api.example.com/v1",
			placeholder: "https://api.example.com/v1",
			rules:       "Rules: must be a valid URL (https://... or http://...).",
		}
	case stateSpecTags:
		return hint{
			title:       fmt.Sprintf("Spec #%d (%s) — Tags", len(m.specs)+1, m.curSpec.Domain),
			description: "Optional. Tags let you filter which specifications are loaded when starting the server.\n\n  When you run:\n    swag2mcp mcp --tags=public,internal\n  only specs with matching tags will be activated.\n\n  Enter comma-separated tags, or leave empty to skip.",
			placeholder: "public, internal (comma-separated)",
		}
	case stateAuthType:
		return hint{
			title:       fmt.Sprintf("Spec #%d (%s) — Auth method", len(m.specs)+1, m.curSpec.Domain),
			description: "Choose an authentication method.\n  Enter the number or type name.\n\n" + authMethodsList() + "\n  Tip: You can use environment variables for any field below.\n  Use $(MY_VAR) syntax — swag2mcp resolves it at runtime.\n  Example: $(API_TOKEN) instead of a hardcoded token.",
			placeholder: "0 (none)",
		}
	case stateAuthField:
		fields := authFieldsFor(m.curSpec.AuthType)
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
			return hint{
				title:       fmt.Sprintf("Spec #%d (%s) — Auth: %s", len(m.specs)+1, m.curSpec.Domain, f.Name),
				description: fmt.Sprintf("Enter the %s for %s authentication.", f.Name, m.curSpec.AuthType),
				placeholder: f.Placeholder,
				rules:       label,
			}
		}
		return hint{}
	case stateCollTitle:
		return hint{
			title:       fmt.Sprintf("Spec #%d (%s), Collection #%d — Title", len(m.specs)+1, m.curSpec.Domain, len(m.curSpec.Collections)+1),
			description: "A name for this collection (a single Swagger/OpenAPI spec file).",
			placeholder: m.curSpec.Domain + " Collection",
			rules:       "Rules: up to 120 characters. Letters, digits, spaces, and basic punctuation allowed.",
		}
	case stateCollLocation:
		return hint{
			title:       fmt.Sprintf("Spec #%d (%s), Collection #%d — Location", len(m.specs)+1, m.curSpec.Domain, len(m.curSpec.Collections)+1),
			description: "Path or URL to the Swagger/OpenAPI spec file.\n  Supports:\n    • Local path:  ./specs/api.json\n    • File URL:    file:///home/user/specs/api.json\n    • HTTP(S):     https://example.com/api/swagger.json",
			placeholder: "./specs/" + m.curSpec.Domain + ".json",
			rules:       "Rules: 5–250 characters. Must point to a valid spec file or URL.",
		}
	default:
		return hint{}
	}
}

// authMethodsList returns a formatted list of available auth methods.
func authMethodsList() string {
	var s string
	for i, m := range availableAuthMethods {
		s += fmt.Sprintf("    %d. %-12s — %s\n", i, m.Type, m.Label)
	}
	return s
}

// View implements tea.Model.
func (m model) View() string {
	var s string

	s += "\n  ╭──────────────────────────────────────────────╮\n"
	s += "  │           swag2mcp — Initialization           │\n"
	s += "  ╰──────────────────────────────────────────────╯\n\n"

	switch m.state {
	case stateConfigPath, stateWorkspaceDir, stateSpecDomain, stateSpecTitle, stateSpecInstruction, stateSpecBaseURL, stateSpecTags, stateAuthType, stateAuthField, stateCollTitle, stateCollLocation:
		h := m.currentHint()
		s += fmt.Sprintf("  %s\n", h.title)
		s += "  " + headerLine(h.title) + "\n\n"
		s += "  " + h.description + "\n\n"
		if h.rules != "" {
			s += "  " + h.rules + "\n\n"
		}
		s += "  ────\n\n"
		s += "  [" + h.placeholder + "]\n"
		s += "  " + m.input.View() + "\n\n"
		s += "  Press Enter to confirm. Leave empty for default.\n"

	case stateAskAddSpec:
		s += "  Add an API specification?\n"
		s += "  ─────────────────────────\n\n"
		s += "  An API specification describes a set of related endpoints.\n"
		s += "  You can add one or more specifications (e.g. petstore, github-api).\n\n"
		s += "  Type y (yes) or n (no), then press Enter.\n"

	case stateAskAddCollection:
		s += fmt.Sprintf("  Spec #%d (%s) — Add a collection?\n", len(m.specs)+1, m.curSpec.Domain)
		s += "  ───────────────────────────────\n\n"
		s += "  A collection points to a single Swagger/OpenAPI spec file.\n"
		s += "  Each specification can have multiple collections.\n\n"
		s += "  Type y (yes) or n (no), then press Enter.\n"

	case stateAskAddAnotherSpec:
		s += "  Add another API specification?\n"
		s += "  ───────────────────────────────\n\n"
		s += "  Type y (yes) or n (no), then press Enter.\n"

	case stateConfirm:
		s += "  Review your configuration:\n"
		s += "  ────────────────────────────\n\n"
		s += fmt.Sprintf("  Config path:     %s\n", m.configPath)
		s += fmt.Sprintf("  Workspace dir:   %s\n", m.workspaceDir)
		s += fmt.Sprintf("  Specifications:  %d\n\n", len(m.specs))
		for i, sp := range m.specs {
			s += fmt.Sprintf("  Spec #%d: %s (%s)\n", i+1, sp.LLMTitle, sp.Domain)
			s += fmt.Sprintf("           Base URL: %s\n", sp.BaseURL)
			if len(sp.Tags) > 0 {
				s += fmt.Sprintf("           Tags:     %s\n", strings.Join(sp.Tags, ", "))
			}
			if sp.AuthType != "" && sp.AuthType != "none" {
				s += fmt.Sprintf("           Auth:     %s\n", sp.AuthType)
			}
			s += fmt.Sprintf("           Collections: %d\n", len(sp.Collections))
			for j, col := range sp.Collections {
				s += fmt.Sprintf("             %d. %s → %s\n", j+1, col.Title, col.Location)
			}
			s += "\n"
		}
		s += "  Write configuration and initialize workspace?\n"
		s += "  Type y (yes) or n (no), then press Enter.\n"

	case stateDone:
		if m.err != nil {
			s += fmt.Sprintf("  ❌ Error: %s\n\n", m.err)
		} else {
			s += fmt.Sprintf("  ✅ Configuration written to: %s\n", m.configPath)
			s += fmt.Sprintf("  ✅ Workspace initialized at: %s\n", m.workspaceDir)
			s += "  Run `swag2mcp mcp` to start the server.\n\n"
		}
	}

	return s
}

// headerLine returns a line of dashes matching the title length.
func headerLine(title string) string {
	n := 0
	for _, r := range title {
		if r > 127 {
			n += 2
		} else {
			n++
		}
	}
	line := ""
	for i := 0; i < n; i++ {
		line += "─"
	}
	return line
}

// RunTUI starts the interactive initialization wizard.
func RunTUI() (configPath, workspaceDir string, specs []SpecInput, err error) {
	p := tea.NewProgram(initialModel())
	final, runErr := p.Run()
	if runErr != nil {
		return "", "", nil, runErr
	}

	m, ok := final.(model)
	if !ok {
		return "", "", nil, fmt.Errorf("unexpected model type")
	}

	return m.configPath, m.workspaceDir, m.specs, m.err
}

// BuildConfigYAML renders the config YAML from the template and collected data.
func BuildConfigYAML(workspaceDir string, specs []SpecInput) ([]byte, error) {
	tmpl, err := template.New("config").Parse(configTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData{
		WorkspaceDir: workspaceDir,
		Specs:        specs,
	}); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	return buf.Bytes(), nil
}

// WriteResult writes the config file and initializes the workspace.
func WriteResult(configPath, workspaceDir string, specs []SpecInput) error {
	ws, err := workspace.New(workspaceDir)
	if err != nil {
		return fmt.Errorf("workspace: %w", err)
	}
	if err := ws.Init(); err != nil {
		return fmt.Errorf("init workspace: %w", err)
	}

	if info, statErr := os.Stat(configPath); statErr == nil && info.IsDir() {
		configPath = filepath.Join(configPath, "swag2mcp.yaml")
	}

	cfgDir := filepath.Dir(configPath)
	if err := os.MkdirAll(cfgDir, 0750); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := BuildConfigYAML(workspaceDir, specs)
	if err != nil {
		return fmt.Errorf("build config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}
