package app

import (
	"github.com/seminhnva/gin-layered-architecture/internal/api/handler"
	"github.com/seminhnva/gin-layered-architecture/internal/repository"
	"github.com/seminhnva/gin-layered-architecture/internal/routes"
	"github.com/seminhnva/gin-layered-architecture/internal/service"
	"github.com/seminhnva/gin-layered-architecture/pkg/cache"
)

type AuthModule struct {
	routes routes.Route
}

func NewAuthModule(deps *ModuleDeps, cache cache.RedisCacheService, env string) *AuthModule {
	authRepo := repository.NewAuthRepo(deps.DB, deps.Queries)
	authSerivce := service.NewAuthService(authRepo, deps.PasswordService, deps.JWTService, cache)
	authHanlder := handler.NewAuthHandler(authSerivce, env)
	authRoutes := routes.NewAuthRoutes(authHanlder)
	return &AuthModule{
		routes: authRoutes,
	}
}

func (am *AuthModule) Routes() routes.Route {
	return am.routes
}
