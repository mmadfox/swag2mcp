package reader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/reader"
	"github.com/stretchr/testify/require"
)

func TestReader_Outline(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{
		"status": "ok",
		"data": [
			{"id": 1, "title": "first"},
			{"id": 2, "title": "second"},
			{"id": 3, "title": "third"}
		]
	}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{MaxDepth: 3, MaxArrayItems: 2})
	require.NoError(t, err)
	require.Equal(t, "object", outline.Type)
	require.Equal(t, int64(len(data)), outline.Size)
	require.Equal(t, []string{"status", "data"}, outline.Keys)
	require.Equal(t, 3, outline.Structure.Structure["data"].ItemCount)
	require.Equal(t, "array", outline.Structure.Structure["data"].Type)

	require.NotEmpty(t, outline.CompressionHints)
}
