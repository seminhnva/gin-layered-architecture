package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/seminhnva/gin-layered-architecture/internal/middleware"
)

type Route interface {
	Register(rg *gin.RouterGroup)
}

func SetUpRouter(corsAllowedOrigins []string, r *gin.Engine, httpLogger *zerolog.Logger, routes ...Route) {
	r.Use(
		middleware.RequestID(),
		middleware.Trace(),
		middleware.Logger(httpLogger),
		middleware.CORS(corsAllowedOrigins),
	)
	api := r.Group("/api/v1")
	for _, route := range routes {
		route.Register(api)
	}
}
