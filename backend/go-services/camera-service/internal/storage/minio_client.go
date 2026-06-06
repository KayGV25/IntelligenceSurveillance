package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/config"
)

type Client struct {
	minioClient *minio.Client
	bucket      string
}

func NewClient(cfg *config.Config) (*Client, error) {
	client, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: cfg.MinIOUseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		minioClient: client,
		bucket:      cfg.MinIOSnapshotsBucket,
	}, nil
}

func (c *Client) UploadSnapshot(
	ctx context.Context,
	objectName string,
	reader io.Reader,
	size int64,
	contentType string,
) (string, string, error) {
	_, err := c.minioClient.PutObject(
		ctx,
		c.bucket,
		objectName,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", "", err
	}

	return objectName, contentType, nil
}
