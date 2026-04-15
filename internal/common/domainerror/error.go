package domainerror

import "errors"

var (
	ErrTokenNotFound         = errors.New("refresh token not found")
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUserNameAlreadyExists = errors.New("username already exists")
	ErrRequired              = errors.New("required")
	ErrInvalidReference      = errors.New("invalid reference")
)
