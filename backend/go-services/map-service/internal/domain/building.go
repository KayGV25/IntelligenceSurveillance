package domain

import (
	"time"

	"github.com/google/uuid"
)

type Building struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Address     *string   `json:"address"`

	Color        string  `json:"color"`
	Opacity      float64 `json:"opacity"`
	Icon         *string `json:"icon"`
	LabelVisible bool    `json:"label_visible"`

	CreatedBy *uuid.UUID `json:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
