package config

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	configValidator *validator.Validate
	// domainRegex for domain validation in config
	domainRegex      = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,60}$`)
	titleRegex       = regexp.MustCompile(`^[\p{L}\p{N} #*_` + "`" + `~>\[\]()|.,!?;:'"\\-]+$`)
	instructionRegex = regexp.MustCompile(`^[\p{L}\p{N}\s#*_` + "`" + `~>\[\]()|.,!?;:'"\\-]+$`)
)

func init() {
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
