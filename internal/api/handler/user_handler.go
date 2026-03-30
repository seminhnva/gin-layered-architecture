package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seminhnva/gin-layered-architecture/internal/common/apperror"
	"github.com/seminhnva/gin-layered-architecture/internal/common/response"
	"github.com/seminhnva/gin-layered-architecture/internal/dto"
	"github.com/seminhnva/gin-layered-architecture/internal/model"
	"github.com/seminhnva/gin-layered-architecture/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) GetUsers(ctx *gin.Context) {
	users, err := uh.service.List(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "users retrieved successfully", toUserListResponse(users))
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var request dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Error(ctx, apperror.BadRequest(fmt.Sprintf("invalid request body: %v", err)))
		return
	}

	user, err := uh.service.Create(ctx.Request.Context(), service.CreateUserInput{
		Name:  request.Name,
		Email: request.Email,
	})
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, http.StatusCreated, "user created successfully", toUserResponse(user))
}

func (uh *UserHandler) GetUserByUUID(ctx *gin.Context) {
	var uri dto.UserURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		response.Error(ctx, apperror.BadRequest(fmt.Sprintf("invalid user id: %v", err)))
		return
	}

	user, err := uh.service.GetByUUID(ctx.Request.Context(), uri.UUID)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "user retrieved successfully", toUserResponse(user))
}

func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	var uri dto.UserURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		response.Error(ctx, apperror.BadRequest(fmt.Sprintf("invalid user id: %v", err)))
		return
	}

	var request dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Error(ctx, apperror.BadRequest(fmt.Sprintf("invalid request body: %v", err)))
		return
	}

	user, err := uh.service.Update(ctx.Request.Context(), uri.UUID, service.UpdateUserInput{
		Name:  request.Name,
		Email: request.Email,
	})
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "user updated successfully", toUserResponse(user))
}

func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	var uri dto.UserURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		response.Error(ctx, apperror.BadRequest(fmt.Sprintf("invalid user id: %v", err)))
		return
	}

	if err := uh.service.Delete(ctx.Request.Context(), uri.UUID); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, "user deleted successfully", nil)
}

func toUserListResponse(users []model.User) []dto.UserResponse {
	responses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, toUserResponse(user))
	}
	return responses
}

func toUserResponse(user model.User) dto.UserResponse {
	return dto.UserResponse{
		UUID:      user.UUID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
