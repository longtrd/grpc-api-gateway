package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Logger   LoggerConfig   `yaml:"logger"`
	Features FeatureConfig  `yaml:"features"`
}

// ServerConfig contains server-specific configuration
type ServerConfig struct {
	Port         string        `yaml:"port"`
	Host         string        `yaml:"host"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	MaxRecvSize  int           `yaml:"max_recv_size"`
	MaxSendSize  int           `yaml:"max_send_size"`
}

// DatabaseConfig contains database configuration
type DatabaseConfig struct {
	Type         string `yaml:"type"` // memory, postgres, mysql
	URL          string `yaml:"url"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

// LoggerConfig contains logging configuration
type LoggerConfig struct {
	Level  string `yaml:"level"`  // debug, info, warn, error
	Format string `yaml:"format"` // json, text
}

// FeatureConfig contains feature flags
type FeatureConfig struct {
	EnableMetrics     bool `yaml:"enable_metrics"`
	EnableTracing     bool `yaml:"enable_tracing"`
	EnableAuth        bool `yaml:"enable_auth"`
	EnableHealthCheck bool `yaml:"enable_health_check"`
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("GRPC_PORT", "50051"),
			Host:         getEnv("GRPC_HOST", "0.0.0.0"),
			ReadTimeout:  getDurationEnv("GRPC_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("GRPC_WRITE_TIMEOUT", 30*time.Second),
			MaxRecvSize:  getIntEnv("GRPC_MAX_RECV_SIZE", 4*1024*1024), // 4MB
			MaxSendSize:  getIntEnv("GRPC_MAX_SEND_SIZE", 4*1024*1024), // 4MB
		},
		Database: DatabaseConfig{
			Type:         getEnv("DATABASE_TYPE", "memory"),
			URL:          getEnv("DATABASE_URL", ""),
			MaxOpenConns: getIntEnv("DATABASE_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getIntEnv("DATABASE_MAX_IDLE_CONNS", 25),
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		Features: FeatureConfig{
			EnableMetrics:     getBoolEnv("ENABLE_METRICS", false),
			EnableTracing:     getBoolEnv("ENABLE_TRACING", false),
			EnableAuth:        getBoolEnv("ENABLE_AUTH", false),
			EnableHealthCheck: getBoolEnv("ENABLE_HEALTH_CHECK", true),
		},
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if c.Database.Type == "postgres" && c.Database.URL == "" {
		return fmt.Errorf("database URL is required for postgres")
	}

	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[c.Logger.Level] {
		return fmt.Errorf("invalid log level: %s", c.Logger.Level)
	}

	validLogFormats := map[string]bool{"json": true, "text": true}
	if !validLogFormats[c.Logger.Format] {
		return fmt.Errorf("invalid log format: %s", c.Logger.Format)
	}

	return nil
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
