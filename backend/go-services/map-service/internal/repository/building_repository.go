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

var ErrBuildingNotFound = errors.New("building not found")

type BuildingRepository struct {
	db *pgxpool.Pool
}

func NewBuildingRepository(db *pgxpool.Pool) *BuildingRepository {
	return &BuildingRepository{db: db}
}

func (r *BuildingRepository) Create(
	ctx context.Context,
	req dto.CreateBuildingRequest,
	createdBy *uuid.UUID,
) (*domain.Building, error) {
	id := uuid.New()

	color := "#2563EB"
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
		INSERT INTO map.buildings (
			id,
			name,
			description,
			address,
			color,
			opacity,
			icon,
			label_visible,
			created_by,
			updated_by
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING
			id,
			name,
			description,
			address,
			color,
			opacity,
			icon,
			label_visible,
			created_by,
			updated_by,
			created_at,
			updated_at
	`

	building := &domain.Building{}

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		req.Name,
		req.Description,
		req.Address,
		color,
		opacity,
		req.Icon,
		labelVisible,
		createdBy,
		createdBy,
	).Scan(
		&building.ID,
		&building.Name,
		&building.Description,
		&building.Address,
		&building.Color,
		&building.Opacity,
		&building.Icon,
		&building.LabelVisible,
		&building.CreatedBy,
		&building.UpdatedBy,
		&building.CreatedAt,
		&building.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return building, nil
}

func (r *BuildingRepository) FindAll(ctx context.Context) ([]domain.Building, error) {
	query := `
		SELECT
			id,
			name,
			description,
			address,
			color,
			opacity,
			icon,
			label_visible,
			created_by,
			updated_by,
			created_at,
			updated_at
		FROM map.buildings
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	buildings := make([]domain.Building, 0)

	for rows.Next() {
		var building domain.Building

		err := rows.Scan(
			&building.ID,
			&building.Name,
			&building.Description,
			&building.Address,
			&building.Color,
			&building.Opacity,
			&building.Icon,
			&building.LabelVisible,
			&building.CreatedBy,
			&building.UpdatedBy,
			&building.CreatedAt,
			&building.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		buildings = append(buildings, building)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return buildings, nil
}

func (r *BuildingRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Building, error) {
	query := `
		SELECT
			id,
			name,
			description,
			address,
			color,
			opacity,
			icon,
			label_visible,
			created_by,
			updated_by,
			created_at,
			updated_at
		FROM map.buildings
		WHERE id = $1
		AND deleted_at IS NULL
	`

	building := &domain.Building{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&building.ID,
		&building.Name,
		&building.Description,
		&building.Address,
		&building.Color,
		&building.Opacity,
		&building.Icon,
		&building.LabelVisible,
		&building.CreatedBy,
		&building.UpdatedBy,
		&building.CreatedAt,
		&building.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrBuildingNotFound
	}

	if err != nil {
		return nil, err
	}

	return building, nil
}

func (r *BuildingRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	req dto.UpdateBuildingRequest,
	updatedBy *uuid.UUID,
) (*domain.Building, error) {
	existing, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.Address != nil {
		existing.Address = req.Address
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
		UPDATE map.buildings
		SET
			name = $2,
			description = $3,
			address = $4,
			color = $5,
			opacity = $6,
			icon = $7,
			label_visible = $8,
			updated_by = $9,
			updated_at = now()
		WHERE id = $1
		AND deleted_at IS NULL
		RETURNING
			id,
			name,
			description,
			address,
			color,
			opacity,
			icon,
			label_visible,
			created_by,
			updated_by,
			created_at,
			updated_at
	`

	building := &domain.Building{}

	err = r.db.QueryRow(
		ctx,
		query,
		id,
		existing.Name,
		existing.Description,
		existing.Address,
		existing.Color,
		existing.Opacity,
		existing.Icon,
		existing.LabelVisible,
		updatedBy,
	).Scan(
		&building.ID,
		&building.Name,
		&building.Description,
		&building.Address,
		&building.Color,
		&building.Opacity,
		&building.Icon,
		&building.LabelVisible,
		&building.CreatedBy,
		&building.UpdatedBy,
		&building.CreatedAt,
		&building.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return building, nil
}

func (r *BuildingRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
	updatedBy *uuid.UUID,
) error {
	query := `
		UPDATE map.buildings
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
		return ErrBuildingNotFound
	}

	return nil
}
