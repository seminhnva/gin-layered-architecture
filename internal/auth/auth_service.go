package auth

import "time"

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

type Hasher interface {
	HashPassword(plainPassword string) (string, error)
	CheckPasswordHash(plainPassword, hashed string) (bool, error)
}

type JWT interface {
	GenerateAccessToken(payload TokenPayload) (string, error)
	VerifyToken(tokenString string) (*TokenClaims, error)
	// 	DescrypAccessTokenPayload(tokenString string) (*EncryptedPayload, error)
	// 	StoreRefreshToken(token RefreshToken) error
	// 	ValidateRefreshToken(token string) (RefreshToken, error)
	// 	RevokeRefreshToken(tokenStr string) error
}
