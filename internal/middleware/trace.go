package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const traceIDContextKey contextKey = "trace_id"

const (
	traceIDHeader = "X-Trace-ID"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(traceIDHeader)
		if traceID == "" {
			traceID = uuid.New().String()
		}
		ctxValue := context.WithValue(c.Request.Context(), traceIDContextKey, traceID)
		c.Request = c.Request.WithContext(ctxValue)
		c.Writer.Header().Set(traceIDHeader, traceID)
		// c.Set(traceIDHeader, traceID)
		c.Next()
	}
}

func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceIDContextKey).(string); ok {
		return traceID
	}
	return ""
}
