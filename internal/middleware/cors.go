package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	allowMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
	allowHeaders = "Content-Type, Authorization, X-Request-ID"
)

func CORS(allowedOrigins []string) gin.HandlerFunc {
	allowedOriginSet := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowedOriginSet[origin] = struct{}{}
	}

	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		if origin == "" {
			ctx.Next()
			return
		}

		headers := ctx.Writer.Header()
		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")

		if _, exists := allowedOriginSet[origin]; !exists {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success":    false,
				"message":    "origin is not allowed",
				"request_id": GetRequestID(ctx),
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "origin is not allowed",
				},
			})
			return
		}

		headers.Set("Access-Control-Allow-Origin", origin)
		headers.Set("Access-Control-Allow-Methods", allowMethods)
		headers.Set("Access-Control-Allow-Headers", allowHeaders)
		headers.Set("Access-Control-Max-Age", "600")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
