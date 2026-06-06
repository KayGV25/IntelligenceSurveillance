package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/dto"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/service"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/common/response"
)

type DiscoveryHandler struct {
	discoveryService *service.DiscoveryService
}

func NewDiscoveryHandler(discoveryService *service.DiscoveryService) *DiscoveryHandler {
	return &DiscoveryHandler{discoveryService: discoveryService}
}

func (h *DiscoveryHandler) Discover(c *gin.Context) {
	var req dto.DiscoverCamerasRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_DISCOVERY_REQUEST", err.Error())
	}

	if req.Method != "CIDR_SCAN" {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "UNSUPPORTED_DISCOVERY_METHOD", "Only CIDR_SCAN is supported in V1")
		return
	}

	devices, err := h.discoveryService.DiscoverCIDR(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "DISCOVERY_FAILED", err.Error())
		return
	}

	response.OK(c, dto.DiscoverCamerasResponse{
		Devices: devices,
		Count:   len(devices),
	})
}

func (h *DiscoveryHandler) GetDiscoveredDevices(c *gin.Context) {
	devices, err := h.discoveryService.GetDiscoveredDevices(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "DISCOVERED_DEVICE_LIST_FAILED", err.Error())
		return
	}

	response.OK(c, gin.H{
		"devices": devices,
		"count":   len(devices),
	})
}

func (h *DiscoveryHandler) GetDiscoveredDeviceByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_DISCOVERED_DEVICE_ID", "Invalid discovered device id")
		return
	}

	device, err := h.discoveryService.GetDiscoveredDeviceByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "DISCOVERED_DEVICE_NOT_FOUND", "Discovered device not found")
		return
	}

	response.OK(c, device)
}
