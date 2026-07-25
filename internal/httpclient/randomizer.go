package httpclient

import (
	"math/rand/v2"
	"os"
	"strings"
)

var (
	userAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:127.0) Gecko/20100101 Firefox/127.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:127.0) Gecko/20100101 Firefox/127.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:128.0) Gecko/20100101 Firefox/128.0",
		"Mozilla/5.0 (X11; Linux i686; rv:127.0) Gecko/20100101 Firefox/127.0",
		"Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.6 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0",
	}

	referers = []string{
		"https://www.google.com/",
		"https://www.google.com/search?q=",
		"https://yandex.ru/",
		"https://yandex.ru/search/?text=",
		"https://duckduckgo.com/",
		"https://duckduckgo.com/?q=",
		"https://www.bing.com/",
		"https://www.bing.com/search?q=",
		"https://www.facebook.com/",
		"https://l.facebook.com/l.php?",
		"https://twitter.com/",
		"https://t.co/",
		"https://www.instagram.com/",
		"https://www.reddit.com/",
		"https://www.linkedin.com/",
		"https://news.ycombinator.com/",
	}

	secChUa = []string{
		`"Chromium";v="125", "Google Chrome";v="125"`,
		`"Chromium";v="126", "Google Chrome";v="126"`,
		`"Chromium";v="127", "Google Chrome";v="127"`,
		`"Chromium";v="125", "Microsoft Edge";v="125"`,
		`"Chromium";v="126", "Microsoft Edge";v="126"`,
		`"Not/A)Brand";v="99", "Google Chrome";v="127"`,
	}

	secChUaPlatform = []string{
		"Windows",
		"macOS",
		"Linux",
	}

	acceptJSON = []string{
		"application/json, text/plain, */*",
		"application/json, application/problem+json, */*",
		"application/json, text/plain;q=0.9, */*;q=0.8",
	}

	acceptHTML = []string{
		"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"text/html,application/xhtml+xml;q=0.9,application/xml;q=0.8,*/*;q=0.7",
	}

	acceptEncodings = []string{
		"gzip, deflate, br",
		"gzip, deflate",
		"gzip, br",
	}

	secFetchSites = []string{
		"same-origin",
		"same-site",
		"cross-site",
		"none",
	}

	secFetchModes = []string{
		"cors",
		"navigate",
		"no-cors",
		"same-origin",
	}

	secFetchDests = []string{
		"document",
		"empty",
		"script",
		"image",
	}

	cachedSystemLang string
)

const (
	langCodeLen   = 2
	langRegionLen = 5
)

// randomUserAgent returns a random browser User-Agent string.
func randomUserAgent() string {
	return userAgents[rand.IntN(len(userAgents))]
}

// randomReferer returns a random Referer URL.
func randomReferer() string {
	return referers[rand.IntN(len(referers))]
}

// randomAccept returns a random Accept header value based on content type.
func randomAccept(isJSON bool) string {
	if isJSON {
		return acceptJSON[rand.IntN(len(acceptJSON))]
	}
	return acceptHTML[rand.IntN(len(acceptHTML))]
}

// randomAcceptEncoding returns a random Accept-Encoding header value.
func randomAcceptEncoding() string {
	return acceptEncodings[rand.IntN(len(acceptEncodings))]
}

// randomSecChUa returns a random Sec-Ch-Ua header value.
func randomSecChUa() string {
	return secChUa[rand.IntN(len(secChUa))]
}

// randomSecChUaPlatform returns a random Sec-Ch-Ua-Platform header value.
func randomSecChUaPlatform() string {
	return secChUaPlatform[rand.IntN(len(secChUaPlatform))]
}

// randomSecFetchSite returns a random Sec-Fetch-Site header value.
func randomSecFetchSite() string {
	return secFetchSites[rand.IntN(len(secFetchSites))]
}

// randomSecFetchMode returns a random Sec-Fetch-Mode header value.
func randomSecFetchMode() string {
	return secFetchModes[rand.IntN(len(secFetchModes))]
}

// randomSecFetchDest returns a random Sec-Fetch-Dest header value.
func randomSecFetchDest() string {
	return secFetchDests[rand.IntN(len(secFetchDests))]
}

// detectSystemLanguage detects the system language from environment variables.
// Falls back to en_US.UTF-8 if none are set. Caches the result.
func detectSystemLanguage() string {
	if cachedSystemLang != "" {
		return cachedSystemLang
	}

	lang := os.Getenv("LANG")
	if lang == "" {
		lang = os.Getenv("MUI_LANG")
	}
	if lang == "" {
		lang = os.Getenv("LC_ALL")
	}
	if lang == "" {
		lang = "en_US.UTF-8"
	}

	cachedSystemLang = lang
	return lang
}

// randomAcceptLanguage returns a random Accept-Language header value based on system language.
func randomAcceptLanguage() string {
	lang := detectSystemLanguage()

	code := "en"
	if len(lang) >= langCodeLen {
		code = strings.ToLower(lang[:2])
	}

	region := strings.ToUpper(code)
	if len(lang) >= langRegionLen {
		region = strings.ToUpper(lang[3:5])
	}

	switch code {
	case "en":
		return "en-US,en;q=0.9"
	case "ru":
		return "ru-RU,ru;q=0.9,en;q=0.8"
	case "de":
		return "de-DE,de;q=0.9,en;q=0.8"
	case "fr":
		return "fr-FR,fr;q=0.9,en;q=0.8"
	case "es":
		return "es-ES,es;q=0.9,en;q=0.8"
	case "pt":
		return "pt-BR,pt;q=0.9,en;q=0.8"
	case "it":
		return "it-IT,it;q=0.9,en;q=0.8"
	case "nl":
		return "nl-NL,nl;q=0.9,en;q=0.8"
	case "pl":
		return "pl-PL,pl;q=0.9,en;q=0.8"
	case "ja":
		return "ja-JP,ja;q=0.9,en;q=0.8"
	case "zh":
		return "zh-CN,zh;q=0.9,en;q=0.8"
	case "ko":
		return "ko-KR,ko;q=0.9,en;q=0.8"
	case "tr":
		return "tr-TR,tr;q=0.9,en;q=0.8"
	case "ar":
		return "ar-SA,ar;q=0.9,en;q=0.8"
	default:
		return code + "-" + region + "," + code + ";q=0.9,en;q=0.8"
	}
}

// randomSecHeaders returns a map of random Sec-* headers.
func randomSecHeaders() map[string]string {
	return map[string]string{
		"Sec-Ch-Ua":          randomSecChUa(),
		"Sec-Ch-Ua-Platform": randomSecChUaPlatform(),
		"Sec-Fetch-Site":     randomSecFetchSite(),
		"Sec-Fetch-Mode":     randomSecFetchMode(),
		"Sec-Fetch-Dest":     randomSecFetchDest(),
	}
}

// RandomizeConfig fills empty fields in cfg with random browser-like values.
// Existing values are never overwritten.
func RandomizeConfig(cfg *Config) {
	if cfg == nil {
		return
	}

	if cfg.UserAgent == "" {
		cfg.UserAgent = randomUserAgent()
	}
	if cfg.Headers == nil {
		cfg.Headers = make(map[string]string)
	}
	if _, ok := cfg.Headers["Accept"]; !ok {
		cfg.Headers["Accept"] = randomAccept(true)
	}
	if _, ok := cfg.Headers["Accept-Language"]; !ok {
		cfg.Headers["Accept-Language"] = randomAcceptLanguage()
	}
	if _, ok := cfg.Headers["Accept-Encoding"]; !ok {
		cfg.Headers["Accept-Encoding"] = randomAcceptEncoding()
	}
	if _, ok := cfg.Headers["Referer"]; !ok {
		cfg.Headers["Referer"] = randomReferer()
	}

	secHeaders := randomSecHeaders()
	for k, v := range secHeaders {
		if _, ok := cfg.Headers[k]; !ok {
			cfg.Headers[k] = v
		}
	}
}
