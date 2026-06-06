package snapshot

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/domain"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/storage"
)

type Service struct {
	storageClient *storage.Client
	timeout       time.Duration
}

func NewService(storageClient *storage.Client, timeout time.Duration) *Service {
	return &Service{
		storageClient: storageClient,
		timeout:       timeout,
	}
}

func (s *Service) Capture(
	ctx context.Context,
	cameraID uuid.UUID,
	contract *domain.CameraConnectionContract,
) (objectPath string, contentType string, err error) {
	if contract.SnapshotURL == nil {
		return "", "", fmt.Errorf("snapshot URL is missing")
	}

	client := http.Client{
		Timeout: s.timeout,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, *contract.SnapshotURL, nil)
	if err != nil {
		return "", "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", "", fmt.Errorf("snapshot request failed with status %d", resp.StatusCode)
	}

	contentType = resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	if !strings.Contains(strings.ToLower(contentType), "image") {
		return "", "", fmt.Errorf("snapshot endpoint did not return image content type: %s", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	ext := "jpg"
	if strings.Contains(strings.ToLower(contentType), "png") {
		ext = "png"
	}

	objectName := fmt.Sprintf(
		"snapshots/%s/%d.%s",
		cameraID.String(),
		time.Now().UTC().UnixMilli(),
		ext,
	)

	return s.storageClient.UploadSnapshot(
		ctx,
		objectName,
		bytes.NewReader(body),
		int64(len(body)),
		contentType,
	)
}
