package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/domain"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/dto"
)

var ErrZoneNotFound = errors.New("zone not found")

type ZoneRepository struct {
	db *pgxpool.Pool
}

func NewZoneRepository(db *pgxpool.Pool) *ZoneRepository {
	return &ZoneRepository{db: db}
}

func polygonToWKT(points []dto.PointRequest) (string, error) {
	if len(points) < 3 {
		return "", errors.New("polygon must have at least 3 points")
	}

	closed := make([]dto.PointRequest, 0, len(points)+1)
	closed = append(closed, points...)

	first := points[0]
	last := points[len(points)-1]

	if first.X != last.X || first.Y != last.Y {
		closed = append(closed, first)
	}

	parts := make([]string, 0, len(closed))
	for _, p := range closed {
		if p.X < 0 || p.X > 1 || p.Y < 0 || p.Y > 1 {
			return "", errors.New("polygon coordinates must be normalized between 0 and 1")
		}

		parts = append(parts, fmt.Sprintf("%s %s",
			strconv.FormatFloat(p.X, 'f', -1, 64),
			strconv.FormatFloat(p.Y, 'f', -1, 64),
		))
	}

	return fmt.Sprintf("POLYGON((%s))", strings.Join(parts, ",")), nil
}

func parsePolygonWKT(wkt string) ([]domain.Point, error) {
	wkt = strings.TrimSpace(wkt)
	wkt = strings.TrimPrefix(wkt, "POLYGON((")
	wkt = strings.TrimSuffix(wkt, "))")

	if wkt == "" {
		return []domain.Point{}, nil
	}

	rawPoints := strings.Split(wkt, ",")
	points := make([]domain.Point, 0, len(rawPoints))

	for i, raw := range rawPoints {
		parts := strings.Fields(strings.TrimSpace(raw))
		if len(parts) != 2 {
			return nil, errors.New("invalid polygon WKT")
		}

		x, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return nil, err
		}

		y, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, err
		}

		if i == len(rawPoints)-1 && len(points) > 0 {
			first := points[0]
			if first.X == x && first.Y == y {
				continue
			}
		}

		points = append(points, domain.Point{
			X: x,
			Y: y,
		})
	}

	return points, nil
}

func scanZone(row pgx.Row) (*domain.Zone, error) {
	zone := &domain.Zone{}
	var polygonWKT string

	err := row.Scan(
		&zone.ID,
		&zone.FloorID,
		&zone.Name,
		&zone.Description,
		&zone.ZoneType,
		&zone.Monitored,
		&polygonWKT,
		&zone.Color,
		&zone.Opacity,
		&zone.Icon,
		&zone.LabelVisible,
		&zone.CreatedBy,
		&zone.UpdatedBy,
		&zone.CreatedAt,
		&zone.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	points, err := parsePolygonWKT(polygonWKT)
	if err != nil {
		return nil, err
	}

	zone.Polygon = points

	return zone, nil
}

func (r *ZoneRepository) Create(ctx context.Context, req dto.CreateZoneRequest, createdBy *uuid.UUID) (*domain.Zone, error) {
	id := uuid.New()

	zoneType := "GENERAL"
	if req.ZoneType != nil && *req.ZoneType != "" {
		zoneType = *req.ZoneType
	}

	monitored := true
	if req.Monitored != nil {
		monitored = *req.Monitored
	}

	color := "#22C55E"
	if req.Color != nil && *req.Color != "" {
		color = *req.Color
	}

	opacity := 0.35
	if req.Opacity != nil {
		opacity = *req.Opacity
	}

	labelVisible := true
	if req.LabelVisible != nil {
		labelVisible = *req.LabelVisible
	}

	polygonWKT, err := polygonToWKT(req.Polygon)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO map.zones (
			id, floor_id, name, description,
			zone_type, monitored, geometry,
			color, opacity, icon, label_visible,
			created_by, updated_by
		)
		VALUES (
			$1,$2,$3,$4,
			$5,$6,ST_GeomFromText($7),
			$8,$9,$10,$11,
			$12,$13
		)
		RETURNING
			id, floor_id, name, description,
			zone_type, monitored, ST_AsText(geometry),
			color, opacity, icon, label_visible,
			created_by, updated_by, created_at, updated_at
	`

	return scanZone(r.db.QueryRow(
		ctx,
		query,
		id,
		req.FloorID,
		req.Name,
		req.Description,
		zoneType,
		monitored,
		polygonWKT,
		color,
		opacity,
		req.Icon,
		labelVisible,
		createdBy,
		createdBy,
	))
}

func (r *ZoneRepository) FindAll(ctx context.Context, floorID *uuid.UUID) ([]domain.Zone, error) {
	query := `
		SELECT
			id, floor_id, name, description,
			zone_type, monitored, ST_AsText(geometry),
			color, opacity, icon, label_visible,
			created_by, updated_by, created_at, updated_at
		FROM map.zones
		WHERE deleted_at IS NULL
		AND ($1::uuid IS NULL OR floor_id = $1)
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, floorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	zones := make([]domain.Zone, 0)

	for rows.Next() {
		zone := &domain.Zone{}
		var polygonWKT string

		err := rows.Scan(
			&zone.ID,
			&zone.FloorID,
			&zone.Name,
			&zone.Description,
			&zone.ZoneType,
			&zone.Monitored,
			&polygonWKT,
			&zone.Color,
			&zone.Opacity,
			&zone.Icon,
			&zone.LabelVisible,
			&zone.CreatedBy,
			&zone.UpdatedBy,
			&zone.CreatedAt,
			&zone.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		points, err := parsePolygonWKT(polygonWKT)
		if err != nil {
			return nil, err
		}

		zone.Polygon = points
		zones = append(zones, *zone)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return zones, nil
}

func (r *ZoneRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Zone, error) {
	query := `
		SELECT
			id, floor_id, name, description,
			zone_type, monitored, ST_AsText(geometry),
			color, opacity, icon, label_visible,
			created_by, updated_by, created_at, updated_at
		FROM map.zones
		WHERE id = $1
		AND deleted_at IS NULL
	`

	zone, err := scanZone(r.db.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrZoneNotFound
	}
	if err != nil {
		return nil, err
	}

	return zone, nil
}

func (r *ZoneRepository) Update(ctx context.Context, id uuid.UUID, req dto.UpdateZoneRequest, updatedBy *uuid.UUID) (*domain.Zone, error) {
	existing, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.FloorID != nil {
		existing.FloorID = *req.FloorID
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.ZoneType != nil {
		existing.ZoneType = *req.ZoneType
	}
	if req.Monitored != nil {
		existing.Monitored = *req.Monitored
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

	polygonReq := make([]dto.PointRequest, 0)

	if len(req.Polygon) > 0 {
		polygonReq = req.Polygon
	} else {
		for _, p := range existing.Polygon {
			polygonReq = append(polygonReq, dto.PointRequest{
				X: p.X,
				Y: p.Y,
			})
		}
	}

	polygonWKT, err := polygonToWKT(polygonReq)
	if err != nil {
		return nil, err
	}

	query := `
		UPDATE map.zones
		SET
			floor_id = $2,
			name = $3,
			description = $4,
			zone_type = $5,
			monitored = $6,
			geometry = ST_GeomFromText($7),
			color = $8,
			opacity = $9,
			icon = $10,
			label_visible = $11,
			updated_by = $12,
			updated_at = now()
		WHERE id = $1
		AND deleted_at IS NULL
		RETURNING
			id, floor_id, name, description,
			zone_type, monitored, ST_AsText(geometry),
			color, opacity, icon, label_visible,
			created_by, updated_by, created_at, updated_at
	`

	return scanZone(r.db.QueryRow(
		ctx,
		query,
		id,
		existing.FloorID,
		existing.Name,
		existing.Description,
		existing.ZoneType,
		existing.Monitored,
		polygonWKT,
		existing.Color,
		existing.Opacity,
		existing.Icon,
		existing.LabelVisible,
		updatedBy,
	))
}

func (r *ZoneRepository) Delete(ctx context.Context, id uuid.UUID, updatedBy *uuid.UUID) error {
	query := `
		UPDATE map.zones
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
		return ErrZoneNotFound
	}

	return nil
}
