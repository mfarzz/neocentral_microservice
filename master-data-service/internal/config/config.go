package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"neocentral-go/pkg/database"
)

// Config holds all configuration for master-data-service.
type Config struct {
	HTTPPort  int
	GRPCPort  int
	DB        database.Config
	NATS      string
	JWTSecret string
}

// Load reads configuration from environment variables (or .env file).
func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		HTTPPort: envInt("HTTP_PORT", 8002),
		GRPCPort: envInt("GRPC_PORT", 9002),
		DB: database.Config{
			Host:     envStr("DB_HOST", "localhost"),
			Port:     envInt("DB_PORT", 3306),
			User:     envStr("DB_USER", "root"),
			Password: envStr("DB_PASSWORD", ""),
			DBName:   envStr("DB_NAME", "neocentral_master"),
			MaxOpen:  envInt("DB_MAX_OPEN", 25),
			MaxIdle:  envInt("DB_MAX_IDLE", 10),
			MaxLife:  time.Duration(envInt("DB_MAX_LIFE_MINUTES", 5)) * time.Minute,
		},
		NATS:      envStr("NATS_URL", "nats://localhost:4222"),
		JWTSecret: envStr("JWT_SECRET", ""),
	}

	if cfg.JWTSecret == "" {
		log.Fatal("❌ JWT_SECRET is required")
	}

	return cfg
}

func envStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	}
	return fallback
}
