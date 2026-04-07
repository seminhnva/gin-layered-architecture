package service

import (
	"github.com/seminhnva/gin-layered-architecture/internal/auth"
	"github.com/seminhnva/gin-layered-architecture/internal/repository"
)

type authService struct {
	repo            repository.UserRepository
	passwordService auth.Hasher
}

func NewAuthService(repo repository.UserRepository, passwordService auth.Hasher) AuthSerivce {
	return &authService{
		repo:            repo,
		passwordService: passwordService,
	}
}
