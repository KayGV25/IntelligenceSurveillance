package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/domain"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/dto"
)

var ErrCameraNotFound = errors.New("camera not found")

type CameraRepository struct {
	db *pgxpool.Pool
}

func NewCameraRepository(db *pgxpool.Pool) *CameraRepository {
	return &CameraRepository{
		db: db,
	}
}

func (r *CameraRepository) Create(
	ctx context.Context,
	req dto.CreateCameraRequest,
	createdBy *uuid.UUID,
) (*domain.Camera, error) {
	id := uuid.New()

	query := `
		INSERT INTO camera.cameras (
			id,
			name,
			description,
			rtsp_url,
			status,
			latitude,
			longitude,
			building_id,
			floor_id,
			zone_id,
			position_x,
			position_y,
			direction_angle,
			fov_angle,
			created_by,
			updated_by
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16
		)
		RETURNING
			id,
			name,
			description,
			rtsp_url,
			status,
			latitude,
			longitude,
			building_id,
			floor_id,
			zone_id,
			position_x,
			position_y,
			direction_angle,
			fov_angle,
			created_by,
			updated_by,
			created_at,
			updated_at
	`

	camera := &domain.Camera{}

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		req.Name,
		req.Description,
		req.RTSPUrl,
		domain.CameraStatusUnknown,
		req.Latitude,
		req.Longitude,
		req.BuildingID,
		req.FloorID,
		req.ZoneID,
		req.PositionX,
		req.PositionY,
		req.DirectionAngle,
		req.FOVAngle,
		createdBy,
		createdBy,
	).Scan(
		&camera.ID,
		&camera.Name,
		&camera.Description,
		&camera.RTSPUrl,
		&camera.Status,
		&camera.Latitude,
		&camera.Longitude,
		&camera.BuildingID,
		&camera.FloorID,
		&camera.ZoneID,
		&camera.PositionX,
		&camera.PositionY,
		&camera.DirectionAngle,
		&camera.FOVAngle,
		&camera.CreatedBy,
		&camera.UpdatedBy,
		&camera.CreatedAt,
		&camera.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return camera, nil
}

func (r *CameraRepository) FindAll(ctx context.Context) ([]domain.Camera, error) {
	query := `
		SELECT
			id,
			name,
			description,
			rtsp_url,
			status,
			latitude,
			longitude,
			building_id,
			floor_id,
			zone_id,
			position_x,
			position_y,
			direction_angle,
			fov_angle,
			created_by,
			updated_by,
			created_at,
			updated_at
		FROM camera.cameras
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cameras := make([]domain.Camera, 0)

	for rows.Next() {
		var camera domain.Camera

		err := rows.Scan(
			&camera.ID,
			&camera.Name,
			&camera.Description,
			&camera.RTSPUrl,
			&camera.Status,
			&camera.Latitude,
			&camera.Longitude,
			&camera.BuildingID,
			&camera.FloorID,
			&camera.ZoneID,
			&camera.PositionX,
			&camera.PositionY,
			&camera.DirectionAngle,
			&camera.FOVAngle,
			&camera.CreatedBy,
			&camera.UpdatedBy,
			&camera.CreatedAt,
			&camera.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		cameras = append(cameras, camera)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cameras, nil
}

func (r *CameraRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Camera, error) {
	query := `
		SELECT
			id,
			name,
			description,
			rtsp_url,
			status,
			latitude,
			longitude,
			building_id,
			floor_id,
			zone_id,
			position_x,
			position_y,
			direction_angle,
			fov_angle,
			created_by,
			updated_by,
			created_at,
			updated_at
		FROM camera.cameras
		WHERE id = $1
		AND deleted_at IS NULL
	`

	camera := &domain.Camera{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&camera.ID,
		&camera.Name,
		&camera.Description,
		&camera.RTSPUrl,
		&camera.Status,
		&camera.Latitude,
		&camera.Longitude,
		&camera.BuildingID,
		&camera.FloorID,
		&camera.ZoneID,
		&camera.PositionX,
		&camera.PositionY,
		&camera.DirectionAngle,
		&camera.FOVAngle,
		&camera.CreatedBy,
		&camera.UpdatedBy,
		&camera.CreatedAt,
		&camera.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCameraNotFound
	}

	if err != nil {
		return nil, err
	}

	return camera, nil
}

func (r *CameraRepository) Delete(ctx context.Context, id uuid.UUID, updatedBy *uuid.UUID) error {
	query := `
		UPDATE camera.cameras
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
		return ErrCameraNotFound
	}

	return nil
}

func (r *CameraRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	req dto.UpdateCameraRequest,
	updatedBy *uuid.UUID,
) (*domain.Camera, error) {
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

	if req.RTSPUrl != nil {
		existing.RTSPUrl = *req.RTSPUrl
	}

	if req.Latitude != nil {
		existing.Latitude = req.Latitude
	}

	if req.Longitude != nil {
		existing.Longitude = req.Longitude
	}

	if req.BuildingID != nil {
		existing.BuildingID = req.BuildingID
	}

	if req.FloorID != nil {
		existing.FloorID = req.FloorID
	}

	if req.ZoneID != nil {
		existing.ZoneID = req.ZoneID
	}

	if req.PositionX != nil {
		existing.PositionX = req.PositionX
	}

	if req.PositionY != nil {
		existing.PositionY = req.PositionY
	}

	if req.DirectionAngle != nil {
		existing.DirectionAngle = req.DirectionAngle
	}

	if req.FOVAngle != nil {
		existing.FOVAngle = req.FOVAngle
	}

	query := `
		UPDATE camera.cameras
		SET
			name = $2,
			description = $3,
			rtsp_url = $4,
			latitude = $5,
			longitude = $6,
			building_id = $7,
			floor_id = $8,
			zone_id = $9,
			position_x = $10,
			position_y = $11,
			direction_angle = $12,
			fov_angle = $13,
			updated_by = $14,
			updated_at = now()
		WHERE id = $1
		AND deleted_at IS NULL
		RETURNING
			id,
			name,
			description,
			rtsp_url,
			status,
			latitude,
			longitude,
			building_id,
			floor_id,
			zone_id,
			position_x,
			position_y,
			direction_angle,
			fov_angle,
			created_by,
			updated_by,
			created_at,
			updated_at
	`

	camera := &domain.Camera{}

	err = r.db.QueryRow(
		ctx,
		query,
		id,
		existing.Name,
		existing.Description,
		existing.RTSPUrl,
		existing.Latitude,
		existing.Longitude,
		existing.BuildingID,
		existing.FloorID,
		existing.ZoneID,
		existing.PositionX,
		existing.PositionY,
		existing.DirectionAngle,
		existing.FOVAngle,
		updatedBy,
	).Scan(
		&camera.ID,
		&camera.Name,
		&camera.Description,
		&camera.RTSPUrl,
		&camera.Status,
		&camera.Latitude,
		&camera.Longitude,
		&camera.BuildingID,
		&camera.FloorID,
		&camera.ZoneID,
		&camera.PositionX,
		&camera.PositionY,
		&camera.DirectionAngle,
		&camera.FOVAngle,
		&camera.CreatedBy,
		&camera.UpdatedBy,
		&camera.CreatedAt,
		&camera.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return camera, nil
}
