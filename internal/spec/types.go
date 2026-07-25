package spec

// Doc is the unified representation of a parsed Swagger/OpenAPI document.
// All versions (Swagger 2.0, OpenAPI 3.0, 3.1) are mapped to this type.
type Doc struct {
	Version     string      `json:"version,omitempty"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	VersionStr  string      `json:"versionStr,omitempty"`
	Servers     []Server    `json:"servers,omitempty"`
	PathItems   []*PathItem `json:"pathItems,omitempty"`
}

type Server struct {
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
}

// PathItem is a single endpoint (path + method + operation).
type PathItem struct {
	Path      string     `json:"path,omitempty"`
	Method    string     `json:"method,omitempty"`
	Operation *Operation `json:"operation,omitempty"`
}

type Operation struct {
	ID          string               `json:"id,omitempty"`
	Tags        []string             `json:"tags,omitempty"`
	Summary     string               `json:"summary,omitempty"`
	Description string               `json:"description,omitempty"`
	Deprecated  bool                 `json:"deprecated,omitempty"`
	Parameters  []*Parameter         `json:"parameters,omitempty"`
	RequestBody *RequestBody         `json:"requestBody,omitempty"`
	Responses   map[string]*Response `json:"responses,omitempty"`
}

type Parameter struct {
	Name        string  `json:"name,omitempty"`
	In          string  `json:"in,omitempty"` // "query", "path", "header", "cookie"
	Description string  `json:"description,omitempty"`
	Required    bool    `json:"required,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
}

type RequestBody struct {
	Description string                `json:"description,omitempty"`
	Required    bool                  `json:"required,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty"`
}

type Response struct {
	Description string                `json:"description,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty"`
}

type MediaType struct {
	Schema *Schema `json:"schema,omitempty"`
}

type Schema struct {
	Type       string             `json:"type,omitempty"`
	Format     string             `json:"format,omitempty"`
	Properties map[string]*Schema `json:"properties,omitempty"`
	Items      *Schema            `json:"items,omitempty"`
	Required   []string           `json:"required,omitempty"`
	Ref        string             `json:"$ref,omitempty"`
	OneOf      []*Schema          `json:"oneOf,omitempty"`
	AnyOf      []*Schema          `json:"anyOf,omitempty"`
	AllOf      []*Schema          `json:"allOf,omitempty"`

	Description string `json:"description,omitempty"`
	Default     any    `json:"default,omitempty"`
	Enum        []any  `json:"enum,omitempty"`
	Example     any    `json:"example,omitempty"`
	Nullable    bool   `json:"nullable,omitempty"`
	ReadOnly    bool   `json:"readOnly,omitempty"`
	WriteOnly   bool   `json:"writeOnly,omitempty"`
}
