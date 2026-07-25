package spec

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-openapi/spec"
)

// parseV2 parses a Swagger 2.0 document into a unified Doc.
func parseV2(data []byte) (*Doc, error) {
	var swag spec.Swagger
	if err := swag.UnmarshalJSON(data); err != nil {
		return nil, fmt.Errorf("swagger 2.0 parse error: %w", err)
	}

	opts := spec.ExpandOptions{ContinueOnError: true}
	if err := spec.ExpandSpec(&swag, &opts); err != nil {
		return nil, fmt.Errorf("swagger 2.0 expand error: %w", err)
	}

	doc := &Doc{
		Version:     specVersion20,
		Title:       swag.Info.Title,
		Description: swag.Info.Description,
		VersionStr:  swag.Info.Version,
	}

	if swag.Host != "" {
		scheme := "https"
		if len(swag.Schemes) > 0 {
			scheme = swag.Schemes[0]
		}
		doc.Servers = append(doc.Servers, Server{
			URL: fmt.Sprintf("%s://%s%s", scheme, swag.Host, swag.BasePath),
		})
	}

	for path, pathItem := range swag.Paths.Paths {
		doc.PathItems = append(doc.PathItems, pathItemToOps(path, pathItem)...)
	}

	return doc, nil
}

// pathItemToOps converts a Swagger path item into a slice of PathItems (one per method).
func pathItemToOps(path string, item spec.PathItem) []*PathItem {
	var out []*PathItem
	type entry struct {
		method string
		op     *spec.Operation
	}
	entries := []entry{
		{http.MethodGet, item.Get},
		{http.MethodPost, item.Post},
		{http.MethodPut, item.Put},
		{http.MethodDelete, item.Delete},
		{http.MethodPatch, item.Patch},
		{http.MethodHead, item.Head},
		{http.MethodOptions, item.Options},
	}
	for _, e := range entries {
		if e.op == nil {
			continue
		}
		out = append(out, &PathItem{
			Path:      path,
			Method:    e.method,
			Operation: swaggerOpToOp(e.op),
		})
	}
	return out
}

// swaggerOpToOp converts a go-openapi Operation to the unified Operation type.
func swaggerOpToOp(op *spec.Operation) *Operation {
	o := &Operation{
		ID:          op.ID,
		Tags:        op.Tags,
		Summary:     op.Summary,
		Description: op.Description,
		Deprecated:  op.Deprecated,
		Parameters:  make([]*Parameter, 0, len(op.Parameters)),
		Responses:   make(map[string]*Response),
	}

	for _, p := range op.Parameters {
		param := &Parameter{
			Name:        p.Name,
			In:          p.In,
			Description: p.Description,
			Required:    p.Required,
			Schema:      swaggerSchemaToSchema(p.Schema),
		}
		if param.Schema == nil && p.Type != "" {
			param.Schema = &Schema{
				Type:    p.Type,
				Format:  p.Format,
				Default: p.Default,
			}
		}
		o.Parameters = append(o.Parameters, param)
	}

	consumes := firstConsumes(op.Consumes)

	// Swagger 2.0: body param from parameters list
	if bodyParam := findBodyParam(op.Parameters); bodyParam != nil {
		o.RequestBody = &RequestBody{
			Required: bodyParam.Required,
			Content: map[string]*MediaType{
				consumes: {Schema: swaggerSchemaToSchema(bodyParam.Schema)},
			},
		}
	}

	if op.Responses != nil {
		for code, resp := range op.Responses.StatusCodeResponses {
			r := &Response{
				Description: resp.Description,
				Content: map[string]*MediaType{
					consumes: {Schema: swaggerSchemaToSchema(resp.Schema)},
				},
			}
			o.Responses[strconv.Itoa(code)] = r
		}
		if op.Responses.Default != nil {
			r := &Response{
				Description: op.Responses.Default.Description,
				Content: map[string]*MediaType{
					consumes: {Schema: swaggerSchemaToSchema(op.Responses.Default.Schema)},
				},
			}
			o.Responses["default"] = r
		}
	}

	return o
}

// findBodyParam finds the body parameter from a Swagger parameter list.
func findBodyParam(params []spec.Parameter) *spec.Parameter {
	for i := range params {
		if params[i].In == "body" {
			return &params[i]
		}
	}
	return nil
}

// firstConsumes returns the first content type from the consumes list, defaulting to JSON.
func firstConsumes(consumes []string) string {
	if len(consumes) > 0 {
		return consumes[0]
	}
	return mediaTypeJSON
}

// swaggerSchemaToSchema converts a go-openapi Schema to the unified Schema type.
func swaggerSchemaToSchema(s *spec.Schema) *Schema {
	if s == nil {
		return nil
	}
	ref := ""
	if s.Ref.String() != "" {
		ref = s.Ref.String()
	}

	typ := ""
	if len(s.Type) > 0 {
		typ = s.Type[0]
	}

	var items *Schema
	if s.Items != nil && s.Items.Schema != nil {
		items = swaggerSchemaToSchema(s.Items.Schema)
	}

	props := make(map[string]*Schema, len(s.Properties))
	for k, v := range s.Properties {
		props[k] = swaggerSchemaToSchema(&v)
	}

	oneOf := make([]*Schema, 0)
	for _, ss := range s.OneOf {
		if s := swaggerSchemaToSchema(&ss); s != nil {
			oneOf = append(oneOf, s)
		}
	}
	anyOf := make([]*Schema, 0)
	for _, ss := range s.AnyOf {
		if s := swaggerSchemaToSchema(&ss); s != nil {
			anyOf = append(anyOf, s)
		}
	}
	allOf := make([]*Schema, 0)
	for _, ss := range s.AllOf {
		if s := swaggerSchemaToSchema(&ss); s != nil {
			allOf = append(allOf, s)
		}
	}

	return &Schema{
		Type:        typ,
		Format:      s.Format,
		Properties:  props,
		Items:       items,
		Required:    s.Required,
		Ref:         ref,
		Description: s.Description,
		Default:     s.Default,
		Nullable:    s.Nullable,
		ReadOnly:    s.ReadOnly,
		Enum:        s.Enum,
		Example:     s.Example,
		OneOf:       oneOf,
		AnyOf:       anyOf,
		AllOf:       allOf,
	}
}
