package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"neocentral-go/auth-service/internal/config"
	"neocentral-go/auth-service/internal/event"
	"neocentral-go/auth-service/internal/grpcserver"
	"neocentral-go/auth-service/internal/handler"
	"neocentral-go/auth-service/internal/repository"
	"neocentral-go/auth-service/internal/service"
	_ "neocentral-go/auth-service/docs" // Swagger docs
	"neocentral-go/pkg/database"
	"neocentral-go/pkg/messaging"
	pb "neocentral-go/proto/auth"
)

// @title           NeoCentral Auth Service API
// @version         1.0
// @description     API Documentation for NeoCentral Auth Service
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@neocentral.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8001
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

	// ── Redis ────────────────────────────────────
	opt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		log.Fatalf("❌ Redis URL parse: %v", err)
	}
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Printf("⚠️  Redis not available (password-reset / verify will fail): %v", err)
	} else {
		log.Println("✅ Connected to Redis")
	}

	// ── NATS ─────────────────────────────────────
	var natsClient *messaging.NATSClient
	var userEvents *event.UserEventPublisher

	natsClient, err = messaging.NewNATSClient(cfg.NATS)
	if err != nil {
		log.Printf("⚠️  NATS not available (events will be skipped): %v", err)
	} else {
		_ = natsClient.EnsureStream("AUTH", []string{"auth.>"})
		userEvents = event.NewUserEventPublisher(natsClient)
		defer natsClient.Close()
	}

	// ── Repositories ─────────────────────────────
	userRepo := repository.NewGormUserRepository(db)

	// ── Services ─────────────────────────────────
	authSvc := service.NewAuthService(userRepo, cfg.JWT, userEvents, rdb)
	profileSvc := service.NewProfileService(userRepo)

	// ── REST Server (Chi) ────────────────────────
	root := chi.NewRouter()
	root.Use(chimw.RequestID)
	root.Use(chimw.RealIP)
	root.Use(chimw.Logger)
	root.Use(chimw.Recoverer)
	root.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Refresh-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Mount service routes
	appRouter := handler.NewRouter(authSvc, profileSvc, cfg.JWT.Secret)
	root.Mount("/", appRouter)

	// Health check
	root.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","service":"auth-service","uptime":"%.0fs"}`,
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

	// ── gRPC Server ──────────────────────────────
	grpcAddr := fmt.Sprintf(":%d", cfg.GRPCPort)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("❌ gRPC listen %s: %v", grpcAddr, err)
	}

	grpcSrv := grpc.NewServer()
	authGRPC := grpcserver.NewAuthGRPCServer(authSvc, userRepo)
	pb.RegisterAuthServiceServer(grpcSrv, authGRPC)
	reflection.Register(grpcSrv)

	// ── Start ────────────────────────────────────
	go func() {
		log.Printf("🌐 REST  → http://localhost%s", httpAddr)
		log.Printf("   Routes: /auth/*, /profile/*")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP: %v", err)
		}
	}()

	go func() {
		log.Printf("⚡ gRPC  → %s", grpcAddr)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("gRPC: %v", err)
		}
	}()

	// ── Graceful Shutdown ────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("🛑 Received %s — shutting down...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	grpcSrv.GracefulStop()
	_ = httpServer.Shutdown(ctx)
	_ = rdb.Close()

	log.Println("✅ auth-service stopped gracefully")
}
