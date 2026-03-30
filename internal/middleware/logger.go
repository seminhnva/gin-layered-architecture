package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	applog "github.com/seminhnva/gin-layered-architecture/pkg/logger"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startedAt := time.Now()
		ctx.Next()

		path := ctx.FullPath()
		if path == "" {
			path = ctx.Request.URL.Path
		}

		applog.Infof(
			"request_id=%s method=%s path=%s status=%d latency=%s client_ip=%s",
			GetRequestID(ctx),
			ctx.Request.Method,
			path,
			ctx.Writer.Status(),
			time.Since(startedAt),
			ctx.ClientIP(),
		)
	}
}
