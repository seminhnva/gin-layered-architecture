package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS(corsAllowedOrigins []string) gin.HandlerFunc {
	allowedOrigins := make(map[string]struct{}, len(corsAllowedOrigins))
	for _, origin := range corsAllowedOrigins {
		allowedOrigins[origin] = struct{}{}
	}

	return func(c *gin.Context) {
		c.Writer.Header().Add("Vary", "Origin")
		c.Writer.Header().Add("Vary", "Access-Control-Request-Method")
		c.Writer.Header().Add("Vary", "Access-Control-Request-Headers")

		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		if _, ok := allowedOrigins[origin]; !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID, X-Trace-ID")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Request-ID, X-Trace-ID")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
