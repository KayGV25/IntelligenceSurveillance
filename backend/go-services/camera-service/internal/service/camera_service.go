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
)

type CameraService struct {
	cameraRepo   *repository.CameraRepository
	deviceRepo   *repository.DiscoveredDeviceRepository
	contractRepo *repository.ConnectionContractRepository
	publisher    event.Publisher
}

func NewCameraService(
	cameraRepo *repository.CameraRepository,
	publisher event.Publisher,
	deviceRepo *repository.DiscoveredDeviceRepository,
	contractRepo *repository.ConnectionContractRepository,
) *CameraService {
	return &CameraService{
		cameraRepo:   cameraRepo,
		publisher:    publisher,
		deviceRepo:   deviceRepo,
		contractRepo: contractRepo,
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
			CameraID:  camera.ID,
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
			EventType: event.CameraCreatedEvent,
			CameraID:  id,
			UserID:    updatedBy,
			Timestamp: time.Now().UTC(),
		}); err != nil {
			log.Printf("failed to publish camera.created event: %v", err)
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
			EventType: event.CameraCreatedEvent,
			CameraID:  camera.ID,
			UserID:    updatedBy,
			Timestamp: time.Now().UTC(),
		}); err != nil {
			log.Printf("failed to publish camera.created event: %v", err)
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
		_ = s.publisher.PublishCameraEvent(ctx, event.CameraEvent{
			EventID:   uuid.New(),
			EventType: event.CameraCreatedEvent,
			CameraID:  camera.ID,
			UserID:    userID,
			Timestamp: time.Now().UTC(),
		})
	}

	return camera, contract, nil
}

func (s *CameraService) GetConnectionByCameraID(
	ctx context.Context,
	cameraID uuid.UUID,
) (*domain.CameraConnectionContract, error) {
	return s.contractRepo.FindByCameraID(ctx, cameraID)
}
