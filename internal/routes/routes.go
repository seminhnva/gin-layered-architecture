package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/seminhnva/gin-layered-architecture/internal/middleware"
)

type Route interface {
	Register(rg *gin.RouterGroup)
}

func SetUpRouter(r *gin.Engine, routes ...Route) {
	r.Use(middleware.Auth())
	api := r.Group("/api/v1")
	for _, route := range routes {
		route.Register(api)
	}
}
