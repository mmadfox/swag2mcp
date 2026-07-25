package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/spec"
)

const (
	schemaTypeObject  = "object"
	schemaTypeArray   = "array"
	contentTypeJSON   = "application/json"
	acceptHeaderJSON  = "application/json, text/plain, */*"
	acceptHeaderOther = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
)

// validateParameters checks that all required parameters are present and that no
// unknown parameters are passed.
func validateParameters(op *spec.Operation, params map[string]any) error {
	paramNames := make(map[string]struct{}, len(op.Parameters))
	for _, p := range op.Parameters {
		paramNames[p.Name] = struct{}{}
	}

	for name := range params {
		if _, exists := paramNames[name]; !exists {
			return fmt.Errorf(
				"unknown parameter %q, all parameters must match the operation schema",
				name,
			)
		}
	}

	var missing []string
	for _, p := range op.Parameters {
		if !p.Required {
			continue
		}
		if _, exists := params[p.Name]; !exists {
			missing = append(missing, p.Name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf(
			"missing required parameters: %s",
			strings.Join(missing, ", "),
		)
	}

	return nil
}

// validateRequestBody validates a request body against the operation's request body schema.
func validateRequestBody(op *spec.Operation, body map[string]any) error {
	if op.RequestBody == nil {
		return nil
	}

	if op.RequestBody.Required && body == nil {
		return errors.New("request body is required for this endpoint")
	}

	if body == nil {
		return nil
	}

	sc := schemaForContentType(op.RequestBody.Content)
	if sc == nil {
		return nil
	}

	return validateSchemaValue(sc, body, "$")
}

// schemaForContentType extracts the JSON schema from a content map, preferring application/json.
func schemaForContentType(ct map[string]*spec.MediaType) *spec.Schema {
	if ct == nil {
		return nil
	}
	mt, exists := ct["application/json"]
	if !exists || mt == nil {
		return nil
	}
	return mt.Schema
}

// validateSchemaValue recursively validates a value against a schema path.
func validateSchemaValue(sc *spec.Schema, value any, path string) error {
	if sc == nil {
		return nil
	}

	switch sc.Type {
	case schemaTypeObject:
		return validateObjectSchema(sc, value, path)
	case schemaTypeArray:
		return validateArraySchema(sc, value, path)
	}

	return nil
}

// validateObjectSchema validates a map value against an object schema.
func validateObjectSchema(sc *spec.Schema, value any, path string) error {
	obj, ok := value.(map[string]any)
	if !ok {
		return nil
	}

	for _, requiredField := range sc.Required {
		if _, exists := obj[requiredField]; !exists {
			return fmt.Errorf("missing required field %q at %s", requiredField, path)
		}
	}

	for key := range obj {
		if _, defined := sc.Properties[key]; !defined {
			return fmt.Errorf(
				"unknown field %q at %s, all fields must match the schema",
				key, path,
			)
		}
	}

	for key, ps := range sc.Properties {
		cv, exists := obj[key]
		if !exists {
			continue
		}
		cp := path + "." + key
		if err := validateSchemaValue(ps, cv, cp); err != nil {
			return err
		}
	}

	return nil
}

// validateArraySchema validates a slice value against an array schema.
func validateArraySchema(sc *spec.Schema, value any, path string) error {
	arr, ok := value.([]any)
	if !ok {
		return nil
	}

	for i, item := range arr {
		cp := fmt.Sprintf("%s[%d]", path, i)
		if err := validateSchemaValue(sc.Items, item, cp); err != nil {
			return err
		}
	}

	return nil
}
