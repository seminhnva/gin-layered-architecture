package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	requestIDContextKey = "request_id"
	RequestIDHeader     = "X-Request-ID"
)

func RequestID() gin.HandlerFunc {
	return (func(ctx *gin.Context) {
		requestID := ctx.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx.Set(requestIDContextKey, requestID)
		ctx.Header(RequestIDHeader, requestID)
		ctx.Next()
	})
}

func GetRequestID(ctx *gin.Context) string {
	value, exists := ctx.Get(requestIDContextKey)
	if !exists {
		return ""
	}

	requestID, ok := value.(string)
	if !ok {
		return ""
	}

	return requestID
}
