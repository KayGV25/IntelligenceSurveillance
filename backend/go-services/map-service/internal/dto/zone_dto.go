package dto

import "github.com/google/uuid"

type PointRequest struct {
	X float64 `json:"x" binding:"required"`
	Y float64 `json:"y" binding:"required"`
}

type CreateZoneRequest struct {
	FloorID uuid.UUID `json:"floor_id" binding:"required"`

	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`

	ZoneType  *string `json:"zone_type"`
	Monitored *bool   `json:"monitored"`

	Polygon []PointRequest `json:"polygon" binding:"required,min=3"`

	Color        *string  `json:"color"`
	Opacity      *float64 `json:"opacity"`
	Icon         *string  `json:"icon"`
	LabelVisible *bool    `json:"label_visible"`
}

type UpdateZoneRequest struct {
	FloorID *uuid.UUID `json:"floor_id"`

	Name        *string `json:"name"`
	Description *string `json:"description"`

	ZoneType  *string `json:"zone_type"`
	Monitored *bool   `json:"monitored"`

	Polygon []PointRequest `json:"polygon"`

	Color        *string  `json:"color"`
	Opacity      *float64 `json:"opacity"`
	Icon         *string  `json:"icon"`
	LabelVisible *bool    `json:"label_visible"`
}
