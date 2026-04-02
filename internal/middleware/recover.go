package middleware

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Recover(recoveryLogger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				stack := debug.Stack()
				stackAt := extractPanicLocation()

				recoveryLogger.Error().
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Str("client_ip", c.ClientIP()).
					Str("request_id", GetRequestID(c)).
					Str("trace_id", GetTraceID(c.Request.Context())).
					Str("panic", fmt.Sprintf("%v", rec)).
					Str("stack_at", stackAt).
					Str("stack", string(stack)).
					Msg("panic recovered")

				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					gin.H{
						"code":    http.StatusInternalServerError,
						"message": "internal server error",
					},
				)
			}
		}()

		c.Next()
	}
}

func extractPanicLocation() string {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		if isAppPanicFrame(frame) {
			return fmt.Sprintf("%s:%d", filepath.ToSlash(frame.File), frame.Line)
		}

		if !more {
			break
		}
	}

	return ""
}

func isAppPanicFrame(frame runtime.Frame) bool {
	file := filepath.ToSlash(frame.File)
	fn := frame.Function

	switch {
	case strings.Contains(fn, "runtime."):
		return false
	case strings.Contains(fn, "runtime/debug."):
		return false
	case strings.Contains(fn, "github.com/gin-gonic/gin."):
		return false
	case strings.HasSuffix(file, "/internal/middleware/recover.go"):
		return false
	default:
		return strings.HasSuffix(file, ".go")
	}
}
