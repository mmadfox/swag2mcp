package spec

// Doc is the unified representation of a parsed Swagger/OpenAPI document.
// All versions (Swagger 2.0, OpenAPI 3.0, 3.1) are mapped to this type.
type Doc struct {
	Version     string // "2.0", "3.0", "3.1", etc.
	Title       string
	Description string
	VersionStr  string // version from info.version
	Servers     []Server
	PathItems   []*PathItem
}

type Server struct {
	URL         string
	Description string
}

// PathItem is a single endpoint (path + method + operation).
type PathItem struct {
	Path      string
	Method    string
	Operation *Operation
}

type Operation struct {
	ID          string
	Tags        []string
	Summary     string
	Description string
	Deprecated  bool
	Parameters  []*Parameter
	RequestBody *RequestBody
	Responses   map[string]*Response
}

type Parameter struct {
	Name            string
	In              string // "query", "path", "header", "cookie"
	Description     string
	Required        bool
	Schema          *Schema
}

type RequestBody struct {
	Description string
	Required    bool
	Content     map[string]*MediaType
}

type Response struct {
	Description string
	Content     map[string]*MediaType
}

type MediaType struct {
	Schema *Schema
}

type Schema struct {
	Type       string
	Format     string
	Properties map[string]*Schema
	Items      *Schema
	Required   []string
	Ref        string
	OneOf      []*Schema
	AnyOf      []*Schema
	AllOf      []*Schema

	Description string
	Default     any
	Enum        []any
	Example     any
	Nullable    bool
	ReadOnly    bool
	WriteOnly   bool
}