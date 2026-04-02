package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/seminhnva/gin-layered-architecture/internal/utils"
)

type HTTPServerConfig struct {
	ServerAddress     string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	ShutdownTimeout   time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMod   string
}

type JWTConfig struct {
	Secret          string
	EncryptKey      string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type LogConfig struct {
	LogLevel      string
	LogFilePath   string
	LogMaxSizeMB  int
	LogMaxBackups int
	LogMaxAgeDays int
	LogCompress   bool
	LocalTime     bool
}

type Config struct {
	AppEnv             string
	Port               string
	CORSAllowedOrigins []string
	JWT                JWTConfig
	DB                 DatabaseConfig
	HTTPServer         HTTPServerConfig
	Logger             LogConfig
}

func NewConfig() (*Config, error) {
	if err := loadEnv(); err != nil {
		return nil, fmt.Errorf("load env file: %w", err)
	}

	appEnv, port, corsAllowedOrigins := loadAppConfig()

	httpConfig, err := loadHTTPServerConfig(port)
	if err != nil {
		return nil, fmt.Errorf("load http config: %w", err)
	}

	jwtConfig, err := loadJWTConfig()
	if err != nil {
		return nil, fmt.Errorf("load jwt config: %w", err)
	}

	databaseConfig, err := loadDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("load database config: %w", err)
	}

	logConfig, err := loadLogConfig()
	if err != nil {
		return nil, fmt.Errorf("load logger config: %w", err)
	}

	cfg := &Config{
		AppEnv:             appEnv,
		Port:               port,
		CORSAllowedOrigins: corsAllowedOrigins,
		JWT:                jwtConfig,
		DB:                 databaseConfig,
		HTTPServer:         httpConfig,
		Logger:             logConfig,
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.HTTPServer.ServerAddress == "" {
		return errors.New("SERVER_ADDRESS cannot be empty")
	}

	if c.Logger.LogFilePath == "" {
		return errors.New("LOG_FILE_PATH cannot be empty")
	}

	if c.Logger.LogMaxSizeMB <= 0 {
		return errors.New("LOG_MAX_SIZE_MB must be greater than 0")
	}

	if c.Logger.LogMaxBackups < 0 {
		return errors.New("LOG_MAX_BACKUPS must be greater than or equal to 0")
	}

	if c.Logger.LogMaxAgeDays < 0 {
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

	if c.HTTPServer.ReadTimeout <= 0 {
		return errors.New("READ_TIMEOUT must be greater than 0")
	}

	if c.HTTPServer.WriteTimeout <= 0 {
		return errors.New("WRITE_TIMEOUT must be greater than 0")
	}

	if c.HTTPServer.IdleTimeout <= 0 {
		return errors.New("IDLE_TIMEOUT must be greater than 0")
	}

	if c.HTTPServer.ReadHeaderTimeout <= 0 {
		return errors.New("READ_HEADER_TIMEOUT must be greater than 0")
	}

	if c.HTTPServer.ShutdownTimeout <= 0 {
		return errors.New("SHUTDOWN_TIMEOUT must be greater than 0")
	}

	return nil
}

func loadEnv() error {
	rootDir, err := utils.GetWorkingDir()
	if err != nil {
		return fmt.Errorf("load root dir: %w", err)
	}

	envPath := filepath.Join(rootDir, ".env")
	if _, err := os.Stat(envPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("stat env file: %w", err)
	}

	if err := godotenv.Load(envPath); err != nil {
		return fmt.Errorf("load env file %s: %w", envPath, err)
	}

	return nil
}

func loadAppConfig() (string, string, []string) {
	appEnv := utils.GetEnv("APP_ENV", "development")
	port := utils.GetEnv("PORT", "8080")
	corsAllowedOrigins := utils.GetSliceEnv("CORS_ALLOWED_ORIGINS")

	return appEnv, port, corsAllowedOrigins
}

func loadHTTPServerConfig(port string) (HTTPServerConfig, error) {
	readTimeout, err := utils.GetDurationEnv("READ_TIMEOUT", 5*time.Second)
	if err != nil {
		return HTTPServerConfig{}, err
	}

	writeTimeout, err := utils.GetDurationEnv("WRITE_TIMEOUT", 10*time.Second)
	if err != nil {
		return HTTPServerConfig{}, err
	}

	idleTimeout, err := utils.GetDurationEnv("IDLE_TIMEOUT", 60*time.Second)
	if err != nil {
		return HTTPServerConfig{}, err
	}

	readHeaderTimeout, err := utils.GetDurationEnv("READ_HEADER_TIMEOUT", 2*time.Second)
	if err != nil {
		return HTTPServerConfig{}, err
	}

	shutdownTimeout, err := utils.GetDurationEnv("SHUTDOWN_TIMEOUT", 10*time.Second)
	if err != nil {
		return HTTPServerConfig{}, err
	}

	return HTTPServerConfig{
		ServerAddress:     utils.GetEnv("SERVER_ADDRESS", ":"+port),
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		ShutdownTimeout:   shutdownTimeout,
	}, nil
}

func loadJWTConfig() (JWTConfig, error) {
	secret, err := utils.MustGetEnv("JWT_SECRET")
	if err != nil {
		return JWTConfig{}, err
	}

	encryptKey, err := utils.MustGetEnv("JWT_ENCRYPT_KEY")
	if err != nil {
		return JWTConfig{}, err
	}

	accessTokenTTL, err := utils.GetDurationEnv("ACCESS_TOKEN_TTL", 15*time.Minute)
	if err != nil {
		return JWTConfig{}, err
	}

	refreshTokenTTL, err := utils.GetDurationEnv("REFRESH_TOKEN_TTL", 168*time.Hour)
	if err != nil {
		return JWTConfig{}, err
	}

	return JWTConfig{
		Secret:          secret,
		EncryptKey:      encryptKey,
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
	}, nil
}

func loadDatabaseConfig() (DatabaseConfig, error) {
	host, err := utils.MustGetEnv("DB_HOST")
	if err != nil {
		return DatabaseConfig{}, err
	}

	user, err := utils.MustGetEnv("DB_USER")
	if err != nil {
		return DatabaseConfig{}, err
	}

	password, err := utils.MustGetEnv("DB_PASSWORD")
	if err != nil {
		return DatabaseConfig{}, err
	}

	dbName, err := utils.MustGetEnv("DB_NAME")
	if err != nil {
		return DatabaseConfig{}, err
	}

	sslMode, err := utils.MustGetEnv("DB_SSLMODE")
	if err != nil {
		return DatabaseConfig{}, err
	}

	return DatabaseConfig{
		Host:     host,
		Port:     utils.GetEnv("DB_PORT", "5432"),
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMod:   sslMode,
	}, nil
}

func loadLogConfig() (LogConfig, error) {
	logMaxSizeMB, err := utils.GetEnvInt("LOG_MAX_SIZE_MB", 10)
	if err != nil {
		return LogConfig{}, err
	}

	logMaxBackups, err := utils.GetEnvInt("LOG_MAX_BACKUPS", 5)
	if err != nil {
		return LogConfig{}, err
	}

	logMaxAgeDays, err := utils.GetEnvInt("LOG_MAX_AGE_DAYS", 30)
	if err != nil {
		return LogConfig{}, err
	}

	logCompress, err := utils.GetEnvBool("LOG_COMPRESS", true)
	if err != nil {
		return LogConfig{}, err
	}

	localTime, err := utils.GetEnvBool("LOG_LOCAL_TIME", false)
	if err != nil {
		return LogConfig{}, err
	}

	return LogConfig{
		LogLevel:      utils.GetEnv("LOG_LEVEL", "info"),
		LogFilePath:   utils.GetEnv("LOG_FILE_PATH", "logs/"),
		LogMaxSizeMB:  logMaxSizeMB,
		LogMaxBackups: logMaxBackups,
		LogMaxAgeDays: logMaxAgeDays,
		LogCompress:   logCompress,
		LocalTime:     localTime,
	}, nil
}
