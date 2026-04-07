package jwtService

import (
	"encoding/json"
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

type Claims struct {
	Data []byte `json:"data"`
	jwt.RegisteredClaims
}

func NewJWTSerivice(jwtSecret string, accessTokenTTL, refreshTokenTTL time.Duration, cache any) auth.JWT {
	return &JWTService{
		secret:          jwtSecret,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		cache:           cache,
	}
}

func (js *JWTService) GenerateAccessToken(payload auth.TokenPayload) (string, error) {

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	claims := &Claims{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(js.accessTokenTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(js.secret))
}

func (js *JWTService) VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// prevent alg:none attack
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(js.secret), nil
	})
	if err != nil || !token.Valid {
		return nil, apperror.NewError("Invalid token", apperror.ErrCodeUnauthorized)
	}
	return claims, nil
}
