package apperror

import (
	"errors"
	"net/http"
)

const (
	CodeBadRequest   = "BAD_REQUEST"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"
	CodeNotFound     = "NOT_FOUND"
	CodeConflict     = "CONFLICT"
	CodeInternal     = "INTERNAL_SERVER_ERROR"
)

type Error struct {
	Status  int
	Code    string
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func New(status int, code, message string) *Error {
	return &Error{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func BadRequest(message string) *Error {
	return New(http.StatusBadRequest, CodeBadRequest, message)
}

func Unauthorized(message string) *Error {
	return New(http.StatusUnauthorized, CodeUnauthorized, message)
}

func Forbidden(message string) *Error {
	return New(http.StatusForbidden, CodeForbidden, message)
}

func NotFound(message string) *Error {
	return New(http.StatusNotFound, CodeNotFound, message)
}

func Conflict(message string) *Error {
	return New(http.StatusConflict, CodeConflict, message)
}

func Internal(message string) *Error {
	return New(http.StatusInternalServerError, CodeInternal, message)
}

func Normalize(err error) *Error {
	if err == nil {
		return nil
	}

	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr
	}

	return Internal("internal server error")
}

func HTTPStatus(err error) int {
	return Normalize(err).Status
}

func Code(err error) string {
	return Normalize(err).Code
}

func Message(err error) string {
	return Normalize(err).Message
}

func IsNotFound(err error) bool {
	return HTTPStatus(err) == http.StatusNotFound
}
