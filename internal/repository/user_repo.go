package repository

import "github.com/seminhnva/gin-layered-architecture/internal/model"

type UserRepo struct {
	users []model.User
}

func NewUserRepo() UserRepository {
	return &UserRepo{
		users: make([]model.User, 0),
	}
}

func (ur *UserRepo) FindUser()   {}
func (ur *UserRepo) FindByUUID() {}
func (ur *UserRepo) Create()     {}
func (ur *UserRepo) Update()     {}
func (ur *UserRepo) Delete()     {}
