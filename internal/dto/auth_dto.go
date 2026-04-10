package dto

import "time"

const RefreshTokenCookieName = "refresh_token"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
type ForgotPasswordRequest struct {
	Email *string `json:"email" binding:"omitempty,email"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type TokenInfo struct {
	AccessToken           string
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

type User struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}
