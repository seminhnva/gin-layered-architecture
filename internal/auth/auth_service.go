package auth

type TokenPayload struct {
	UserID   string
	Email    string
	UserName string
	Name     string
}
type Hasher interface {
	HashPassword(plain_password string) (string, error)
	CheckPasswordHash(plain_password, hashed string) (bool, error)
}

type JWT interface {
	GenerateAccessToken(payload TokenPayload) (string, error)
	// 	GenerateRefreshToken(user sqlc.User) (RefreshToken, error)
	// 	ParseToken(tokenString string) (*Claims, error)
	// 	DescrypAccessTokenPayload(tokenString string) (*EncryptedPayload, error)
	// 	StoreRefreshToken(token RefreshToken) error
	// 	ValidateRefreshToken(token string) (RefreshToken, error)
	// 	RevokeRefreshToken(tokenStr string) error
}
