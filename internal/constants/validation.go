package constants

type ValidationMessageKey string

const (
	ValidationRequired       ValidationMessageKey = "validation.required"
	ValidationUUID           ValidationMessageKey = "validation.uuid"
	ValidationEmail          ValidationMessageKey = "validation.email"
	ValidationGT             ValidationMessageKey = "validation.gt"
	ValidationGTE            ValidationMessageKey = "validation.gte"
	ValidationLT             ValidationMessageKey = "validation.lt"
	ValidationLTE            ValidationMessageKey = "validation.lte"
	ValidationMin            ValidationMessageKey = "validation.min"
	ValidationMax            ValidationMessageKey = "validation.max"
	ValidationLen            ValidationMessageKey = "validation.len"
	ValidationAlpha          ValidationMessageKey = "validation.alpha"
	ValidationNumeric        ValidationMessageKey = "validation.numeric"
	ValidationAlphanum       ValidationMessageKey = "validation.alphanum"
	ValidationInvalid        ValidationMessageKey = "validation.invalid"
	ValidationOneOf          ValidationMessageKey = "validation.oneof"
	ValidationEmailAdvanced  ValidationMessageKey = "validation.email_advanced"
	ValidationPasswordStrong ValidationMessageKey = "validation.passwordStrong"
)
