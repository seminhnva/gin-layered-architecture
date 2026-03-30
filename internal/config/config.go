package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv             string
	Port               string
	ServerAddress      string
	LogLevel           string
	LogFilePath        string
	LogMaxSizeMB       int
	LogMaxBackups      int
	LogMaxAgeDays      int
	LogCompress        bool
	JWTSecret          string
	CORSAllowedOrigins []string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	ReadHeaderTimeout  time.Duration
	ShutdownTimeout    time.Duration
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	serverAddress := getEnv("SERVER_ADDRESS", ":"+port)
	jwtSecret, err := mustGetEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}
	logMaxSizeMB, err := getIntEnv("LOG_MAX_SIZE_MB", 10)
	if err != nil {
		return nil, err
	}
	logMaxBackups, err := getIntEnv("LOG_MAX_BACKUPS", 5)
	if err != nil {
		return nil, err
	}
	logMaxAgeDays, err := getIntEnv("LOG_MAX_AGE_DAYS", 30)
	if err != nil {
		return nil, err
	}
	logCompress, err := getBoolEnv("LOG_COMPRESS", true)
	if err != nil {
		return nil, err
	}

	readTimeout, err := getDurationEnv("READ_TIMEOUT", 5*time.Second)
	if err != nil {
		return nil, err
	}

	writeTimeout, err := getDurationEnv("WRITE_TIMEOUT", 10*time.Second)
	if err != nil {
		return nil, err
	}

	idleTimeout, err := getDurationEnv("IDLE_TIMEOUT", 60*time.Second)
	if err != nil {
		return nil, err
	}

	readHeaderTimeout, err := getDurationEnv("READ_HEADER_TIMEOUT", 2*time.Second)
	if err != nil {
		return nil, err
	}

	shutdownTimeout, err := getDurationEnv("SHUTDOWN_TIMEOUT", 10*time.Second)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		Port:               port,
		ServerAddress:      serverAddress,
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		LogFilePath:        getEnv("LOG_FILE_PATH", "logs/app.log"),
		LogMaxSizeMB:       logMaxSizeMB,
		LogMaxBackups:      logMaxBackups,
		LogMaxAgeDays:      logMaxAgeDays,
		LogCompress:        logCompress,
		JWTSecret:          jwtSecret,
		CORSAllowedOrigins: getSliceEnv("CORS_ALLOWED_ORIGINS"),
		ReadTimeout:        readTimeout,
		WriteTimeout:       writeTimeout,
		IdleTimeout:        idleTimeout,
		ReadHeaderTimeout:  readHeaderTimeout,
		ShutdownTimeout:    shutdownTimeout,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.ServerAddress == "" {
		return errors.New("SERVER_ADDRESS cannot be empty")
	}

	if c.LogFilePath == "" {
		return errors.New("LOG_FILE_PATH cannot be empty")
	}

	if c.LogMaxSizeMB <= 0 {
		return errors.New("LOG_MAX_SIZE_MB must be greater than 0")
	}

	if c.LogMaxBackups < 0 {
		return errors.New("LOG_MAX_BACKUPS must be greater than or equal to 0")
	}

	if c.LogMaxAgeDays < 0 {
		return errors.New("LOG_MAX_AGE_DAYS must be greater than or equal to 0")
	}

	for _, origin := range c.CORSAllowedOrigins {
		if origin == "*" {
			return errors.New("CORS_ALLOWED_ORIGINS cannot contain wildcard '*'; use explicit origins")
		}

		parsedOrigin, err := url.Parse(origin)
		if err != nil || parsedOrigin.Scheme == "" || parsedOrigin.Host == "" {
			return fmt.Errorf("invalid CORS origin: %s", origin)
		}
	}

	if _, err := strconv.Atoi(c.Port); err != nil {
		return fmt.Errorf("PORT must be numeric: %w", err)
	}

	if c.ReadTimeout <= 0 {
		return errors.New("READ_TIMEOUT must be greater than 0")
	}

	if c.WriteTimeout <= 0 {
		return errors.New("WRITE_TIMEOUT must be greater than 0")
	}

	if c.IdleTimeout <= 0 {
		return errors.New("IDLE_TIMEOUT must be greater than 0")
	}

	if c.ReadHeaderTimeout <= 0 {
		return errors.New("READ_HEADER_TIMEOUT must be greater than 0")
	}

	if c.ShutdownTimeout <= 0 {
		return errors.New("SHUTDOWN_TIMEOUT must be greater than 0")
	}

	return nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func mustGetEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return value, nil
}

func getSliceEnv(key string) []string {
	value := os.Getenv(key)
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		values = append(values, trimmed)
	}

	return values
}

func getIntEnv(key string, fallback int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be numeric: %w", key, err)
	}

	return parsed, nil
}

func getBoolEnv(key string, fallback bool) (bool, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean: %w", key, err)
	}

	return parsed, nil
}

func getDurationEnv(key string, fallback time.Duration) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration: %w", key, err)
	}

	return parsed, nil
}
