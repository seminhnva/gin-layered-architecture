package apperror

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	errorCodeContextKey    = "error_code"
	errorMessageContextKey = "error_message"
	errorDetailsContextKey = "error_details"
	rootErrorContextKey    = "root_error"
)

type ErrorCode string

type LogContext struct {
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Details   any    `json:"details,omitempty"`
	RootError string `json:"root_error,omitempty"`
}

type Error struct {
	Status  int
	Code    ErrorCode
	Message string
	Details any
	Err     error
}

func (ed *Error) Error() string {
	return ed.Message
}

const (
	ErrCodeBadRequest      ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeConflict        ErrorCode = "CONFLICT"
	ErrCodeValidation      ErrorCode = "VALIDATION_ERROR"
	ErrCodeInternal        ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrCodeUserNotFound    ErrorCode = "USER_NOT_FOUND"
	ErrCodeEmailExists     ErrorCode = "EMAIL_ALREADY_EXISTS"
	ErrCodeInvalidPassword ErrorCode = "INVALID_PASSWORD"
	ErrCodeTooManyRequest  ErrorCode = "TOO_MANY_REQUEST"
)

var errorCodeToHTTPStatus = map[ErrorCode]int{
	ErrCodeBadRequest:      http.StatusBadRequest,          // 400
	ErrCodeValidation:      http.StatusUnprocessableEntity, // 422
	ErrCodeUnauthorized:    http.StatusUnauthorized,        // 401
	ErrCodeInvalidPassword: http.StatusUnauthorized,        // 401
	ErrCodeForbidden:       http.StatusForbidden,           // 403
	ErrCodeNotFound:        http.StatusNotFound,            // 404
	ErrCodeUserNotFound:    http.StatusNotFound,            // 404
	ErrCodeConflict:        http.StatusConflict,            // 409
	ErrCodeEmailExists:     http.StatusConflict,            // 409
	ErrCodeInternal:        http.StatusInternalServerError, // 500
	ErrCodeTooManyRequest:  http.StatusTooManyRequests,     // 429
}

func HttpStatusCodeFromCode(code ErrorCode) int {
	if statusCode, exist := errorCodeToHTTPStatus[code]; exist {
		return statusCode
	}
	return http.StatusInternalServerError
}

func NewError(message string, code ErrorCode) error {
	return &Error{
		Status:  HttpStatusCodeFromCode(code),
		Message: message,
		Code:    code,
	}
}

func WrapError(err error, message string, code ErrorCode) error {
	return &Error{
		Status:  HttpStatusCodeFromCode(code),
		Message: message,
		Code:    code,
		Err:     err,
	}
}

func NewValidationError(message string, details any) error {
	return &Error{
		Status:  HttpStatusCodeFromCode(ErrCodeValidation),
		Code:    ErrCodeValidation,
		Message: message,
		Details: details,
	}
}

func SetErrorContext(c *gin.Context, appErr *Error) {
	if appErr == nil {
		return
	}

	c.Set(errorCodeContextKey, string(appErr.Code))
	c.Set(errorMessageContextKey, appErr.Message)

	if appErr.Details != nil {
		c.Set(errorDetailsContextKey, appErr.Details)
	}

	if appErr.Err != nil {
		c.Set(rootErrorContextKey, appErr.Err.Error())
	}
}

func GetLogContext(c *gin.Context) (*LogContext, bool) {
	logCtx := &LogContext{}

	if value, ok := c.Get(errorCodeContextKey); ok {
		if code, ok := value.(string); ok {
			logCtx.Code = code
		}
	}

	if value, ok := c.Get(errorMessageContextKey); ok {
		if message, ok := value.(string); ok {
			logCtx.Message = message
		}
	}

	if value, ok := c.Get(errorDetailsContextKey); ok {
		logCtx.Details = value
	}

	if value, ok := c.Get(rootErrorContextKey); ok {
		if rootErr, ok := value.(string); ok {
			logCtx.RootError = rootErr
		}
	}

	if logCtx.Code == "" && logCtx.Message == "" && logCtx.Details == nil && logCtx.RootError == "" {
		return nil, false
	}

	return logCtx, true
}
