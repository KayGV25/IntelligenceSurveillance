package domain

import (
	"time"

	"github.com/google/uuid"
)

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Zone struct {
	ID      uuid.UUID `json:"id"`
	FloorID uuid.UUID `json:"floor_id"`

	Name        string  `json:"name"`
	Description *string `json:"description"`

	ZoneType  string `json:"zone_type"`
	Monitored bool   `json:"monitored"`

	Polygon []Point `json:"polygon"`

	Color        string  `json:"color"`
	Opacity      float64 `json:"opacity"`
	Icon         *string `json:"icon"`
	LabelVisible bool    `json:"label_visible"`

	CreatedBy *uuid.UUID `json:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}