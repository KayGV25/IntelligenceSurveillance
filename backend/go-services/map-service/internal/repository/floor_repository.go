package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/domain"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/dto"
)

var ErrFloorNotFound = errors.New("floor not found")

type FloorRepository struct {
	db *pgxpool.Pool
}

func NewFloorRepository(db *pgxpool.Pool) *FloorRepository {
	return &FloorRepository{db: db}
}

func (r *FloorRepository) Create(ctx context.Context, req dto.CreateFloorRequest, createdBy *uuid.UUID) (*domain.Floor, error) {
	id := uuid.New()

	unit := "meters"
	if req.Unit != nil && *req.Unit != "" {
		unit = *req.Unit
	}

	color := "#64748B"
	if req.Color != nil && *req.Color != "" {
		color = *req.Color
	}

	opacity := 1.0
	if req.Opacity != nil {
		opacity = *req.Opacity
	}

	labelVisible := true
	if req.LabelVisible != nil {
		labelVisible = *req.LabelVisible
	}

	query := `
		INSERT INTO map.floors (
			id, building_id, name, description, floor_number,
			width, height, unit, color, opacity, icon, label_visible,
			created_by, updated_by
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		RETURNING
			id, building_id, name, description, floor_number,
			width, height, unit, color, opacity, icon, label_visible,
			created_by, updated_by, created_at, updated_at
	`

	floor := &domain.Floor{}

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		req.BuildingID,
		req.Name,
		req.Description,
		req.FloorNumber,
		req.Width,
		req.Height,
		unit,
		color,
		opacity,
		req.Icon,
		labelVisible,
		createdBy,
		createdBy,
	).Scan(
		&floor.ID,
		&floor.BuildingID,
		&floor.Name,
		&floor.Description,
		&floor.FloorNumber,
		&floor.Width,
		&floor.Height,
		&floor.Unit,
		&floor.Color,
		&floor.Opacity,
		&floor.Icon,
		&floor.LabelVisible,
		&floor.CreatedBy,
		&floor.UpdatedBy,
		&floor.CreatedAt,
		&floor.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return floor, nil
}

func (r *FloorRepository) FindAll(ctx context.Context) ([]domain.Floor, error) {
	query := `
		SELECT
			id, building_id, name, description, floor_number,
			width, height, unit, color, opacity, icon, label_visible,
			created_by, updated_by, created_at, updated_at
		FROM map.floors
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	floors := make([]domain.Floor, 0)

	for rows.Next() {
		var floor domain.Floor

		err := rows.Scan(
			&floor.ID,
			&floor.BuildingID,
			&floor.Name,
			&floor.Description,
			&floor.FloorNumber,
			&floor.Width,
			&floor.Height,
			&floor.Unit,
			&floor.Color,
			&floor.Opacity,
			&floor.Icon,
			&floor.LabelVisible,
			&floor.CreatedBy,
			&floor.UpdatedBy,
			&floor.CreatedAt,
			&floor.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		floors = append(floors, floor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return floors, nil
}

func (r *FloorRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Floor, error) {
	query := `
		SELECT
			id, building_id, name, description, floor_number,
			width, height, unit, color, opacity, icon, label_visible,
			created_by, updated_by, created_at, updated_at
		FROM map.floors
		WHERE id = $1
		AND deleted_at IS NULL
	`

	floor := &domain.Floor{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&floor.ID,
		&floor.BuildingID,
		&floor.Name,
		&floor.Description,
		&floor.FloorNumber,
		&floor.Width,
		&floor.Height,
		&floor.Unit,
		&floor.Color,
		&floor.Opacity,
		&floor.Icon,
		&floor.LabelVisible,
		&floor.CreatedBy,
		&floor.UpdatedBy,
		&floor.CreatedAt,
		&floor.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrFloorNotFound
	}

	if err != nil {
		return nil, err
	}

	return floor, nil
}

func (r *FloorRepository) Update(ctx context.Context, id uuid.UUID, req dto.UpdateFloorRequest, updatedBy *uuid.UUID) (*domain.Floor, error) {
	existing, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.BuildingID != nil {
		existing.BuildingID = *req.BuildingID
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.FloorNumber != nil {
		existing.FloorNumber = req.FloorNumber
	}
	if req.Width != nil {
		existing.Width = req.Width
	}
	if req.Height != nil {
		existing.Height = req.Height
	}
	if req.Unit != nil {
		existing.Unit = *req.Unit
	}
	if req.Color != nil {
		existing.Color = *req.Color
	}
	if req.Opacity != nil {
		existing.Opacity = *req.Opacity
	}
	if req.Icon != nil {
		existing.Icon = req.Icon
	}
	if req.LabelVisible != nil {
		existing.LabelVisible = *req.LabelVisible
	}

	query := `
		UPDATE map.floors
		SET
			building_id = $2,
			name = $3,
			description = $4,
			floor_number = $5,
			width = $6,
			height = $7,
			unit = $8,
			color = $9,
			opacity = $10,
			icon = $11,
			label_visible = $12,
			updated_by = $13,
			updated_at = now()
		WHERE id = $1
		AND deleted_at IS NULL
		RETURNING
			id, building_id, name, description, floor_number,
			width, height, unit, color, opacity, icon, label_visible,
			created_by, updated_by, created_at, updated_at
	`

	floor := &domain.Floor{}

	err = r.db.QueryRow(
		ctx,
		query,
		id,
		existing.BuildingID,
		existing.Name,
		existing.Description,
		existing.FloorNumber,
		existing.Width,
		existing.Height,
		existing.Unit,
		existing.Color,
		existing.Opacity,
		existing.Icon,
		existing.LabelVisible,
		updatedBy,
	).Scan(
		&floor.ID,
		&floor.BuildingID,
		&floor.Name,
		&floor.Description,
		&floor.FloorNumber,
		&floor.Width,
		&floor.Height,
		&floor.Unit,
		&floor.Color,
		&floor.Opacity,
		&floor.Icon,
		&floor.LabelVisible,
		&floor.CreatedBy,
		&floor.UpdatedBy,
		&floor.CreatedAt,
		&floor.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return floor, nil
}

func (r *FloorRepository) Delete(ctx context.Context, id uuid.UUID, updatedBy *uuid.UUID) error {
	query := `
		UPDATE map.floors
		SET
			deleted_at = now(),
			updated_by = $2,
			updated_at = now()
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := r.db.Exec(ctx, query, id, updatedBy)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrFloorNotFound
	}

	return nil
}
