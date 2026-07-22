package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"neocentral-go/document-service/internal/dto"
	"neocentral-go/document-service/internal/model"
	"neocentral-go/document-service/internal/repository"
	"neocentral-go/document-service/pkg/storage"

	"github.com/google/uuid"
)

type DocumentService struct {
	docRepo  repository.DocumentRepository
	typeRepo repository.DocumentTypeRepository
	minio    *storage.MinioService
	bucket   string
}

func NewDocumentService(
	docRepo repository.DocumentRepository,
	typeRepo repository.DocumentTypeRepository,
	minio *storage.MinioService,
	bucket string,
) *DocumentService {
	return &DocumentService{
		docRepo:  docRepo,
		typeRepo: typeRepo,
		minio:    minio,
		bucket:   bucket,
	}
}

func (s *DocumentService) UploadDocument(
	ctx context.Context,
	userID *string,
	docTypeID *string,
	fileName string,
	fileSize int64,
	contentType string,
	reader io.Reader,
) (*dto.DocumentUploadResponse, error) {
	docID := uuid.New().String()
	objectName := fmt.Sprintf("%s/%s_%s", time.Now().Format("2006/01/02"), docID, fileName)

	// Upload to MinIO
	_, err := s.minio.UploadFile(ctx, objectName, reader, fileSize, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to storage: %w", err)
	}

	// Save to DB
	doc := &model.Document{
		ID:             docID,
		UserID:         userID,
		DocumentTypeID: docTypeID,
		S3Bucket:       s.bucket,
		S3ObjectName:   objectName,
		FileName:       &fileName,
	}

	err = s.docRepo.Create(ctx, doc)
	if err != nil {
		// Clean up MinIO on db failure
		_ = s.minio.DeleteFile(ctx, objectName)
		return nil, fmt.Errorf("failed to save document record: %w", err)
	}

	return &dto.DocumentUploadResponse{
		ID: docID,
	}, nil
}

func (s *DocumentService) GetDownloadURL(ctx context.Context, docID string) (*dto.DocumentUploadResponse, error) {
	doc, err := s.docRepo.FindByID(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	url, err := s.minio.GetPresignedURL(ctx, doc.S3ObjectName, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return &dto.DocumentUploadResponse{
		ID:           doc.ID,
		PresignedURL: url.String(),
	}, nil
}
