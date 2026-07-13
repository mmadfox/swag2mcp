package service

import (
	"embed"
	"errors"
	"fmt"
	"strings"
)

const (
	// Name is the name of the service.
	Name = "swag2mcp"
	// CollectionByID is the name of the collection_by_id tool.
	CollectionByID = "collection_by_id"
	// CollectionBySpec is the name of the collection_by_spec tool.
	CollectionBySpec = "collection_by_spec"
	// SpecByID is the name of the spec_by_id tool.
	SpecByID = "spec_by_id"
	// SpecList is the name of the spec_list tool.
	SpecList = "spec_list"
	// Inspect is the name of the inspect tool.
	Inspect = "inspect"
	// Search is the name of the search tool.
	Search = "search"
	// TagByID is the name of the tag_by_id tool.
	TagByID = "tag_by_id"
	// TagByCollection is the name of the tag_by_collection tool.
	TagByCollection = "tag_by_collection"
	// TagBySpec is the name of the tag_by_spec tool.
	TagBySpec = "tag_by_spec"
	// Invoke is the name of the invoke tool.
	Invoke = "invoke"
	// EndpointByID is the name of the endpoint_by_id tool.
	EndpointByID = "endpoint_by_id"
	// EndpointByTag is the name of the endpoint_by_tag tool.
	EndpointByTag = "endpoint_by_tag"
	// EndpointByCollection is the name of the endpoint_by_collection tool.
	EndpointByCollection = "endpoint_by_collection"
	// EndpointBySpec is the name of the endpoint_by_spec tool.
	EndpointBySpec = "endpoint_by_spec"
	// Auth is the name of the auth tool.
	Auth = "auth"
	// Info is the name of the info tool.
	Info = "info"
)

//go:embed definitions/*.md
var toolDefsFS embed.FS

// Tool represents a single MCP tool definition.
type Tool struct {
	Name        string `json:"name"        jsonschema:"required,Unique identifier for the tool"`
	Description string `json:"description" jsonschema:"required,Detailed description of what the tool does, when to use it, and what arguments it expects"`
}

// ToolDefinitions represents the complete set of MCP tools with their descriptions.
type ToolDefinitions struct {
	Instruction string `json:"instruction" jsonschema:"required,Instruction for the LLM about when to use each tool"`
	Tools       []Tool `json:"tools"       jsonschema:"required,List of available MCP tools with their detailed descriptions"`
}

// MakeToolDefinitions loads tool descriptions from embedded markdown files.
func (s *Service) MakeToolDefinitions() (ToolDefinitions, error) {
	entries, err := toolDefsFS.ReadDir("definitions")
	if err != nil {
		return ToolDefinitions{}, fmt.Errorf("failed to read tool definitions: %w", err)
	}

	var tools []Tool
	var instruction string

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		// Check if this is the instruction file.
		if entry.Name() == "instruction.md" {
			instruction, err = loadInstructionFromEmbed()
			if err != nil {
				return ToolDefinitions{}, fmt.Errorf("failed to load instruction: %w", err)
			}
			continue
		}

		loadedTool, loadErr := loadToolFromEmbed(entry.Name())
		if loadErr != nil {
			return ToolDefinitions{}, fmt.Errorf("failed to load tool from %s: %w", entry.Name(), loadErr)
		}
		if loadedTool.Name == Auth && s.disableLLMAuth.Load() {
			continue
		}
		tools = append(tools, loadedTool)
	}

	availableSpecs := s.makeAvaliablesSpecs()

	if err != nil {
		return ToolDefinitions{}, fmt.Errorf("failed to make available specs: %w", err)
	}

	return ToolDefinitions{
		Instruction: strings.Join([]string{instruction, availableSpecs}, " \n "),
		Tools:       tools,
	}, nil
}

// loadInstructionFromEmbed loads the instruction.md file from [embed.FS].
func loadInstructionFromEmbed() (string, error) {
	content, err := toolDefsFS.ReadFile("definitions/instruction.md")
	if err != nil {
		return "", fmt.Errorf("failed to read instruction file: %w", err)
	}

	return strings.TrimSpace(string(content)), nil
}

// loadToolFromEmbed parses a markdown file from [embed.FS] and returns a Tool.
func loadToolFromEmbed(filename string) (Tool, error) {
	content, err := toolDefsFS.ReadFile("definitions/" + filename)
	if err != nil {
		return Tool{}, fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) < 3 || lines[0] != "---" {
		return Tool{}, errors.New("invalid markdown format: missing frontmatter delimiter")
	}

	// Find the end of frontmatter.
	frontmatterEnd := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			frontmatterEnd = i
			break
		}
	}

	if frontmatterEnd == -1 {
		return Tool{}, errors.New("invalid markdown format: missing closing frontmatter delimiter")
	}

	// Extract name from frontmatter.
	var name string
	for i := 1; i < frontmatterEnd; i++ {
		if after, ok := strings.CutPrefix(lines[i], "name:"); ok {
			name = strings.TrimSpace(after)
			break
		}
	}

	if name == "" {
		return Tool{}, errors.New("invalid markdown format: missing name in frontmatter")
	}

	// Extract description (everything after frontmatter).
	var description strings.Builder
	inDescription := false
	for i := frontmatterEnd + 1; i < len(lines); i++ {
		if !inDescription && strings.TrimSpace(lines[i]) != "" {
			inDescription = true
		}
		if inDescription {
			description.WriteString(lines[i])
			description.WriteString("\n")
		}
	}

	return Tool{
		Name:        name,
		Description: strings.TrimSpace(description.String()),
	}, nil
}

// makeAvaliablesSpecs generates a structured, readable string of available specs and their collections.
func (s *Service) makeAvaliablesSpecs() string {
	var sb strings.Builder
	specs := s.index.AllSpecs()
	if len(specs) == 0 {
		return ""
	}

	const maxCollectionPerSpec = 10

	sb.WriteString("\n# Available specs\n")
	for _, spec := range specs {
		sb.WriteString("\n---\n")

		sb.WriteString("specID: ")
		sb.WriteString(spec.ID)
		sb.WriteString("\n")

		sb.WriteString("domain: ")
		sb.WriteString(spec.Domain)
		sb.WriteString("\n")

		sb.WriteString("title: ")
		sb.WriteString(spec.LLMTitle)
		sb.WriteString("\n")

		if len(spec.LLMInstruction) > 0 {
			sb.WriteString("instruction: ")
			sb.WriteString(strings.ReplaceAll(spec.LLMInstruction, "\n", " ")) // Normalize to a single line.
			sb.WriteString("\n")
		}

		collections, err := s.index.CollectionsBySpec(spec.ID)
		if err != nil {
			sb.WriteString("\n **No available collections**")
			continue
		}

		var hasMoreCollections bool
		if len(collections) > maxCollectionPerSpec {
			collections = collections[:maxCollectionPerSpec]
			hasMoreCollections = true
		}

		sb.WriteString("\n**Available collections:**")
		for _, collection := range collections {
			sb.WriteString("\n collectionID: ")
			sb.WriteString(collection.ID)
			sb.WriteString("\n title: ")
			sb.WriteString(collection.LLMTitle)

			if len(collection.LLMInstruction) > 0 {
				sb.WriteString("\n instruction: ")
				sb.WriteString(strings.ReplaceAll(collection.LLMInstruction, "\n", " "))
			}

			sb.WriteString("\n")
		}

		if hasMoreCollections {
			sb.WriteString(
				"\n **more than 10 collections are available for this spec; use the `collection_by_spec` tool to retrieve the full list.** \n",
			)
		}
	}

	return sb.String()
}
