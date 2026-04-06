package service

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/seminhnva/gin-layered-architecture/internal/auth"
	"github.com/seminhnva/gin-layered-architecture/internal/common/apperror"
	"github.com/seminhnva/gin-layered-architecture/internal/common/domainerror"
	"github.com/seminhnva/gin-layered-architecture/internal/db/sqlc"
	v1dto "github.com/seminhnva/gin-layered-architecture/internal/dto/v1"
	"github.com/seminhnva/gin-layered-architecture/internal/repository"
	"github.com/seminhnva/gin-layered-architecture/internal/utils"
	"golang.org/x/sync/errgroup"
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
func (us *userService) GetUsers(ctx context.Context, param v1dto.ListUsersQuery) (*v1dto.PaginationResponse[v1dto.UserDTO], error) {
	filter := repository.UserListFilter{
		Page:   param.Page,
		Limit:  param.Limit,
		Search: param.Search,
		Order:  param.Order,
		SortBy: param.SortBy,
	}

	var (
		users []sqlc.ListUsersRow
		total int64
	)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		// users, err = us.repo.GetUsers(ctx, query)
		users, err = us.repo.GetUsersV2(ctx, filter)
		return err
	})
	g.Go(func() error {
		var err error
		total, err = us.repo.CountUser(ctx, sqlc.CountUsersParams{
			Search: param.Search,
		})
		return err
	})
	if err := g.Wait(); err != nil {
		return nil, apperror.NewError(fmt.Sprintf("Get user: %s", err), apperror.ErrCodeInternal)
	}

	data := make([]v1dto.UserDTO, len(users))
	for i, u := range users {
		data[i] = *v1dto.ToUserDTOFromList(u)
	}
	totalPages := int32(math.Ceil(float64(total) / float64(param.Limit)))

	return &v1dto.PaginationResponse[v1dto.UserDTO]{
		Data:       data,
		Total:      total,
		Page:       int(param.Page),
		Limit:      int(param.Limit),
		TotalPages: totalPages,
	}, nil

}
func (us *userService) CreateUser(ctx context.Context, req v1dto.CreateUserRequest) (sqlc.CreateUserRow, error) {
	hashedPw, err := us.passwordService.HashPassword(req.Password)
	if err != nil {
		return sqlc.CreateUserRow{}, apperror.WrapError(err, "Internal server error", apperror.ErrCodeInternal)
	}

	user, err := us.repo.Create(ctx, sqlc.CreateUserParams{
		UserName:     req.UserName,
		Email:        utils.NormalizeString(req.Email),
		Name:         req.Name,
		PasswordHash: hashedPw,
	})
	if err != nil {
		if errors.Is(err, domainerror.ErrEmailAlreadyExists) {
			return sqlc.CreateUserRow{}, apperror.NewError("email already exists", apperror.ErrCodeConflict)
		}
		if errors.Is(err, domainerror.ErrUserNameAlreadyExists) {
			return sqlc.CreateUserRow{}, apperror.NewError("username already exists", apperror.ErrCodeConflict)
		}

		return sqlc.CreateUserRow{}, apperror.WrapError(err, "Fail to create user", apperror.ErrCodeInternal)
	}

	return user, nil
}
func (us *userService) GetUserByUUID(ctx context.Context, ID uuid.UUID) (sqlc.FindUserByIDRow, error) {
	user, err := us.repo.FindByUUID(ctx, ID)
	if err != nil {
		if errors.Is(err, domainerror.ErrUserNotFound) {
			return sqlc.FindUserByIDRow{}, apperror.NewError("User not found", apperror.ErrCodeNotFound)
		}
		return sqlc.FindUserByIDRow{}, apperror.WrapError(err, "Fail to find user", apperror.ErrCodeInternal)
	}
	return user, nil
}
func (us *userService) UpdateUser(ctx context.Context, params sqlc.UpdateUserByIDParams) (sqlc.UpdateUserByIDRow, error) {
	if params.PasswordHash != nil {
		hashedPw, err := us.passwordService.HashPassword(*params.PasswordHash)
		if err != nil {
			return sqlc.UpdateUserByIDRow{}, apperror.WrapError(err, "Internal server error", apperror.ErrCodeInternal)
		}
		params.PasswordHash = &hashedPw
	}

	user, err := us.repo.Update(ctx, params)
	if err != nil {
		if errors.Is(err, domainerror.ErrUserNotFound) {
			return sqlc.UpdateUserByIDRow{}, apperror.NewError("User not found", apperror.ErrCodeNotFound)
		}
		if errors.Is(err, domainerror.ErrEmailAlreadyExists) {
			return sqlc.UpdateUserByIDRow{}, apperror.NewError("email already exists", apperror.ErrCodeConflict)
		}
		if errors.Is(err, domainerror.ErrUserNameAlreadyExists) {
			return sqlc.UpdateUserByIDRow{}, apperror.NewError("username already exists", apperror.ErrCodeConflict)
		}
		return sqlc.UpdateUserByIDRow{}, apperror.WrapError(err, "Fail to update user", apperror.ErrCodeInternal)
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
