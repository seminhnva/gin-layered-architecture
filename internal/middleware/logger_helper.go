package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

func formatSize(size int64) string {
	switch {
	case size >= 1<<20:
		return fmt.Sprintf("%.2f MB", float64(size)/(1<<20))
	case size >= 1<<10:
		return fmt.Sprintf("%.2f KB", float64(size)/(1<<10))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

func truncate(s string) string {
	const limit = 2000
	if len(s) <= limit {
		return s
	}

	return s[:limit] + "... [TRUNCATED]"
}

func compactText(s string) string {
	if strings.TrimSpace(s) == "" {
		return ""
	}

	return truncate(strings.Join(strings.Fields(s), " "))
}

func parseLoggedPayload(raw []byte) any {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil
	}

	var parsed any
	if err := json.Unmarshal(trimmed, &parsed); err == nil {
		return parsed
	}

	return compactText(string(trimmed))
}

func sanitizeValue(value any, sensitiveFields map[string]bool) any {
	switch val := value.(type) {
	case map[string]any:
		result := make(map[string]any, len(val))
		for k, v := range val {
			if sensitiveFields[strings.ToLower(k)] {
				result[k] = "[REDACTED]"
				continue
			}
			result[k] = sanitizeValue(v, sensitiveFields)
		}
		return result
	case []any:
		arr := make([]any, len(val))
		for i, item := range val {
			arr[i] = sanitizeValue(item, sensitiveFields)
		}
		return arr
	default:
		return value
	}
}
