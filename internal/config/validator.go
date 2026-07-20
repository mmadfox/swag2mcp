package config

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/mmadfox/swag2mcp/internal/cache"
	"github.com/mmadfox/swag2mcp/internal/spec"
)

var (
	configValidator   *validator.Validate //nolint:gochecknoglobals // Lazily initialized singleton.
	configValidatorMu sync.Mutex          //nolint:gochecknoglobals // Guards lazy validator initialization.
	domainRegex       = regexp.MustCompile(`^[a-z0-9_-]{1,60}$`)
	titleRegex        = regexp.MustCompile(`^[\p{L}\p{N} #*_` + "`" + `~>\[\]()|.,!?;:'"\\–—\-]+$`)
	instructionRegex  = regexp.MustCompile(`^[\p{L}\p{N}\s#*_` + "`" + `~>\[\]()|.,!?;:'"\\–—\-]+$`)
)

// validationError describes a single validation issue.
type validationError struct {
	field      string
	message    string
	spec       string
	collection string
	location   string
	errType    string
}

// validationErrors collects multiple validation errors.
type validationErrors []validationError

// Error returns a formatted string listing all validation errors.
func (ve validationErrors) Error() string {
	if len(ve) == 0 {
		return "no validation errors"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Configuration validation failed with %d error(s):\n", len(ve))
	for i, e := range ve {
		prefix := e.errType
		if prefix == "" {
			prefix = "config"
		}
		msg := e.message
		if e.spec != "" {
			msg = fmt.Sprintf("Spec %q", e.spec)
			if e.collection != "" {
				msg += fmt.Sprintf(", Collection %q", e.collection)
			}
			msg += ": " + e.message
		}
		fmt.Fprintf(&b, "  %d. [%s] %s\n", i+1, prefix, msg)
	}
	return b.String()
}

// ValidateOptions holds optional dependencies for comprehensive validation.
type ValidateOptions struct {
	Cache *cache.Cache
	Tags  []string
}

// ValidateConfig performs comprehensive validation of the configuration.
// Returns nil if valid, or an error listing all issues.
func ValidateConfig(cfg *Config, opts ValidateOptions) error {
	var errs validationErrors

	cfg.HTTPClient.SetDefaults()

	filter := NewFilter(opts.Tags)

	if err := cfg.Validate(filter); err != nil {
		errs = append(errs, validationError{
			errType: "config",
			message: err.Error(),
		})
	}

	errs = append(errs, validateDuplicateDomains(cfg)...)
	if cfg.MockEnabled {
		errs = append(errs, validateMockPorts(cfg, filter)...)
	}
	errs = append(errs, validateSpecLocations(cfg, filter, opts.Cache)...)
	errs = append(errs, validateGlobalHTTPClient(cfg)...)

	if len(errs) == 0 {
		return nil
	}
	return errs
}

// validateDuplicateDomains checks that no two active specs share the same domain.
func validateDuplicateDomains(cfg *Config) []validationError {
	var errs []validationError
	seen := make(map[string]int)
	for i, sp := range cfg.Specs {
		if sp.Disable {
			continue
		}
		if j, ok := seen[sp.Domain]; ok {
			errs = append(errs, validationError{
				spec:    sp.Domain,
				errType: "config",
				message: fmt.Sprintf("duplicate domain %q (specs #%d and #%d)",
					sp.Domain, j+1, i+1),
			})
		} else {
			seen[sp.Domain] = i
		}
	}
	return errs
}

// validateMockPorts checks for duplicate ports across mock auth config and collection base_mock_url values.
func validateMockPorts(cfg *Config, filter *Filter) []validationError {
	var errs []validationError
	usedPorts := make(map[int]string)

	if cfg.MockAuth != nil {
		if cfg.MockAuth.OAuth2Port > 0 {
			usedPorts[cfg.MockAuth.OAuth2Port] = "mock_auth.oauth2_port"
		}
		if cfg.MockAuth.DigestPort > 0 {
			usedPorts[cfg.MockAuth.DigestPort] = "mock_auth.digest_port"
		}
	}

	for _, sp := range cfg.Specs {
		if sp.Disable {
			continue
		}
		if !filter.MatchSpec(sp.Tags...) {
			continue
		}

		for _, col := range sp.Collections {
			if col.Disable {
				continue
			}
			if col.BaseMockURL != "" {
				port := extractPort(col.BaseMockURL)
				if port > 0 {
					label := sp.Domain + "/" + col.LLMTitle
					if existing, ok := usedPorts[port]; ok {
						errs = append(errs, validationError{
							spec:       sp.Domain,
							collection: col.LLMTitle,
							errType:    "config",
							message:    fmt.Sprintf("duplicate mock port %d: used by %q and %q", port, existing, label),
						})
					} else {
						usedPorts[port] = label
					}
				}
			}
		}
	}
	return errs
}

// validateSpecLocations checks that all collection spec locations are accessible
// and contain a valid OpenAPI/Swagger/Postman document.
func validateSpecLocations(cfg *Config, filter *Filter, cacheInstance *cache.Cache) []validationError {
	var errs []validationError
	for _, sp := range cfg.Specs {
		if sp.Disable {
			continue
		}
		if !filter.MatchSpec(sp.Tags...) {
			continue
		}

		for _, col := range sp.Collections {
			if col.Disable {
				continue
			}

			if cacheInstance == nil {
				continue
			}

			loc := col.Location
			err := cacheInstance.Exists(context.Background(), loc)
			if err != nil {
				ve := validationError{
					spec:       sp.Domain,
					collection: col.LLMTitle,
					location:   loc,
					errType:    "file",
					message:    err.Error(),
				}
				var locErr *cache.LocationError
				if errors.As(err, &locErr) {
					ve.errType = locErr.Type
					if locErr.Type == "url" {
						ve.message += "\n    Ensure the URL points to a raw OpenAPI/Swagger JSON or YAML file (e.g. https://example.com/openapi.json)"
					}
				}
				errs = append(errs, ve)
				continue
			}

			specPath, rErr := cacheInstance.Resolve(context.Background(), loc)
			if rErr != nil {
				continue
			}

			data, rErr := os.ReadFile(specPath)
			if rErr != nil {
				continue
			}

			if _, pErr := spec.Parse(data); pErr != nil {
				errs = append(errs, validationError{
					spec:       sp.Domain,
					collection: col.LLMTitle,
					location:   loc,
					errType:    "file",
					message:    fmt.Sprintf("location does not appear to be a valid OpenAPI/Swagger spec — expected OpenAPI 3.x, Swagger 2.0, or Postman collection: %s", pErr),
				})
			}
		}
	}
	return errs
}

// validateGlobalHTTPClient validates the global HTTP client configuration.
func validateGlobalHTTPClient(cfg *Config) []validationError {
	var errs []validationError
	if cfg.HTTPClient == nil {
		return nil
	}

	errs = append(errs, collectStructErrors("http_client", *cfg.HTTPClient)...)

	if cfg.HTTPClient.Proxy != nil {
		errs = append(errs, collectStructErrors("http_client.proxy", *cfg.HTTPClient.Proxy)...)
	}

	return errs
}

// extractPort extracts the port number from a "host:port" or "host:port/path" string.
func extractPort(addr string) int {
	// Handle "host:port/path" — cut at the first slash after the port
	_, portStr, found := strings.Cut(addr, ":")
	if !found {
		return 0
	}
	if idx := strings.IndexByte(portStr, '/'); idx >= 0 {
		portStr = portStr[:idx]
	}
	var port int
	if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
		return 0
	}
	return port
}

// getValidator returns the singleton config validator, initializing it on first call.
func getValidator() (*validator.Validate, error) {
	configValidatorMu.Lock()
	defer configValidatorMu.Unlock()
	if configValidator != nil {
		return configValidator, nil
	}
	configValidator = validator.New(
		validator.WithRequiredStructEnabled(),
	)
	if err := configValidator.RegisterValidation("domain_format", domainFormatValidation); err != nil {
		return nil, fmt.Errorf("register domain_format validation: %w", err)
	}
	if err := configValidator.RegisterValidation("title_format", titleFormatValidation); err != nil {
		return nil, fmt.Errorf("register title_format validation: %w", err)
	}
	if err := configValidator.RegisterValidation("instruction_format", instructionFormatValidation); err != nil {
		return nil, fmt.Errorf("register instruction_format validation: %w", err)
	}
	if err := configValidator.RegisterValidation("mock_addr_format", mockAddrFormatValidation); err != nil {
		return nil, fmt.Errorf("register mock_addr_format validation: %w", err)
	}
	if err := configValidator.RegisterValidation("proxy_url_format", proxyURLFormatValidation); err != nil {
		return nil, fmt.Errorf("register proxy_url_format validation: %w", err)
	}
	return configValidator, nil
}

// proxyURLFormatValidation validates that the URL has a supported proxy scheme.
func proxyURLFormatValidation(fl validator.FieldLevel) bool {
	u := fl.Field().String()
	if u == "" {
		return true
	}
	parsed, err := url.Parse(u)
	if err != nil {
		return false
	}
	switch parsed.Scheme {
	case "http", "https", "socks5", "socks5h":
		return true
	default:
		return false
	}
}

// mockAddrFormatValidation validates that the address is in format "host:port"
// or "host:port/path", where host is localhost, 127.0.0.1, or 0.0.0.0.
func mockAddrFormatValidation(fl validator.FieldLevel) bool {
	addr := fl.Field().String()
	if addr == "" {
		return true
	}

	// Try to parse as URL first (handles "host:port/path")
	if strings.Contains(addr, "://") {
		u, err := url.Parse(addr)
		if err != nil {
			return false
		}
		addr = u.Host
	}

	// Strip path suffix: "host:port/path" → "host:port"
	if hostPort, _, found := strings.Cut(addr, "/"); found {
		addr = hostPort
	}

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}

	if host != "localhost" && host != "127.0.0.1" && host != "0.0.0.0" {
		return false
	}

	if port == "" {
		return false
	}

	for _, c := range port {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

// domainFormatValidation validates that a domain matches the allowed pattern.
func domainFormatValidation(fl validator.FieldLevel) bool {
	return domainRegex.MatchString(fl.Field().String())
}

// titleFormatValidation validates that a title contains only allowed characters.
func titleFormatValidation(fl validator.FieldLevel) bool {
	return titleRegex.MatchString(fl.Field().String())
}

// instructionFormatValidation validates that an instruction contains only allowed characters.
func instructionFormatValidation(fl validator.FieldLevel) bool {
	return instructionRegex.MatchString(fl.Field().String())
}

// humanReadableError converts a validator field error into a human-readable message.
func humanReadableError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return requiredFieldError(fe.Field())
	case "min":
		return minFieldError(fe.Field(), fe.Param())
	case "max":
		return maxFieldError(fe.Field(), fe.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL — provide a full URL starting with http:// or https://", fe.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", fe.Field(), fe.Param())
	case "domain_format":
		return "Domain must be 1-60 lowercase characters using only letters, digits, hyphens, and underscores"
	case "title_format":
		return "LLMTitle contains invalid characters — use letters, digits, spaces, and basic punctuation only"
	case "instruction_format":
		return "LLMInstruction contains invalid characters — use letters, digits, spaces, and basic punctuation only"
	case "mock_addr_format":
		return fmt.Sprintf("%s must be in format 'host:port' or 'host:port/path' where host is localhost, 127.0.0.1, or 0.0.0.0 (e.g. 'localhost:8080' or '127.0.0.1:9000/v1/api')", fe.Field())
	case "proxy_url_format":
		return "Proxy URL must use a supported scheme: http, https, socks5, or socks5h (e.g. 'socks5h://127.0.0.1:1080')"
	default:
		return fe.Error()
	}
}

// requiredFieldError returns a human-readable message for a required field validation error.
func requiredFieldError(field string) string {
	switch field {
	case "Domain":
		return "Domain is required — provide a unique identifier for this API (e.g. 'meteo', 'github-api')"
	case "LLMTitle":
		return "LLMTitle is required — provide a human-readable name the LLM will use to reference this API"
	case "BaseURL":
		return "BaseURL is required — provide the base URL for all API requests (e.g. 'https://api.example.com/v1')"
	case "Location":
		return "Location is required — provide a path or URL to the Swagger/OpenAPI spec file"
	default:
		return fmt.Sprintf("%s is required", field)
	}
}

// minFieldError returns a human-readable message for a min-length validation error.
func minFieldError(field, param string) string {
	switch field {
	case "LLMTitle":
		return fmt.Sprintf("LLMTitle must be at least %s characters — provide a more descriptive name", param)
	case "Location":
		return fmt.Sprintf("Location must be at least %s characters — the path or URL is too short", param)
	case "Timeout":
		return "Timeout must be at least 1 second — set a reasonable timeout for HTTP requests"
	case "MaxRedirects":
		return "MaxRedirects must be at least 0 — set to 0 to disable redirects"
	case "MaxResponseSize":
		return "MaxResponseSize must be at least 256 bytes — the minimum response size limit"
	default:
		return fmt.Sprintf("%s must be at least %s", field, param)
	}
}

// maxFieldError returns a human-readable message for a max-length validation error.
func maxFieldError(field, param string) string {
	switch field {
	case "LLMTitle":
		return fmt.Sprintf("LLMTitle must be at most %s characters — the name is too long", param)
	case "LLMInstruction":
		return fmt.Sprintf("LLMInstruction must be at most %s characters — the instruction is too long", param)
	case "Location":
		return fmt.Sprintf("Location must be at most %s characters — the path or URL is too long", param)
	case "Timeout":
		return "Timeout must be at most 5 minutes (300 seconds) — set a reasonable timeout for HTTP requests"
	case "MaxRedirects":
		return "MaxRedirects must be at most 50 — too many redirects may cause infinite loops"
	case "MaxResponseSize":
		return "MaxResponseSize must be at most 10 MB (10485760 bytes) — the maximum response size limit"
	default:
		return fmt.Sprintf("%s must be at most %s", field, param)
	}
}

// collectStructErrors runs the validator on a struct and collects all field errors.
func collectStructErrors(prefix string, v any) []validationError {
	var errs []validationError
	val, err := getValidator()
	if err != nil {
		errs = append(errs, validationError{
			field:   prefix,
			message: fmt.Sprintf("validator initialization failed: %s", err),
		})
		return errs
	}
	if err := val.Struct(v); err != nil {
		//nolint:errorlint // validator.Struct returns ValidationErrors directly.
		fe, ok := err.(validator.ValidationErrors)
		if !ok {
			return errs
		}
		for _, f := range fe {
			fieldPath := prefix + "." + f.Field()
			errs = append(errs, validationError{
				field:   fieldPath,
				message: humanReadableError(f),
			})
		}
	}
	return errs
}
