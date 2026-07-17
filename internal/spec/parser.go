package spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"go.yaml.in/yaml/v3"
)

// toJSON converts raw bytes (JSON or YAML) to normalized JSON.
func toJSON(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty document")
	}
	// Fast path: already JSON
	trimmed := strings.TrimLeft(string(data), " \t\r\n")
	if trimmed != "" && trimmed[0] == '{' {
		// Postman collections are always JSON. Check before parsing.
		if isPostman(data) {
			return data, nil
		}
		return data, nil
	}
	// YAML path — Postman collections are never YAML, so skip postman check.
	var raw any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("not valid YAML: %w", err)
	}
	clean := cleanupYAML(raw)
	jsonData, err := json.Marshal(clean)
	if err != nil {
		return nil, fmt.Errorf("yaml to json conversion: %w", err)
	}
	return jsonData, nil
}

// Parse parses a Swagger 2.0, OpenAPI 3.x, or Postman collection document (JSON or YAML)
// and returns a unified Doc.
func Parse(data []byte) (*Doc, error) {
	jsonData, err := toJSON(data)
	if err != nil {
		return nil, fmt.Errorf("spec parse: %w", err)
	}

	// Postman detection first — Postman collections don't have swagger/openapi fields.
	if isPostman(jsonData) {
		return parsePostman(jsonData)
	}

	version, err := detectVersion(jsonData)
	if err != nil {
		return nil, fmt.Errorf("spec parse: %w", err)
	}

	switch {
	case version == specVersion20:
		return parseV2(jsonData)
	case strings.HasPrefix(version, "3."):
		jsonData = preprocessV3(jsonData)
		return parseV3(jsonData)
	default:
		return nil, fmt.Errorf("unsupported spec version %q", version)
	}
}

// cleanupYAML converts yaml.v3's map[interface{}]interface{} to map[string]interface{}.
func cleanupYAML(v any) any {
	switch vv := v.(type) {
	case map[string]any:
		for k, val := range vv {
			vv[k] = cleanupYAML(val)
		}
		return vv
	case map[any]any:
		m := make(map[string]any, len(vv))
		for k, val := range vv {
			key := fmt.Sprint(k)
			m[key] = cleanupYAML(val)
		}
		return m
	case []any:
		for i, val := range vv {
			vv[i] = cleanupYAML(val)
		}
	}
	return v
}

// preprocessV3 applies workarounds for known kin-openapi limitations.
func preprocessV3(data []byte) []byte {
	// OpenAPI 3.1 allows items: false (meaning "no items allowed").
	// kin-openapi expects items: {...} or items: null.
	// Replace `"items":false` and `"items": true` with a no-op empty schema ref.
	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		return data
	}
	fixItems(doc)
	fixed, err := json.Marshal(doc)
	if err != nil {
		return data
	}
	return fixed
}

// fixItems replaces boolean items values with empty schema objects for kin-openapi compat.
func fixItems(v any) {
	switch node := v.(type) {
	case map[string]any:
		if items, ok := node["items"]; ok {
			if _, isBool := items.(bool); isBool {
				node["items"] = map[string]any{}
			}
		}
		for _, val := range node {
			fixItems(val)
		}
	case []any:
		for _, val := range node {
			fixItems(val)
		}
	}
}

type versionDoc struct {
	Swagger string `json:"swagger"`
	OpenAPI string `json:"openapi"`
}

// detectVersion reads the swagger or openapi field to determine the spec version.
func detectVersion(data []byte) (string, error) {
	var v versionDoc
	if err := json.Unmarshal(data, &v); err != nil {
		return "", fmt.Errorf("version detection: %w", err)
	}
	// Prefer swagger (backward compat), then openapi
	switch {
	case v.Swagger != "":
		return v.Swagger, nil
	case v.OpenAPI != "":
		return v.OpenAPI, nil
	default:
		return "", errors.New("cannot detect spec version: missing 'swagger' or 'openapi' field")
	}
}
