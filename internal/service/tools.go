package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"embed"
	"errors"
	"fmt"
	"strings"
)

//go:embed definitions/*.md
var toolDefsFS embed.FS

type toolsService struct {
	index           IndexReader
	llmAuthDisabled func() bool
}

func newToolsService(index IndexReader, llmAuthDisabled func() bool) *toolsService {
	return &toolsService{
		index:           index,
		llmAuthDisabled: llmAuthDisabled,
	}
}

// MakeToolDefinitions loads tool descriptions from embedded markdown files
// and returns the complete set of MCP tools with their descriptions.
func (ts *toolsService) MakeToolDefinitions() (ToolDefinitions, error) {
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
		if loadedTool.Name == Auth && ts.llmAuthDisabled() {
			continue
		}
		tools = append(tools, loadedTool)
	}

	availableSpecs := ts.makeAvailableSpecs()

	return ToolDefinitions{
		Instruction: strings.Join([]string{instruction, availableSpecs}, " \n "),
		Tools:       tools,
	}, nil
}

func (ts *toolsService) makeAvailableSpecs() string {
	var sb strings.Builder
	specs := ts.index.AllSpecs()
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
			sb.WriteString(strings.ReplaceAll(spec.LLMInstruction, "\n", " "))
			sb.WriteString("\n")
		}

		collections, err := ts.index.CollectionsBySpec(spec.ID)
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
