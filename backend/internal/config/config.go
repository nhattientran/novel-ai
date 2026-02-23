package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration
type Config struct {
	Neo4jURI      string
	Neo4jUser     string
	Neo4jPassword string
	JWTSecret     string
	Port          string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Neo4jURI:      getEnv("NEO4J_URI", "bolt://localhost:7687"),
		Neo4jUser:     getEnv("NEO4J_USER", "neo4j"),
		Neo4jPassword: getEnv("NEO4J_PASSWORD", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		Port:          getEnv("PORT", "8080"),
	}

	if cfg.Neo4jPassword == "" {
		return nil, fmt.Errorf("NEO4J_PASSWORD is required")
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
