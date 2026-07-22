package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

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

func specFileName(domain, title, location string) string {
	base := title
	if base == "" {
		base = specFileNameBase(location)
	}

	ext := filepath.Ext(base)
	if ext == "" {
		ext = ".yaml"
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

	if sanitized == domain {
		return fmt.Sprintf("%s%s", domain, ext)
	}

	return fmt.Sprintf("%s-%s%s", domain, sanitized, ext)
}

func removeDiacritics(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.Predicate(unicode.IsMark)), norm.NFC)
	result, _, _ := transform.String(t, s)
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
