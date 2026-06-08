package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/domain"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/dto"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/repository"
)

type FloorService struct {
	floorRepo *repository.FloorRepository
}

func NewFloorService(floorRepo *repository.FloorRepository) *FloorService {
	return &FloorService{floorRepo: floorRepo}
}

func (s *FloorService) Create(ctx context.Context, req dto.CreateFloorRequest, createdBy *uuid.UUID) (*domain.Floor, error) {
	return s.floorRepo.Create(ctx, req, createdBy)
}

func (s *FloorService) GetAll(ctx context.Context) ([]domain.Floor, error) {
	return s.floorRepo.FindAll(ctx)
}

func (s *FloorService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Floor, error) {
	return s.floorRepo.FindByID(ctx, id)
}

func (s *FloorService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateFloorRequest, updatedBy *uuid.UUID) (*domain.Floor, error) {
	return s.floorRepo.Update(ctx, id, req, updatedBy)
}

func (s *FloorService) Delete(ctx context.Context, id uuid.UUID, updatedBy *uuid.UUID) error {
	return s.floorRepo.Delete(ctx, id, updatedBy)
}
