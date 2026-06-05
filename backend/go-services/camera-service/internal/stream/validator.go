package stream

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/domain"
)

type ValidationResult struct {
	Status     domain.CameraStatus
	CheckedURL string
	Message    string
}

type Validator struct {
	timeout time.Duration
}

func NewValidator(timeout time.Duration) *Validator {
	return &Validator{timeout: timeout}
}

func (v *Validator) Validate(contract *domain.CameraConnectionContract) ValidationResult {
	switch contract.ConnectionType {
	case domain.ConnectionTypeHTTPMJPEG:
		return v.validateHTTPMJPEG(contract)

	case domain.ConnectionTypeHTTPSnapshot:
		return v.validateHTTPSnapshot(contract)

	case domain.ConnectionTypeRTSPManual, domain.ConnectionTypeONVIF:
		return v.validateRTSP(contract)

	default:
		return ValidationResult{
			Status:  domain.CameraStatusUnknown,
			Message: "Unsupported connection type",
		}
	}
}

func (v *Validator) validateHTTPMJPEG(contract *domain.CameraConnectionContract) ValidationResult {
	if contract.MainStreamURL == nil {
		return ValidationResult{
			Status:  domain.CameraStatusOffline,
			Message: "Main stream URL is missing",
		}
	}

	url := *contract.MainStreamURL

	client := http.Client{
		Timeout: v.timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return ValidationResult{
			Status:     domain.CameraStatusOffline,
			CheckedURL: url,
			Message:    err.Error(),
		}
	}
	defer resp.Body.Close()

	contentType := strings.ToLower(resp.Header.Get("Content-Type"))

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if strings.Contains(contentType, "multipart") ||
			strings.Contains(contentType, "image/jpeg") ||
			strings.Contains(contentType, "video") ||
			contentType == "" {
			return ValidationResult{
				Status:     domain.CameraStatusOnline,
				CheckedURL: url,
				Message:    "HTTP MJPEG stream is reachable",
			}
		}
	}

	return ValidationResult{
		Status:     domain.CameraStatusDegraded,
		CheckedURL: url,
		Message:    fmt.Sprintf("Unexpected response: status=%d content_type=%s", resp.StatusCode, contentType),
	}
}

func (v *Validator) validateHTTPSnapshot(contract *domain.CameraConnectionContract) ValidationResult {
	if contract.SnapshotURL == nil {
		return ValidationResult{
			Status:  domain.CameraStatusOffline,
			Message: "Snapshot URL is missing",
		}
	}

	url := *contract.SnapshotURL

	client := http.Client{
		Timeout: v.timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return ValidationResult{
			Status:     domain.CameraStatusOffline,
			CheckedURL: url,
			Message:    err.Error(),
		}
	}
	defer resp.Body.Close()

	contentType := strings.ToLower(resp.Header.Get("Content-Type"))

	if resp.StatusCode >= 200 && resp.StatusCode < 300 &&
		(strings.Contains(contentType, "image/jpeg") || strings.Contains(contentType, "image/png")) {
		return ValidationResult{
			Status:     domain.CameraStatusOnline,
			CheckedURL: url,
			Message:    "Snapshot endpoint is reachable",
		}
	}

	return ValidationResult{
		Status:     domain.CameraStatusDegraded,
		CheckedURL: url,
		Message:    fmt.Sprintf("Unexpected response: status=%d content_type=%s", resp.StatusCode, contentType),
	}
}

func (v *Validator) validateRTSP(contract *domain.CameraConnectionContract) ValidationResult {
	var url string

	if contract.MainStreamURL != nil {
		url = *contract.MainStreamURL
	} else if contract.RTSPUrl != nil {
		url = *contract.RTSPUrl
	} else {
		return ValidationResult{
			Status:  domain.CameraStatusOffline,
			Message: "RTSP URL is missing",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), v.timeout)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"ffprobe",
		"-v", "error",
		"-rtsp_transport", "tcp",
		"-i", url,
	)

	if err := cmd.Run(); err != nil {
		return ValidationResult{
			Status:     domain.CameraStatusOffline,
			CheckedURL: url,
			Message:    err.Error(),
		}
	}

	return ValidationResult{
		Status:     domain.CameraStatusOnline,
		CheckedURL: url,
		Message:    "RTSP stream is reachable",
	}
}
