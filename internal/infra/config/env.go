package config

import (
	"os"
)

func GetEnv(key string, fallback ...string) string {
	value := os.Getenv(key)

	if value != "" {
		return value
	}

	if len(fallback) > 0 {
		return fallback[0]
	}

	return ""
}
