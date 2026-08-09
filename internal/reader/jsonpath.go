/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package reader

import (
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// rootPath is the gjson path that addresses the root value of a document.
const rootPath = "@this"

// NormalizeJSONPath resolves "-" (last element) segments to concrete array
// indices, because gjson does not support "-" for the last element of an
// array. It resolves "-" in any position, not just the terminal segment.
// It returns the path unchanged when there is no "-" or when it cannot be
// resolved against the data.
//
// Examples (data = [{"id":1},{"id":2},{"id":3}]):
//
//	"-"            -> "2"          (last element of the root array)
//	"data.-"       -> "data.2"     (last element of the data array)
//	"data.-.name"  -> "data.2.name" (field of the last element)
func NormalizeJSONPath(data []byte, path string) string {
	if path == "" || !strings.Contains(path, "-") {
		return path
	}

	segments := strings.Split(path, ".")
	for i, seg := range segments {
		if seg != "-" {
			continue
		}
		parent := strings.Join(segments[:i], ".")
		idx, ok := lastArrayIndex(data, parent)
		if !ok {
			return path
		}
		segments[i] = strconv.Itoa(idx)
	}
	return strings.Join(segments, ".")
}

// lastArrayIndex returns the index of the last element of the array at the
// given dotted path. An empty parent refers to the root value. It returns
// ok=false when the value at parent is not a non-empty array.
func lastArrayIndex(data []byte, parent string) (int, bool) {
	query := parent
	if query == "" {
		query = rootPath
	}
	res := gjson.GetBytes(data, query)
	if !res.IsArray() {
		return 0, false
	}
	arr := res.Array()
	if len(arr) == 0 {
		return 0, false
	}
	return len(arr) - 1, true
}
