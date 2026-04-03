package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func MustGetEnv(key string) (string, error) {
	val := GetEnv(key, "")
	if val == "" {
		return "", fmt.Errorf("%s is not set", key)
	}
	return val, nil
}
func GetWorkingDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Unable to get working dir: %w", err)

	}
	return dir, nil
}

func GetDurationEnv(key string, fallback time.Duration) (time.Duration, error) {
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

func GetEnvInt(key string, defaultValue int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	intV, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue, fmt.Errorf("GetEnvInt: invalid value for key %q: %w", key, err)
	}
	return intV, nil
}

func GetEnvBool(key string, defaultValue bool) (bool, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue, fmt.Errorf("GetEnvBool: invalid value for key %q: %w", key, err)
	}
	return parsed, nil
}

func GetSliceEnv(key string) []string {
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

func NormalizeString(text string) string {
	return strings.ToLower(strings.TrimSpace(text))
}
