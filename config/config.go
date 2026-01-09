package config

import (
	"os"
	"strconv"
)

type Limiter struct {
	RPS      float64
	Burst    int
	Enabbled bool
}

type Config struct {
	Env         string
	ServerAddr  string
	DatabaseURL string
	MongoDBURL  string
	JWTSecret   string
	Limiter     Limiter
}

func Load() *Config {
	rps, _ := strconv.ParseFloat(getEnv("LIMITER_RPS", "2"), 64)
	burst, _ := strconv.Atoi(getEnv("LIMITER_BURST", "4"))
	enabled, _ := strconv.ParseBool(getEnv("LIMITER_ENABLED", "True"))

	appLimiter := Limiter{
		RPS:      rps,
		Burst:    burst,
		Enabbled: enabled,
	}

	return &Config{
		Env:         getEnv("ENV", "development"),
		ServerAddr:  getEnv("SERVER_ADDRESS", ":8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/habits?sslmode=disable"),
		MongoDBURL:  getEnv("MONGO_DB_URL", "mongodb://localhost:27017"),
		JWTSecret:   getEnv("JWT_SECRET", "9b36f2a2-f8a1-4826-90a6-71d16ca14932"),
		Limiter:     appLimiter,
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
