package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/seminhnva/gin-layered-architecture/internal/common/apperror"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimitPolicy struct {
	Name  string
	Limit rate.Limit
	Burst int
}

var (
	DefaultRateLimitPolicy = RateLimitPolicy{
		Name:  "default",
		Limit: 5,
		Burst: 10,
	}

	LoginRateLimitPolicy = RateLimitPolicy{
		Name:  "login",
		Limit: 1,
		Burst: 5,
	}

	RefreshRateLimitPolicy = RateLimitPolicy{
		Name:  "refresh",
		Limit: 1,
		Burst: 5,
	}

	ForgotPasswordRateLimitPolicy = RateLimitPolicy{
		Name:  "forgot_password",
		Limit: 0.2,
		Burst: 3,
	}
)
var (
	mu      sync.Mutex
	clients = make(map[string]*Client)
)

func RateLimiter(rateLimiterLogger *zerolog.Logger, policy RateLimitPolicy) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c)
		limiter := getRateLimiter(policy, ip)
		clientKey := policy.Name + ":" + ip

		if !limiter.Allow() {
			if ShouldlogRateLimit(clientKey) {
				rateLimiterLogger.Warn().
					Str("policy", policy.Name).
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Str("client_ip", c.ClientIP()).
					Str("user_agent", c.Request.UserAgent()).
					Msg("ratelimiter exceed")
			}
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message":    "Too many requests. Please try again later.",
				"request_id": GetRequestID(c),
				"error": gin.H{
					"code": apperror.ErrCodeTooManyRequest,
				},
			})

			return
		}
		c.Next()
	}
}

func CleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

func getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}
	return ip
}

func getRateLimiter(policy RateLimitPolicy, ip string) *rate.Limiter {
	key := policy.Name + ":" + ip

	mu.Lock()

	defer mu.Unlock()

	client, exist := clients[key]
	if !exist {
		limiter := rate.NewLimiter(policy.Limit, policy.Burst)
		newClient := &Client{
			limiter,
			time.Now(),
		}
		clients[key] = newClient
		return limiter
	}
	client.lastSeen = time.Now()
	return client.limiter
}

var rateLimitLogCache = sync.Map{}

const rateLimitingLogTTL = 10 * time.Second

func ShouldlogRateLimit(key string) bool {
	now := time.Now()
	if val, exist := rateLimitLogCache.Load(key); exist {
		if t, ok := val.(time.Time); ok && now.Sub(t) < rateLimitingLogTTL {
			return false
		}
	}
	rateLimitLogCache.Store(key, now)
	return true
}
