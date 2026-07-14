package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"neocentral-go/pkg/auth"
	"neocentral-go/pkg/database"
)

// Config holds all configuration for auth-service.
type Config struct {
	HTTPPort int
	GRPCPort int
	DB       database.Config
	Redis    RedisConfig
	JWT      auth.JWTConfig
	NATS     string
	Microsoft MicrosoftConfig
	SMTP     SMTPConfig
	BaseURL  string
	AppName  string
}

type RedisConfig struct {
	URL string
}

type MicrosoftConfig struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	RedirectURI  string
}

type SMTPConfig struct {
	Host string
	Port int
	User string
	Pass string
	From string
}

// Load reads configuration from environment variables (or .env file).
func Load() *Config {
	// Load .env file if it exists (dev mode)
	_ = godotenv.Load()

	cfg := &Config{
		HTTPPort: envInt("HTTP_PORT", 8001),
		GRPCPort: envInt("GRPC_PORT", 9001),
		DB: database.Config{
			Host:     envStr("DB_HOST", "localhost"),
			Port:     envInt("DB_PORT", 3306),
			User:     envStr("DB_USER", "root"),
			Password: envStr("DB_PASSWORD", ""),
			DBName:   envStr("DB_NAME", "neocentral_auth"),
			MaxOpen:  envInt("DB_MAX_OPEN", 25),
			MaxIdle:  envInt("DB_MAX_IDLE", 10),
			MaxLife:  time.Duration(envInt("DB_MAX_LIFE_MINUTES", 5)) * time.Minute,
		},
		Redis: RedisConfig{
			URL: envStr("REDIS_URL", "redis://localhost:6379"),
		},
		JWT: auth.JWTConfig{
			Secret:        envStr("JWT_SECRET", ""),
			AccessExpiry:  parseDuration(envStr("JWT_EXPIRES_IN", "7d")),
			RefreshSecret: envStr("REFRESH_TOKEN_SECRET", ""),
			RefreshExpiry: parseDuration(envStr("REFRESH_TOKEN_EXPIRES_IN", "30d")),
		},
		NATS: envStr("NATS_URL", "nats://localhost:4222"),
		Microsoft: MicrosoftConfig{
			ClientID:     envStr("CLIENT_ID", ""),
			ClientSecret: envStr("CLIENT_SECRET", ""),
			TenantID:     envStr("TENANT_ID", ""),
			RedirectURI:  envStr("REDIRECT_URI", ""),
		},
		SMTP: SMTPConfig{
			Host: envStr("SMTP_HOST", "smtp.gmail.com"),
			Port: envInt("SMTP_PORT", 587),
			User: envStr("SMTP_USER", ""),
			Pass: envStr("SMTP_PASS", ""),
			From: envStr("SMTP_FROM", ""),
		},
		BaseURL: envStr("BASE_URL", "http://localhost:8001"),
		AppName: envStr("APP_NAME", "NeoCentral API"),
	}

	// Validate required vars
	if cfg.JWT.Secret == "" {
		log.Fatal("❌ JWT_SECRET is required")
	}
	if cfg.JWT.RefreshSecret == "" {
		log.Fatal("❌ REFRESH_TOKEN_SECRET is required")
	}

	return cfg
}

// ── Helpers ─────────────────────────────────────────────────────

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

// parseDuration converts strings like "7d", "15m", "24h" to time.Duration.
func parseDuration(s string) time.Duration {
	if len(s) == 0 {
		return 7 * 24 * time.Hour
	}
	unit := s[len(s)-1]
	numStr := s[:len(s)-1]
	n, err := strconv.Atoi(numStr)
	if err != nil {
		return 7 * 24 * time.Hour
	}
	switch unit {
	case 'd':
		return time.Duration(n) * 24 * time.Hour
	case 'h':
		return time.Duration(n) * time.Hour
	case 'm':
		return time.Duration(n) * time.Minute
	case 's':
		return time.Duration(n) * time.Second
	default:
		return 7 * 24 * time.Hour
	}
}
