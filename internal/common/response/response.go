package response

import (
	"github.com/gin-gonic/gin"
	"github.com/seminhnva/gin-layered-architecture/internal/common/apperror"
	"github.com/seminhnva/gin-layered-architecture/internal/middleware"
)

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type APIResponse struct {
	Success   bool         `json:"success"`
	Message   string       `json:"message,omitempty"`
	Data      any          `json:"data,omitempty"`
	Error     *ErrorDetail `json:"error,omitempty"`
	RequestID string       `json:"request_id,omitempty"`
}

func Success(ctx *gin.Context, status int, message string, data any) {
	ctx.JSON(status, APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		RequestID: middleware.GetRequestID(ctx),
	})
}

func Error(ctx *gin.Context, err error) {
	status := apperror.HTTPStatus(err)
	ctx.AbortWithStatusJSON(status, APIResponse{
		Success: false,
		Message: apperror.Message(err),
		Error: &ErrorDetail{
			Code:    apperror.Code(err),
			Message: apperror.Message(err),
		},
		RequestID: middleware.GetRequestID(ctx),
	})
}
