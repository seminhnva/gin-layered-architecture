package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/seminhnva/gin-layered-architecture/internal/service"
)

type AuthHandler struct {
	service service.AuthSerivce
}

func NewAuthHandler(service service.AuthSerivce) *AuthHandler {
	return &AuthHandler{
		service: service}
}

func (as *AuthHandler) Login(c *gin.Context)          {}
func (as *AuthHandler) Logout(c *gin.Context)         {}
func (as *AuthHandler) RefreshToken(c *gin.Context)   {}
func (as *AuthHandler) ForgotPassword(c *gin.Context) {}
func (as *AuthHandler) ResetPassword(c *gin.Context)  {}
