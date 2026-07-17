package spec

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// parseV3 parses an OpenAPI 3.x document into a unified Doc.
func parseV3(data []byte) (*Doc, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(data)
	if err != nil {
		return nil, fmt.Errorf("openapi 3 parse error: %w", err)
	}

	return openapi3DocToDoc(doc), nil
}

// openapi3DocToDoc converts a kin-openapi T to the unified Doc type.
func openapi3DocToDoc(doc *openapi3.T) *Doc {
	out := &Doc{
		Version: doc.OpenAPI,
	}

	if doc.Info != nil {
		out.Title = doc.Info.Title
		out.Description = doc.Info.Description
		out.VersionStr = doc.Info.Version
	}

	for _, s := range doc.Servers {
		out.Servers = append(out.Servers, Server{
			URL:         s.URL,
			Description: s.Description,
		})
	}

	for path, pathItem := range doc.Paths.Map() {
		ops := openapi3PathItemToOps(path, pathItem)
		out.PathItems = append(out.PathItems, ops...)
	}

	return out
}

// openapi3PathItemToOps converts a kin-openapi PathItem into a slice of PathItems (one per method).
func openapi3PathItemToOps(path string, item *openapi3.PathItem) []*PathItem {
	var out []*PathItem
	type entry struct {
		method string
		op     *openapi3.Operation
	}
	entries := []entry{
		{http.MethodGet, item.Get},
		{http.MethodPost, item.Post},
		{http.MethodPut, item.Put},
		{http.MethodDelete, item.Delete},
		{http.MethodPatch, item.Patch},
		{http.MethodHead, item.Head},
		{http.MethodOptions, item.Options},
		{http.MethodTrace, item.Trace},
	}
	for _, e := range entries {
		if e.op == nil {
			continue
		}
		out = append(out, &PathItem{
			Path:      path,
			Method:    e.method,
			Operation: openapi3OpToOp(e.op),
		})
	}
	return out
}

// openapi3OpToOp converts a kin-openapi Operation to the unified Operation type.
func openapi3OpToOp(op *openapi3.Operation) *Operation {
	o := &Operation{
		ID:          op.OperationID,
		Tags:        op.Tags,
		Summary:     op.Summary,
		Description: op.Description,
		Deprecated:  op.Deprecated,
		Parameters:  make([]*Parameter, 0, len(op.Parameters)),
		Responses:   make(map[string]*Response, op.Responses.Len()),
	}

	for _, pref := range op.Parameters {
		if pref == nil || pref.Value == nil {
			continue
		}
		p := pref.Value
		param := &Parameter{
			Name:        p.Name,
			In:          p.In,
			Description: p.Description,
			Required:    p.Required,
			Schema:      schemaRefToSchema(p.Schema),
		}
		o.Parameters = append(o.Parameters, param)
	}

	if op.RequestBody != nil && op.RequestBody.Value != nil {
		rb := op.RequestBody.Value
		o.RequestBody = &RequestBody{
			Description: rb.Description,
			Required:    rb.Required,
			Content:     openapi3ContentToContent(rb.Content),
		}
	}

	for code, rref := range op.Responses.Map() {
		if rref == nil || rref.Value == nil {
			continue
		}
		desc := ""
		if rref.Value.Description != nil {
			desc = *rref.Value.Description
		}
		r := &Response{
			Description: desc,
			Content:     openapi3ContentToContent(rref.Value.Content),
		}
		o.Responses[code] = r
	}

	return o
}

// openapi3ContentToContent converts kin-openapi Content to the unified MediaType map.
func openapi3ContentToContent(content openapi3.Content) map[string]*MediaType {
	if len(content) == 0 {
		return nil
	}
	out := make(map[string]*MediaType, len(content))
	for ct, mt := range content {
		out[ct] = &MediaType{
			Schema: schemaRefToSchema(mt.Schema),
		}
	}
	return out
}

// schemaRefToSchema converts a kin-openapi SchemaRef to the unified Schema type.
func schemaRefToSchema(sref *openapi3.SchemaRef) *Schema {
	if sref == nil || sref.Value == nil {
		return nil
	}
	s := sref.Value

	return &Schema{
		Type:        extractSchemaType(s),
		Format:      s.Format,
		Properties:  extractSchemaProperties(s),
		Items:       schemaRefToSchema(s.Items),
		Required:    s.Required,
		Ref:         sref.Ref,
		Description: s.Description,
		Default:     s.Default,
		Nullable:    s.Nullable,
		ReadOnly:    s.ReadOnly,
		WriteOnly:   s.WriteOnly,
		Example:     s.Example,
		OneOf:       extractSchemaComposition(s.OneOf),
		AnyOf:       extractSchemaComposition(s.AnyOf),
		AllOf:       extractSchemaComposition(s.AllOf),
	}
}

// extractSchemaType returns the first non-null type from a schema's type list.
func extractSchemaType(s *openapi3.Schema) string {
	if s.Type == nil {
		return ""
	}
	for _, t := range s.Type.Slice() {
		if !strings.EqualFold(t, "null") {
			return t
		}
	}
	return ""
}

// extractSchemaProperties converts a schema's property map to the unified Schema map.
func extractSchemaProperties(s *openapi3.Schema) map[string]*Schema {
	props := make(map[string]*Schema, len(s.Properties))
	for k, vref := range s.Properties {
		props[k] = schemaRefToSchema(vref)
	}
	return props
}

// extractSchemaComposition converts a slice of SchemaRefs to a slice of unified Schemas.
func extractSchemaComposition(refs []*openapi3.SchemaRef) []*Schema {
	out := make([]*Schema, 0, len(refs))
	for _, ss := range refs {
		if s := schemaRefToSchema(ss); s != nil {
			out = append(out, s)
		}
	}
	return out
}
