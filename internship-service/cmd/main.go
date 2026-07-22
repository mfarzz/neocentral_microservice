package main

import (
	"log"
	"net/http"

	"neocentral-go/internship-service/internal/config"
	"neocentral-go/internship-service/internal/handler"
	_ "neocentral-go/internship-service/docs" // Swagger docs
	"neocentral-go/pkg/database"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// @title           NeoCentral Internship Service API
// @version         1.0
// @description     API Documentation for NeoCentral Internship Service
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@neocentral.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8005
// @BasePath  /api/v1/internships

func main() {
	cfg := config.LoadConfig()

	if err := database.EnsureDatabaseFromDSN(cfg.DBURL); err != nil {
		log.Fatalf("Failed to ensure database: %v", err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DBURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	router := chi.NewRouter()

	handler.SetupRoutes(router, db, cfg)

	log.Printf("Starting Internship Service on port %s...", cfg.AppPort)
	if err := http.ListenAndServe(":"+cfg.AppPort, router); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
