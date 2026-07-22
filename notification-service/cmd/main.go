package main

import (
	"log"
	"net/http"
	"os"

	"neocentral-go/notification-service/internal/config"
	"neocentral-go/notification-service/internal/handler"
	"neocentral-go/notification-service/internal/repository"
	"neocentral-go/notification-service/internal/service"
	"neocentral-go/notification-service/pkg/broker"
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

	// 2. Connect to NATS
	natsBroker, err := broker.NewNATSBroker(cfg.NatsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsBroker.Close()

	// 3. Init Repositories
	notifRepo := repository.NewNotificationRepository(db)

	// 4. Init Services
	sseManager := service.NewSSEManager()
	notifService := service.NewNotificationService(notifRepo, sseManager)

	// 5. Start NATS Consumer
	natsConsumer := handler.NewNatsConsumer(notifService)
	natsConsumer.Start(natsBroker)
	log.Println("NATS Consumer started")

	// 6. Init HTTP Handlers
	notifHandler := handler.NewNotificationHandler(notifService, sseManager)

	// 7. Setup Router
	r := handler.SetupRouter(cfg, notifHandler)

	// 8. Start Server
	port := cfg.AppPort
	if port == "" {
		port = "8004"
	}

	log.Printf("Notification Service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
