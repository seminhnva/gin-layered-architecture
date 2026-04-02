package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/seminhnva/gin-layered-architecture/internal/api/validation"
	"github.com/seminhnva/gin-layered-architecture/internal/common/apperror"
	"github.com/seminhnva/gin-layered-architecture/internal/common/response"
	v1dto "github.com/seminhnva/gin-layered-architecture/internal/dto/v1"
	"github.com/seminhnva/gin-layered-architecture/internal/service/v1"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service}
}
func (uh *UserHandler) GetUsers(c *gin.Context) {
	var params v1dto.ListUsersQuery

	if err := c.ShouldBindQuery(&params); err != nil {
		response.Error(c, validation.HandleValidationError(err))
		return
	}
	params.Normalize()
	arg := v1dto.ListUsersQuery{
		Page:   params.Page,
		Limit:  params.Limit,
		SortBy: params.SortBy,
		Order:  params.Order,
		Search: params.Search,
	}
	uh.service.GetUsers(c.Request.Context(), arg)

}
func (uh *UserHandler) CreateUser(c *gin.Context) {
	var req v1dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validation.HandleValidationError(err))
		return
	}

	uh.service.CreateUser()
}
func (uh *UserHandler) GetUserByUUID(c *gin.Context) {
	var params v1dto.GetUserIdParams
	if err := c.ShouldBindUri(&params); err != nil {
		response.Error(c, validation.HandleValidationError(err))
		return
	}
	UserID, err := uuid.Parse(params.UUID)
	if err != nil {
		response.Error(c, apperror.NewError("invalid user id", apperror.ErrCodeBadRequest))
		return
	}

	uh.service.GetUserByUUID(UserID)
}
func (uh *UserHandler) UpdateUser(c *gin.Context) {
	var req v1dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validation.HandleValidationError(err))
		return
	}
	uh.service.UpdateUser()
}
func (uh *UserHandler) DeleteUser(c *gin.Context) {
	var params v1dto.GetUserIdParams
	if err := c.ShouldBindUri(&params); err != nil {
		response.Error(c, validation.HandleValidationError(err))
		return
	}
	userID, err := uuid.Parse(params.UUID)
	if err != nil {
		response.Error(c, apperror.NewError("invalid user id", apperror.ErrCodeBadRequest))
		return
	}
	uh.service.DeleteUser(userID)
}
