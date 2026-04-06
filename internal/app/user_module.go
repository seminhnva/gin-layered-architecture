package app

import (
	"github.com/seminhnva/gin-layered-architecture/internal/api/handler/v1"
	"github.com/seminhnva/gin-layered-architecture/internal/repository"
	"github.com/seminhnva/gin-layered-architecture/internal/routes"
	"github.com/seminhnva/gin-layered-architecture/internal/service/v1"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule(deps *ModuleDeps) *UserModule {
	userRepo := repository.NewUserRepo(deps.db, deps.Queries)
	userSerivce := service.NewUserService(userRepo, deps.PasswordService)
	userHanlder := handler.NewUserHandler(userSerivce)
	userRoutes := routes.NewUserRoutes(userHanlder)
	return &UserModule{
		routes: userRoutes,
	}
}

func (um *UserModule) Routes() routes.Route {
	return um.routes
}
