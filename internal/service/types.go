package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"time"

	"github.com/mmadfox/swag2mcp/internal/spec"
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
	// ResponseOutline is the name of the response_outline tool.
	ResponseOutline = "response_outline"
	// ResponseCompress is the name of the response_compress tool.
	ResponseCompress = "response_compress"
	// ResponseSlice is the name of the response_slice tool.
	ResponseSlice = "response_slice"
)

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

// TagListItem represents a tag with its method count for display in lists.
type TagListItem struct {
	ID           string `json:"id"           jsonschema:"required,Unique identifier for the tag"`
	Title        string `json:"title"        jsonschema:"required,Human-readable title of the tag"`
	CountMethods int    `json:"countMethods" jsonschema:"required,Number of methods in the tag"`
}

// ToolInfo represents a single MCP tool with its name and description.
type ToolInfo struct {
	Name        string `json:"name"        jsonschema:"required,Name of the tool"`
	Description string `json:"description" jsonschema:"required,Description of the tool"`
}

// EndpointSearchItem represents an endpoint in the spec.
type EndpointSearchItem struct {
	ID              string `json:"id"              jsonschema:"required,Unique identifier for the endpoint"`
	TagID           string `json:"tagId"           jsonschema:"required,Unique identifier for the tag"`
	TagName         string `json:"tagName"         jsonschema:"required,Human-readable name of the tag"`
	CollectionID    string `json:"collectionId"    jsonschema:"required,Unique identifier for the collection"`
	CollectionTitle string `json:"collectionTitle" jsonschema:"required,Human-readable title of the collection"`
	SpecID          string `json:"specId"          jsonschema:"required,Unique identifier for the spec"`
	SpecDomain      string `json:"specDomain"      jsonschema:"required,Domain or category of the spec"`
	Method          string `json:"method"          jsonschema:"required,HTTP method (GET, POST, etc.)"`
	Path            string `json:"path"            jsonschema:"required,API path"`
	Summary         string `json:"summary"         jsonschema:"required,Human-readable summary of the endpoint"`
}

// EndpointTagItem represents an endpoint within a tag context.
type EndpointTagItem struct {
	ID      string `json:"id"      jsonschema:"required,Unique identifier for the endpoint"`
	Method  string `json:"method"  jsonschema:"required,HTTP method (GET, POST, etc.)"`
	Path    string `json:"path"    jsonschema:"required,API path"`
	Summary string `json:"summary" jsonschema:"required,Human-readable summary of the endpoint"`
}

// EndpointCollectionItem represents an endpoint within a collection context.
type EndpointCollectionItem struct {
	ID      string `json:"id"      jsonschema:"required,Unique identifier for the endpoint"`
	TagID   string `json:"tagId"   jsonschema:"required,Unique identifier for the tag"`
	TagName string `json:"tagName" jsonschema:"required,Human-readable name of the tag"`
	Method  string `json:"method"  jsonschema:"required,HTTP method (GET, POST, etc.)"`
	Path    string `json:"path"    jsonschema:"required,API path"`
	Summary string `json:"summary" jsonschema:"required,Human-readable summary of the endpoint"`
}

// Endpoint represents a minimal endpoint with method, path, and summary.
type Endpoint struct {
	ID      string `json:"id"      jsonschema:"required,Unique identifier for the endpoint"`
	Method  string `json:"method"  jsonschema:"required,HTTP method (GET, POST, etc.)"`
	Path    string `json:"path"    jsonschema:"required,API path"`
	Summary string `json:"summary" jsonschema:"required,Human-readable summary of the endpoint"`
}

// Spec is a specification like Openapi or Swagger.
type Spec struct {
	ID     string `json:"id"     jsonschema:"required,A unique 32-character MD5 hash identifier for the spec,pattern=^[0-9a-f]{32}$"`
	Domain string `json:"domain" jsonschema:"required,The domain or category of the spec,minLength=1"`
}

// SpecItem is a specification like Openapi or Swagger.
type SpecItem struct {
	ID     string `json:"id"     jsonschema:"required,A unique 32-character MD5 hash identifier for the spec,pattern=^[0-9a-f]{32}$"`
	Domain string `json:"domain" jsonschema:"required,The domain or category of the spec,minLength=1"`
}

// CollectionItem represents a collection in the spec.
type CollectionItem struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	LLMTitle     string `json:"llmTitle,omitempty"`
	CountTags    int    `json:"countTags"`
	CountMethods int    `json:"countMethods"`
}

// Collection represents a collection in the spec.
type Collection struct {
	ID           string `json:"id"           jsonschema:"required,Unique identifier for the collection"`
	Title        string `json:"title"        jsonschema:"required,Human-readable title of the collection"`
	CountMethods int    `json:"countMethods" jsonschema:"required,Number of methods in the collection"`
}

// SpecByIDRequest contains the spec ID used to look up a specific specification.
type SpecByIDRequest struct {
	ID string `json:"id" validate:"required,md5" jsonschema:"required,A unique 32-character MD5 hash identifier for the spec,pattern=^[0-9a-f]{32}$"`
}

// SpecByIDResponse contains the requested spec and its associated collections.
type SpecByIDResponse struct {
	Spec        Spec             `json:"spec"        jsonschema:"required,Specification"`
	Collections []CollectionItem `json:"collections" jsonschema:"required,List of collections associated with the spec"`
}

// SpecsResponse contains the list of all available specifications.
type SpecsResponse struct {
	Specs []SpecItem `json:"specs" jsonschema:"required,List of specifications"`
}

// CollectionsRequest represents a request to list all collections for a given spec.
type CollectionsRequest struct {
	SpecID string `json:"specId" jsonschema:"required," validate:"required,md5"`
}

// CollectionsResponse represents a response to list all collections for a given spec.
type CollectionsResponse struct {
	Spec        Spec             `json:"spec"        jsonschema:"required,Specification"`
	Collections []CollectionItem `json:"collections" jsonschema:"List of collections associated with the spec,required"`
}

// CollectionByIDRequest represents a request to get a collection by its ID.
type CollectionByIDRequest struct {
	ID string `json:"id" validate:"required,md5" jsonschema:"Unique identifier for the collection,required"`
}

// CollectionByIDResponse represents a response to get a collection by its ID.
type CollectionByIDResponse struct {
	Spec       Spec          `json:"spec"       jsonschema:"required,Specification"`
	Collection Collection    `json:"collection" jsonschema:"required,Collection"`
	Tags       []TagListItem `json:"tags"       jsonschema:"List of tags associated with the collection,required"`
}

// TagsByCollectionRequest represents a request to list all tags for a given collection.
type TagsByCollectionRequest struct {
	CollectionID string `json:"collectionId" jsonschema:"required," validate:"required,md5"`
}

// TagsByCollectionResponse represents a response to list all tags for a given collection.
type TagsByCollectionResponse struct {
	Spec       Spec          `json:"spec"       jsonschema:"required,Specification"`
	Collection Collection    `json:"collection" jsonschema:"required,Collection"`
	Tags       []TagListItem `json:"tags"       jsonschema:"required,List of tags associated with the collection"`
}

// TagByIDRequest represents a request to get a tag by its ID.
type TagByIDRequest struct {
	ID string `json:"id" validate:"required,md5" jsonschema:"required,Unique identifier for the tag"`
}

// TagByIDResponse represents a response to get a tag by its ID.
type TagByIDResponse struct {
	Tag TagListItem `json:"tag" jsonschema:"required,"`
}

// TagsBySpecRequest represents a request to list all tags for a given spec.
type TagsBySpecRequest struct {
	SpecID string `json:"specId" jsonschema:"required," validate:"required,md5"`
}

// TagsBySpecResponse represents a response to list all tags for a given spec.
type TagsBySpecResponse struct {
	Tags []TagListItem `json:"tags" jsonschema:"required,List of tags associated with the spec"`
}

// EndpointsByTagRequest contains the tag ID used to look up endpoints.
type EndpointsByTagRequest struct {
	TagID string `json:"tagId" jsonschema:"required," validate:"required,md5"`
}

// EndpointsByTagResponse contains the spec, collection, tag, and endpoints.
type EndpointsByTagResponse struct {
	Spec       Spec              `json:"spec"       jsonschema:"required,Specification"`
	Collection Collection        `json:"collection" jsonschema:"required,Collection"`
	Tag        TagListItem       `json:"tag"        jsonschema:"required,Tag"`
	Endpoints  []EndpointTagItem `json:"endpoints"  jsonschema:"required,List of endpoints associated with the tag"`
}

// EndpointsByCollectionRequest contains the collection ID used to look up endpoints.
type EndpointsByCollectionRequest struct {
	CollectionID string `json:"collectionId" jsonschema:"required," validate:"required,md5"`
}

// EndpointsByCollectionResponse contains the spec, collection, and endpoints.
type EndpointsByCollectionResponse struct {
	Spec       Spec                     `json:"spec"       jsonschema:"required,Specification"`
	Collection Collection               `json:"collection" jsonschema:"required,Collection"`
	Endpoints  []EndpointCollectionItem `json:"endpoints"  jsonschema:"required,List of endpoints associated with the collection"`
}

// EndpointsBySpecRequest contains the spec ID used to look up all endpoints.
type EndpointsBySpecRequest struct {
	SpecID string `json:"specId" jsonschema:"required," validate:"required,md5"`
}

// EndpointsBySpecResponse contains the list of endpoints associated with the spec.
type EndpointsBySpecResponse struct {
	Endpoints []EndpointSearchItem `json:"endpoints" jsonschema:"required,List of endpoints associated with the spec"`
}

// EndpointByIDRequest contains the unique endpoint ID to look up a single endpoint.
type EndpointByIDRequest struct {
	ID string `json:"id" validate:"required,md5" jsonschema:"required,Unique identifier for the endpoint"`
}

// EndpointByIDResponse contains the spec, collection, tag, and endpoint details.
type EndpointByIDResponse struct {
	Spec       Spec        `json:"spec"       jsonschema:"required,Specification"`
	Collection Collection  `json:"collection" jsonschema:"required,Collection"`
	Tag        TagListItem `json:"tag"        jsonschema:"required,Tag"`
	Endpoint   Endpoint    `json:"endpoint"   jsonschema:"required,"`
}

// SearchRequest contains the search query and result limit.
type SearchRequest struct {
	Query string `json:"query" jsonschema:"required,"                                    validate:"required"`
	Limit int    `json:"limit" jsonschema:"required,Maximum number of results to return" validate:"required,min=1,max=50"`
}

// SearchResponse contains the list of endpoints that matched the search query.
type SearchResponse struct {
	Endpoints []EndpointSearchItem `json:"endpoints" jsonschema:"required,List of endpoints matching the search query"`
}

// InspectRequest contains the endpoint ID used to retrieve full endpoint details.
type InspectRequest struct {
	EndpointID string `json:"endpointId" validate:"required,md5" jsonschema:"required,The 32-character MD5 hash ID of the endpoint to inspect"`
}

// InspectResponse contains the full details of an endpoint.
type InspectResponse struct {
	ID           string          `json:"id"                jsonschema:"required,Unique identifier for the endpoint"`
	TagID        string          `json:"tagId"             jsonschema:"required,Unique identifier for the tag"`
	CollectionID string          `json:"collectionId"      jsonschema:"required,Unique identifier for the collection"`
	SpecID       string          `json:"specId"            jsonschema:"required,Unique identifier for the spec"`
	SpecDomain   string          `json:"specDomain"        jsonschema:"required,Domain of the spec"`
	Method       string          `json:"method"            jsonschema:"required,HTTP method (GET, POST, etc.)"`
	Path         string          `json:"path"              jsonschema:"required,API path"`
	BaseURL      string          `json:"baseUrl"           jsonschema:"required,Base URL of the API"`
	FullURL      string          `json:"fullUrl"           jsonschema:"required,Full URL of the endpoint"`
	Operation    *spec.Operation `json:"operation"         jsonschema:"required,Operation details"`
}

// AuthRequest contains the parameters needed to retrieve authentication information.
type AuthRequest struct {
	SpecID string `json:"specId" validate:"required,md5" jsonschema:"required,The 32-character MD5 hash ID of the spec/domain to get an auth token for"`
}

// AuthResponse contains the authentication token, headers, and query parameters.
type AuthResponse struct {
	Token       string            `json:"token"`
	Headers     map[string]string `json:"headers,omitempty"`
	QueryParams map[string]string `json:"queryParams,omitempty"`
}

// InfoSnapshot is a point-in-time snapshot of the service state,
// computed once after Bootstrap and served on every Info call.
type InfoSnapshot struct {
	Version    string         `json:"version"`
	Workspace  string         `json:"workspace"`
	Uptime     time.Duration  `json:"-"`
	Specs      SpecsSummary   `json:"specs"`
	HTTPClient HTTPClientInfo `json:"http_client"`
	MCP        MCPInfo        `json:"mcp"`
	Auth       AuthInfo       `json:"auth"`
	Mock       MockInfo       `json:"mock"`
}

// InfoResponse holds the complete runtime and configuration summary.
type InfoResponse struct {
	Version    string         `json:"version,omitempty"`
	Workspace  string         `json:"workspace"`
	Uptime     string         `json:"uptime,omitempty"`
	Specs      SpecsSummary   `json:"specs"`
	HTTPClient HTTPClientInfo `json:"http_client"`
	MCP        MCPInfo        `json:"mcp"`
	Auth       AuthInfo       `json:"auth"`
	Mock       MockInfo       `json:"mock"`
}

// SpecsSummary holds aggregate spec statistics.
type SpecsSummary struct {
	Total       int `json:"total"`
	Active      int `json:"active"`
	Disabled    int `json:"disabled"`
	Collections int `json:"collections"`
	Endpoints   int `json:"endpoints"`
}

// HTTPClientInfo holds the effective HTTP client configuration.
type HTTPClientInfo struct {
	Randomize       bool              `json:"randomize"`
	UserAgent       string            `json:"user_agent,omitempty"`
	Timeout         string            `json:"timeout,omitempty"`
	FollowRedirects *bool             `json:"follow_redirects,omitempty"`
	MaxRedirects    *int              `json:"max_redirects,omitempty"`
	MaxResponseSize string            `json:"max_response_size"`
	Proxy           *ProxyInfo        `json:"proxy,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
	Cookies         []CookieInfo      `json:"cookies,omitempty"`
}

// ProxyInfo holds proxy configuration details.
type ProxyInfo struct {
	URL      string   `json:"url"`
	Username string   `json:"username,omitempty"`
	Bypass   []string `json:"bypass,omitempty"`
}

// CookieInfo holds a single cookie configuration.
type CookieInfo struct {
	Name     string `json:"name"`
	Domain   string `json:"domain,omitempty"`
	Path     string `json:"path,omitempty"`
	Secure   bool   `json:"secure"`
	HTTPOnly bool   `json:"http_only"`
}

// MCPInfo holds the MCP server configuration.
type MCPInfo struct {
	Transport   string `json:"transport"`
	Addr        string `json:"addr,omitempty"`
	Path        string `json:"path,omitempty"`
	AuthEnabled bool   `json:"auth_enabled"`
}

// AuthInfo holds authentication method information.
type AuthInfo struct {
	Methods []string `json:"methods,omitempty"`
}

// MockInfo holds mock server configuration.
type MockInfo struct {
	Enabled bool `json:"enabled"`
}
