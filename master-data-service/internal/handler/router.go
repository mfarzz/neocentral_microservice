package handler

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"neocentral-go/master-data-service/internal/service"
)

// NewRouter builds the Chi router for master-data-service.
func NewRouter(
	ayService *service.AcademicYearService,
	roomService *service.RoomService,
	sgService *service.ScienceGroupService,
	thesisService *service.ThesisTopicService,
	internshipHolidayService *service.InternshipHolidayService,
	jwtSecret string,
) chi.Router {
	r := chi.NewRouter()

	ayH := NewAcademicYearHandler(ayService)
	roomH := NewRoomHandler(roomService)
	sgH := NewScienceGroupHandler(sgService)
	thesisH := NewThesisHandler(thesisService)
	holidayH := NewInternshipHolidayHandler(internshipHolidayService)

	// All master-data routes are protected
	r.Route("/master-data", func(r chi.Router) {
		r.Use(AuthGuardMiddleware(jwtSecret))

		// ── Academic Years ────────────────────────
		r.Route("/academic-years", func(r chi.Router) {
			r.Get("/", ayH.GetAll)
			r.Get("/active", ayH.GetActive)
			r.Get("/{id}", ayH.GetByID)
			
			r.Group(func(r chi.Router) {
				r.Use(RequireRoleMiddleware("Admin", "GKM"))
				r.Post("/", ayH.Create)
				r.Patch("/{id}", ayH.Update)
				r.Delete("/{id}", ayH.Delete)
			})
		})

		// ── Rooms ────────────────────────────────
		r.Route("/rooms", func(r chi.Router) {
			r.Get("/", roomH.GetAll)
			r.Get("/{id}", roomH.GetByID)

			r.Group(func(r chi.Router) {
				r.Use(RequireRoleMiddleware("Admin", "GKM"))
				r.Post("/", roomH.Create)
				r.Patch("/{id}", roomH.Update)
				r.Delete("/{id}", roomH.Delete)
			})
		})

		// ── Science Groups ───────────────────────
		r.Route("/science-groups", func(r chi.Router) {
			r.Get("/", sgH.GetAll)
			r.Get("/{id}", sgH.GetByID)

			r.Group(func(r chi.Router) {
				r.Use(RequireRoleMiddleware("Admin", "GKM"))
				r.Post("/", sgH.Create)
				r.Patch("/{id}", sgH.Update)
				r.Delete("/{id}", sgH.Delete)
			})
		})

		// ── Thesis Topics ────────────────────────
		r.Route("/thesis-topics", func(r chi.Router) {
			r.Get("/", thesisH.GetAllTopics)
			r.Get("/{id}", thesisH.GetTopicByID)

			r.Group(func(r chi.Router) {
				r.Use(RequireRoleMiddleware("Admin", "GKM"))
				r.Post("/", thesisH.CreateTopic)
				r.Patch("/{id}", thesisH.UpdateTopic)
				r.Delete("/{id}", thesisH.DeleteTopic)
			})
		})

		// ── Thesis Statuses ──────────────────────
		r.Get("/thesis-statuses", thesisH.GetAllStatuses)

		// ── Internship Holidays ──────────────────
		r.Route("/internship-holidays", func(r chi.Router) {
			r.Get("/", holidayH.GetAll)
			
			r.Group(func(r chi.Router) {
				r.Use(RequireRoleMiddleware("Admin", "GKM"))
				r.Post("/", holidayH.Create)
				r.Put("/{id}", holidayH.Update)
				r.Delete("/{id}", holidayH.Delete)
			})
		})
	})

	// ── Swagger UI ─────────────────────────────
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	return r
}
