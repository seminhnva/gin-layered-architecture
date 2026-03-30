package repository

import (
	"context"

	"github.com/seminhnva/gin-layered-architecture/internal/model"
)

type UserRepository interface {
	List(ctx context.Context) ([]model.User, error)
	Create(ctx context.Context, user model.User) (model.User, error)
	FindByUUID(ctx context.Context, uuid string) (model.User, error)
	FindByEmail(ctx context.Context, email string) (model.User, error)
	Update(ctx context.Context, user model.User) (model.User, error)
	Delete(ctx context.Context, uuid string) error
}
