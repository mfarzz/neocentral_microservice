package handler

import (
	"neocentral-go/internship-service/internal/config"
	"neocentral-go/internship-service/internal/repository"
	"neocentral-go/internship-service/internal/service"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"gorm.io/gorm"
)

func SetupRoutes(router *chi.Mux, db *gorm.DB, cfg *config.Config) {
	// Repositories
	proposalRepo := repository.NewProposalRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	internshipRepo := repository.NewInternshipRepository(db)
	letterRepo := repository.NewSupervisorLetterRepository(db)
	logbookRepo := repository.NewLogbookRepository(db)
	bimbinganRepo := repository.NewBimbinganRepository(db)
	seminarRepo := repository.NewSeminarRepository(db)
	penilaianRepo := repository.NewPenilaianRepository(db)
	
	// Services
	pendaftaranService := service.NewPendaftaranService(proposalRepo, companyRepo)
	penunjukanService := service.NewPenunjukanService(internshipRepo, letterRepo)
	pelaksanaanService := service.NewPelaksanaanService(internshipRepo, logbookRepo)
	bimbinganService := service.NewBimbinganService(bimbinganRepo, internshipRepo)
	seminarService := service.NewSeminarService(seminarRepo)
	penilaianService := service.NewPenilaianService(penilaianRepo)
	monitoringService := service.NewMonitoringService(internshipRepo)

	// Handlers
	pendaftaranHandler := NewPendaftaranHandler(pendaftaranService)
	penunjukanHandler := NewPenunjukanHandler(penunjukanService)
	pelaksanaanHandler := NewPelaksanaanHandler(pelaksanaanService)
	bimbinganHandler := NewBimbinganHandler(bimbinganService)
	seminarHandler := NewSeminarHandler(seminarService)
	penilaianHandler := NewPenilaianHandler(penilaianService)
	monitoringHandler := NewMonitoringHandler(monitoringService)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The url pointing to API definition
	))

	router.Route("/api/v1/internships", func(r chi.Router) {
		// Modul 1: Pendaftaran
		r.Route("/pendaftaran", func(r chi.Router) {
			r.Post("/proposals", pendaftaranHandler.CreateProposal)
			r.Get("/proposals", pendaftaranHandler.ListStudentProposals)
			r.Get("/proposals/{id}", pendaftaranHandler.GetProposalByID)
			r.Put("/proposals/{id}", pendaftaranHandler.UpdateProposal)
			r.Patch("/proposals/{id}/status", pendaftaranHandler.UpdateProposalStatus)
			r.Delete("/proposals/{id}", pendaftaranHandler.DeleteProposal)
			r.Get("/companies", pendaftaranHandler.ListCompanies)
		})

		// Modul 2: Penunjukan Pembimbing
		r.Route("/penunjukan-pembimbing", func(r chi.Router) {
			r.Post("/assign-bulk", penunjukanHandler.BulkAssignSupervisor)
			r.Get("/workload", penunjukanHandler.GetLecturersWorkload)
			r.Get("/letters/{supervisorId}", penunjukanHandler.GetSupervisorLetter)
			r.Put("/letters/{supervisorId}", penunjukanHandler.UpdateSupervisorLetter)
		})

		// Modul 3: Pelaksanaan
		r.Route("/pelaksanaan", func(r chi.Router) {
			r.Get("/history", pelaksanaanHandler.GetInternshipHistory)
			r.Get("/{internshipId}/logbooks", pelaksanaanHandler.GetLogbooks)
			r.Post("/{internshipId}/logbooks", pelaksanaanHandler.CreateLogbook)
			r.Put("/logbooks/{logbookId}", pelaksanaanHandler.UpdateLogbook)
			r.Put("/{internshipId}/details", pelaksanaanHandler.UpdateInternshipDetails)
			r.Post("/{internshipId}/documents/{docType}", pelaksanaanHandler.SubmitDocument)
		})

		// Modul 4: Bimbingan
		r.Route("/bimbingan", func(r chi.Router) {
			r.Get("/questions", bimbinganHandler.GetQuestions)
			r.Post("/questions", bimbinganHandler.CreateQuestion)
			r.Get("/criteria", bimbinganHandler.GetLecturerCriteria)
			r.Post("/criteria", bimbinganHandler.CreateLecturerCriteria)
			r.Get("/student", bimbinganHandler.GetStudentGuidance)
			r.Post("/student", bimbinganHandler.SubmitStudentGuidance)
			r.Post("/evaluation", bimbinganHandler.SubmitLecturerEvaluation)
		})

		// Modul 5: Seminar
		r.Route("/seminar", func(r chi.Router) {
			r.Get("/upcoming", seminarHandler.GetUpcomingSeminars)
			r.Post("/", seminarHandler.RegisterSeminar)
			r.Post("/{id}/approve-bulk", seminarHandler.ApproveSeminarBulk)
			
			// Audience
			r.Post("/{id}/audience", seminarHandler.RegisterAudience)
			r.Post("/{id}/audience/validate-bulk", seminarHandler.ValidateAudienceBulk)
		})

		// Modul 6: Penilaian
		r.Route("/penilaian", func(r chi.Router) {
			r.Get("/cpmks", penilaianHandler.GetCPMKs)
			r.Post("/cpmks", penilaianHandler.CreateCPMK)
			r.Post("/rubrics", penilaianHandler.CreateRubric)
			
			r.Post("/{internshipId}/lecturer-assessment", penilaianHandler.SubmitLecturerAssessment)
			r.Post("/field-assessment/{token}", penilaianHandler.SubmitFieldAssessment)
		})

		// Modul 7: Monitoring (Sekdep/Admin)
		r.Route("/monitoring", func(r chi.Router) {
			r.Get("/", monitoringHandler.GetInternshipList)
			r.Get("/{id}", monitoringHandler.GetInternshipDetail)
			r.Get("/grade-recap", monitoringHandler.GetGradeRecap)
			r.Get("/stats", monitoringHandler.GetMonitoringStats)
			r.Post("/{id}/verify-document", monitoringHandler.VerifyDocument)
			r.Post("/{id}/reject-report", monitoringHandler.RejectFinalReport)
		})
	})
}
