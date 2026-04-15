package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seminhnva/gin-layered-architecture/internal/api/validation"
	"github.com/seminhnva/gin-layered-architecture/internal/auth"
	"github.com/seminhnva/gin-layered-architecture/internal/common/apperror"
	"github.com/seminhnva/gin-layered-architecture/internal/common/response"
	"github.com/seminhnva/gin-layered-architecture/internal/dto"
	"github.com/seminhnva/gin-layered-architecture/internal/service"
)

type AuthHandler struct {
	service service.AuthSerivce
	env     string
}

func NewAuthHandler(service service.AuthSerivce, env string) *AuthHandler {
	return &AuthHandler{
		service: service,
		env:     env,
	}
}

func (as *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validation.HandleValidationError(err))
		return
	}

	tokenInfo, err := as.service.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	maxAge := int(time.Until(tokenInfo.RefreshTokenExpiresAt).Seconds())

	c.SetCookie(
		dto.RefreshTokenCookieName, tokenInfo.RefreshToken,
		maxAge,
		"/",
		"",
		as.env == "production",
		true,
	)

	tokenData := dto.LoginResponse{
		AccessToken: tokenInfo.AccessToken,
		TokenType:   "Bearer",
	}
	response.Success(c, http.StatusOK, tokenData)
}

func (as *AuthHandler) Logout(c *gin.Context) {
	accessToken, err := auth.GetBearerToken(c.Request.Header)
	if err != nil {
		response.Error(c, apperror.NewError("Missing Authorization header", apperror.ErrCodeBadRequest))
		return
	}

	refreshToken, err := c.Cookie(dto.RefreshTokenCookieName)
	if err != nil {
		response.Error(c, err)
		return
	}
	err = as.service.Logout(c, accessToken, refreshToken)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.SetCookie(
		dto.RefreshTokenCookieName, "",
		-1,
		"/",
		"",
		as.env == "production",
		true,
	)
	response.Success(c, http.StatusOK, nil)

}
func (as *AuthHandler) RefreshToken(c *gin.Context) {
	rawRefreshToken, err := c.Cookie(dto.RefreshTokenCookieName)
	if err != nil {
		response.Error(c, err)
		return
	}

	tokenInfo, err := as.service.RefreshToken(c.Request.Context(), rawRefreshToken)
	if err != nil {
		response.Error(c, err)
		return
	}
	maxAge := int(time.Until(tokenInfo.RefreshTokenExpiresAt).Seconds())

	c.SetCookie(
		dto.RefreshTokenCookieName, tokenInfo.RefreshToken,
		maxAge,
		"/",
		"",
		as.env == "production",
		true,
	)

	tokenData := dto.LoginResponse{
		AccessToken: tokenInfo.AccessToken,
		TokenType:   "Bearer",
	}
	response.Success(c, http.StatusOK, tokenData)
}

func (as *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validation.HandleValidationError(err))
		return
	}
	err := as.service.ForgotPassword(c.Request.Context(), req.Email)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, gin.H{
		"message": "A password reset link has been sent.",
	})

}
func (as *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validation.HandleValidationError(err))
		return
	}

	if err := as.service.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, gin.H{
		"message": "Password reset successful.",
	})
}
