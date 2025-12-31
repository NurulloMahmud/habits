package config

import "os"

type Config struct {
	Env         string
	ServerAddr  string
	DatabaseURL string
	JWTSecret   string
}

func Load() *Config {
	return &Config{
		Env:         getEnv("ENV", "development"),
		ServerAddr:  getEnv("SERVER_ADDRESS", ":8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/habits?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "9b36f2a2-f8a1-4826-90a6-71d16ca14932"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}
