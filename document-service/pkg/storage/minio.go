package storage

import (
	"context"
	"io"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioService struct {
	client *minio.Client
	bucket string
}

func NewMinioService(endpoint, accessKey, secretKey string, useSSL bool, bucket string) (*MinioService, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		log.Printf("Bucket %s does not exist, creating it...", bucket)
		err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &MinioService{
		client: minioClient,
		bucket: bucket,
	}, nil
}

func (s *MinioService) UploadFile(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (minio.UploadInfo, error) {
	info, err := s.client.PutObject(ctx, s.bucket, objectName, reader, objectSize, minio.PutObjectOptions{ContentType: contentType})
	return info, err
}

func (s *MinioService) GetFileStream(ctx context.Context, objectName string) (*minio.Object, error) {
	return s.client.GetObject(ctx, s.bucket, objectName, minio.GetObjectOptions{})
}

func (s *MinioService) GetPresignedURL(ctx context.Context, objectName string, expires time.Duration) (*url.URL, error) {
	reqParams := make(url.Values)
	return s.client.PresignedGetObject(ctx, s.bucket, objectName, expires, reqParams)
}

func (s *MinioService) DeleteFile(ctx context.Context, objectName string) error {
	return s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
}
