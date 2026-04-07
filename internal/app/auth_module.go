package app

import (
	"github.com/seminhnva/gin-layered-architecture/internal/api/handler/v1"
	"github.com/seminhnva/gin-layered-architecture/internal/repository"
	"github.com/seminhnva/gin-layered-architecture/internal/routes"
	"github.com/seminhnva/gin-layered-architecture/internal/service/v1"
)

type AuthModule struct {
	routes routes.Route
}

func NewAuthModule(deps *ModuleDeps) *AuthModule {
	authRepo := repository.NewAuthRepo(deps.db, deps.Queries)
	authSerivce := service.NewUserService(authRepo, deps.PasswordService)
	authHanlder := handler.NewUserHandler(authSerivce)
	authRoutes := routes.NewUserRoutes(authHanlder)
	return &AuthModule{
		routes: authRoutes,
	}
}

func (am *AuthModule) Routes() routes.Route {
	return am.routes
}
