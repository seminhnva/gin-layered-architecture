package app

import (
	"github.com/seminhnva/gin-layered-architecture/internal/api/handler"
	"github.com/seminhnva/gin-layered-architecture/internal/repository"
	"github.com/seminhnva/gin-layered-architecture/internal/routes"
	"github.com/seminhnva/gin-layered-architecture/internal/service"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule() *UserModule {
	userRepo := repository.NewUserRepo()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userRoutes := routes.NewUserRoutes(userHandler)

	return &UserModule{
		routes: userRoutes,
	}
}

func (um *UserModule) Routes() routes.Route {
	return um.routes
}
