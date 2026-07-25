package reader

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

// jsonType constants describe JSON value kinds used across the reader package.
const (
	jsonTypeObject = "object"
	jsonTypeArray  = "array"
	jsonTypeString = "string"
	jsonTypeNumber = "number"
	jsonTypeBool   = "bool"
	jsonTypeTrue   = "true"
	jsonTypeFalse  = "false"
	jsonTypeNull   = "null"
	jsonTypeValue  = "value"
)

// Outline describes the high-level structure of a JSON response file.
type Outline struct {
	Type             string      `json:"type"`
	Size             int64       `json:"size"`
	LineCount        int         `json:"lineCount"`
	Depth            int         `json:"depth"`
	Structure        outlineNode `json:"structure"`
	SchemaHint       string      `json:"schemaHint"`
	Keys             []string    `json:"keys,omitempty"`
	ItemCount        int         `json:"itemCount,omitempty"`
	ItemType         string      `json:"itemType,omitempty"`
	CompressionHints []string    `json:"compressionHints"`
	NavigationHints  navigation  `json:"navigationHints"`
}

// OutlineOptions controls how the outline is built.
type OutlineOptions struct {
	MaxDepth      int `json:"maxDepth,omitempty"`
	MaxArrayItems int `json:"maxArrayItems,omitempty"`
}

// outlineNode is a recursive description of a JSON value.
type outlineNode struct {
	Type        string                 `json:"type"`
	Key         string                 `json:"key,omitempty"`
	Value       string                 `json:"value,omitempty"`
	Keys        []string               `json:"keys,omitempty"`
	Structure   map[string]outlineNode `json:"structure,omitempty"`
	ItemCount   int                    `json:"itemCount,omitempty"`
	ItemType    string                 `json:"itemType,omitempty"`
	SampleItems []outlineNode          `json:"sampleItems,omitempty"`
}

// navigation helps LLM move through the file logically.
type navigation struct {
	TopLevelPaths []topLevelPath `json:"topLevelPaths"`
	ArrayPaths    []arrayPath    `json:"arrayPaths"`
}

// topLevelPath names a root-level key.
type topLevelPath struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

// arrayPath describes an array found in the outline.
type arrayPath struct {
	Path     string `json:"path"`
	Length   int    `json:"length"`
	ItemType string `json:"itemType"`
}

// CompressMode selects how a large JSON value is reduced.
type CompressMode string

const (
	// CompressFirstOfArray keeps only the first element of a homogeneous array.
	CompressFirstOfArray CompressMode = "first_of_array"
	// CompressSampleArray keeps a head and tail sample of an array.
	CompressSampleArray CompressMode = "sample_array"
	// CompressTruncateStrings shortens every string value.
	CompressTruncateStrings CompressMode = "truncate_strings"
	// CompressKeysOnly replaces object values with their type names.
	CompressKeysOnly CompressMode = "keys_only"
	// CompressSelectKeys keeps only the selected keys for every object in an array.
	CompressSelectKeys CompressMode = "select_keys"
)

// CompressOptions controls how compression is performed.
type CompressOptions struct {
	JSONPath   string       `json:"jsonPath,omitempty"`
	Mode       CompressMode `json:"mode"`
	ArrayHead  int          `json:"arrayHead,omitempty"`
	ArrayTail  int          `json:"arrayTail,omitempty"`
	StringLen  int          `json:"stringLen,omitempty"`
	SelectKeys []string     `json:"selectKeys,omitempty"`
	Limit      int          `json:"limit,omitempty"`
}

// CompressResult is the outcome of compression.
type CompressResult struct {
	Body     any    `json:"body,omitempty"`
	TooLarge bool   `json:"tooLarge"`
	Hint     string `json:"hint,omitempty"`
}

// Slice describes a fragment extracted from a JSON file.
type Slice struct {
	Lines      [2]int `json:"lines"`
	Fragment   string `json:"fragment,omitempty"`
	Value      any    `json:"value"`
	JSONPath   string `json:"jsonPath,omitempty"`
	Context    string `json:"context"`
	IsComplete bool   `json:"isComplete"`
	NextLine   int    `json:"nextLine"`
	PrevLine   int    `json:"prevLine"`
	NextPath   string `json:"nextPath,omitempty"`
	PrevPath   string `json:"prevPath,omitempty"`
}

// SliceOptions controls fragment extraction.
type SliceOptions struct {
	JSONPath string `json:"jsonPath,omitempty"`
	Line     int    `json:"line,omitempty"`
	Range    string `json:"range,omitempty"`
	Around   int    `json:"around,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}
