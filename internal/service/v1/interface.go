package service

import (
	"context"

	"github.com/google/uuid"
	v1dto "github.com/seminhnva/gin-layered-architecture/internal/dto/v1"
)

type UserService interface {
	GetUsers(ctx context.Context, query v1dto.ListUsersQuery)
	CreateUser()
	GetUserByUUID(ID uuid.UUID)
	UpdateUser()
	DeleteUser(ID uuid.UUID)
}
