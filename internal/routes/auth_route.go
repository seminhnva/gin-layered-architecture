package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/seminhnva/gin-layered-architecture/internal/api/handler"
)

type AuthRoute struct {
	handler *handler.AuthHandler
}

func NewAuthRoutes(handler *handler.AuthHandler) *AuthRoute {
	return &AuthRoute{
		handler: handler,
	}
}

func (ar *AuthRoute) Register(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", ar.handler.Login)
		auth.POST("/logout", ar.handler.Logout)
		auth.POST("/refresh", ar.handler.RefreshToken)
		auth.POST("/forgot-password", ar.handler.ForgotPassword)
		auth.POST("/reset-password", ar.handler.ResetPassword)
	}
}
