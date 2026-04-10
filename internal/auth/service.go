package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TokenPayload struct {
	UserID   string
	Email    string
	UserName string
	Name     string
	Roles    []string
}

type TokenClaims struct {
	UserID    string
	Email     string
	UserName  string
	Name      string
	Roles     []string
	TokenID   string
	Issuer    string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type RefreshToken struct {
	TokenHash string
	UserID    uuid.UUID
	IssuedAt  time.Time
	ExpiresAt time.Time
	Revoked   bool
	RevokedAt *time.Time
}

type Hasher interface {
	HashPassword(plainPassword string) (string, error)
	CheckPasswordHash(plainPassword, hashed string) (bool, error)
}

type JWT interface {
	GenerateAccessToken(payload TokenPayload) (string, error)
	VerifyToken(tokenString string) (*TokenClaims, error)
	GenerateRefreshToken(userID uuid.UUID) (rawToken string, stored RefreshToken, err error)
	// 	ValidateRefreshToken(token string) (RefreshToken, error)
	// 	RevokeRefreshToken(tokenStr string) error
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	return parts[1], nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "ApiKey" {
		return "", errors.New("invalid authorization header")
	}
	return parts[1], nil

}
