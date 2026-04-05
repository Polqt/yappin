package util

import (
	"os"
	"strings"
)

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func GetEnvList(key string, defaults []string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaults
	}

	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			items = append(items, item)
		}
	}

	if len(items) == 0 {
		return defaults
	}

	return items
}
