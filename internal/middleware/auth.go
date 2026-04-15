package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/seminhnva/gin-layered-architecture/internal/auth"
	"github.com/seminhnva/gin-layered-architecture/internal/common/apperror"
	"github.com/seminhnva/gin-layered-architecture/internal/constants"
	"github.com/seminhnva/gin-layered-architecture/pkg/cache"
)

func Auth(jwtService auth.JWT, cache cache.RedisCacheService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			abortAuthError(c, apperror.ErrCodeUnauthorized, "Missing or invalid authorization header")
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			abortAuthError(c, apperror.ErrCodeUnauthorized, "Missing or invalid authorization header")
			return

		}
		tokenString := parts[1]

		claims, err := jwtService.VerifyAcessToken(tokenString)
		if err != nil {
			abortAuthError(c, apperror.ErrCodeUnauthorized, "Invalid or expired access token")
			return
		}
		blackLstKey := constants.BlackListCachePrefix + claims.TokenID
		exist, err := cache.Exists(blackLstKey)
		if err != nil {
			abortAuthError(c, apperror.ErrCodeInternal, "Failed to validate access token")
			return
		}
		if exist {
			abortAuthError(c, apperror.ErrCodeUnauthorized, "Access token has been revoked")
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_roles", claims.Roles)
		c.Next()
	}
}

func abortAuthError(c *gin.Context, code apperror.ErrorCode, message string) {
	c.AbortWithStatusJSON(apperror.HttpStatusCodeFromCode(code), gin.H{
		"message":    message,
		"request_id": GetRequestID(c),
		"error": gin.H{
			"code": code,
		},
	})
}
