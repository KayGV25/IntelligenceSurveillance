package domain

import (
	"time"

	"github.com/google/uuid"
)

type CameraStatus string

const (
	CameraStatusOnline   CameraStatus = "ONLINE"
	CameraStatusOffline  CameraStatus = "OFFLINE"
	CameraStatusDegraded CameraStatus = "DEGRADED"
	CameraStatusUnknown  CameraStatus = "UNKNOWN"
)

type Camera struct {
	ID             uuid.UUID    `json:"id"`
	Name           string       `json:"name"`
	Description    *string      `json:"description"`
	RTSPUrl        *string      `json:"rtsp_url"`
	Status         CameraStatus `json:"status"`
	Latitude       *float64     `json:"latitude"`
	Longitude      *float64     `json:"longitude"`
	BuildingID     *uuid.UUID   `json:"building_id"`
	FloorID        *uuid.UUID   `json:"floor_id"`
	ZoneID         *uuid.UUID   `json:"zone_id"`
	PositionX      *float64     `json:"position_x"`
	PositionY      *float64     `json:"position_y"`
	DirectionAngle *float64     `json:"direction_angle"`
	FOVAngle       *float64     `json:"fov_angle"`
	CreatedBy      *uuid.UUID   `json:"created_by"`
	UpdatedBy      *uuid.UUID   `json:"updated_by"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}
