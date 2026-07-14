package handler

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"neocentral-go/auth-service/internal/service"
)

// NewRouter builds the Chi router for auth-service.
func NewRouter(
	authSvc *service.AuthService,
	profileSvc *service.ProfileService,
	jwtSecret string,
) chi.Router {
	r := chi.NewRouter()

	authH := NewAuthHandler(authSvc)
	profileH := NewProfileHandler(profileSvc)

	// ── Public auth routes ─────────────────────────
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authH.Login)
		r.Post("/refresh", authH.RefreshToken)
		r.Post("/forgot-password", authH.RequestPasswordReset)
		r.Post("/reset-password", authH.ResetPassword)
		r.Post("/verify", authH.VerifyAccount)
		r.Get("/verify", authH.VerifyAccount) // support GET with ?token=
		r.Post("/request-verification", authH.RequestVerification)
		r.Post("/microsoft/callback", authH.MicrosoftCallback)
		r.Post("/microsoft/mobile", authH.MicrosoftMobileLogin)

		// Protected auth routes
		r.Group(func(r chi.Router) {
			r.Use(AuthGuardMiddleware(jwtSecret))
			r.Post("/logout", authH.Logout)
			r.Post("/change-password", authH.ChangePassword)
		})
	})

	// ── Protected profile routes ───────────────────
	r.Route("/profile", func(r chi.Router) {
		r.Use(AuthGuardMiddleware(jwtSecret))
		r.Get("/me", profileH.GetMe)
	})

	// ── Swagger UI ─────────────────────────────────
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The url pointing to API definition
	))

	return r
}
