package service

import "github.com/seminhnva/gin-layered-architecture/internal/repository"

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}
func (us *userService) GetUser() {
	us.repo.FindUser()
}
func (us *userService) CreateUser() {
	us.repo.Create()

}
func (us *userService) GetUserByUUID() {
	us.repo.FindByUUID()

}
func (us *userService) UpdateUser() {
	us.repo.Update()

}
func (us *userService) DeleteUser() {
	us.repo.Delete()

}
