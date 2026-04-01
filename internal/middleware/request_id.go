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
	return (func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set(requestIDContextKey, requestID)
		c.Header(RequestIDHeader, requestID)
		c.Next()
	})
}

func GetRequestID(c *gin.Context) string {
	value, exists := c.Get(requestIDContextKey)
	if !exists {
		return ""
	}

	requestID, ok := value.(string)
	if !ok {
		return ""
	}

	return requestID
}
