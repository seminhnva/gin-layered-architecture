package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seminhnva/gin-layered-architecture/internal/common/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (hh *HealthHandler) Health(ctx *gin.Context) {
	response.Success(ctx, http.StatusOK, "service is healthy", gin.H{"status": "ok"})
}

func (hh *HealthHandler) Ready(ctx *gin.Context) {
	response.Success(ctx, http.StatusOK, "service is ready", gin.H{"status": "ready"})
}
