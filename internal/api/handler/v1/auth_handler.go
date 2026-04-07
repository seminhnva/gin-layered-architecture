package handler

import "github.com/seminhnva/gin-layered-architecture/internal/service/v1"

type AuthHandler struct {
	service service.UserService
}

func NewAuthHandler(service service.UserService) *AuthHandler {
	return &AuthHandler{
		service: service}
}
