package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*Client)
)

func getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}
	return ip
}

func getRateLimiter(ip string) *rate.Limiter {
	mu.Lock()

	defer mu.Unlock()

	client, exist := clients[ip]
	if !exist {
		limiter := rate.NewLimiter(5, 10)
		newClient := &Client{
			limiter,
			time.Now(),
		}
		clients[ip] = newClient
		return limiter
	}
	client.lastSeen = time.Now()
	return client.limiter
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

func RateLimiter(rateLimiterLogger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c)
		limiter := getRateLimiter(ip)

		if !limiter.Allow() {
			if ShouldlogRateLimit(ip) {
				rateLimiterLogger.Warn().
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Str("client_ip", c.ClientIP()).
					Str("user_agent", c.Request.UserAgent()).
					Msg("ratelimiter exceed")
			}

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"errror": "Too many request",
			})
			return
		}
		c.Next()
	}
}

var rateLimitLogCache = sync.Map{}

const rateLimitingLogTTL = 10 * time.Second

func ShouldlogRateLimit(ip string) bool {
	now := time.Now()
	if val, exist := rateLimitLogCache.Load(ip); exist {
		if t, ok := val.(time.Time); ok && now.Sub(t) < rateLimitingLogTTL {
			return false
		}
	}
	rateLimitLogCache.Store(ip, now)
	return true
}
