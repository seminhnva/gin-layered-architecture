package service

import (
	"context"

	"github.com/seminhnva/gin-layered-architecture/internal/dto"
)

type AuthSerivce interface {
	Login(ctx context.Context, params dto.LoginRequest) (dto.TokenInfo, error)
	Logout(ctx context.Context, accessToken, rawRefreshToken string) error
	RefreshToken(ctx context.Context) error
	ForgotPassword(ctx context.Context) error
	ResetPassword(ctx context.Context) error
}
