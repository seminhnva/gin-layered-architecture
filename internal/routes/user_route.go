package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/seminhnva/gin-layered-architecture/internal/api/handler"
	"github.com/seminhnva/gin-layered-architecture/internal/middleware"
)

type UserRoute struct {
	handler *handler.UserHandler
}

func NewUserRoutes(handler *handler.UserHandler) *UserRoute {
	return &UserRoute{
		handler: handler,
	}
}

func (ur *UserRoute) Register(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.Use(middleware.Auth())
	{
		users.GET("", ur.handler.GetUsers)
		users.POST("", ur.handler.CreateUser)
		users.GET("/:uuid", ur.handler.GetUserByUUID)
		users.PUT("/:uuid", ur.handler.UpdateUser)
		users.DELETE("/:uuid", ur.handler.DeleteUser)
	}
}
