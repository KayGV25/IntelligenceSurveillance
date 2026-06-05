package dto

import "github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/domain"

type DiscoverCamerasRequest struct {
	Method      string `json:"method" binding:"required"`
	NetworkCIDR string `json:"network_cidr" binding:"required"`
	Ports       []int  `json:"ports"`
	TimeoutMs   int    `json:"timeout_ms"`
	MaxWorkers  int    `json:"max_workers"`
}

type DiscoverCamerasResponse struct {
	Devices []domain.DiscoveredDevice `json:"devices"`
	Count   int                       `json:"count"`
}
