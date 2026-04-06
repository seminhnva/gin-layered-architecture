package routes

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/seminhnva/gin-layered-architecture/internal/middleware"
)

type Route interface {
	Register(rg *gin.RouterGroup)
}

func SetUpRouter(corsAllowedOrigins []string, r *gin.Engine, httpLogger, recoveryLogger *zerolog.Logger, routes ...Route) {
	r.Use(
		middleware.Recover(recoveryLogger),
		middleware.RequestID(),
		middleware.Trace(),
		gzip.Gzip(gzip.DefaultCompression),
		middleware.Logger(httpLogger),
		middleware.CORS(corsAllowedOrigins),
	)
	api := r.Group("/api/v1")
	for _, route := range routes {
		route.Register(api)
	}
}
