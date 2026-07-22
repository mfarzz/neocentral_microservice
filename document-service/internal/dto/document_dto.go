package dto

import "time"

type DocumentResponse struct {
	ID             string    `json:"id"`
	UserID         *string   `json:"userId"`
	DocumentTypeID *string   `json:"documentTypeId"`
	FileName       *string   `json:"fileName"`
	CreatedAt      time.Time `json:"createdAt"`
}

type DocumentUploadResponse struct {
	ID           string `json:"id"`
	PresignedURL string `json:"presignedUrl,omitempty"`
}

type DocumentTypeRequest struct {
	Name string `json:"name" validate:"required"`
}

type DocumentTypeResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}
