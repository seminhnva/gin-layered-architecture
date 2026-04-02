package bootstrap

import (
	"github.com/seminhnva/gin-layered-architecture/internal/config"
	"github.com/seminhnva/gin-layered-architecture/pkg/logger"
)

func NewLoggerOptions(cfg config.LogConfig) logger.Options {
	return logger.Options{
		Dir:        cfg.LogFilePath,
		Level:      cfg.LogLevel,
		MaxSizeMB:  cfg.LogMaxSizeMB,
		MaxBackups: cfg.LogMaxBackups,
		MaxAgeDays: cfg.LogMaxAgeDays,
		Compress:   cfg.LogCompress,
		LocalTime:  cfg.LocalTime,
	}
}
