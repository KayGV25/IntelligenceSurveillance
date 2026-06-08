package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/domain"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/dto"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/repository"
)

type ZoneService struct {
	zoneRepo *repository.ZoneRepository
}

func NewZoneService(zoneRepo *repository.ZoneRepository) *ZoneService {
	return &ZoneService{zoneRepo: zoneRepo}
}

func (s *ZoneService) Create(ctx context.Context, req dto.CreateZoneRequest, createdBy *uuid.UUID) (*domain.Zone, error) {
	return s.zoneRepo.Create(ctx, req, createdBy)
}

func (s *ZoneService) GetAll(ctx context.Context, floorID *uuid.UUID) ([]domain.Zone, error) {
	return s.zoneRepo.FindAll(ctx, floorID)
}

func (s *ZoneService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Zone, error) {
	return s.zoneRepo.FindByID(ctx, id)
}

func (s *ZoneService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateZoneRequest, updatedBy *uuid.UUID) (*domain.Zone, error) {
	return s.zoneRepo.Update(ctx, id, req, updatedBy)
}

func (s *ZoneService) Delete(ctx context.Context, id uuid.UUID, updatedBy *uuid.UUID) error {
	return s.zoneRepo.Delete(ctx, id, updatedBy)
}
