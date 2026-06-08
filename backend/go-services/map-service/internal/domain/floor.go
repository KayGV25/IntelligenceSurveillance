package domain

import (
	"time"

	"github.com/google/uuid"
)

type Floor struct {
	ID         uuid.UUID `json:"id"`
	BuildingID uuid.UUID `json:"building_id"`

	Name        string  `json:"name"`
	Description *string `json:"description"`
	FloorNumber *int    `json:"floor_number"`

	Width  *float64 `json:"width"`
	Height *float64 `json:"height"`
	Unit   string   `json:"unit"`

	Color        string  `json:"color"`
	Opacity      float64 `json:"opacity"`
	Icon         *string `json:"icon"`
	LabelVisible bool    `json:"label_visible"`

	CreatedBy *uuid.UUID `json:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
