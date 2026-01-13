// Package config provides centralized configuration management for the application.
// It loads and validates all configuration from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"chat-application/internal/constants"
)

// Config holds all application configuration values.
type Config struct {
	// Environment settings
	Environment string

	// Server settings
	ServerPort       string
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPIdleTimeout  time.Duration

	// Database settings
	DatabaseURL    string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBMaxOpenConns int
	DBMaxIdleConns int
	DBConnLifetime time.Duration

	// JWT settings
	JWTSecretKey   string
	JWTTokenExpiry time.Duration
	JWTCookieName  string

	// Rate limiting
	DefaultRateLimit  int
	AuthRateLimit     int
	RoomCreationLimit int
	RateLimitWindow   time.Duration

	// Room settings
	MaxRooms            int
	RoomDefaultExpiry   time.Duration
	RoomCleanupInterval time.Duration
	MaxRoomHistory      int

	// External services
	RedditClientID     string
	RedditClientSecret string

	// Allowed origins for CORS and WebSocket
	AllowedOrigins []string
}

// Load creates a new Config instance from environment variables.
// It uses sensible defaults where appropriate.
func Load() (*Config, error) {
	env := getEnv("ENVIRONMENT", "development")

	cfg := &Config{
		// Environment
		Environment: env,

		// Server
		ServerPort:       getEnv("PORT", "8080"),
		HTTPReadTimeout:  constants.HTTPServerTimeout,
		HTTPWriteTimeout: constants.HTTPServerTimeout,
		HTTPIdleTimeout:  constants.HTTPServerTimeout * 2,

		// Database
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "password"),
		DBName:         getEnv("DB_NAME", "chat_app"),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		DBMaxOpenConns: constants.DBMaxOpenConns,
		DBMaxIdleConns: constants.DBMaxIdleConns,
		DBConnLifetime: constants.DBConnMaxLifetime,

		// JWT
		JWTSecretKey:   getEnv("JWT_SECRET_KEY", ""),
		JWTTokenExpiry: constants.JWTTokenExpiry,
		JWTCookieName:  constants.JWTCookieName,

		// Rate limiting
		DefaultRateLimit:  constants.DefaultRateLimit,
		AuthRateLimit:     constants.AuthRateLimit,
		RoomCreationLimit: constants.RoomCreationLimit,
		RateLimitWindow:   constants.RateLimitWindow,

		// Room settings
		MaxRooms:            getEnvInt("MAX_ROOMS", constants.DefaultRoomLimit),
		RoomDefaultExpiry:   constants.RoomDefaultExpiry,
		RoomCleanupInterval: constants.RoomCleanupInterval,
		MaxRoomHistory:      constants.MaxRoomHistory,

		// External services
		RedditClientID:     getEnv("REDDIT_CLIENT_ID", ""),
		RedditClientSecret: getEnv("REDDIT_CLIENT_SECRET", ""),

		// Allowed origins
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://localhost:5174",
			"https://yappin.chat",
		},
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that all required configuration values are present.
func (c *Config) Validate() error {
	if c.Environment == "production" {
		if c.JWTSecretKey == "" {
			return fmt.Errorf("JWT_SECRET_KEY is required in production")
		}
		if c.DatabaseURL == "" {
			return fmt.Errorf("DATABASE_URL is required in production")
		}
	}
	return nil
}

// IsDevelopment returns true if running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.Environment != "production"
}

// IsProduction returns true if running in production mode.
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// GetDSN returns the database connection string based on environment.
func (c *Config) GetDSN() string {
	if c.IsProduction() && c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

// Helper functions

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
