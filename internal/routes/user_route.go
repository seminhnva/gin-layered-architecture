package routes

import (
	"github.com/gin-gonic/gin"
	v1handler "github.com/seminhnva/gin-layered-architecture/internal/api/handler/v1"
)

type UserRoute struct {
	handler *v1handler.UserHandler
}

func NewUserRoutes(handler *v1handler.UserHandler) *UserRoute {
	return &UserRoute{
		handler: handler,
	}
}

func (ur *UserRoute) Register(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.GET("", ur.handler.GetUsers)
		users.POST("", ur.handler.CreateUser)
		users.GET("/:uuid", ur.handler.GetUserByUUID)
		users.PUT("/:uuid", ur.handler.UpdateUser)
		users.DELETE("/:uuid", ur.handler.DeleteUser)
	}
}
