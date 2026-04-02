package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/seminhnva/gin-layered-architecture/internal/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(logFilePath string, logConfig config.LogConfig) (*zerolog.Logger, error) {
	logger, err := NewFileLogger(
		filepath.Join(logConfig.LogFilePath, logFilePath),
		logConfig.LogLevel,
		logConfig.LogMaxSizeMB,
		logConfig.LogMaxBackups,
		logConfig.LogMaxAgeDays,
		logConfig.LogCompress,
		logConfig.LocalTime,
	)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	return logger, nil
}

func NewFileLogger(logFilePath, level string, maxSize, maxBackups, maxAge int, compress, localTime bool) (*zerolog.Logger, error) {
	resolvedPath, err := resolvePath(logFilePath)
	if err != nil {
		return nil, err
	}

	if err := ensureParentDir(resolvedPath); err != nil {
		return nil, err
	}

	zerolog.TimeFieldFormat = time.RFC3339
	lv, err := zerolog.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	zerolog.SetGlobalLevel(lv)

	rollingFile := &lumberjack.Logger{
		Filename:   resolvedPath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
		LocalTime:  localTime,
	}
	var write io.Writer
	write = zerolog.MultiLevelWriter(os.Stdout, rollingFile)
	logger := zerolog.New(write).With().Timestamp().Logger()

	return &logger, nil
}
