package reader

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
)

const (
	compressedTypeField   = "type"
	compressedTypeArray   = "array"
	compressedNoteField   = "note"
	compressedItemField   = "item"
	compressedLengthField = "length"
)

// Compress reduces the JSON value selected by jsonPath and mode.
func (r *reader) Compress(path string, opts CompressOptions) (CompressResult, error) {
	if err := r.validatePath(path); err != nil {
		return CompressResult{}, err
	}

	opts = normalizeCompressOptions(opts)

	data, err := os.ReadFile(path)
	if err != nil {
		return CompressResult{}, fmt.Errorf("read file: %w", err)
	}

	selected := data
	if opts.JSONPath != "" {
		result := gjson.GetBytes(data, opts.JSONPath)
		if !result.Exists() {
			return CompressResult{}, fmt.Errorf("%w: %s", ErrPathNotFound, opts.JSONPath)
		}
		selected = []byte(result.Raw)
	}

	compressed, err := r.compressValue(selected, opts)
	if err != nil {
		return CompressResult{}, fmt.Errorf("compress value: %w", err)
	}

	if opts.Limit > 0 && len(compressed) > opts.Limit {
		return CompressResult{
			TooLarge: true,
			Hint: fmt.Sprintf(
				"Compressed result still exceeds the %d byte limit. Use response_slice with a deeper jsonPath.",
				opts.Limit),
		}, nil
	}

	var body any
	if err := json.Unmarshal(compressed, &body); err != nil {
		body = string(compressed)
	}

	return CompressResult{
		Body: body,
		Hint: fmt.Sprintf("Compressed using mode %s%s.", opts.Mode, jsonPathHint(opts.JSONPath)),
	}, nil
}

// jsonPathHint returns a human-readable suffix when jsonPath is set.
func jsonPathHint(path string) string {
	if path == "" {
		return ""
	}
	return " at " + path
}

// compressValue applies the selected compression mode to raw JSON bytes.
func (r *reader) compressValue(data []byte, opts CompressOptions) ([]byte, error) {
	switch opts.Mode {
	case CompressFirstOfArray:
		return compressFirstOfArray(data)
	case CompressSampleArray:
		return compressSampleArray(data, opts.ArrayHead, opts.ArrayTail)
	case CompressTruncateStrings:
		return compressTruncateStrings(data, opts.StringLen)
	case CompressKeysOnly:
		return compressKeysOnly(data)
	case CompressSelectKeys:
		return compressSelectKeys(data, opts.SelectKeys)
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidJSONPath, opts.Mode)
	}
}

// compressFirstOfArray keeps only the first element of an array.
func compressFirstOfArray(data []byte) ([]byte, error) {
	var head []any
	decoder := json.NewDecoder(bytes.NewReader(data))
	token, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrNotJSON, err)
	}
	delim, ok := token.(json.Delim)
	if !ok || delim != '[' {
		return nil, errors.New("first_of_array requires an array value")
	}

	if !decoder.More() {
		return []byte("[]"), nil
	}

	var first any
	if err := decoder.Decode(&first); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrNotJSON, err)
	}
	head = append(head, first)

	wrapped := map[string]any{
		"compressed": map[string]any{
			compressedTypeField:   compressedTypeArray,
			compressedLengthField: 1,
			compressedNoteField:   "Only the first item of the original array is kept.",
			compressedItemField:   head[0],
		},
	}
	return json.Marshal(wrapped)
}

// compressSampleArray keeps head and tail items of an array.
func compressSampleArray(data []byte, head, tail int) ([]byte, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	token, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrNotJSON, err)
	}
	delim, ok := token.(json.Delim)
	if !ok || delim != '[' {
		return nil, errors.New("sample_array requires an array value")
	}

	var items []any
	skipped := 0
	for decoder.More() {
		var item any
		if err := decoder.Decode(&item); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrNotJSON, err)
		}
		items = append(items, item)
	}

	total := len(items)
	var sample []any
	if total <= head+tail {
		sample = items
	} else {
		sample = make([]any, 0, head+tail)
		sample = append(sample, items[:head]...)
		sample = append(sample, items[total-tail:]...)
		skipped = total - head - tail
	}

	wrapped := map[string]any{
		"compressed": map[string]any{
			compressedTypeField:   compressedTypeArray,
			compressedLengthField: len(sample),
			"original":            total,
			"skipped":             skipped,
			"sample":              sample,
			"sampleRange":         fmt.Sprintf("items [0-%d] and [%d-%d]", head-1, total-tail, total-1),
		},
	}
	return json.Marshal(wrapped)
}

// compressTruncateStrings shortens every string value in JSON.
func compressTruncateStrings(data []byte, maxLen int) ([]byte, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	value, err := decodeAndTruncate(decoder, maxLen)
	if err != nil {
		return nil, err
	}
	return json.Marshal(value)
}

// decodeAndTruncate recursively decodes and truncates strings.
func decodeAndTruncate(decoder *json.Decoder, maxLen int) (any, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	switch t := token.(type) {
	case json.Delim:
		return decodeAndTruncateContainer(decoder, t, maxLen)
	case string:
		return truncateString(t, maxLen), nil
	case json.Number:
		return t, nil
	default:
		return t, nil
	}
}

// decodeAndTruncateContainer handles object/array containers for decodeAndTruncate.
func decodeAndTruncateContainer(decoder *json.Decoder, delim json.Delim, maxLen int) (any, error) {
	if delim == '{' {
		return decodeAndTruncateObject(decoder, maxLen)
	}
	return decodeAndTruncateArray(decoder, maxLen)
}

// decodeAndTruncateObject decodes an object and truncates its string values.
func decodeAndTruncateObject(decoder *json.Decoder, maxLen int) (map[string]any, error) {
	obj := make(map[string]any)
	for decoder.More() {
		key, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		k, ok := key.(string)
		if !ok {
			return nil, ErrNotJSON
		}
		val, err := decodeAndTruncate(decoder, maxLen)
		if err != nil {
			return nil, err
		}
		obj[k] = val
	}
	if _, err := decoder.Token(); err != nil {
		return nil, err
	}
	return obj, nil
}

// decodeAndTruncateArray decodes an array and truncates its string values.
func decodeAndTruncateArray(decoder *json.Decoder, maxLen int) ([]any, error) {
	arr := []any{}
	for decoder.More() {
		val, err := decodeAndTruncate(decoder, maxLen)
		if err != nil {
			return nil, err
		}
		arr = append(arr, val)
	}
	if _, err := decoder.Token(); err != nil {
		return nil, err
	}
	return arr, nil
}

// truncateString shortens a string, adding an ellipsis.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// compressKeysOnly replaces object values with type-only placeholders.
func compressKeysOnly(data []byte) ([]byte, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	value, err := decodeKeysOnly(decoder)
	if err != nil {
		return nil, err
	}
	return json.Marshal(value)
}

// decodeKeysOnly recursively replaces values with type names.
func decodeKeysOnly(decoder *json.Decoder) (any, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	switch t := token.(type) {
	case json.Delim:
		return decodeKeysOnlyContainer(decoder, t)
	case string:
		return jsonTypeString, nil
	case json.Number:
		return jsonTypeNumber, nil
	case bool:
		return jsonTypeBool, nil
	case nil:
		return jsonTypeNull, nil
	}
	return jsonTypeValue, nil
}

// decodeKeysOnlyContainer handles object/array containers for decodeKeysOnly.
func decodeKeysOnlyContainer(decoder *json.Decoder, delim json.Delim) (any, error) {
	if delim == '{' {
		return decodeKeysOnlyObject(decoder)
	}
	return decodeKeysOnlyArray(decoder)
}

// decodeKeysOnlyObject decodes an object replacing values with type names.
func decodeKeysOnlyObject(decoder *json.Decoder) (map[string]any, error) {
	obj := make(map[string]any)
	for decoder.More() {
		key, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		k, ok := key.(string)
		if !ok {
			return nil, ErrNotJSON
		}
		val, err := decodeKeysOnly(decoder)
		if err != nil {
			return nil, err
		}
		obj[k] = val
	}
	if _, err := decoder.Token(); err != nil {
		return nil, err
	}
	return obj, nil
}

// decodeKeysOnlyArray decodes an array discarding values.
func decodeKeysOnlyArray(decoder *json.Decoder) (string, error) {
	for decoder.More() {
		if _, err := decodeKeysOnly(decoder); err != nil {
			return "", err
		}
	}
	if _, err := decoder.Token(); err != nil {
		return "", err
	}
	return compressedTypeArray, nil
}

// compressSelectKeys keeps only selected keys for every object in an array or object.
func compressSelectKeys(data []byte, selectKeys []string) ([]byte, error) {
	if len(selectKeys) == 0 {
		return nil, errors.New("select_keys requires at least one selectKey")
	}

	set := make(map[string]struct{}, len(selectKeys))
	for _, k := range selectKeys {
		set[k] = struct{}{}
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	value, err := decodeSelectKeys(decoder, set)
	if err != nil {
		return nil, err
	}
	return json.Marshal(value)
}

// decodeSelectKeys recursively keeps only selected keys in objects.
func decodeSelectKeys(decoder *json.Decoder, keep map[string]struct{}) (any, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	switch t := token.(type) {
	case json.Delim:
		return decodeSelectKeysContainer(decoder, t, keep)
	default:
		return t, nil
	}
}

// decodeSelectKeysContainer handles object/array containers for decodeSelectKeys.
func decodeSelectKeysContainer(decoder *json.Decoder, delim json.Delim, keep map[string]struct{}) (any, error) {
	if delim == '{' {
		return decodeSelectKeysObject(decoder, keep)
	}
	return decodeSelectKeysArray(decoder, keep)
}

// decodeSelectKeysObject decodes an object keeping only selected keys.
func decodeSelectKeysObject(decoder *json.Decoder, keep map[string]struct{}) (map[string]any, error) {
	obj := make(map[string]any)
	for decoder.More() {
		key, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		k, ok := key.(string)
		if !ok {
			return nil, ErrNotJSON
		}
		if _, keepKey := keep[k]; keepKey {
			val, err := decodeSelectKeys(decoder, keep)
			if err != nil {
				return nil, err
			}
			obj[k] = val
			continue
		}
		if err := skipValue(decoder); err != nil {
			return nil, err
		}
	}
	if _, err := decoder.Token(); err != nil {
		return nil, err
	}
	return obj, nil
}

// decodeSelectKeysArray decodes an array applying select-key filtering to each item.
func decodeSelectKeysArray(decoder *json.Decoder, keep map[string]struct{}) ([]any, error) {
	arr := []any{}
	for decoder.More() {
		val, err := decodeSelectKeys(decoder, keep)
		if err != nil {
			return nil, err
		}
		arr = append(arr, val)
	}
	if _, err := decoder.Token(); err != nil {
		return nil, err
	}
	return arr, nil
}
