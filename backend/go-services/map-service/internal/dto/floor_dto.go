package dto

import "github.com/google/uuid"

type CreateFloorRequest struct {
	BuildingID uuid.UUID `json:"building_id" binding:"required"`

	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	FloorNumber *int    `json:"floor_number"`

	Width  *float64 `json:"width"`
	Height *float64 `json:"height"`
	Unit   *string  `json:"unit"`

	Color        *string  `json:"color"`
	Opacity      *float64 `json:"opacity"`
	Icon         *string  `json:"icon"`
	LabelVisible *bool    `json:"label_visible"`
}

type UpdateFloorRequest struct {
	BuildingID *uuid.UUID `json:"building_id"`

	Name        *string `json:"name"`
	Description *string `json:"description"`
	FloorNumber *int    `json:"floor_number"`

	Width  *float64 `json:"width"`
	Height *float64 `json:"height"`
	Unit   *string  `json:"unit"`

	Color        *string  `json:"color"`
	Opacity      *float64 `json:"opacity"`
	Icon         *string  `json:"icon"`
	LabelVisible *bool    `json:"label_visible"`
}
	