package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/domain"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/dto"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/repository"
)

type BuildingService struct {
	buildingRepo *repository.BuildingRepository
}

func NewBuildingService(buildingRepo *repository.BuildingRepository) *BuildingService {
	return &BuildingService{buildingRepo: buildingRepo}
}

func (s *BuildingService) Create(
	ctx context.Context,
	req dto.CreateBuildingRequest,
	createdBy *uuid.UUID,
) (*domain.Building, error) {
	return s.buildingRepo.Create(ctx, req, createdBy)
}

func (s *BuildingService) GetAll(ctx context.Context) ([]domain.Building, error) {
	return s.buildingRepo.FindAll(ctx)
}

func (s *BuildingService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Building, error) {
	return s.buildingRepo.FindByID(ctx, id)
}

func (s *BuildingService) Update(
	ctx context.Context,
	id uuid.UUID,
	req dto.UpdateBuildingRequest,
	updatedBy *uuid.UUID,
) (*domain.Building, error) {
	return s.buildingRepo.Update(ctx, id, req, updatedBy)
}

func (s *BuildingService) Delete(
	ctx context.Context,
	id uuid.UUID,
	updatedBy *uuid.UUID,
) error {
	return s.buildingRepo.Delete(ctx, id, updatedBy)
}
