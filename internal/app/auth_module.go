package app

import (
	"github.com/rs/zerolog"
	"github.com/seminhnva/gin-layered-architecture/internal/api/handler"
	"github.com/seminhnva/gin-layered-architecture/internal/repository"
	"github.com/seminhnva/gin-layered-architecture/internal/routes"
	"github.com/seminhnva/gin-layered-architecture/internal/service"
	"github.com/seminhnva/gin-layered-architecture/pkg/cache"
)

type AuthModule struct {
	routes routes.Route
}

func NewAuthModule(deps *ModuleDeps, cache cache.RedisCacheService, env string, rateLimiterLogger *zerolog.Logger) *AuthModule {
	userRepo := repository.NewUserRepo(deps.DB, deps.Queries)
	authSerivce := service.NewAuthService(userRepo, deps.PasswordService, deps.JWTService, cache)
	authHanlder := handler.NewAuthHandler(authSerivce, env)
	authRoutes := routes.NewAuthRoutes(authHanlder, rateLimiterLogger)
	return &AuthModule{
		routes: authRoutes,
	}
}

func (am *AuthModule) Routes() routes.Route {
	return am.routes
}
