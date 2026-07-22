package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort       string
	DBURL         string
	JWTSecret     string
	MinioEndpoint string
	MinioAccess   string
	MinioSecret   string
	MinioSSL      bool
	MinioBucket   string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	return &Config{
		AppPort:       getEnv("APP_PORT", "8003"),
		DBURL:         getEnv("DB_URL", "root:rootpassword@tcp(localhost:3306)/neocentral_document?charset=utf8mb4&parseTime=True&loc=Local"),
		JWTSecret:     getEnv("JWT_SECRET", "supersecretkey"),
		MinioEndpoint: getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccess:   getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecret:   getEnv("MINIO_SECRET_KEY", "minioadminpassword"),
		MinioSSL:      getEnv("MINIO_USE_SSL", "false") == "true",
		MinioBucket:   getEnv("MINIO_BUCKET", "neocentral-documents"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
