package reader

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Outline returns a structural summary of the JSON file at path.
func (r *reader) Outline(path string, opts OutlineOptions) (Outline, error) {
	if err := r.validatePath(path); err != nil {
		return Outline{}, err
	}

	opts = normalizeOutlineOptions(opts)

	info, err := os.Stat(path)
	if err != nil {
		return Outline{}, fmt.Errorf("stat file: %w", err)
	}

	lineCount, err := countLines(path)
	if err != nil {
		return Outline{}, fmt.Errorf("count lines: %w", err)
	}

	node, depth, err := r.buildOutlineNode(path, "", opts)
	if err != nil {
		return Outline{}, fmt.Errorf("build outline: %w", err)
	}

	out := Outline{
		Type:       node.Type,
		Size:       info.Size(),
		LineCount:  lineCount,
		Depth:      depth,
		Structure:  node,
		SchemaHint: buildSchemaHint("", node),
		Keys:       node.Keys,
		ItemCount:  node.ItemCount,
		ItemType:   node.ItemType,
	}

	out.CompressionHints = buildCompressionHints("", node)
	out.NavigationHints = buildNavigationHints("", node)
	return out, nil
}

// countLines counts newline-separated lines in a file.
func countLines(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	count := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return count, nil
}

// buildOutlineNode builds a recursive structural summary of the JSON value at path.
func (r *reader) buildOutlineNode(path, key string, opts OutlineOptions) (outlineNode, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return outlineNode{}, 0, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	decoder.UseNumber()

	node, depth, err := r.decodeNode(decoder, key, opts, opts.MaxDepth)
	if err != nil {
		return outlineNode{}, 0, err
	}
	return node, depth, nil
}

// decodeNode recursively decodes the next JSON value and describes it.
func (r *reader) decodeNode(
	decoder *json.Decoder,
	key string,
	opts OutlineOptions,
	depth int,
) (outlineNode, int, error) {
	token, err := decoder.Token()
	if err != nil {
		return outlineNode{}, 0, fmt.Errorf("%w: %w", ErrNotJSON, err)
	}

	node := outlineNode{Key: key}
	maxDepth := 1

	switch t := token.(type) {
	case json.Delim:
		switch t {
		case '{':
			node.Type = jsonTypeObject
			objDepth, err := r.decodeObject(decoder, &node, opts, depth)
			if err != nil {
				return outlineNode{}, 0, err
			}
			if objDepth > maxDepth {
				maxDepth = objDepth
			}
		case '[':
			node.Type = jsonTypeArray
			arrDepth, err := r.decodeArray(decoder, &node, opts, depth)
			if err != nil {
				return outlineNode{}, 0, err
			}
			if arrDepth > maxDepth {
				maxDepth = arrDepth
			}
		}
	case string:
		node.Type = jsonTypeString
		node.Value = truncateValue(t, defaultValuePreviewLen)
	case json.Number:
		node.Type = jsonTypeNumber
		node.Value = t.String()
	case bool:
		if t {
			node.Type = jsonTypeTrue
		} else {
			node.Type = jsonTypeFalse
		}
		node.Value = strconv.FormatBool(t)
	case nil:
		node.Type = jsonTypeNull
		node.Value = jsonTypeNull
	}

	return node, maxDepth, nil
}

// decodeObject fills node with object key/type information.
func (r *reader) decodeObject(decoder *json.Decoder, node *outlineNode, opts OutlineOptions, depth int) (int, error) {
	maxDepth := 1
	node.Structure = make(map[string]outlineNode)
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return 0, fmt.Errorf("%w: %w", ErrNotJSON, err)
		}
		childKey, ok := token.(string)
		if !ok {
			return 0, ErrNotJSON
		}
		node.Keys = append(node.Keys, childKey)

		if depth <= 1 {
			if err := skipValue(decoder); err != nil {
				return 0, err
			}
			continue
		}
		child, childDepth, err := r.decodeNode(decoder, childKey, opts, depth-1)
		if err != nil {
			return 0, err
		}
		if childDepth+1 > maxDepth {
			maxDepth = childDepth + 1
		}
		node.Structure[childKey] = child
	}

	if err := decodeObjectEnd(decoder); err != nil {
		return 0, err
	}
	return maxDepth, nil
}

// decodeArray fills node with array length and item type information.
func (r *reader) decodeArray(decoder *json.Decoder, node *outlineNode, opts OutlineOptions, depth int) (int, error) {
	maxDepth := 1
	count := 0
	state := &arrayDecodeState{}

	for decoder.More() {
		count++
		if depth <= 1 || count > opts.MaxArrayItems {
			if err := skipValue(decoder); err != nil {
				return 0, err
			}
			continue
		}
		child, childDepth, err := r.decodeNode(decoder, "", opts, depth-1)
		if err != nil {
			return 0, err
		}
		if childDepth+1 > maxDepth {
			maxDepth = childDepth + 1
		}
		state.add(child)
	}

	if err := decodeArrayEnd(decoder); err != nil {
		return 0, err
	}

	node.ItemCount = count
	if state.homogeneous && state.itemType != "" {
		node.ItemType = state.itemType
	}
	if count <= opts.MaxArrayItems || state.homogeneous {
		node.SampleItems = state.sampleItems
	}
	return maxDepth, nil
}

// arrayDecodeState tracks type info while decoding array samples.
type arrayDecodeState struct {
	sampleItems []outlineNode
	itemType    string
	homogeneous bool
}

// add incorporates a decoded array item into the state.
func (s *arrayDecodeState) add(child outlineNode) {
	if s.sampleItems == nil {
		s.homogeneous = true
	}
	s.sampleItems = append(s.sampleItems, child)
	if s.itemType == "" {
		s.itemType = child.Type
		return
	}
	if s.itemType != child.Type {
		s.homogeneous = false
	}
}

// skipValue consumes the next complete JSON value without building a structure.
func skipValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrNotJSON, err)
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delim {
	case '{':
		for decoder.More() {
			if _, err := decoder.Token(); err != nil {
				return fmt.Errorf("%w: %w", ErrNotJSON, err)
			}
			if err := skipValue(decoder); err != nil {
				return err
			}
		}
	case '[':
		for decoder.More() {
			if err := skipValue(decoder); err != nil {
				return err
			}
		}
	}
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("%w: %w", ErrNotJSON, err)
	}
	return nil
}

// decodeObjectEnd reads the closing brace of an object.
func decodeObjectEnd(decoder *json.Decoder) error {
	end, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrNotJSON, err)
	}
	if end != json.Delim('}') {
		return ErrNotJSON
	}
	return nil
}

// decodeArrayEnd reads the closing bracket of an array.
func decodeArrayEnd(decoder *json.Decoder) error {
	end, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrNotJSON, err)
	}
	if end != json.Delim(']') {
		return ErrNotJSON
	}
	return nil
}

// truncateValue shortens a string to maxLen for display in outlines.
func truncateValue(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// buildSchemaHint creates a one-line description of the node structure.
func buildSchemaHint(path string, node outlineNode) string {
	switch node.Type {
	case jsonTypeObject:
		return fmt.Sprintf("Object with keys [%s]", strings.Join(node.Keys, ", "))
	case jsonTypeArray:
		if node.ItemType != "" {
			return fmt.Sprintf("Array of %s with %d items", node.ItemType, node.ItemCount)
		}
		return fmt.Sprintf("Array with %d items", node.ItemCount)
	default:
		if path == "" {
			return fmt.Sprintf("%s value", node.Type)
		}
		return fmt.Sprintf("%s value at %s", node.Type, path)
	}
}

// buildCompressionHints suggests compression strategies based on structure.
func buildCompressionHints(path string, node outlineNode) []string {
	var hints []string
	prefix := ""
	if path != "" {
		prefix = path + "."
	}

	switch node.Type {
	case jsonTypeArray:
		if node.ItemCount > 1 {
			hints = append(hints, fmt.Sprintf(
				"Use response_compress with mode 'first_of_array' and jsonPath '%s' to inspect one item.", path))
			hints = append(hints, fmt.Sprintf(
				"Use response_compress with mode 'sample_array' and jsonPath '%s' to inspect head/tail items.", path))
		}
	case jsonTypeObject:
		for _, key := range node.Keys {
			hints = append(hints, fmt.Sprintf(
				"Use response_slice with jsonPath '%s%s' to inspect this field.", prefix, key))
		}
	}

	for _, child := range node.SampleItems {
		hints = append(hints, buildCompressionHints(pathForChild(path, child.Key), child)...)
	}

	if node.Structure != nil {
		for key, child := range node.Structure {
			hints = append(hints, buildCompressionHints(pathForChild(path, key), child)...)
		}
	}

	return hints
}

// buildNavigationHints collects top-level paths and arrays for LLM navigation.
func buildNavigationHints(path string, node outlineNode) navigation {
	var nav navigation
	if path == "" {
		switch node.Type {
		case jsonTypeObject:
			for _, key := range node.Keys {
				nav.TopLevelPaths = append(nav.TopLevelPaths, topLevelPath{
					Path: key,
					Type: "unknown",
				})
			}
		case jsonTypeArray:
			nav.ArrayPaths = append(nav.ArrayPaths, arrayPath{
				Path:     "",
				Length:   node.ItemCount,
				ItemType: node.ItemType,
			})
		}
	}
	for _, child := range node.SampleItems {
		if child.Type == jsonTypeArray {
			nav.ArrayPaths = append(nav.ArrayPaths, arrayPath{
				Path:     pathForChild(path, child.Key),
				Length:   child.ItemCount,
				ItemType: child.ItemType,
			})
		}
	}
	if node.Structure != nil {
		for key, child := range node.Structure {
			if child.Type == jsonTypeArray {
				nav.ArrayPaths = append(nav.ArrayPaths, arrayPath{
					Path:     pathForChild(path, key),
					Length:   child.ItemCount,
					ItemType: child.ItemType,
				})
			}
		}
	}
	return nav
}

// pathForChild builds a dotted path to a child node.
func pathForChild(parent, child string) string {
	if parent == "" {
		return child
	}
	return parent + "." + child
}
