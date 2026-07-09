package spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	mediaTypeJSON = "application/json"
	paramTypeStr  = "string"
	paramTypeFile = "file"
	paramTypeObj  = "object"
	paramInQuery  = "query"
	paramInPath   = "path"
	bodyModeRaw   = "raw"
)

type postmanCollection struct {
	Info postmanInfo   `json:"info"`
	Item []postmanItem `json:"item"`
}

type postmanInfo struct {
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

type postmanItem struct {
	Name        string          `json:"name"`
	Description json.RawMessage `json:"description"`
	Request     *postmanRequest `json:"request,omitempty"`
	Item        []postmanItem   `json:"item,omitempty"`
}

type postmanRequest struct {
	Method string          `json:"method"`
	URL    json.RawMessage `json:"url"`
	Header []postmanHeader `json:"header,omitempty"`
	Body   *postmanBody    `json:"body,omitempty"`
}

type postmanHeader struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Disabled    bool   `json:"disabled,omitempty"`
	Description string `json:"description,omitempty"`
}

type postmanURL struct {
	Raw      string            `json:"raw"`
	Protocol string            `json:"protocol"`
	Host     []string          `json:"host"`
	Path     []json.RawMessage `json:"path"`
	Port     string            `json:"port"`
	Query    []postmanQueryVar `json:"query,omitempty"`
	Variable []postmanVariable `json:"variable,omitempty"`
}

type postmanQueryVar struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Disabled    bool   `json:"disabled,omitempty"`
	Description string `json:"description,omitempty"`
}

type postmanVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type postmanBody struct {
	Mode     string              `json:"mode"`
	Raw      string              `json:"raw,omitempty"`
	FormData []postmanFormData   `json:"formdata,omitempty"`
	URLEnc   []postmanURLEncoded `json:"urlencoded,omitempty"`
	Options  map[string]any      `json:"options,omitempty"`
	GraphQL  map[string]any      `json:"graphql,omitempty"`
}

type postmanFormData struct {
	Key         string `json:"key"`
	Value       string `json:"value,omitempty"`
	Type        string `json:"type,omitempty"`
	Src         any    `json:"src,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
	Description string `json:"description,omitempty"`
}

type postmanURLEncoded struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Disabled    bool   `json:"disabled,omitempty"`
	Description string `json:"description,omitempty"`
}

func parsePostman(data []byte) (*Doc, error) {
	var col postmanCollection
	if err := json.Unmarshal(data, &col); err != nil {
		return nil, fmt.Errorf("postman parse: %w", err)
	}

	doc := &Doc{
		Version:   "2.x",
		Title:     col.Info.Name,
		PathItems: make([]*PathItem, 0),
	}

	if err := flattenPostmanItems(nil, col.Item, doc); err != nil {
		return nil, err
	}

	if len(doc.PathItems) == 0 {
		return nil, errors.New("postman collection has no requests")
	}

	return doc, nil
}

func flattenPostmanItems(folderNames []string, items []postmanItem, doc *Doc) error {
	for _, item := range items {
		if item.Request != nil {
			pi := postmanItemToPathItem(item, folderNames)
			doc.PathItems = append(doc.PathItems, pi)
		}
		if len(item.Item) > 0 {
			names := append(folderNames, item.Name) //nolint:gocritic // intentional: avoid aliasing parent slice
			if err := flattenPostmanItems(names, item.Item, doc); err != nil {
				return err
			}
		}
	}
	return nil
}

func postmanItemToPathItem(item postmanItem, folderNames []string) *PathItem {
	req := item.Request

	method := strings.ToUpper(req.Method)
	if method == "" {
		method = http.MethodGet
	}

	apiPath := extractPostmanPath(req.URL)
	if apiPath == "" {
		apiPath = "/"
	}

	tag := postmanTag(item.Name, folderNames)

	op := &Operation{
		Summary:    item.Name,
		Tags:       []string{tag},
		Parameters: make([]*Parameter, 0),
		Responses:  make(map[string]*Response),
	}

	appendPostmanURLParams(req.URL, op)
	appendPostmanHeaders(req.Header, op)
	appendPostmanBody(req.Body, op, method)

	// Default response
	op.Responses["200"] = &Response{
		Description: "OK",
	}

	return &PathItem{
		Path:      apiPath,
		Method:    method,
		Operation: op,
	}
}

func extractPostmanPath(rawURL json.RawMessage) string {
	if rawURL == nil {
		return "/"
	}

	var rawStr string
	if err := json.Unmarshal(rawURL, &rawStr); err == nil {
		return extractPathFromURLString(rawStr)
	}

	var u postmanURL
	if err := json.Unmarshal(rawURL, &u); err != nil {
		return "/"
	}

	// Prefer structured path over raw — it handles path variables properly.
	if len(u.Path) > 0 {
		var segments []string
		for _, seg := range u.Path {
			var s string
			if err := json.Unmarshal(seg, &s); err == nil {
				segments = append(segments, s)
				continue
			}
			var pv struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			}
			if err := json.Unmarshal(seg, &pv); err == nil && pv.Value != "" {
				segments = append(segments, "{"+pv.Value+"}")
			}
		}
		if len(segments) > 0 {
			return "/" + strings.Join(segments, "/")
		}
	}

	// Fallback to raw
	if u.Raw != "" {
		return extractPathFromURLString(u.Raw)
	}

	return "/"
}

func extractPathFromURLString(rawURL string) string {
	if !strings.Contains(rawURL, "://") {
		rawURL = "http://" + rawURL
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	path := parsed.Path
	if path == "" {
		return "/"
	}
	// Convert Postman :param syntax to {param}
	segments := strings.Split(path, "/")
	for i, seg := range segments {
		if strings.HasPrefix(seg, ":") {
			segments[i] = "{" + seg[1:] + "}"
		}
	}
	return strings.Join(segments, "/")
}

func appendPostmanURLParams(rawURL json.RawMessage, op *Operation) {
	var u postmanURL
	if err := json.Unmarshal(rawURL, &u); err != nil {
		return
	}

	for _, v := range u.Variable {
		op.Parameters = append(op.Parameters, &Parameter{
			Name:     v.Key,
			In:       paramInPath,
			Required: true,
			Schema:   &Schema{Type: paramTypeStr},
		})
	}

	for _, q := range u.Query {
		if q.Disabled {
			continue
		}
		param := &Parameter{
			Name:        q.Key,
			In:          paramInQuery,
			Description: q.Description,
			Schema:      &Schema{Type: paramTypeStr},
		}
		op.Parameters = append(op.Parameters, param)
	}
}

func appendPostmanHeaders(headers []postmanHeader, op *Operation) {
	for _, h := range headers {
		if h.Disabled {
			continue
		}
		op.Parameters = append(op.Parameters, &Parameter{
			Name:        h.Key,
			In:          "header",
			Description: h.Description,
			Schema:      &Schema{Type: paramTypeStr, Default: h.Value},
		})
	}
}

func appendPostmanBody(body *postmanBody, op *Operation, method string) {
	if body == nil {
		return
	}

	switch method {
	case http.MethodGet, http.MethodHead, http.MethodDelete, http.MethodOptions:
		return
	}

	rb := &RequestBody{
		Content: make(map[string]*MediaType),
	}

	switch body.Mode {
	case bodyModeRaw:
		ct := guessPostmanContentType(body)
		rb.Content[ct] = &MediaType{
			Schema: &Schema{
				Type:        paramTypeStr,
				Description: body.Raw,
				Example:     body.Raw,
			},
		}

	case "urlencoded":
		props := make(map[string]*Schema)
		for _, f := range body.URLEnc {
			if f.Disabled {
				continue
			}
			props[f.Key] = &Schema{Type: paramTypeStr, Default: f.Value}
		}
		rb.Content["application/x-www-form-urlencoded"] = &MediaType{
			Schema: &Schema{
				Type:       "object",
				Properties: props,
			},
		}

	case "formdata":
		props := make(map[string]*Schema)
		for _, f := range body.FormData {
			if f.Disabled {
				continue
			}
			typ := paramTypeStr
			if f.Type == paramTypeFile {
				typ = paramTypeFile
			}
			props[f.Key] = &Schema{Type: typ, Default: f.Value}
		}
		rb.Content["multipart/form-data"] = &MediaType{
			Schema: &Schema{
				Type:       "object",
				Properties: props,
			},
		}

	case "graphql":
		rb.Content["application/json"] = &MediaType{
			Schema: &Schema{
				Type:        paramTypeObj,
				Description: "GraphQL query",
			},
		}
	}

	if len(rb.Content) > 0 {
		op.RequestBody = rb
	}
}

func postmanTag(itemName string, folders []string) string {
	if len(folders) > 0 {
		return sanitizePostmanTag(folders[len(folders)-1])
	}
	return sanitizePostmanTag(itemName)
}

func sanitizePostmanTag(name string) string {
	tag := strings.ToLower(name)
	tag = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}
		return '-'
	}, tag)
	return strings.Trim(tag, "-")
}

func guessPostmanContentType(body *postmanBody) string {
	if body == nil || body.Raw == "" {
		return mediaTypeJSON
	}
	raw := strings.TrimSpace(body.Raw)
	if strings.HasPrefix(raw, "{") || strings.HasPrefix(raw, "[") {
		return mediaTypeJSON
	}
	if strings.HasPrefix(raw, "<") {
		return "application/xml"
	}
	return "text/plain"
}

// isPostman checks if bytes represent a Postman collection.
func isPostman(data []byte) bool {
	var probe struct {
		Info struct {
			Schema string `json:"schema"`
		} `json:"info"`
		Item []any `json:"item"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return false
	}
	return probe.Info.Schema != "" &&
		strings.Contains(probe.Info.Schema, "getpostman.com") &&
		len(probe.Item) > 0
}
