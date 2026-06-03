package dto

import "github.com/google/uuid"

type CreateCameraRequest struct {
	Name           string     `json:"name" binding:"required"`
	Description    *string    `json:"description"`
	RTSPUrl        string     `json:"rtsp_url" binding:"required"`
	Latitude       *float64   `json:"latitude"`
	Longitude      *float64   `json:"longitude"`
	BuildingID     *uuid.UUID `json:"building_id"`
	FloorID        *uuid.UUID `json:"floor_id"`
	ZoneID         *uuid.UUID `json:"zone_id"`
	PositionX      *float64   `json:"position_x"`
	PositionY      *float64   `json:"position_y"`
	DirectionAngle *float64   `json:"direction_angle"`
	FOVAngle       *float64   `json:"fov_angle"`
}

type UpdateCameraRequest struct {
	Name           *string    `json:"name"`
	Description    *string    `json:"description"`
	RTSPUrl        *string    `json:"rtsp_url"`
	Latitude       *float64   `json:"latitude"`
	Longitude      *float64   `json:"longitude"`
	BuildingID     *uuid.UUID `json:"building_id"`
	FloorID        *uuid.UUID `json:"floor_id"`
	ZoneID         *uuid.UUID `json:"zone_id"`
	PositionX      *float64   `json:"position_x"`
	PositionY      *float64   `json:"position_y"`
	DirectionAngle *float64   `json:"direction_angle"`
	FOVAngle       *float64   `json:"fov_angle"`
}
