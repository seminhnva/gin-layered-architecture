package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/seminhnva/gin-layered-architecture/internal/auth"
	"github.com/seminhnva/gin-layered-architecture/internal/common/apperror"
	"github.com/seminhnva/gin-layered-architecture/internal/common/domainerror"
	"github.com/seminhnva/gin-layered-architecture/internal/db/sqlc"
	v1dto "github.com/seminhnva/gin-layered-architecture/internal/dto/v1"
	"github.com/seminhnva/gin-layered-architecture/internal/repository"
	"github.com/seminhnva/gin-layered-architecture/internal/utils"
)

type userService struct {
	repo            repository.UserRepository
	passwordService auth.Hasher
}

func NewUserService(repo repository.UserRepository, passwordService auth.Hasher) UserService {
	return &userService{
		repo:            repo,
		passwordService: passwordService,
	}
}
func (us *userService) GetUsers(ctx context.Context, query v1dto.ListUsersQuery) {
}
func (us *userService) CreateUser(ctx context.Context, req v1dto.CreateUserRequest) (sqlc.User, error) {
	hashedPw, err := us.passwordService.HashPassword(req.Password)
	if err != nil {
		return sqlc.User{}, apperror.WrapError(err, "Internal server error", apperror.ErrCodeInternal)
	}

	user, err := us.repo.Create(ctx, sqlc.CreateUserParams{
		UserName:     req.UserName,
		Email:        utils.NormalizeString(req.Email),
		Name:         req.Name,
		PasswordHash: hashedPw,
	})
	if err != nil {
		if errors.Is(err, domainerror.ErrEmailAlreadyExists) {
			return sqlc.User{}, apperror.NewError("email already exists", apperror.ErrCodeConflict)
		}
		if errors.Is(err, domainerror.ErrUserNameAlreadyExists) {
			return sqlc.User{}, apperror.NewError("username already exists", apperror.ErrCodeConflict)
		}

		return sqlc.User{}, apperror.WrapError(err, "Fail to create user", apperror.ErrCodeInternal)
	}

	return user, nil
}
func (us *userService) GetUserByUUID(ctx context.Context, ID uuid.UUID) (sqlc.User, error) {
	user, err := us.repo.FindByUUID(ctx, ID)
	if err != nil {
		if errors.Is(err, domainerror.ErrUserNotFound) {
			return sqlc.User{}, apperror.NewError("User not found", apperror.ErrCodeNotFound)
		}
		return sqlc.User{}, apperror.WrapError(err, "Fail to find user", apperror.ErrCodeInternal)
	}
	return user, nil
}
func (us *userService) UpdateUser(ctx context.Context, params sqlc.UpdateUserByIDParams) (sqlc.User, error) {
	if params.PasswordHash != nil {
		hashedPw, err := us.passwordService.HashPassword(*params.PasswordHash)
		if err != nil {
			return sqlc.User{}, apperror.WrapError(err, "Internal server error", apperror.ErrCodeInternal)
		}
		params.PasswordHash = &hashedPw
	}

	user, err := us.repo.Update(ctx, params)
	if err != nil {
		if errors.Is(err, domainerror.ErrUserNotFound) {
			return sqlc.User{}, apperror.NewError("User not found", apperror.ErrCodeNotFound)
		}
		if errors.Is(err, domainerror.ErrEmailAlreadyExists) {
			return sqlc.User{}, apperror.NewError("email already exists", apperror.ErrCodeConflict)
		}
		if errors.Is(err, domainerror.ErrUserNameAlreadyExists) {
			return sqlc.User{}, apperror.NewError("username already exists", apperror.ErrCodeConflict)
		}
		return sqlc.User{}, apperror.WrapError(err, "Fail to update user", apperror.ErrCodeInternal)
	}
	return user, nil

}

func (us *userService) DeleteUser(ctx context.Context, ID uuid.UUID) error {
	err := us.repo.SoftDelete(ctx, ID)
	if err != nil {
		if errors.Is(err, domainerror.ErrUserNotFound) {
			return apperror.NewError("User not found", apperror.ErrCodeNotFound)
		}
		return apperror.WrapError(err, "Fail to delete user", apperror.ErrCodeInternal)
	}
	return nil
}
