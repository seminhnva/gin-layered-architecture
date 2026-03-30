package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/seminhnva/gin-layered-architecture/internal/api/handler"
)

func RegisterHealthRoutes(r *gin.Engine) {
	healthHandler := handler.NewHealthHandler()

	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)
}
