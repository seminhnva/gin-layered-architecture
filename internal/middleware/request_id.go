package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const (
	RequestIDHeader = "X-Request-ID"
	requestIDKey    = "request_id"
)

func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = newRequestID()
		}

		ctx.Set(requestIDKey, requestID)
		ctx.Writer.Header().Set(RequestIDHeader, requestID)
		ctx.Next()
	}
}

func GetRequestID(ctx *gin.Context) string {
	value, exists := ctx.Get(requestIDKey)
	if !exists {
		return ""
	}

	requestID, ok := value.(string)
	if !ok {
		return ""
	}

	return requestID
}

func newRequestID() string {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return "request-id-unavailable"
	}

	return hex.EncodeToString(buffer)
}
