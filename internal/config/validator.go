package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/mmadfox/swag2mcp/internal/cache"
)

var (
	configValidator   *validator.Validate //nolint:gochecknoglobals // lazily initialized
	configValidatorMu sync.Mutex          //nolint:gochecknoglobals // guards validator init
	domainRegex       = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,60}$`)
	titleRegex        = regexp.MustCompile(`^[\p{L}\p{N} #*_` + "`" + `~>\[\]()|.,!?;:'"\\-]+$`)
	instructionRegex  = regexp.MustCompile(`^[\p{L}\p{N}\s#*_` + "`" + `~>\[\]()|.,!?;:'"\\-]+$`)
	scriptDomainRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,60}$`)
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

	filter := NewFilter(opts.Tags)

	if err := cfg.Validate(filter); err != nil {
		errs = append(errs, validationError{
			errType: "config",
			message: err.Error(),
		})
	}

	seen := make(map[string]int)
	for i, spec := range cfg.Specs {
		if spec.Disable {
			continue
		}
		if j, ok := seen[spec.Domain]; ok {
			errs = append(errs, validationError{
				spec:    spec.Domain,
				errType: "config",
				message: fmt.Sprintf("duplicate domain %q (specs #%d and #%d)", spec.Domain, j+1, i+1),
			})
		} else {
			seen[spec.Domain] = i
		}
	}

	for _, spec := range cfg.Specs {
		if spec.Disable {
			continue
		}
		if !filter.MatchSpec(spec.Tags...) {
			continue
		}

		for _, col := range spec.Collections {
			if col.Disable {
				continue
			}

			if opts.Cache == nil {
				continue
			}

			loc := col.Location
			err := opts.Cache.Exists(loc)
			if err != nil {
				ve := validationError{
					spec:       spec.Domain,
					collection: col.LLMTitle,
					location:   loc,
					errType:    "file",
					message:    err.Error(),
				}
				var locErr *cache.LocationError
				if errors.As(err, &locErr) {
					ve.errType = locErr.Type
				}
				errs = append(errs, ve)
			}
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return errs
}

func getValidator() *validator.Validate {
	configValidatorMu.Lock()
	defer configValidatorMu.Unlock()
	if configValidator != nil {
		return configValidator
	}
	configValidator = validator.New(
		validator.WithRequiredStructEnabled(),
	)
	if err := configValidator.RegisterValidation("domain_format", domainFormatValidation); err != nil {
		panic(err)
	}
	if err := configValidator.RegisterValidation("title_format", titleFormatValidation); err != nil {
		panic(err)
	}
	if err := configValidator.RegisterValidation("instruction_format", instructionFormatValidation); err != nil {
		panic(err)
	}
	if err := configValidator.RegisterValidation("script_domain_format", scriptDomainFormatValidation); err != nil {
		panic(err)
	}
	return configValidator
}

func domainFormatValidation(fl validator.FieldLevel) bool {
	return domainRegex.MatchString(fl.Field().String())
}

func titleFormatValidation(fl validator.FieldLevel) bool {
	return titleRegex.MatchString(fl.Field().String())
}

func instructionFormatValidation(fl validator.FieldLevel) bool {
	return instructionRegex.MatchString(fl.Field().String())
}

func scriptDomainFormatValidation(fl validator.FieldLevel) bool {
	return scriptDomainRegex.MatchString(fl.Field().String())
}

// humanReadableError translates a validator.FieldError into a human-readable message.
func humanReadableError(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()
	param := fe.Param()

	switch tag {
	case "required":
		switch field {
		case "Domain":
			return "Domain is required — provide a unique identifier for this API (e.g. 'petstore', 'github-api')"
		case "LLMTitle":
			return "LLMTitle is required — provide a human-readable name the LLM will use to reference this API"
		case "BaseURL":
			return "BaseURL is required — provide the base URL for all API requests (e.g. 'https://api.example.com/v1')"
		case "Location":
			return "Location is required — provide a path or URL to the Swagger/OpenAPI spec file"
		default:
			return fmt.Sprintf("%s is required", field)
		}
	case "min":
		switch field {
		case "LLMTitle":
			return fmt.Sprintf("LLMTitle must be at least %s characters — provide a more descriptive name", param)
		case "Location":
			return fmt.Sprintf("Location must be at least %s characters — the path or URL is too short", param)
		default:
			return fmt.Sprintf("%s must be at least %s characters", field, param)
		}
	case "max":
		switch field {
		case "LLMTitle":
			return fmt.Sprintf("LLMTitle must be at most %s characters — the name is too long", param)
		case "LLMInstruction":
			return fmt.Sprintf("LLMInstruction must be at most %s characters — the instruction is too long", param)
		case "Location":
			return fmt.Sprintf("Location must be at most %s characters — the path or URL is too long", param)
		default:
			return fmt.Sprintf("%s must be at most %s characters", field, param)
		}
	case "url":
		return fmt.Sprintf("%s must be a valid URL — provide a full URL starting with http:// or https://", field)
	case "domain_format":
		return "Domain must be 1-60 characters using only letters, digits, hyphens, and underscores"
	case "title_format":
		return "LLMTitle contains invalid characters — use letters, digits, spaces, and basic punctuation only"
	case "instruction_format":
		return "LLMInstruction contains invalid characters — use letters, digits, spaces, and basic punctuation only"
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	default:
		return fe.Error()
	}
}

// collectStructErrors runs the validator on a struct and collects all field errors.
func collectStructErrors(prefix string, v any) []validationError {
	var errs []validationError
	if err := getValidator().Struct(v); err != nil {
		//nolint:errorlint // validator.Struct returns ValidationErrors directly
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
