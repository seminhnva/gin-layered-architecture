package service

import (
	"context"

	"github.com/google/uuid"
	v1dto "github.com/seminhnva/gin-layered-architecture/internal/dto/v1"
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
func (us *userService) GetUsers(ctx context.Context, query v1dto.ListUsersQuery) {
	us.repo.FindUser()
}
func (us *userService) CreateUser() {
	us.repo.Create()

}
func (us *userService) GetUserByUUID(ID uuid.UUID) {
	us.repo.FindByUUID()

}
func (us *userService) UpdateUser() {
	us.repo.Update()

}
func (us *userService) DeleteUser(ID uuid.UUID) {
	us.repo.Delete()

}
