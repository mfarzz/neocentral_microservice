package handler

import (
	"net/http"

	"neocentral-go/document-service/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetupRouter(cfg *config.Config, docHandler *DocumentHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Adjust in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Document Service OK"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Public or token-optional routes
		
		r.Group(func(r chi.Router) {
			r.Use(AuthGuardMiddleware(cfg.JWTSecret))
			
			r.Route("/documents", func(r chi.Router) {
				r.Post("/upload", docHandler.Upload)
				r.Get("/{id}/download-url", docHandler.GetDownloadURL)
			})
		})
	})

	return r
}
