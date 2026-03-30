package service

import (
	"context"

	"github.com/seminhnva/gin-layered-architecture/internal/model"
)

type CreateUserInput struct {
	Name  string
	Email string
}

type UpdateUserInput struct {
	Name  *string
	Email *string
}

type UserService interface {
	List(ctx context.Context) ([]model.User, error)
	Create(ctx context.Context, input CreateUserInput) (model.User, error)
	GetByUUID(ctx context.Context, uuid string) (model.User, error)
	Update(ctx context.Context, uuid string, input UpdateUserInput) (model.User, error)
	Delete(ctx context.Context, uuid string) error
}
