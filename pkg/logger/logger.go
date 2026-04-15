package logger

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Options struct {
	Dir        string
	Level      string
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
	Compress   bool
	LocalTime  bool
}

func InitLogger(fileName string, opts Options) (*zerolog.Logger, error) {
	logger, err := NewFileLogger(
		filepath.Join(opts.Dir, fileName),
		opts.Level,
		opts.MaxSizeMB,
		opts.MaxBackups,
		opts.MaxAgeDays,
		opts.Compress,
		opts.LocalTime,
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
	rollingFile := &lumberjack.Logger{
		Filename:   resolvedPath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
		LocalTime:  localTime,
	}
	var write io.Writer
	write = zerolog.MultiLevelWriter(rollingFile)
	logger := zerolog.New(write).Level(lv).With().Timestamp().Logger()

	return &logger, nil
}
