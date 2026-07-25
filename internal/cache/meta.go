package cache

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

func (m fileMeta) IsExpired() bool {
	return time.Since(m.CachedAt) > time.Duration(m.TTLSec)*time.Second
}

func readMeta(path string) (fileMeta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return fileMeta{}, err
	}
	var m fileMeta
	if uErr := json.Unmarshal(data, &m); uErr != nil {
		return fileMeta{}, uErr
	}
	return m, nil
}

func writeMeta(path string, m fileMeta) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal meta: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}
