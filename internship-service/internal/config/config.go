package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort            string
	DBURL              string
	JWTSecret          string
	NatsURL            string
	DocumentServiceURL string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	return &Config{
		AppPort:            getEnv("APP_PORT", "8005"),
		DBURL:              getEnv("DB_URL", "root:rootpassword@tcp(localhost:3306)/neocentral_internship?charset=utf8mb4&parseTime=True&loc=Local"),
		JWTSecret:          getEnv("JWT_SECRET", "supersecretkey"),
		NatsURL:            getEnv("NATS_URL", "nats://localhost:4222"),
		DocumentServiceURL: getEnv("DOCUMENT_SERVICE_URL", "http://localhost:8003"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
