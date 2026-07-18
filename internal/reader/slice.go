package reader

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// Slice extracts a JSON fragment by jsonPath or by line range.
func (r *reader) Slice(path string, opts SliceOptions) (Slice, error) {
	if err := r.validatePath(path); err != nil {
		return Slice{}, err
	}

	opts = normalizeSliceOptions(opts)

	if opts.JSONPath != "" {
		return r.sliceByJSONPath(path, opts)
	}

	return r.sliceByLines(path, opts)
}

// sliceByJSONPath extracts the value at a gjson path and its surrounding context.
func (r *reader) sliceByJSONPath(path string, opts SliceOptions) (Slice, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Slice{}, fmt.Errorf("read file: %w", err)
	}

	result := gjson.GetBytes(data, opts.JSONPath)
	if !result.Exists() {
		return Slice{}, fmt.Errorf("%w: %s", ErrPathNotFound, opts.JSONPath)
	}

	raw := []byte(result.Raw)
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		value = result.Value()
	}

	context := inferContext(result)
	lines, err := locateLines(data, raw)
	if err != nil {
		lines = [2]int{0, 0}
	}

	slice := Slice{
		Lines:      lines,
		Value:      value,
		JSONPath:   opts.JSONPath,
		Context:    context,
		IsComplete: true,
		NextLine:   lines[1] + 1,
		PrevLine:   lines[0] - 1,
		NextPath:   adjacentPath(opts.JSONPath, 1),
		PrevPath:   adjacentPath(opts.JSONPath, -1),
	}

	if opts.Limit == 0 || len(raw) <= opts.Limit/2 {
		slice.Fragment = string(raw)
	}

	if opts.Limit > 0 && len(raw) > opts.Limit {
		return slice, fmt.Errorf("%w: fragment exceeds limit", ErrPathNotFound)
	}

	return slice, nil
}

// inferContext determines whether the gjson result is an object, array, or scalar.
func inferContext(result gjson.Result) string {
	raw := strings.TrimSpace(result.Raw)
	if len(raw) == 0 {
		return jsonTypeValue
	}
	switch raw[0] {
	case '{':
		return jsonTypeObject
	case '[':
		return jsonTypeArray
	default:
		return jsonTypeValue
	}
}

// adjacentPath returns a sibling path by index adjustment, best effort.
func adjacentPath(path string, delta int) string {
	if path == "" {
		return ""
	}
	idx := strings.LastIndexAny(path, "[.")
	if idx == -1 {
		return ""
	}
	if path[idx] == '[' {
		end := strings.Index(path[idx:], "]")
		if end == -1 {
			return ""
		}
		end += idx
		n, err := strconv.Atoi(path[idx+1 : end])
		if err != nil {
			return ""
		}
		return path[:idx+1] + strconv.Itoa(n+delta) + path[end:]
	}
	n, err := strconv.Atoi(path[idx+1:])
	if err != nil {
		return ""
	}
	return path[:idx+1] + strconv.Itoa(n+delta)
}

// locateLines finds the 1-based line range of a raw JSON substring in the full data.
func locateLines(data, raw []byte) ([2]int, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return [2]int{}, errors.New("empty fragment")
	}
	before, _, found := bytes.Cut(data, trimmed)
	if !found {
		return [2]int{}, errors.New("fragment not found in source")
	}

	startLine := bytes.Count(before, []byte{'\n'}) + 1
	endLine := startLine + bytes.Count(trimmed, []byte{'\n'})
	return [2]int{startLine, endLine}, nil
}

// sliceByLines returns a raw fragment around the requested lines.
func (r *reader) sliceByLines(path string, opts SliceOptions) (Slice, error) {
	start, end, err := parseLineRange(opts)
	if err != nil {
		return Slice{}, err
	}

	f, err := os.Open(path)
	if err != nil {
		return Slice{}, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		if lineNo >= start && lineNo <= end {
			lines = append(lines, scanner.Text())
		}
		if lineNo >= end {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return Slice{}, fmt.Errorf("scan file: %w", err)
	}

	if len(lines) == 0 {
		return Slice{}, fmt.Errorf("%w: no lines in range %d-%d", ErrInvalidLineRange, start, end)
	}

	fragment := strings.Join(lines, "\n")
	var value any
	if err := json.Unmarshal([]byte(fragment), &value); err != nil {
		value = nil
	}

	context := inferContextFromRaw(fragment)
	slice := Slice{
		Lines:      [2]int{start, start + len(lines) - 1},
		Fragment:   fragment,
		Value:      value,
		Context:    context,
		IsComplete: value != nil,
		NextLine:   start + len(lines),
		PrevLine:   start - 1,
	}

	if opts.Limit > 0 && len(fragment) > opts.Limit {
		return slice, fmt.Errorf("%w: fragment exceeds limit", ErrInvalidLineRange)
	}

	return slice, nil
}

// parseLineRange converts options into an absolute 1-based line range.
func parseLineRange(opts SliceOptions) (int, int, error) {
	if opts.Line > 0 {
		start := max(opts.Line-opts.Around, 1)
		return start, opts.Line + opts.Around, nil
	}

	if opts.Range != "" {
		const rangeParts = 2
		parts := strings.Split(opts.Range, "-")
		if len(parts) != rangeParts {
			return 0, 0, fmt.Errorf("%w: range must be start-end", ErrInvalidLineRange)
		}
		start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return 0, 0, fmt.Errorf("%w: invalid start line", ErrInvalidLineRange)
		}
		end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return 0, 0, fmt.Errorf("%w: invalid end line", ErrInvalidLineRange)
		}
		if start < 1 || end < start {
			return 0, 0, fmt.Errorf("%w: invalid range", ErrInvalidLineRange)
		}
		return start, end, nil
	}

	return 0, 0, fmt.Errorf("%w: either line or range must be specified", ErrInvalidLineRange)
}

// inferContextFromRaw guesses whether the raw fragment is object/array/value.
func inferContextFromRaw(fragment string) string {
	trimmed := strings.TrimSpace(fragment)
	if len(trimmed) == 0 {
		return jsonTypeValue
	}
	switch trimmed[0] {
	case '{':
		return jsonTypeObject
	case '[':
		return jsonTypeArray
	default:
		return jsonTypeValue
	}
}
