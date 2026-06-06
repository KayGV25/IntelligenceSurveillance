package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/domain"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/dto"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/event"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/repository"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/snapshot"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/stream"
)

type CameraService struct {
	cameraRepo      *repository.CameraRepository
	deviceRepo      *repository.DiscoveredDeviceRepository
	contractRepo    *repository.ConnectionContractRepository
	publisher       event.Publisher
	streamValidator *stream.Validator
	snapshotService *snapshot.Service
}

func NewCameraService(
	cameraRepo *repository.CameraRepository,
	publisher event.Publisher,
	deviceRepo *repository.DiscoveredDeviceRepository,
	contractRepo *repository.ConnectionContractRepository,
	streamValidator *stream.Validator,
	snapshotService *snapshot.Service,
) *CameraService {
	return &CameraService{
		cameraRepo:      cameraRepo,
		publisher:       publisher,
		deviceRepo:      deviceRepo,
		contractRepo:    contractRepo,
		streamValidator: streamValidator,
		snapshotService: snapshotService,
	}
}

func (s *CameraService) Create(
	ctx context.Context,
	req dto.CreateCameraRequest,
	createdBy *uuid.UUID,
) (*domain.Camera, error) {
	camera, err := s.cameraRepo.Create(ctx, req, createdBy)
	if err != nil {
		return nil, err
	}

	if s.publisher != nil {
		if err := s.publisher.PublishCameraEvent(ctx, event.CameraEvent{
			EventID:   uuid.New(),
			EventType: event.CameraCreatedEvent,
			CameraID:  &camera.ID,
			UserID:    createdBy,
			Timestamp: time.Now().UTC(),
		}); err != nil {
			log.Printf("failed to publish camera.created event: %v", err)
		}
	}

	return camera, nil
}

func (s *CameraService) GetAll(ctx context.Context) ([]domain.Camera, error) {
	return s.cameraRepo.FindAll(ctx)
}

func (s *CameraService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Camera, error) {
	return s.cameraRepo.FindByID(ctx, id)
}

func (s *CameraService) Delete(ctx context.Context, id uuid.UUID, updatedBy *uuid.UUID) error {
	if err := s.cameraRepo.Delete(ctx, id, updatedBy); err != nil {
		return err
	}

	if s.publisher != nil {
		if err := s.publisher.PublishCameraEvent(ctx, event.CameraEvent{
			EventID:   uuid.New(),
			EventType: event.CameraDeletedEvent,
			CameraID:  &id,
			UserID:    updatedBy,
			Timestamp: time.Now().UTC(),
		}); err != nil {
			log.Printf("failed to publish camera.deleted event: %v", err)
		}
	}

	return nil
}

func (s *CameraService) Update(
	ctx context.Context,
	id uuid.UUID,
	req dto.UpdateCameraRequest,
	updatedBy *uuid.UUID,
) (*domain.Camera, error) {
	camera, err := s.cameraRepo.Update(ctx, id, req, updatedBy)
	if err != nil {
		return nil, err
	}

	if s.publisher != nil {
		if err := s.publisher.PublishCameraEvent(ctx, event.CameraEvent{
			EventID:   uuid.New(),
			EventType: event.CameraUpdatedEvent,
			CameraID:  &camera.ID,
			UserID:    updatedBy,
			Timestamp: time.Now().UTC(),
		}); err != nil {
			log.Printf("failed to publish camera.updated event: %v", err)
		}
	}

	return camera, nil
}

func (s *CameraService) ConnectDiscoveredDevice(
	ctx context.Context,
	discoveredDeviceID uuid.UUID,
	req dto.ConnectDiscoveredDeviceRequest,
	userID *uuid.UUID,
) (*domain.Camera, *domain.CameraConnectionContract, error) {
	device, err := s.deviceRepo.FindByID(ctx, discoveredDeviceID)
	if err != nil {
		return nil, nil, err
	}

	existingContract, err := s.contractRepo.FindByDiscoveredDeviceID(ctx, discoveredDeviceID)
	if err != nil {
		return nil, nil, err
	}

	if existingContract != nil {
		existingCamera, err := s.cameraRepo.FindByID(ctx, existingContract.CameraID)
		if err != nil {
			return nil, nil, err
		}

		return existingCamera, existingContract, nil
	}

	var rtspURL *string
	var mainStreamURL *string
	var snapshotURL *string

	if device.RTSPSupported && device.RTSPPort != nil {
		url := "rtsp://" + device.IPAddress + ":554"
		if *device.RTSPPort != 554 {
			url = "rtsp://" + device.IPAddress + ":" + fmt.Sprintf("%d", *device.RTSPPort)
		}
		rtspURL = &url
		mainStreamURL = &url
	}

	if device.HTTPSupported && device.HTTPPort != nil {
		baseURL := "http://" + device.IPAddress + ":" + fmt.Sprintf("%d", *device.HTTPPort)

		videoURL := baseURL + "/video"
		snapURL := baseURL + "/shot.jpg"

		mainStreamURL = &videoURL
		snapshotURL = &snapURL
	}

	cameraReq := dto.CreateCameraRequest{
		Name:        req.Name,
		Description: req.Description,
		RTSPUrl:     rtspURL,
	}

	camera, err := s.cameraRepo.Create(ctx, cameraReq, userID)
	if err != nil {
		return nil, nil, err
	}

	connectionType := domain.ConnectionTypeRTSPManual
	if device.ONVIFSupported {
		connectionType = domain.ConnectionTypeONVIF
	} else if device.RTSPSupported {
		connectionType = domain.ConnectionTypeRTSPManual
	} else if device.HTTPSupported {
		connectionType = domain.ConnectionTypeHTTPMJPEG
	}

	contract, err := s.contractRepo.Create(ctx, domain.CameraConnectionContract{
		CameraID:           camera.ID,
		DiscoveredDeviceID: &device.ID,
		ConnectionType:     connectionType,
		IPAddress:          device.IPAddress,
		RTSPUrl:            rtspURL,
		MainStreamURL:      mainStreamURL,
		SnapshotURL:        snapshotURL,
		Username:           req.Username,
		EncryptedPassword:  req.Password,
		CreatedBy:          userID,
		UpdatedBy:          userID,
	})
	if err != nil {
		return nil, nil, err
	}

	if s.publisher != nil {
		if err := s.publisher.PublishCameraEvent(ctx, event.CameraEvent{
			EventID:   uuid.New(),
			EventType: event.CameraCreatedEvent,
			CameraID:  &camera.ID,
			UserID:    userID,
			Timestamp: time.Now().UTC(),
		}); err != nil {
			log.Printf("failed to publish camera.created event: %v", err)
		}

		if err := s.publisher.PublishCameraEvent(ctx, event.CameraEvent{
			EventID:            uuid.New(),
			EventType:          event.CameraConnectedEvent,
			CameraID:           &camera.ID,
			DiscoveredDeviceID: &device.ID,
			IPAddress:          device.IPAddress,
			UserID:             userID,
			Timestamp:          time.Now().UTC(),
		}); err != nil {
			log.Printf("failed to publish camera.connected event: %v", err)
		}
	}

	return camera, contract, nil
}

func (s *CameraService) GetConnectionByCameraID(
	ctx context.Context,
	cameraID uuid.UUID,
) (*domain.CameraConnectionContract, error) {
	return s.contractRepo.FindByCameraID(ctx, cameraID)
}

func (s *CameraService) ValidateStream(
	ctx context.Context,
	cameraID uuid.UUID,
	userID *uuid.UUID,
) (*dto.ValidateStreamResponse, error) {
	camera, err := s.cameraRepo.FindByID(ctx, cameraID)
	if err != nil {
		return nil, err
	}

	contract, err := s.contractRepo.FindByCameraID(ctx, camera.ID)
	if err != nil {
		return nil, err
	}

	result := s.streamValidator.Validate(contract)

	if err := s.cameraRepo.UpdateStatus(ctx, camera.ID, result.Status, userID); err != nil {
		return nil, err
	}

	if s.publisher != nil {
		eventType := event.CameraOfflineEvent
		if result.Status == domain.CameraStatusOnline {
			eventType = event.CameraOnlineEvent
		}

		_ = s.publisher.PublishCameraEvent(ctx, event.CameraEvent{
			EventID:   uuid.New(),
			EventType: eventType,
			CameraID:  &camera.ID,
			UserID:    userID,
			Timestamp: time.Now().UTC(),
		})
	}

	return &dto.ValidateStreamResponse{
		CameraID:       camera.ID.String(),
		Status:         string(result.Status),
		ConnectionType: string(contract.ConnectionType),
		CheckedURL:     result.CheckedURL,
		Message:        result.Message,
	}, nil
}

func (s *CameraService) CaptureSnapshot(
	ctx context.Context,
	cameraID uuid.UUID,
) (*dto.SnapshotResponse, error) {
	camera, err := s.cameraRepo.FindByID(ctx, cameraID)
	if err != nil {
		return nil, err
	}

	contract, err := s.contractRepo.FindByCameraID(ctx, camera.ID)
	if err != nil {
		return nil, err
	}

	objectPath, contentType, err := s.snapshotService.Capture(ctx, camera.ID, contract)
	if err != nil {
		return nil, err
	}

	return &dto.SnapshotResponse{
		CameraID:    camera.ID.String(),
		ObjectPath:  objectPath,
		ContentType: contentType,
		Message:     "Snapshot captured successfully",
	}, nil
}

func (s *CameraService) GetHealth(
	ctx context.Context,
	cameraID uuid.UUID,
) (*dto.CameraHealthResponse, error) {
	camera, err := s.cameraRepo.FindByID(ctx, cameraID)
	if err != nil {
		return nil, err
	}

	contract, err := s.contractRepo.FindByCameraID(ctx, camera.ID)
	if err != nil {
		return nil, err
	}

	message := "Camera status is unknown"
	switch camera.Status {
	case domain.CameraStatusOnline:
		message = "Camera is online"
	case domain.CameraStatusOffline:
		message = "Camera is offline"
	case domain.CameraStatusDegraded:
		message = "Camera is degraded"
	}

	return &dto.CameraHealthResponse{
		CameraID:       camera.ID.String(),
		Status:         string(camera.Status),
		ConnectionType: string(contract.ConnectionType),
		MainStreamURL:  contract.MainStreamURL,
		SnapshotURL:    contract.SnapshotURL,
		Message:        message,
	}, nil
}

func (s *CameraService) GetStreamInfo(
	ctx context.Context,
	cameraID uuid.UUID,
) (*dto.StreamInfoResponse, error) {
	camera, err := s.cameraRepo.FindByID(ctx, cameraID)
	if err != nil {
		return nil, err
	}

	contract, err := s.contractRepo.FindByCameraID(ctx, camera.ID)
	if err != nil {
		return nil, err
	}

	var streamMode string
	var streamURL *string
	requiresProxy := false
	message := "Stream information available"

	switch contract.ConnectionType {
	case domain.ConnectionTypeHTTPMJPEG:
		streamMode = "MJPEG"
		streamURL = contract.MainStreamURL
		requiresProxy = false

	case domain.ConnectionTypeHTTPSnapshot:
		streamMode = "SNAPSHOT"
		streamURL = contract.SnapshotURL
		requiresProxy = false

	case domain.ConnectionTypeRTSPManual:
		streamMode = "RTSP"
		if contract.MainStreamURL != nil {
			streamURL = contract.MainStreamURL
		} else {
			streamURL = contract.RTSPUrl
		}
		requiresProxy = true
		message = "RTSP stream requires backend proxy/transcoding for browser playback"

	case domain.ConnectionTypeONVIF:
		streamMode = "ONVIF_RTSP"
		if contract.MainStreamURL != nil {
			streamURL = contract.MainStreamURL
		} else {
			streamURL = contract.RTSPUrl
		}
		requiresProxy = true
		message = "ONVIF/RTSP stream requires backend proxy/transcoding for browser playback"

	default:
		streamMode = "UNKNOWN"
		message = "Unsupported or unknown stream type"
	}

	return &dto.StreamInfoResponse{
		CameraID:       camera.ID.String(),
		ConnectionType: string(contract.ConnectionType),
		StreamMode:     streamMode,
		StreamURL:      streamURL,
		SnapshotURL:    contract.SnapshotURL,
		RequiresProxy:  requiresProxy,
		Message:        message,
	}, nil
}
