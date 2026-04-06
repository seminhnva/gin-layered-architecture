package handler

import (
	"net/http"

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
	users, err := uh.service.GetUsers(c.Request.Context(), params)
	if err != nil {
		response.Error(c, err)
		return
	}
	meta := &response.Meta{
		Total:      users.Total,
		Page:       users.Page,
		Limit:      users.Limit,
		TotalPages: users.TotalPages,
	}
	links :=
		response.BuildPaginationLinks(
			c,
			int(users.Page),
			int(users.Limit),
			users.TotalPages,
			v1dto.ListUsersAllowedQueryKeys)
	response.Paginated(c, http.StatusOK, users.Data, meta, links)
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
	var req v1dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validation.HandleValidationError(err))
		return
	}
	req.ToCreateUserParams()
	user, err := uh.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusCreated, v1dto.ToUserDTOFromCreate(user))
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

	user, err := uh.service.GetUserByUUID(c.Request.Context(), UserID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, v1dto.ToUserDTOFromFind(user))
}
func (uh *UserHandler) UpdateUser(c *gin.Context) {
	var params v1dto.GetUserIdParams
	if err := c.ShouldBindUri(&params); err != nil {
		response.Error(c, validation.HandleValidationError(err))
		return
	}
	var req v1dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validation.HandleValidationError(err))
		return
	}

	UserID, err := uuid.Parse(params.UUID)
	if err != nil {
		response.Error(c, apperror.NewError("invalid user id", apperror.ErrCodeBadRequest))
		return
	}

	updateParams := req.ToUpdateUserParams(UserID)

	user, err := uh.service.UpdateUser(c.Request.Context(), updateParams)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, v1dto.ToUserDTOFromUpdate(user))
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
	if err := uh.service.DeleteUser(c.Request.Context(), userID); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, nil)
}
