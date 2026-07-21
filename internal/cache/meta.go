package cache

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type fileMeta struct {
	Source     string    `json:"source"`
	SourceType string    `json:"source_type"` // "url" or "local"
	CachedAt   time.Time `json:"cached_at"`
	ModTime    time.Time `json:"mod_time"`
	TTLSec     int       `json:"ttl_sec"`
}

// IsExpired reports whether the cached file has exceeded its TTL.
func (m fileMeta) IsExpired() bool {
	return time.Since(m.CachedAt) > time.Duration(m.TTLSec)*time.Second
}

// readMeta reads and parses a fileMeta from the given path.
func readMeta(path string) (fileMeta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return fileMeta{}, err
	}
	var m fileMeta
	if err := json.Unmarshal(data, &m); err != nil {
		return fileMeta{}, err
	}
	return m, nil
}

// writeMeta serializes a fileMeta and writes it to the given path.
func writeMeta(path string, m fileMeta) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal meta: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}
