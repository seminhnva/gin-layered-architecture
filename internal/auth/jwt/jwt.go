package jwtService

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/seminhnva/gin-layered-architecture/internal/auth"
	"github.com/seminhnva/gin-layered-architecture/internal/common/apperror"
)

const (
	Issuer = "server"
)

type JWTService struct {
	secret          string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	cache           any
}

type claims struct {
	UserID   string   `json:"sub"`
	Email    string   `json:"email"`
	UserName string   `json:"user_name"`
	Name     string   `json:"name"`
	Roles    []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

func NewJWTService(jwtSecret string, accessTokenTTL, refreshTokenTTL time.Duration, cache any) auth.JWT {
	return &JWTService{
		secret:          jwtSecret,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		cache:           cache,
	}
}

func NewJWTSerivice(jwtSecret string, accessTokenTTL, refreshTokenTTL time.Duration, cache any) auth.JWT {
	return NewJWTService(jwtSecret, accessTokenTTL, refreshTokenTTL, cache)
}

func (js *JWTService) GenerateAccessToken(payload auth.TokenPayload) (string, error) {
	now := time.Now().UTC()

	tokenClaims := &claims{
		UserID:   payload.UserID,
		Email:    payload.Email,
		UserName: payload.UserName,
		Name:     payload.Name,
		Roles:    payload.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    Issuer,
			Subject:   payload.UserID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(js.accessTokenTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	return token.SignedString([]byte(js.secret))
}

func (js *JWTService) VerifyToken(tokenString string) (*auth.TokenClaims, error) {
	tokenClaims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, tokenClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(js.secret), nil
	})
	if err != nil || !token.Valid {
		return nil, apperror.NewError("Invalid token", apperror.ErrCodeUnauthorized)
	}

	verifiedClaims := &auth.TokenClaims{
		UserID:   tokenClaims.UserID,
		Email:    tokenClaims.Email,
		UserName: tokenClaims.UserName,
		Name:     tokenClaims.Name,
		Roles:    tokenClaims.Roles,
		TokenID:  tokenClaims.ID,
		Issuer:   tokenClaims.Issuer,
	}

	if tokenClaims.IssuedAt != nil {
		verifiedClaims.IssuedAt = tokenClaims.IssuedAt.Time
	}
	if tokenClaims.ExpiresAt != nil {
		verifiedClaims.ExpiresAt = tokenClaims.ExpiresAt.Time
	}

	return verifiedClaims, nil
}

func (js *JWTService) GenerateRefreshToken(userID uuid.UUID) (rawToken string, stored auth.RefreshToken, err error) {
	tokenBytes := make([]byte, 32)
	if _, err = rand.Read(tokenBytes); err != nil {
		return
	}

	rawToken = base64.URLEncoding.EncodeToString(tokenBytes)

	now := time.Now().UTC()
	stored = auth.RefreshToken{
		TokenHash: HashToken(rawToken),
		UserID:    userID,
		IssuedAt:  now,
		ExpiresAt: now.Add(js.refreshTokenTTL),
		Revoked:   false,
	}
	return
}

func HashToken(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])
	return tokenHash
}
