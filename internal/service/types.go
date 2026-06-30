package service

type TagListItem struct {
	ID           string `json:"id" jsonschema:"required,Unique identifier for the tag"`
	Title        string `json:"title" jsonschema:"required,Human-readable title of the tag"`
	CountMethods int    `json:"countMethods" jsonschema:"required,Number of methods in the tag"`
}

type ToolInfo struct {
	Name        string `json:"name" jsonschema:"required,Name of the tool"`
	Description string `json:"description" jsonschema:"required,Description of the tool"`
}

// EndpointItem represents an endpoint in the spec.
type EndpointItem struct {
	ID           string `json:"id" jsonschema:"required,Unique identifier for the endpoint"`
	TagID        string `json:"tagId" jsonschema:"required,Unique identifier for the tag"`
	CollectionID string `json:"collectionId" jsonschema:"required,Unique identifier for the collection"`
	SpecID       string `json:"specId" jsonschema:"required,Unique identifier for the spec"`
	Method       string `json:"method" jsonschema:"required,HTTP method (GET, POST, etc.)"`
	Path         string `json:"path" jsonschema:"required,API path"`
	Summary      string `json:"summary" jsonschema:"required,Human-readable summary of the endpoint"`
	Deprecated   bool   `json:"deprecated" jsonschema:"required,Whether the endpoint is deprecated"`
}

// Spec is a specification like Openapi or Swagger.
type Spec struct {
	ID     string `json:"id" jsonschema:"required,A unique 32-character MD5 hash identifier for the spec,pattern=^[0-9a-f]{32}$"`
	Domain string `json:"domain" jsonschema:"required,The domain or category of the spec,minLength=1"`
}

// SpecItem is a specification like Openapi or Swagger.
type SpecItem struct {
	ID     string `json:"id" jsonschema:"required,A unique 32-character MD5 hash identifier for the spec,pattern=^[0-9a-f]{32}$"`
	Domain string `json:"domain" jsonschema:"required,The domain or category of the spec,minLength=1"`
}

// CollectionItem represents a collection in the spec.
type CollectionItem struct {
	ID           string `json:"id" jsonschema:"required,Unique identifier for the collection"`
	Title        string `json:"title" jsonschema:"required,Human-readable title of the collection"`
	CountTags    int    `json:"countTags" jsonschema:"required,Number of tags in the collection"`
	CountMethods int    `json:"countMethods" jsonschema:"required,Number of methods in the collection"`
}

// Collection represents a collection in the spec.
type Collection struct {
	ID           string `json:"id" jsonschema:"required,Unique identifier for the collection"`
	Title        string `json:"title" jsonschema:"required,Human-readable title of the collection"`
	CountMethods int    `json:"countMethods" jsonschema:"required,Number of methods in the collection"`
}
