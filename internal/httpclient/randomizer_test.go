package httpclient

import (
	"strings"
	"testing"
)

func TestRandomUserAgent(t *testing.T) {
	ua := randomUserAgent()
	if ua == "" {
		t.Fatal("randomUserAgent() returned empty")
	}
	if !strings.HasPrefix(ua, "Mozilla/5.0") {
		t.Errorf("randomUserAgent() = %q, want Mozilla/5.0 prefix", ua)
	}
}

func TestRandomReferer(t *testing.T) {
	ref := randomReferer()
	if ref == "" {
		t.Fatal("randomReferer() returned empty")
	}
	if !strings.HasPrefix(ref, "https://") {
		t.Errorf("randomReferer() = %q, want https:// prefix", ref)
	}
}

func TestRandomAccept_JSON(t *testing.T) {
	accept := randomAccept(true)
	if accept == "" {
		t.Fatal("randomAccept(true) returned empty")
	}
	if !strings.Contains(accept, "application/json") {
		t.Errorf("randomAccept(true) = %q, want application/json", accept)
	}
}

func TestRandomAccept_HTML(t *testing.T) {
	accept := randomAccept(false)
	if accept == "" {
		t.Fatal("randomAccept(false) returned empty")
	}
	if !strings.Contains(accept, "text/html") {
		t.Errorf("randomAccept(false) = %q, want text/html", accept)
	}
}

func TestRandomAcceptEncoding(t *testing.T) {
	enc := randomAcceptEncoding()
	if enc == "" {
		t.Fatal("randomAcceptEncoding() returned empty")
	}
}

func TestRandomSecChUa(t *testing.T) {
	ua := randomSecChUa()
	if ua == "" {
		t.Fatal("randomSecChUa() returned empty")
	}
}

func TestRandomSecChUaPlatform(t *testing.T) {
	p := randomSecChUaPlatform()
	if p == "" {
		t.Fatal("randomSecChUaPlatform() returned empty")
	}
}

func TestRandomSecFetchSite(t *testing.T) {
	s := randomSecFetchSite()
	if s == "" {
		t.Fatal("randomSecFetchSite() returned empty")
	}
}

func TestRandomSecFetchMode(t *testing.T) {
	m := randomSecFetchMode()
	if m == "" {
		t.Fatal("randomSecFetchMode() returned empty")
	}
}

func TestRandomSecFetchDest(t *testing.T) {
	d := randomSecFetchDest()
	if d == "" {
		t.Fatal("randomSecFetchDest() returned empty")
	}
}

func TestDetectSystemLanguage(t *testing.T) {
	t.Setenv("LANG", "ru_RU.UTF-8")
	t.Setenv("MUI_LANG", "")
	t.Setenv("LC_ALL", "")
	cachedSystemLang = ""

	lang := detectSystemLanguage()
	if lang != "ru_RU.UTF-8" {
		t.Errorf("detectSystemLanguage() = %q, want ru_RU.UTF-8", lang)
	}
}

func TestDetectSystemLanguage_Fallback(t *testing.T) {
	t.Setenv("LANG", "")
	t.Setenv("MUI_LANG", "")
	t.Setenv("LC_ALL", "")
	cachedSystemLang = ""

	lang := detectSystemLanguage()
	if lang != "en_US.UTF-8" {
		t.Errorf("detectSystemLanguage() = %q, want en_US.UTF-8", lang)
	}
}

func TestRandomAcceptLanguage_Russian(t *testing.T) {
	t.Setenv("LANG", "ru_RU.UTF-8")
	cachedSystemLang = ""

	al := randomAcceptLanguage()
	if !strings.Contains(al, "ru-RU") {
		t.Errorf("randomAcceptLanguage() = %q, want ru-RU", al)
	}
}

func TestRandomAcceptLanguage_English(t *testing.T) {
	t.Setenv("LANG", "en_US.UTF-8")
	cachedSystemLang = ""

	al := randomAcceptLanguage()
	if !strings.Contains(al, "en-US") {
		t.Errorf("randomAcceptLanguage() = %q, want en-US", al)
	}
}

func TestRandomSecHeaders(t *testing.T) {
	h := randomSecHeaders()
	if len(h) == 0 {
		t.Fatal("randomSecHeaders() returned empty map")
	}
	if _, ok := h["Sec-Ch-Ua"]; !ok {
		t.Error("missing Sec-Ch-Ua")
	}
	if _, ok := h["Sec-Ch-Ua-Platform"]; !ok {
		t.Error("missing Sec-Ch-Ua-Platform")
	}
	if _, ok := h["Sec-Fetch-Site"]; !ok {
		t.Error("missing Sec-Fetch-Site")
	}
	if _, ok := h["Sec-Fetch-Mode"]; !ok {
		t.Error("missing Sec-Fetch-Mode")
	}
	if _, ok := h["Sec-Fetch-Dest"]; !ok {
		t.Error("missing Sec-Fetch-Dest")
	}
}

func TestDetectSystemLanguage_Cache(t *testing.T) {
	t.Setenv("LANG", "de_DE.UTF-8")
	cachedSystemLang = ""

	_ = detectSystemLanguage()

	t.Setenv("LANG", "fr_FR.UTF-8")

	lang := detectSystemLanguage()
	if lang != "de_DE.UTF-8" {
		t.Errorf("detectSystemLanguage() = %q, want de_DE.UTF-8 (cached)", lang)
	}
}
