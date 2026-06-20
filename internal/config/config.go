package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort      string
	RedisHost       string
	RedisPort       string
	BaseURL         string
	RateLimit       int
	RateLimitWindow time.Duration
	URLTTL          time.Duration
}

func Load() *Config {
	return &Config{
		ServerPort:      getEnv("SERVER_PORT", "3000"),
		RedisHost:       getEnv("REDIS_HOST", "localhost"),
		RedisPort:       getEnv("REDIS_PORT", "6379"),
		BaseURL:         getEnv("BASE_URL", "http://localhost:3000"),
		RateLimit:       getEnvInt("RATE_LIMIT_REQUESTS", 5),
		RateLimitWindow: getEnvDuration("RATE_LIMIT_WINDOW", time.Hour),
		URLTTL:          getEnvDuration("URL_TTL", 24*time.Hour),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}
