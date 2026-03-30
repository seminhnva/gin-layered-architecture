package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/seminhnva/gin-layered-architecture/internal/config"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type Level int

const (
	levelDebug Level = iota
	levelInfo
	levelWarn
	levelError
)

var (
	mu           sync.RWMutex
	currentLevel = levelInfo
	std          = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC)
)

func Setup(cfg *config.Config) error {
	parsedLevel, err := parseLevel(cfg.LogLevel)
	if err != nil {
		return err
	}

	logDir := filepath.Dir(cfg.LogFilePath)
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return fmt.Errorf("create log directory: %w", err)
	}

	rotatingWriter := &lumberjack.Logger{
		Filename:   cfg.LogFilePath,
		MaxSize:    cfg.LogMaxSizeMB,
		MaxBackups: cfg.LogMaxBackups,
		MaxAge:     cfg.LogMaxAgeDays,
		Compress:   cfg.LogCompress,
	}

	flags := log.Ldate | log.Ltime | log.LUTC
	writer := io.MultiWriter(os.Stdout, rotatingWriter)

	log.SetFlags(flags)
	log.SetOutput(writer)

	mu.Lock()
	currentLevel = parsedLevel
	std = log.New(writer, "", flags)
	mu.Unlock()

	return nil
}

func Debugf(format string, args ...any) {
	logf(levelDebug, "DEBUG", format, args...)
}

func Infof(format string, args ...any) {
	logf(levelInfo, "INFO", format, args...)
}

func Warnf(format string, args ...any) {
	logf(levelWarn, "WARN", format, args...)
}

func Errorf(format string, args ...any) {
	logf(levelError, "ERROR", format, args...)
}

func Fatalf(format string, args ...any) {
	getLogger().Fatalf("[FATAL] %s", fmt.Sprintf(format, args...))
}

func logf(level Level, label, format string, args ...any) {
	if !enabled(level) {
		return
	}

	getLogger().Printf("[%s] %s", label, fmt.Sprintf(format, args...))
}

func enabled(level Level) bool {
	mu.RLock()
	defer mu.RUnlock()

	return level >= currentLevel
}

func getLogger() *log.Logger {
	mu.RLock()
	defer mu.RUnlock()

	return std
}

func parseLevel(value string) (Level, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return levelDebug, nil
	case "info", "":
		return levelInfo, nil
	case "warn", "warning":
		return levelWarn, nil
	case "error":
		return levelError, nil
	default:
		return levelInfo, fmt.Errorf("LOG_LEVEL must be one of: debug, info, warn, error")
	}
}
