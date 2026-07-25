package httpclient

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomUserAgent(t *testing.T) {
	ua := randomUserAgent()
	require.NotEmpty(t, ua)
	assert.True(t, strings.HasPrefix(ua, "Mozilla/5.0"))
}

func TestRandomReferer(t *testing.T) {
	ref := randomReferer()
	require.NotEmpty(t, ref)
	assert.True(t, strings.HasPrefix(ref, "https://"))
}

func TestRandomAccept_JSON(t *testing.T) {
	accept := randomAccept(true)
	require.NotEmpty(t, accept)
	assert.Contains(t, accept, "application/json")
}

func TestRandomAccept_HTML(t *testing.T) {
	accept := randomAccept(false)
	require.NotEmpty(t, accept)
	assert.Contains(t, accept, "text/html")
}

func TestRandomAcceptEncoding(t *testing.T) {
	enc := randomAcceptEncoding()
	require.NotEmpty(t, enc)
}

func TestRandomSecChUa(t *testing.T) {
	ua := randomSecChUa()
	require.NotEmpty(t, ua)
}

func TestRandomSecChUaPlatform(t *testing.T) {
	p := randomSecChUaPlatform()
	require.NotEmpty(t, p)
}

func TestRandomSecFetchSite(t *testing.T) {
	s := randomSecFetchSite()
	require.NotEmpty(t, s)
}

func TestRandomSecFetchMode(t *testing.T) {
	m := randomSecFetchMode()
	require.NotEmpty(t, m)
}

func TestRandomSecFetchDest(t *testing.T) {
	d := randomSecFetchDest()
	require.NotEmpty(t, d)
}

func TestDetectSystemLanguage(t *testing.T) {
	t.Setenv("LANG", "ru_RU.UTF-8")
	t.Setenv("MUI_LANG", "")
	t.Setenv("LC_ALL", "")
	cachedSystemLang = ""

	lang := detectSystemLanguage()
	assert.Equal(t, "ru_RU.UTF-8", lang)
}

func TestDetectSystemLanguage_Fallback(t *testing.T) {
	t.Setenv("LANG", "")
	t.Setenv("MUI_LANG", "")
	t.Setenv("LC_ALL", "")
	cachedSystemLang = ""

	lang := detectSystemLanguage()
	assert.Equal(t, "en_US.UTF-8", lang)
}

func TestRandomAcceptLanguage_Russian(t *testing.T) {
	t.Setenv("LANG", "ru_RU.UTF-8")
	cachedSystemLang = ""

	al := randomAcceptLanguage()
	assert.Contains(t, al, "ru-RU")
}

func TestRandomAcceptLanguage_English(t *testing.T) {
	t.Setenv("LANG", "en_US.UTF-8")
	cachedSystemLang = ""

	al := randomAcceptLanguage()
	assert.Contains(t, al, "en-US")
}

func TestRandomSecHeaders(t *testing.T) {
	h := randomSecHeaders()
	require.NotEmpty(t, h)
	assert.Contains(t, h, "Sec-Ch-Ua")
	assert.Contains(t, h, "Sec-Ch-Ua-Platform")
	assert.Contains(t, h, "Sec-Fetch-Site")
	assert.Contains(t, h, "Sec-Fetch-Mode")
	assert.Contains(t, h, "Sec-Fetch-Dest")
}

func TestDetectSystemLanguage_Cache(t *testing.T) {
	t.Setenv("LANG", "de_DE.UTF-8")
	cachedSystemLang = ""

	_ = detectSystemLanguage()

	t.Setenv("LANG", "fr_FR.UTF-8")

	lang := detectSystemLanguage()
	assert.Equal(t, "de_DE.UTF-8", lang)
}
