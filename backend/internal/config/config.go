package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port         string
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	Schema       string
	SSLMode      string
	MaxIdleConns int
	MaxOpenConns int
}

type JWTConfig struct {
	Secret             string
	PlatformExpiry     time.Duration
	SystemExpiry       time.Duration
	RefreshExpiry      time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type CORSConfig struct {
	AllowedOrigins []string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8000"),
			Mode:         getEnv("SERVER_MODE", "development"),
			ReadTimeout:  getDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			User:         getEnv("DB_USER", "postgres"),
			Password:     getEnv("DB_PASSWORD", "postgres"),
			DBName:       getEnv("DB_NAME", "gerege_db"),
			Schema:       getEnv("DB_SCHEMA", "gerege_base"),
			SSLMode:      getEnv("DB_SSLMODE", "disable"),
			MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 100),
		},
		JWT: JWTConfig{
			Secret:         getEnv("JWT_SECRET", "your-super-secret-key-change-in-production"),
			PlatformExpiry: getDuration("JWT_PLATFORM_EXPIRY", 24*time.Hour),
			SystemExpiry:   getDuration("JWT_SYSTEM_EXPIRY", 8*time.Hour),
			RefreshExpiry:  getDuration("JWT_REFRESH_EXPIRY", 168*time.Hour),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:3001"}),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return splitAndTrim(value, ",")
	}
	return defaultValue
}

func splitAndTrim(s string, sep string) []string {
	var result []string
	for _, part := range splitString(s, sep) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitString(s string, sep string) []string {
	var result []string
	var current string
	for _, r := range s {
		if string(r) == sep {
			result = append(result, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
