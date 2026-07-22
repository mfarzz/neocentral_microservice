package handler

import (
	"net/http"

	"neocentral-go/document-service/internal/service"
	"neocentral-go/pkg/response"

	"github.com/go-chi/chi/v5"
)

type DocumentHandler struct {
	svc *service.DocumentService
}

func NewDocumentHandler(svc *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{svc: svc}
}

func (h *DocumentHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// 32 MB max memory
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "File is required")
		return
	}
	defer file.Close()

	docTypeID := r.FormValue("documentTypeId")
	var docTypePtr *string
	if docTypeID != "" {
		docTypePtr = &docTypeID
	}

	userID := GetUserIDFromContext(r.Context())
	var userPtr *string
	if userID != "" {
		userPtr = &userID
	}

	resp, err := h.svc.UploadDocument(
		r.Context(),
		userPtr,
		docTypePtr,
		header.Filename,
		header.Size,
		header.Header.Get("Content-Type"),
		file,
	)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Document uploaded successfully",
		"data":    resp,
	})
}

func (h *DocumentHandler) GetDownloadURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	resp, err := h.svc.GetDownloadURL(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, resp)
}
