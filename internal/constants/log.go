package constants

type logFilePath string

const (
	AppLogFilePath      logFilePath = "app.log"
	HttpLogFilePath     logFilePath = "http.log"
	RecoveryLogFilePath logFilePath = "recovery.log"
	RateLimiterFilePath logFilePath = "rate_limiter.log"
)
