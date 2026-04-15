package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/seminhnva/gin-layered-architecture/internal/api/handler"
	"github.com/seminhnva/gin-layered-architecture/internal/middleware"
)

type AuthRoute struct {
	handler           *handler.AuthHandler
	rateLimiterLogger *zerolog.Logger
}

func NewAuthRoutes(handler *handler.AuthHandler, rateLimiterLogger *zerolog.Logger) *AuthRoute {
	return &AuthRoute{
		handler:           handler,
		rateLimiterLogger: rateLimiterLogger,
	}
}

func (ar *AuthRoute) Register(rg *gin.RouterGroup) {
	loginLimiter := middleware.RateLimiter(ar.rateLimiterLogger, middleware.LoginRateLimitPolicy)
	refreshLimiter := middleware.RateLimiter(ar.rateLimiterLogger, middleware.RefreshRateLimitPolicy)
	forgotPasswordLimiter := middleware.RateLimiter(ar.rateLimiterLogger, middleware.ForgotPasswordRateLimitPolicy)
	auth := rg.Group("/auth")
	{
		auth.POST("/login", loginLimiter, ar.handler.Login)
		auth.POST("/logout", ar.handler.Logout)
		auth.POST("/refresh", refreshLimiter, ar.handler.RefreshToken)
		auth.POST("/forgot-password", forgotPasswordLimiter, ar.handler.ForgotPassword)
		auth.POST("/reset-password", ar.handler.ResetPassword)
	}
}
