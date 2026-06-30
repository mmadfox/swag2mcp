package spec

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func parseV3(data []byte) (*Doc, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(data)
	if err != nil {
		return nil, fmt.Errorf("openapi 3 parse error: %w", err)
	}

	return openapi3DocToDoc(doc), nil
}

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

func openapi3PathItemToOps(path string, item *openapi3.PathItem) []*PathItem {
	var out []*PathItem
	type entry struct {
		method string
		op     *openapi3.Operation
	}
	entries := []entry{
		{"GET", item.Get},
		{"POST", item.Post},
		{"PUT", item.Put},
		{"DELETE", item.Delete},
		{"PATCH", item.Patch},
		{"HEAD", item.Head},
		{"OPTIONS", item.Options},
		{"TRACE", item.Trace},
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

func schemaRefToSchema(sref *openapi3.SchemaRef) *Schema {
	if sref == nil || sref.Value == nil {
		return nil
	}
	s := sref.Value

	// Type: *Types -> pick first non-null type
	typ := ""
	if s.Type != nil {
		types := s.Type.Slice()
		for _, t := range types {
			if !strings.EqualFold(t, "null") {
				typ = t
				break
			}
		}
	}

	props := make(map[string]*Schema, len(s.Properties))
	for k, vref := range s.Properties {
		props[k] = schemaRefToSchema(vref)
	}

	oneOf := make([]*Schema, 0, len(s.OneOf))
	for _, ss := range s.OneOf {
		if s := schemaRefToSchema(ss); s != nil {
			oneOf = append(oneOf, s)
		}
	}
	anyOf := make([]*Schema, 0, len(s.AnyOf))
	for _, ss := range s.AnyOf {
		if s := schemaRefToSchema(ss); s != nil {
			anyOf = append(anyOf, s)
		}
	}
	allOf := make([]*Schema, 0, len(s.AllOf))
	for _, ss := range s.AllOf {
		if s := schemaRefToSchema(ss); s != nil {
			allOf = append(allOf, s)
		}
	}

	return &Schema{
		Type:        typ,
		Format:      s.Format,
		Properties:  props,
		Items:       schemaRefToSchema(s.Items),
		Required:    s.Required,
		Ref:         sref.Ref,
		Description: s.Description,
		Default:     s.Default,
		Nullable:    s.Nullable,
		ReadOnly:    s.ReadOnly,
		WriteOnly:   s.WriteOnly,
		Example:     s.Example,
		OneOf:       oneOf,
		AnyOf:       anyOf,
		AllOf:       allOf,
	}
}