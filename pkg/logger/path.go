package logger

import (
	"fmt"
	"os"
	"path/filepath"
)

func resolvePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get cwd: %w", err)
	}

	return filepath.Join(cwd, path), nil
}

func ensureParentDir(path string) error {
	logDir := filepath.Dir(path)
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return fmt.Errorf("create log directory: %w", err)
	}

	return nil
}
