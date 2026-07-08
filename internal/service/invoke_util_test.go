package service

import (
	"net/http"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/types"
)

func TestResolveMaxResponseSize_NilConfig(t *testing.T) {
	t.Parallel()

	size := resolveMaxResponseSize(nil)
	if size != defaultMaxResponseSize {
		t.Errorf("got %d, want %d", size, defaultMaxResponseSize)
	}
}

func TestResolveMaxResponseSize_NilField(t *testing.T) {
	t.Parallel()

	size := resolveMaxResponseSize(&types.HTTPClientConfig{})
	if size != defaultMaxResponseSize {
		t.Errorf("got %d, want %d", size, defaultMaxResponseSize)
	}
}

func TestResolveMaxResponseSize_Custom(t *testing.T) {
	t.Parallel()

	val := 4096
	size := resolveMaxResponseSize(&types.HTTPClientConfig{MaxResponseSize: &val})
	if size != 4096 {
		t.Errorf("got %d, want %d", size, 4096)
	}
}

func TestResolveMaxResponseSize_ExceedsMax(t *testing.T) {
	t.Parallel()

	val := 2 * 1024 * 1024 // 2 MB
	size := resolveMaxResponseSize(&types.HTTPClientConfig{MaxResponseSize: &val})
	if size != maxMaxResponseSize {
		t.Errorf("got %d, want %d", size, maxMaxResponseSize)
	}
}

func TestResolveMaxResponseSize_Zero(t *testing.T) {
	t.Parallel()

	val := 0
	size := resolveMaxResponseSize(&types.HTTPClientConfig{MaxResponseSize: &val})
	if size != defaultMaxResponseSize {
		t.Errorf("got %d, want %d", size, defaultMaxResponseSize)
	}
}

func TestResolveMaxResponseSize_Negative(t *testing.T) {
	t.Parallel()

	val := -100
	size := resolveMaxResponseSize(&types.HTTPClientConfig{MaxResponseSize: &val})
	if size != defaultMaxResponseSize {
		t.Errorf("got %d, want %d", size, defaultMaxResponseSize)
	}
}

func TestOpenCommand_Darwin(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("only on darwin")
	}

	cmd := openCommand("/tmp/test.json")
	if cmd != "open /tmp/test.json" {
		t.Errorf("got %q, want %q", cmd, "open /tmp/test.json")
	}
}

func TestOpenCommand_Linux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("only on linux")
	}

	cmd := openCommand("/tmp/test.json")
	if cmd != "xdg-open /tmp/test.json" {
		t.Errorf("got %q, want %q", cmd, "xdg-open /tmp/test.json")
	}
}

func TestOpenCommand_Windows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("only on windows")
	}

	cmd := openCommand("C:\\test.json")
	if cmd != "start C:\\test.json" {
		t.Errorf("got %q, want %q", cmd, "start C:\\test.json")
	}
}

func TestFormatSize_Bytes(t *testing.T) {
	t.Parallel()

	if s := formatSize(500); s != "500 B" {
		t.Errorf("got %q, want %q", s, "500 B")
	}
}

func TestFormatSize_KB(t *testing.T) {
	t.Parallel()

	if s := formatSize(2048); s != "2.0 KB" {
		t.Errorf("got %q, want %q", s, "2.0 KB")
	}
}

func TestFormatSize_MB(t *testing.T) {
	t.Parallel()

	if s := formatSize(1048576); s != "1.0 MB" {
		t.Errorf("got %q, want %q", s, "1.0 MB")
	}
}

func TestFormatSize_GB(t *testing.T) {
	t.Parallel()

	if s := formatSize(1073741824); s != "1.0 GB" {
		t.Errorf("got %q, want %q", s, "1.0 GB")
	}
}

func TestFormatSize_Zero(t *testing.T) {
	t.Parallel()

	if s := formatSize(0); s != "0 B" {
		t.Errorf("got %q, want %q", s, "0 B")
	}
}

func TestRandomSuffix_Length(t *testing.T) {
	t.Parallel()

	suffix := randomSuffix(6)
	if len(suffix) != 6 {
		t.Errorf("len = %d, want %d", len(suffix), 6)
	}
}

func TestRandomSuffix_HexChars(t *testing.T) {
	t.Parallel()

	suffix := randomSuffix(12)
	for _, c := range suffix {
		if !strings.ContainsRune("0123456789abcdef", c) {
			t.Errorf("unexpected char %c in suffix %q", c, suffix)
		}
	}
}

func TestRandomSuffix_Unique(t *testing.T) {
	t.Parallel()

	s1 := randomSuffix(6)
	s2 := randomSuffix(6)
	if s1 == s2 {
		t.Error("two random suffixes are identical")
	}
}

func TestSaveLargeResponse(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	body := make([]byte, 10000)
	response := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}

	endpoint := &types.Endpoint{
		Name: "GET",
		Path: "/test",
	}

	resp, err := svc.saveLargeResponse(response, body, t.Name(), endpoint, 2048)
	if err != nil {
		t.Fatalf("saveLargeResponse() = %v", err)
	}

	if resp.FileRef == nil {
		t.Fatal("FileRef is nil")
	}
	if resp.FileRef.Size != 10000 {
		t.Errorf("Size = %d, want %d", resp.FileRef.Size, 10000)
	}
	if resp.FileRef.SizeHint == "" {
		t.Error("SizeHint is empty")
	}
	if resp.FileRef.MaxSizeHint == "" {
		t.Error("MaxSizeHint is empty")
	}
	if resp.FileRef.Message == "" {
		t.Error("Message is empty")
	}
	if resp.FileRef.OpenCmd == "" {
		t.Error("OpenCmd is empty")
	}
	if !strings.HasPrefix(resp.FileRef.Path, svc.ws.ResponsesDir()) {
		t.Errorf("Path %q not in responses dir %q", resp.FileRef.Path, svc.ws.ResponsesDir())
	}

	if _, statErr := os.Stat(resp.FileRef.Path); os.IsNotExist(statErr) {
		t.Error("response file was not created on disk")
	}
}

func TestSaveLargeResponse_FileContent(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	body := []byte(`{"key": "value"}`)
	response := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{},
	}

	endpoint := &types.Endpoint{
		Name: "GET",
		Path: "/test",
	}

	resp, err := svc.saveLargeResponse(response, body, t.Name(), endpoint, 100)
	if err != nil {
		t.Fatalf("saveLargeResponse() = %v", err)
	}

	data, err := os.ReadFile(resp.FileRef.Path)
	if err != nil {
		t.Fatalf("ReadFile() = %v", err)
	}
	if string(data) != string(body) {
		t.Errorf("file content = %q, want %q", string(data), string(body))
	}
}

func TestSaveLargeResponse_StatusCode(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	body := []byte("test")
	response := &http.Response{
		StatusCode: http.StatusNotFound,
		Header:     http.Header{},
	}

	endpoint := &types.Endpoint{
		Name: "GET",
		Path: "/test",
	}

	resp, err := svc.saveLargeResponse(response, body, t.Name(), endpoint, 100)
	if err != nil {
		t.Fatalf("saveLargeResponse() = %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}
