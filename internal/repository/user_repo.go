package repository

import (
	"context"
	"strings"
	"sync"

	"github.com/seminhnva/gin-layered-architecture/internal/common/apperror"
	"github.com/seminhnva/gin-layered-architecture/internal/model"
)

type UserRepo struct {
	mu    sync.RWMutex
	users []model.User
}

func NewUserRepo() UserRepository {
	return &UserRepo{
		users: make([]model.User, 0),
	}
}

func (ur *UserRepo) List(context.Context) ([]model.User, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	users := make([]model.User, len(ur.users))
	copy(users, ur.users)

	return users, nil
}

func (ur *UserRepo) FindByUUID(_ context.Context, uuid string) (model.User, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	for _, user := range ur.users {
		if user.UUID == uuid {
			return user, nil
		}
	}

	return model.User{}, apperror.NotFound("user not found")
}

func (ur *UserRepo) FindByEmail(_ context.Context, email string) (model.User, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	for _, user := range ur.users {
		if strings.EqualFold(user.Email, email) {
			return user, nil
		}
	}

	return model.User{}, apperror.NotFound("user not found")
}

func (ur *UserRepo) Create(_ context.Context, user model.User) (model.User, error) {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	ur.users = append(ur.users, user)

	return user, nil
}

func (ur *UserRepo) Update(_ context.Context, user model.User) (model.User, error) {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	for index, currentUser := range ur.users {
		if currentUser.UUID == user.UUID {
			ur.users[index] = user
			return user, nil
		}
	}

	return model.User{}, apperror.NotFound("user not found")
}

func (ur *UserRepo) Delete(_ context.Context, uuid string) error {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	for index, user := range ur.users {
		if user.UUID == uuid {
			ur.users = append(ur.users[:index], ur.users[index+1:]...)
			return nil
		}
	}

	return apperror.NotFound("user not found")
}
