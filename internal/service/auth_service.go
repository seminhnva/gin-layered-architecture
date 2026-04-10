package service

import (
	"context"
	"time"

	"github.com/seminhnva/gin-layered-architecture/internal/auth"
	jwtService "github.com/seminhnva/gin-layered-architecture/internal/auth/jwt"
	"github.com/seminhnva/gin-layered-architecture/internal/common/apperror"
	"github.com/seminhnva/gin-layered-architecture/internal/constants"
	"github.com/seminhnva/gin-layered-architecture/internal/db/sqlc"
	"github.com/seminhnva/gin-layered-architecture/internal/dto"
	"github.com/seminhnva/gin-layered-architecture/internal/repository"
	"github.com/seminhnva/gin-layered-architecture/internal/utils"
	"github.com/seminhnva/gin-layered-architecture/pkg/cache"
)

type authService struct {
	repo            repository.AuthRepository
	passwordService auth.Hasher
	jwtService      auth.JWT
	cache           cache.RedisCacheService
}

func NewAuthService(repo repository.AuthRepository, passwordService auth.Hasher, jwtService auth.JWT, cache cache.RedisCacheService) AuthSerivce {
	return &authService{
		repo:            repo,
		passwordService: passwordService,
		jwtService:      jwtService,
		cache:           cache,
	}
}

func (as *authService) Login(ctx context.Context, params dto.LoginRequest) (dto.TokenInfo, error) {
	params.Email = utils.NormalizeString(params.Email)

	user, err := as.repo.GetByEmail(ctx, sqlc.GetByEmailParams{
		Email: params.Email,
	})
	if err != nil {
		return dto.TokenInfo{}, err
	}

	matchPw, err := as.passwordService.CheckPasswordHash(params.Password, user.PasswordHash)
	if err != nil {
		return dto.TokenInfo{}, err
	}
	if !matchPw {
		return dto.TokenInfo{}, apperror.NewError("Invalid email or password", apperror.ErrCodeUnauthorized)
	}

	accessToken, err := as.jwtService.GenerateAccessToken(auth.TokenPayload{
		UserID:   user.UserID.String(),
		Email:    user.UserName,
		Name:     user.Name,
		UserName: user.UserName,
	})
	if err != nil {
		return dto.TokenInfo{}, err
	}

	refreshToken, storeToken, err := as.jwtService.GenerateRefreshToken(user.UserID)
	if err != nil {
		return dto.TokenInfo{}, err
	}

	tokenInfo := dto.TokenInfo{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: storeToken.ExpiresAt,
	}

	cacheKey := constants.RefreshTokenCachePrefix + jwtService.HashToken(tokenInfo.RefreshToken)
	if err := as.cache.Set(cacheKey, user.UserID.String(), time.Until(storeToken.ExpiresAt)); err != nil {
		return dto.TokenInfo{}, err
	}
	return tokenInfo, nil
}
func (as *authService) Logout(ctx context.Context, accessToken, rawRefreshToken string) error {
	_, err := as.jwtService.VerifyToken(accessToken)
	if err != nil {
		return apperror.NewError("Invalid access token", apperror.ErrCodeUnauthorized)
	}

	cacheKey := constants.RefreshTokenCachePrefix + jwtService.HashToken(rawRefreshToken)
	as.cache.Del(cacheKey)

	return nil
}
func (as *authService) RefreshToken(ctx context.Context) error {
	return nil
}
func (as *authService) ForgotPassword(ctx context.Context) error {
	return nil
}
func (as *authService) ResetPassword(ctx context.Context) error {
	return nil
}
