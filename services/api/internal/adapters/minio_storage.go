package adapters

import (
	"context"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorageAdapter struct {
	client     *minio.Client
	bucketName string
}

// GetReadLink generates a signed URL that allows temporary public access to a private object.
func (a *MinioStorageAdapter) GetReadLink(ctx context.Context, key string, expires time.Duration) (string, error) {
	// We can also set response headers here (e.g., forcing a specific filename)
	reqParams := make(url.Values)

	u, err := a.client.PresignedGetObject(ctx, a.bucketName, key, expires, reqParams)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

// NewMinioStorageAdapter initializes the client and ensures the bucket exists.
func NewMinioStorageAdapter(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinioStorageAdapter, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinioStorageAdapter{
		client:     minioClient,
		bucketName: bucket,
	}, nil
}

// GenerateUploadLink creates a presigned URL for a PUT request.
func (a *MinioStorageAdapter) GenerateUploadLink(ctx context.Context, key string, maxSize int64, expires time.Duration) (string, error) {
	// Set policy constraints (like max file size)
	//reqParams := make(url.Values)

	// MinIO/S3 uses PresignedPutObject for direct uploads
	u, err := a.client.PresignedPutObject(ctx, a.bucketName, key, expires)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

// DeleteObject removes the file from the bucket.
func (a *MinioStorageAdapter) DeleteObject(ctx context.Context, key string) error {
	return a.client.RemoveObject(ctx, a.bucketName, key, minio.RemoveObjectOptions{})
}
