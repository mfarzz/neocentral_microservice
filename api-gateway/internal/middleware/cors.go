package middleware

import (
	"github.com/go-chi/cors"
)

// CORS returns a middleware that handles Cross-Origin Resource Sharing.
func CORS() *cors.Cors {
	return cors.New(cors.Options{
		// Set to "*" for development. In production, change to specific frontend domains.
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}
