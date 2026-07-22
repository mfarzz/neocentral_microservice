package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"neocentral-go/api-gateway/internal/middleware"
	"neocentral-go/api-gateway/internal/proxy"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load environment variables
	_ = godotenv.Load(".env") // Ignore error if .env doesn't exist, rely on system env

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:8001"
	}

	masterDataServiceURL := os.Getenv("MASTER_DATA_SERVICE_URL")
	if masterDataServiceURL == "" {
		masterDataServiceURL = "http://localhost:8002"
	}

	documentServiceURL := os.Getenv("DOCUMENT_SERVICE_URL")
	if documentServiceURL == "" {
		documentServiceURL = "http://localhost:8003"
	}

	// 2. Setup Router & Global Middlewares
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(middleware.CORS().Handler)

	// 3. Health Check
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"service": "api-gateway", "status": "online"}`))
	})

	// 4. Configure Upstream Proxies
	// Auth Service → /api/v1/auth/*, /api/v1/profile/*
	authProxy := proxy.NewReverseProxy(authServiceURL)
	r.Mount("/api/v1/auth", http.StripPrefix("/api/v1", authProxy))
	r.Mount("/api/v1/profile", http.StripPrefix("/api/v1", authProxy))

	// Master Data Service → /api/v1/master-data/*
	masterProxy := proxy.NewReverseProxy(masterDataServiceURL)
	r.Mount("/api/v1/master-data", http.StripPrefix("/api/v1", masterProxy))

	// Document Service → /api/v1/documents/*
	documentProxy := proxy.NewReverseProxy(documentServiceURL)
	r.Mount("/api/v1/documents", http.StripPrefix("/api/v1", documentProxy))

	// 5. Start Server with Graceful Shutdown
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("🚀 API Gateway is running on http://localhost:%s", port)
		log.Printf("   Routing /api/v1/auth         -> %s", authServiceURL)
		log.Printf("   Routing /api/v1/master-data  -> %s", masterDataServiceURL)
		log.Printf("   Routing /api/v1/documents    -> %s", documentServiceURL)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Listen error: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Received interrupt — shutting down API Gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server Shutdown Failed:%+v", err)
	}
	log.Println("✅ API Gateway stopped gracefully")
}
