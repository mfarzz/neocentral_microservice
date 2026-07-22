package main

import (
	"log"
	"net/http"
	"os"

	"neocentral-go/document-service/internal/config"
	"neocentral-go/document-service/internal/handler"
	"neocentral-go/document-service/internal/repository"
	"neocentral-go/document-service/internal/service"
	"neocentral-go/document-service/pkg/storage"
	"neocentral-go/pkg/database"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()

	// 1. Ensure database exists, then connect
	if err := database.EnsureDatabaseFromDSN(cfg.DBURL); err != nil {
		log.Fatalf("Failed to ensure database: %v", err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DBURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connection established")

	// 2. Init MinIO
	minioSvc, err := storage.NewMinioService(
		cfg.MinioEndpoint,
		cfg.MinioAccess,
		cfg.MinioSecret,
		cfg.MinioSSL,
		cfg.MinioBucket,
	)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}
	log.Println("MinIO connection established")

	// 3. Init Repositories
	docRepo := repository.NewDocumentRepository(db)
	typeRepo := repository.NewDocumentTypeRepository(db)

	// 4. Init Services
	docService := service.NewDocumentService(docRepo, typeRepo, minioSvc, cfg.MinioBucket)

	// 5. Init Handlers
	docHandler := handler.NewDocumentHandler(docService)

	// 6. Setup Router
	r := handler.SetupRouter(cfg, docHandler)

	// 7. Start Server
	port := cfg.AppPort
	if port == "" {
		port = "8003"
	}

	log.Printf("Document Service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
