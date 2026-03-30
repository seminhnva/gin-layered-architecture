package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/seminhnva/gin-layered-architecture/internal/common/apperror"
	"github.com/seminhnva/gin-layered-architecture/internal/model"
	"github.com/seminhnva/gin-layered-architecture/internal/repository"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) List(ctx context.Context) ([]model.User, error) {
	return us.repo.List(ctx)
}

func (us *userService) Create(ctx context.Context, input CreateUserInput) (model.User, error) {
	name := strings.TrimSpace(input.Name)
	email := normalizeEmail(input.Email)

	if name == "" {
		return model.User{}, apperror.BadRequest("name is required")
	}

	if email == "" {
		return model.User{}, apperror.BadRequest("email is required")
	}

	existingUser, err := us.repo.FindByEmail(ctx, email)
	if err == nil && existingUser.UUID != "" {
		return model.User{}, apperror.Conflict("email is already in use")
	}

	if err != nil && !apperror.IsNotFound(err) {
		return model.User{}, err
	}

	uuid, err := newUUID()
	if err != nil {
		return model.User{}, apperror.Internal("failed to generate user id")
	}

	now := time.Now().UTC()
	user := model.User{
		UUID:      uuid,
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return us.repo.Create(ctx, user)
}

func (us *userService) GetByUUID(ctx context.Context, uuid string) (model.User, error) {
	return us.repo.FindByUUID(ctx, uuid)
}

func (us *userService) Update(ctx context.Context, uuid string, input UpdateUserInput) (model.User, error) {
	if input.Name == nil && input.Email == nil {
		return model.User{}, apperror.BadRequest("at least one field must be provided")
	}

	currentUser, err := us.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return model.User{}, err
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return model.User{}, apperror.BadRequest("name cannot be empty")
		}
		currentUser.Name = name
	}

	if input.Email != nil {
		email := normalizeEmail(*input.Email)
		if email == "" {
			return model.User{}, apperror.BadRequest("email cannot be empty")
		}

		foundUser, findErr := us.repo.FindByEmail(ctx, email)
		if findErr == nil && foundUser.UUID != currentUser.UUID {
			return model.User{}, apperror.Conflict("email is already in use")
		}

		if findErr != nil && !apperror.IsNotFound(findErr) {
			return model.User{}, findErr
		}

		currentUser.Email = email
	}

	currentUser.UpdatedAt = time.Now().UTC()

	return us.repo.Update(ctx, currentUser)
}

func (us *userService) Delete(ctx context.Context, uuid string) error {
	return us.repo.Delete(ctx, uuid)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func newUUID() (string, error) {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	buffer[6] = (buffer[6] & 0x0f) | 0x40
	buffer[8] = (buffer[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		buffer[0:4],
		buffer[4:6],
		buffer[6:8],
		buffer[8:10],
		buffer[10:16],
	), nil
}
