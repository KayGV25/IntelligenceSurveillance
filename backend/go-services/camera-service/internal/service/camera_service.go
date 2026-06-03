package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/domain"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/dto"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/repository"
)

type CameraService struct {
	cameraRepo *repository.CameraRepository
}

func NewCameraService(cameraRepo *repository.CameraRepository) *CameraService {
	return &CameraService{
		cameraRepo: cameraRepo,
	}
}

func (s *CameraService) Create(
	ctx context.Context,
	req dto.CreateCameraRequest,
	createdBy *uuid.UUID,
) (*domain.Camera, error) {
	return s.cameraRepo.Create(ctx, req, createdBy)
}

func (s *CameraService) GetAll(ctx context.Context) ([]domain.Camera, error) {
	return s.cameraRepo.FindAll(ctx)
}

func (s *CameraService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Camera, error) {
	return s.cameraRepo.FindByID(ctx, id)
}

func (s *CameraService) Delete(ctx context.Context, id uuid.UUID, updatedBy *uuid.UUID) error {
	return s.cameraRepo.Delete(ctx, id, updatedBy)
}

func (s *CameraService) Update(
	ctx context.Context,
	id uuid.UUID,
	req dto.UpdateCameraRequest,
	updatedBy *uuid.UUID,
) (*domain.Camera, error) {
	return s.cameraRepo.Update(ctx, id, req, updatedBy)
}
