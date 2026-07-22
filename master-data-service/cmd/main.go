package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"neocentral-go/master-data-service/internal/config"
	"neocentral-go/master-data-service/internal/handler"
	"neocentral-go/master-data-service/internal/repository"
	"neocentral-go/master-data-service/internal/service"
	// _ "neocentral-go/master-data-service/docs" // Swagger docs (uncomment after running: make swagger)
	"neocentral-go/pkg/database"
)

// @title           NeoCentral Master Data Service API
// @version         1.0
// @description     API Documentation for NeoCentral Master Data Service
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@neocentral.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8002
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	cfg := config.Load()
	startTime := time.Now()

	// ── MySQL (GORM) ─────────────────────────────
	db, err := database.NewMySQLConnection(cfg.DB)
	if err != nil {
		log.Fatalf("❌ Database: %v", err)
	}

	// ── Repositories ─────────────────────────────
	ayRepo := repository.NewGormAcademicYearRepo(db)
	roomRepo := repository.NewGormRoomRepo(db)
	sgRepo := repository.NewGormScienceGroupRepo(db)
	topicRepo := repository.NewGormThesisTopicRepo(db)
	statusRepo := repository.NewGormThesisStatusRepo(db)
	holidayRepo := repository.NewGormInternshipHolidayRepo(db)

	// ── Services ─────────────────────────────────
	aySvc := service.NewAcademicYearService(ayRepo)
	roomSvc := service.NewRoomService(roomRepo)
	sgSvc := service.NewScienceGroupService(sgRepo)
	thesisSvc := service.NewThesisTopicService(topicRepo, statusRepo)
	holidaySvc := service.NewInternshipHolidayService(holidayRepo)

	// ── REST Server (Chi) ────────────────────────
	root := chi.NewRouter()
	root.Use(chimw.RequestID)
	root.Use(chimw.RealIP)
	root.Use(chimw.Logger)
	root.Use(chimw.Recoverer)
	root.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Mount service routes
	appRouter := handler.NewRouter(aySvc, roomSvc, sgSvc, thesisSvc, holidaySvc, cfg.JWTSecret)
	root.Mount("/", appRouter)

	// Health check
	root.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","service":"master-data-service","uptime":"%.0fs"}`,
			time.Since(startTime).Seconds())
	})

	httpAddr := fmt.Sprintf(":%d", cfg.HTTPPort)
	httpServer := &http.Server{
		Addr:         httpAddr,
		Handler:      root,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ── Start ────────────────────────────────────
	go func() {
		log.Printf("🌐 REST  → http://localhost%s", httpAddr)
		log.Printf("   Routes: /master-data/*")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP: %v", err)
		}
	}()

	// ── Graceful Shutdown ────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("🛑 Received %s — shutting down...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_ = httpServer.Shutdown(ctx)
	log.Println("✅ master-data-service stopped gracefully")
}
