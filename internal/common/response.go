package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seminhnva/gin-layered-architecture/internal/common/apperror"
	"github.com/seminhnva/gin-layered-architecture/internal/middleware"
)

type ErrorResponse struct {
	Code    apperror.ErrorCode `json:"code,omitempty"`
	Details any                `json:"details,omitempty"`
}

type APIResponse struct {
	Success    bool           `json:"success"`
	StatusCode int            `json:"status_code,omitempty"`
	Message    string         `json:"message,omitempty"`
	Data       any            `json:"data,omitempty"`
	Error      *ErrorResponse `json:"error,omitempty"`
	RequestID  string         `json:"request_id,omitempty"`
}

func Success(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		RequestID: middleware.GetRequestID(c),
	})

}

func Error(c *gin.Context, err error) {
	if appErr, ok := err.(*apperror.Error); ok {
		apperror.SetErrorContext(c, appErr)

		errRes := APIResponse{
			Success:   false,
			Message:   appErr.Message,
			RequestID: middleware.GetRequestID(c),
			Error: &ErrorResponse{
				Code: appErr.Code,
			},
		}

		if appErr.Code == apperror.ErrCodeValidation && appErr.Details != nil {
			errRes.Error.Details = appErr.Details
		}

		c.JSON(apperror.HttpStatusCodeFromCode(appErr.Code), errRes)
		return
	}

	internalErr := &apperror.Error{
		Status:  http.StatusInternalServerError,
		Code:    apperror.ErrCodeInternal,
		Message: "internal server error",
		Err:     err,
	}

	apperror.SetErrorContext(c, internalErr)

	c.JSON(http.StatusInternalServerError, APIResponse{
		Success:   false,
		Message:   internalErr.Message,
		RequestID: middleware.GetRequestID(c),
		Error: &ErrorResponse{
			Code: internalErr.Code,
		},
	})
}
