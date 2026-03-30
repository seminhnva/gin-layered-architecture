package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/seminhnva/gin-layered-architecture/internal/config"
	"github.com/seminhnva/gin-layered-architecture/internal/middleware"
)

type Route interface {
	Register(rg *gin.RouterGroup)
}

func SetUpRouter(r *gin.Engine, cfg *config.Config, routes ...Route) {
	r.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recovery(),
		middleware.CORS(cfg.CORSAllowedOrigins),
	)

	RegisterHealthRoutes(r)

	api := r.Group("/api/v1")
	for _, route := range routes {
		route.Register(api)
	}
}
