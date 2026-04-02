package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/seminhnva/gin-layered-architecture/internal/common/apperror"
	"github.com/seminhnva/gin-layered-architecture/internal/constants"
)

type ValidationError struct {
	Key   constants.ValidationMessageKey
	Param string
}

var validationKeys = map[string]constants.ValidationMessageKey{
	"required":       constants.ValidationRequired,
	"uuid":           constants.ValidationUUID,
	"email":          constants.ValidationEmail,
	"gt":             constants.ValidationGT,
	"gte":            constants.ValidationGTE,
	"lt":             constants.ValidationLT,
	"lte":            constants.ValidationLTE,
	"min":            constants.ValidationMin,
	"max":            constants.ValidationMax,
	"len":            constants.ValidationLen,
	"alpha":          constants.ValidationAlpha,
	"numeric":        constants.ValidationNumeric,
	"alphanum":       constants.ValidationAlphanum,
	"oneof":          constants.ValidationOneOf,
	"email_advanced": constants.ValidationEmailAdvanced,
	"passwordStrong": constants.ValidationPasswordStrong,
}

var defaultValidator = map[constants.ValidationMessageKey]func(param string) string{
	constants.ValidationRequired:      func(_ string) string { return "this field is required" },
	constants.ValidationUUID:          func(_ string) string { return "must be a valid UUID" },
	constants.ValidationEmail:         func(_ string) string { return "must be a valid email address" },
	constants.ValidationGT:            func(p string) string { return "must be greater than " + p },
	constants.ValidationGTE:           func(p string) string { return "must be greater than or equal to " + p },
	constants.ValidationLT:            func(p string) string { return "must be less than " + p },
	constants.ValidationLTE:           func(p string) string { return "must be less than or equal to " + p },
	constants.ValidationMin:           func(p string) string { return "must be at least " + p + " characters" },
	constants.ValidationMax:           func(p string) string { return "must be at most " + p + " characters" },
	constants.ValidationLen:           func(p string) string { return "must be exactly " + p + " characters" },
	constants.ValidationAlpha:         func(_ string) string { return "must contain only letters" },
	constants.ValidationNumeric:       func(_ string) string { return "must contain only numbers" },
	constants.ValidationAlphanum:      func(_ string) string { return "must contain only letters and numbers" },
	constants.ValidationInvalid:       func(_ string) string { return "invalid value" },
	constants.ValidationOneOf:         func(p string) string { return "must be one of: " + p },
	constants.ValidationEmailAdvanced: func(p string) string { return "invalid email" },
	constants.ValidationPasswordStrong: func(p string) string {
		return "password must contain a number, a lowercase letter, an uppercase letter, and a special character"
	},
}

func parseValidationError(e validator.FieldError) ValidationError {
	key, exist := validationKeys[e.Tag()]
	if !exist {
		key = constants.ValidationInvalid
	}
	return ValidationError{Key: key, Param: e.Param()}
}

func HandleValidationError(err error) error {
	validationErrs, exists := err.(validator.ValidationErrors)
	if !exists {
		return apperror.NewValidationError("validation error", "invalid request payload")
	}
	errors := make(map[string]string, len(validationErrs))
	for _, fieldErr := range validationErrs {
		field := fieldErr.Field()
		validationError := parseValidationError(fieldErr)

		msgFn, exists := defaultValidator[validationError.Key]
		if exists {
			errors[field] = msgFn(validationError.Param)
		} else {
			errors[field] = "invalid value"
		}
	}
	return apperror.NewValidationError("validation error", errors)
}
