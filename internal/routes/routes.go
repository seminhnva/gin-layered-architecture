package routes

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/seminhnva/gin-layered-architecture/internal/auth"
	"github.com/seminhnva/gin-layered-architecture/internal/middleware"
	"github.com/seminhnva/gin-layered-architecture/pkg/cache"
)

type Route interface {
	Register(rg *gin.RouterGroup)
}

func SetUpRouter(corsAllowedOrigins []string, r *gin.Engine, jwtService auth.JWT, cache cache.RedisCacheService, httpLogger, recoveryLogger, rateLimiterLogger *zerolog.Logger, routes ...Route) {
	r.Use(
		middleware.Recover(recoveryLogger),
		middleware.RequestID(),
		middleware.Trace(),
		gzip.Gzip(gzip.DefaultCompression),
		middleware.RateLimiter(rateLimiterLogger, middleware.DefaultRateLimitPolicy),
		middleware.Logger(httpLogger),
		middleware.CORS(corsAllowedOrigins),
	)
	api := r.Group("/api")
	apiv1 := r.Group("/api/v1")
	protected := apiv1.Group("")
	protected.Use(
		middleware.Auth(jwtService, cache),
	)
	for _, route := range routes {
		switch route.(type) {
		case *AuthRoute:
			route.Register(api)
		default:
			route.Register(protected)
		}

	}
}
