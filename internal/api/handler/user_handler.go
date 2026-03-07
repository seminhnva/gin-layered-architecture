package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/seminhnva/gin-layered-architecture/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service}
}
func (uh *UserHandler) GetUsers(ctx *gin.Context) {
	uh.service.GetUser()
}
func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	uh.service.CreateUser()
}
func (uh *UserHandler) GetUserByUUID(ctx *gin.Context) {
	uh.service.GetUserByUUID()
}
func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	uh.service.UpdateUser()
}
func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	uh.service.DeleteUser()
}
