package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/seminhnva/gin-layered-architecture/internal/auth"
	"github.com/seminhnva/gin-layered-architecture/internal/common/apperror"
	"github.com/seminhnva/gin-layered-architecture/internal/common/domainerror"
	"github.com/seminhnva/gin-layered-architecture/internal/constants"
	"github.com/seminhnva/gin-layered-architecture/internal/db/sqlc"
	"github.com/seminhnva/gin-layered-architecture/internal/dto"
	"github.com/seminhnva/gin-layered-architecture/internal/repository"
	"github.com/seminhnva/gin-layered-architecture/internal/utils"
	"github.com/seminhnva/gin-layered-architecture/pkg/cache"
)

type authService struct {
	userRepo        repository.UserRepository
	passwordService auth.Hasher
	jwtService      auth.JWT
	cache           cache.RedisCacheService
}

func NewAuthService(userRepo repository.UserRepository, passwordService auth.Hasher, jwtService auth.JWT, cache cache.RedisCacheService) AuthSerivce {
	return &authService{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
		cache:           cache,
	}
}

func (as *authService) Login(ctx context.Context, params dto.LoginRequest) (dto.TokenInfo, error) {
	params.Email = utils.NormalizeString(params.Email)

	user, err := as.userRepo.GetByEmail(ctx, sqlc.GetByEmailParams{
		Email: params.Email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.TokenInfo{}, apperror.NewError("Invalid email or password", apperror.ErrCodeUnauthorized)
		}
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
		Email:    user.Email,
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

	if err := as.storeRefreshToken(storeToken); err != nil {
		return dto.TokenInfo{}, err

	}
	return tokenInfo, nil
}

func (as *authService) Logout(ctx context.Context, accessToken, rawRefreshToken string) error {
	claims, err := as.jwtService.VerifyAcessToken(accessToken)
	if err != nil {
		return apperror.NewError("Invalid access token", apperror.ErrCodeUnauthorized)
	}
	blackLstKey := constants.BlackListCachePrefix + claims.TokenID
	if err := as.cache.Set(blackLstKey, "revoke", time.Until(claims.ExpiresAt)); err != nil {
		return apperror.NewError("Internal server error", apperror.ErrCodeInternal)

	}

	if err := as.deleteRefreshToken(rawRefreshToken); err != nil {
		return err
	}
	return nil
}

func (as *authService) RefreshToken(ctx context.Context, rawRefreshToken string) (dto.TokenInfo, error) {
	storedRefreshToken, err := as.verifyRefreshToken(rawRefreshToken)
	if err != nil {
		return dto.TokenInfo{}, err
	}

	user, err := as.userRepo.FindByUUID(ctx, storedRefreshToken.UserID)
	if err != nil {
		if errors.Is(err, domainerror.ErrUserNotFound) {
			return dto.TokenInfo{}, apperror.NewError("User not found", apperror.ErrCodeNotFound)
		}
		return dto.TokenInfo{}, apperror.WrapError(err, "Fail to find user", apperror.ErrCodeInternal)
	}
	accessToken, err := as.jwtService.GenerateAccessToken(auth.TokenPayload{
		UserID:   user.UserID.String(),
		Email:    user.Email,
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

	if err := as.storeRefreshToken(storeToken); err != nil {
		return dto.TokenInfo{}, err

	}

	if err := as.deleteRefreshToken(rawRefreshToken); err != nil {
		return dto.TokenInfo{}, err
	}
	return tokenInfo, nil
}
func (as *authService) ForgotPassword(ctx context.Context, email string) error {
	email = utils.NormalizeString(email)

	user, err := as.userRepo.GetByEmail(ctx, sqlc.GetByEmailParams{
		Email: email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	randomCode, err := utils.GenerateRandomString(16)
	if err != nil {
		return apperror.NewError("Internal server error", apperror.ErrCodeInternal)
	}
	resetLink := fmt.Sprintf("https://yourdomain.com/reset-password?code=%s", randomCode)
	log.Println(resetLink)
	cacheKey := constants.ResetPasswordPrefix + randomCode
	resetPasswordInfo := auth.ResetPassword{
		UserID:   user.UserID,
		IssuedAt: time.Now(),
	}
	err = as.cache.Set(cacheKey, resetPasswordInfo, 5*time.Minute)

	if err != nil {
		return apperror.NewError("Internal server error", apperror.ErrCodeInternal)
	}
	return nil
}
func (as *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	cacheKey := constants.ResetPasswordPrefix + token
	var resetPasswordInfo auth.ResetPassword
	if err := as.cache.Get(cacheKey, &resetPasswordInfo); err != nil {
		return apperror.NewError("Reset code is wrong", apperror.ErrCodeBadRequest)
	}
	hashPassword, err := as.passwordService.HashPassword(newPassword)
	if err != nil {
		return apperror.NewError("Internal server error", apperror.ErrCodeInternal)

	}
	_, err = as.userRepo.ResetPassword(ctx, sqlc.ResetPasswordParams{
		UserID:       resetPasswordInfo.UserID,
		PasswordHash: hashPassword,
	})
	if err != nil {
		return apperror.WrapError(err, "Fail to update password", apperror.ErrCodeInternal)

	}
	if err = as.cache.Del(cacheKey); err != nil {
		return apperror.NewError("Internal server error", apperror.ErrCodeInternal)
	}

	return nil
}

func (as *authService) storeRefreshToken(token auth.RefreshToken) error {
	cacheKey := constants.RefreshTokenCachePrefix + token.TokenHash
	return as.cache.Set(cacheKey, token, time.Until(token.ExpiresAt))
}

func (as *authService) deleteRefreshToken(rawRefreshToken string) error {
	cacheKey := constants.RefreshTokenCachePrefix + auth.HashToken(rawRefreshToken)
	return as.cache.Del(cacheKey)
}

func (as *authService) verifyRefreshToken(rawRefreshToken string) (auth.RefreshToken, error) {
	cacheKey := constants.RefreshTokenCachePrefix + auth.HashToken(rawRefreshToken)
	var storedRefreshToken auth.RefreshToken
	if err := as.cache.Get(cacheKey, &storedRefreshToken); err != nil {
		return auth.RefreshToken{}, apperror.NewError("Refresh token is invalid or expired", apperror.ErrCodeUnauthorized)
	}
	if storedRefreshToken.ExpiresAt.Before(time.Now()) {
		return auth.RefreshToken{}, apperror.NewError("Refresh token is invalid or expired", apperror.ErrCodeUnauthorized)
	}
	return storedRefreshToken, nil
}
