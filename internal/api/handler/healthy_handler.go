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

func (hh *HealthHandler) Health(c *gin.Context) {
	response.Success(c, http.StatusOK, "service is healthy", gin.H{"status": "ok"})
}

func (hh *HealthHandler) Ready(c *gin.Context) {
	response.Success(c, http.StatusOK, "service is ready", gin.H{"status": "ready"})
}
