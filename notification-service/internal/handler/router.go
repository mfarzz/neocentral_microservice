package handler

import (
	"net/http"

	"neocentral-go/notification-service/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetupRouter(cfg *config.Config, notifHandler *NotificationHandler) *chi.Mux {
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
		w.Write([]byte("Notification Service OK"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(AuthGuardMiddleware(cfg.JWTSecret))
			
			r.Route("/notifications", func(r chi.Router) {
				r.Get("/", notifHandler.GetHistory)
				r.Get("/stream", notifHandler.Stream)
				r.Patch("/read-all", notifHandler.MarkAllAsRead)
				r.Patch("/{id}/read", notifHandler.MarkAsRead)
			})
		})
	})

	return r
}
