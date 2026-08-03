/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const defaultSpecName = "spec"

const (
	maxTitleLen    = 20
	maxPathPartLen = 40
)

type specFilter struct {
	domains map[string]struct{}
}

func makeFilter(domains []string) *specFilter {
	f := &specFilter{domains: make(map[string]struct{}, len(domains))}
	for _, d := range domains {
		f.domains[strings.TrimSpace(d)] = struct{}{}
	}
	return f
}

func (f *specFilter) match(domain string) bool {
	if len(f.domains) == 0 {
		return true
	}
	_, ok := f.domains[domain]
	return ok
}

func isURL(s string) bool {
	return strings.HasPrefix(s, "https://") || strings.HasPrefix(s, "http://")
}

func extFromLocation(location string) string {
	if !isURL(location) {
		return ""
	}
	u, err := url.Parse(location)
	if err != nil {
		return ""
	}
	return filepath.Ext(u.Path)
}

func truncateByLastHyphen(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	truncated := s[:maxLen]
	if lastHyphen := strings.LastIndex(truncated, "-"); lastHyphen > 0 {
		return truncated[:lastHyphen]
	}
	return truncated
}

func specFileName(domain, title, location, pathPart string) string {
	base := title
	if base == "" {
		base = specFileNameBase(location)
	}

	ext := filepath.Ext(base)
	if ext == "" {
		ext = extFromLocation(location)
	}
	if ext == "" {
		ext = filepath.Ext(location)
	}
	base = strings.TrimSuffix(base, ext)
	base = strings.TrimSuffix(base, ".yml")

	sanitized := strings.ToLower(base)
	sanitized = strings.NewReplacer(
		" ", "-",
		"_", "-",
		".", "-",
	).Replace(sanitized)
	sanitized = removeDiacritics(sanitized)
	sanitized = truncateByLastHyphen(sanitized, maxTitleLen)

	if pathPart != "" {
		pathPart = truncateByLastHyphen(pathPart, maxPathPartLen)
		sanitized = sanitized + "-" + pathPart
	}

	if sanitized == domain {
		return fmt.Sprintf("%s%s", domain, ext)
	}

	return fmt.Sprintf("%s-%s%s", domain, sanitized, ext)
}

// pathPartFromLocation extracts the directory path from a URL and converts it
// to a dash-separated string. For "https://example.com/api/v1/spec.yaml"
// it returns "api-v1". Returns empty string for non-URL locations.
func pathPartFromLocation(location string) string {
	if !strings.HasPrefix(location, "http://") && !strings.HasPrefix(location, "https://") {
		return ""
	}
	u, err := url.Parse(location)
	if err != nil {
		return ""
	}
	dir := strings.TrimSuffix(u.Path, filepath.Base(u.Path))
	dir = strings.Trim(dir, "/")
	if dir == "" {
		return ""
	}
	return strings.ReplaceAll(dir, "/", "-")
}

func removeDiacritics(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.Predicate(unicode.IsMark)), norm.NFC)
	result, _, err := transform.String(t, s)
	if err != nil {
		return s
	}
	return result
}

func specFileNameBase(location string) string {
	if strings.HasPrefix(location, "http://") || strings.HasPrefix(location, "https://") {
		u, err := url.Parse(location)
		if err == nil && u.Path != "" && u.Path != "/" {
			return filepath.Base(u.Path)
		}
		return defaultSpecName
	}
	return filepath.Base(location)
}
