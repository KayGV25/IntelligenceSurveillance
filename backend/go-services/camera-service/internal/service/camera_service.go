package service

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/domain"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/dto"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/event"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/repository"
)

type CameraService struct {
	cameraRepo *repository.CameraRepository
	publisher  event.Publisher
}

func NewCameraService(
	cameraRepo *repository.CameraRepository,
	publisher event.Publisher,
) *CameraService {
	return &CameraService{
		cameraRepo: cameraRepo,
		publisher:  publisher,
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
