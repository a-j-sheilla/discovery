package configs

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server ServerConfig
	TMDB   TMDBConfig
	OMDB   OMDBConfig
	Cache  CacheConfig
	Rate   RateLimitConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Host string
}

// TMDBConfig holds TMDB API configuration
type TMDBConfig struct {
	APIKey  string
	BaseURL string
}

// OMDBConfig holds OMDB API configuration
type OMDBConfig struct {
	APIKey  string
	BaseURL string
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Duration time.Duration
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "localhost"),
		},
		TMDB: TMDBConfig{
			APIKey:  getEnv("TMDB_API_KEY", ""),
			BaseURL: getEnv("TMDB_BASE_URL", "https://api.themoviedb.org/3"),
		},
		OMDB: OMDBConfig{
			APIKey:  getEnv("OMDB_API_KEY", ""),
			BaseURL: getEnv("OMDB_BASE_URL", "http://www.omdbapi.com"),
		},
		Cache: CacheConfig{
			Duration: time.Duration(getEnvAsInt("CACHE_DURATION_MINUTES", 30)) * time.Minute,
		},
		Rate: RateLimitConfig{
			RequestsPerMinute: getEnvAsInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 60),
		},
	}

	return config, nil
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
