/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

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

	// Propagate top-level security to operations that do not define their own.
	// OpenAPI 3.x allows a global "security" that applies to all operations.
	applyTopLevelSecurity(out.PathItems, securityRequirementsToMaps(doc.Security))

	return out
}

// securityRequirementsToMaps converts kin-openapi SecurityRequirements to the
// unified []map[string][]string representation.
func securityRequirementsToMaps(reqs openapi3.SecurityRequirements) []map[string][]string {
	if len(reqs) == 0 {
		return nil
	}
	out := make([]map[string][]string, len(reqs))
	for i, req := range reqs {
		out[i] = req
	}
	return out
}

// openapi3PathItemToOps converts a kin-openapi PathItem into a slice of PathItems (one per method).
// Path-level parameters are merged into each operation. Operation-level parameters with the same
// name+in override path-level parameters, as per the OpenAPI spec.
func openapi3PathItemToOps(path string, item *openapi3.PathItem) []*PathItem {
	pathParams := openapi3ParamsToParams(item.Parameters)

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
		op := openapi3OpToOp(e.op)
		if len(pathParams) > 0 {
			op.Parameters = mergeParameters(pathParams, op.Parameters)
		}
		out = append(out, &PathItem{
			Path:      path,
			Method:    e.method,
			Operation: op,
		})
	}
	return out
}

// openapi3ParamsToParams converts a slice of kin-openapi ParameterRef to unified Parameters.
func openapi3ParamsToParams(refs []*openapi3.ParameterRef) []*Parameter {
	out := make([]*Parameter, 0, len(refs))
	for _, pref := range refs {
		if pref == nil || pref.Value == nil {
			continue
		}
		p := pref.Value
		out = append(out, &Parameter{
			Name:        p.Name,
			In:          p.In,
			Description: p.Description,
			Required:    p.Required,
			Schema:      schemaRefToSchema(p.Schema),
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

	if op.Security != nil {
		o.securityDeclared = true
		if len(*op.Security) > 0 {
			o.Security = make([]map[string][]string, len(*op.Security))
			for i, sec := range *op.Security {
				o.Security[i] = sec
			}
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
	return schemaRefToSchemaVisited(sref, make(map[*openapi3.Schema]struct{}))
}

// schemaRefToSchemaVisited converts a SchemaRef to the unified Schema type,
// tracking visited schema pointers to break circular $ref cycles. kin-openapi
// inlines references, so a self-referencing schema yields the same pointer
// repeatedly; without cycle detection this recurses until stack overflow.
func schemaRefToSchemaVisited(sref *openapi3.SchemaRef, visited map[*openapi3.Schema]struct{}) *Schema {
	if sref == nil || sref.Value == nil {
		return nil
	}
	s := sref.Value
	if _, ok := visited[s]; ok {
		// Circular reference: emit a schema carrying only the ref so the
		// structure is preserved without recursing further.
		return &Schema{Ref: sref.Ref}
	}
	// Mark the schema as in-progress for the duration of this subtree, then
	// unmark it on the way out (DFS backtracking). This detects cycles along
	// the current path while still allowing the same schema to be expanded
	// again from a different, non-cyclic branch.
	visited[s] = struct{}{}
	defer delete(visited, s)

	return &Schema{
		Type:        extractSchemaType(s),
		Format:      s.Format,
		Properties:  extractSchemaProperties(s, visited),
		Items:       schemaRefToSchemaVisited(s.Items, visited),
		Required:    s.Required,
		Ref:         sref.Ref,
		Description: s.Description,
		Default:     s.Default,
		Nullable:    s.Nullable,
		ReadOnly:    s.ReadOnly,
		WriteOnly:   s.WriteOnly,
		Example:     s.Example,
		OneOf:       extractSchemaComposition(s.OneOf, visited),
		AnyOf:       extractSchemaComposition(s.AnyOf, visited),
		AllOf:       extractSchemaComposition(s.AllOf, visited),
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
func extractSchemaProperties(s *openapi3.Schema, visited map[*openapi3.Schema]struct{}) map[string]*Schema {
	props := make(map[string]*Schema, len(s.Properties))
	for k, vref := range s.Properties {
		props[k] = schemaRefToSchemaVisited(vref, visited)
	}
	return props
}

// extractSchemaComposition converts a slice of SchemaRefs to a slice of unified Schemas.
func extractSchemaComposition(refs []*openapi3.SchemaRef, visited map[*openapi3.Schema]struct{}) []*Schema {
	out := make([]*Schema, 0, len(refs))
	for _, ss := range refs {
		if s := schemaRefToSchemaVisited(ss, visited); s != nil {
			out = append(out, s)
		}
	}
	return out
}

// mergeParameters merges base (path-level) parameters into op (operation-level) parameters.
// Operation-level parameters with the same name+in override base parameters.
func mergeParameters(base, op []*Parameter) []*Parameter {
	seen := make(map[string]int, len(op))
	for i, p := range op {
		seen[p.Name+"\x00"+p.In] = i
	}

	merged := make([]*Parameter, 0, len(base)+len(op))
	for _, p := range base {
		if _, exists := seen[p.Name+"\x00"+p.In]; exists {
			continue
		}
		merged = append(merged, p)
	}
	merged = append(merged, op...)
	return merged
}
