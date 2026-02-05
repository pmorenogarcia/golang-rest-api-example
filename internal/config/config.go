package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	PokeAPI  PokeAPIConfig
	Logging  LoggingConfig
	CORS     CORSConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// PokeAPIConfig holds PokeAPI client configuration
type PokeAPIConfig struct {
	BaseURL string
	Timeout time.Duration
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins string
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Load .env file if it exists (optional, won't fail if not found)
	_ = godotenv.Load()

	// Configure viper
	viper.AutomaticEnv()

	// Set default values
	setDefaults()

	// Parse configuration
	cfg := &Config{
		Server: ServerConfig{
			Port:         viper.GetString("SERVER_PORT"),
			ReadTimeout:  viper.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout: viper.GetDuration("SERVER_WRITE_TIMEOUT"),
			IdleTimeout:  viper.GetDuration("SERVER_IDLE_TIMEOUT"),
		},
		PokeAPI: PokeAPIConfig{
			BaseURL: viper.GetString("POKEAPI_BASE_URL"),
			Timeout: viper.GetDuration("POKEAPI_TIMEOUT"),
		},
		Logging: LoggingConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
		},
		CORS: CORSConfig{
			AllowedOrigins: viper.GetString("CORS_ALLOWED_ORIGINS"),
		},
	}

	// Validate configuration
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// setDefaults sets default values for configuration
func setDefaults() {
	// Server defaults
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_READ_TIMEOUT", "10s")
	viper.SetDefault("SERVER_WRITE_TIMEOUT", "10s")
	viper.SetDefault("SERVER_IDLE_TIMEOUT", "120s")

	// PokeAPI defaults
	viper.SetDefault("POKEAPI_BASE_URL", "https://pokeapi.co/api/v2")
	viper.SetDefault("POKEAPI_TIMEOUT", "30s")

	// Logging defaults
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")

	// CORS defaults
	viper.SetDefault("CORS_ALLOWED_ORIGINS", "*")
}

// validate validates the configuration
func (c *Config) validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("SERVER_PORT is required")
	}

	if c.PokeAPI.BaseURL == "" {
		return fmt.Errorf("POKEAPI_BASE_URL is required")
	}

	if c.Logging.Level == "" {
		return fmt.Errorf("LOG_LEVEL is required")
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.Logging.Level] {
		return fmt.Errorf("invalid LOG_LEVEL: must be one of debug, info, warn, error")
	}

	return nil
}
