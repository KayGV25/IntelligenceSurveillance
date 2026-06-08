package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/common/response"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/common/security"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/dto"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/repository"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/service"
)

type ZoneHandler struct {
	zoneService *service.ZoneService
}

func NewZoneHandler(zoneService *service.ZoneService) *ZoneHandler {
	return &ZoneHandler{zoneService: zoneService}
}

func (h *ZoneHandler) Create(c *gin.Context) {
	var req dto.CreateZoneRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ZONE_REQUEST", err.Error())
		return
	}

	userCtx := security.FromGin(c)

	var userID *uuid.UUID
	if userCtx != nil {
		userID = userCtx.UserID
	}

	zone, err := h.zoneService.Create(c.Request.Context(), req, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "ZONE_CREATE_FAILED", err.Error())
		return
	}

	response.Created(c, zone)
}

func (h *ZoneHandler) GetAll(c *gin.Context) {
	var floorID *uuid.UUID

	floorIDQuery := c.Query("floor_id")
	if floorIDQuery != "" {
		parsed, err := uuid.Parse(floorIDQuery)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_FLOOR_ID", "Invalid floor id")
			return
		}

		floorID = &parsed
	}

	zones, err := h.zoneService.GetAll(c.Request.Context(), floorID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "ZONE_LIST_FAILED", err.Error())
		return
	}

	response.OK(c, zones)
}

func (h *ZoneHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ZONE_ID", "Invalid zone id")
		return
	}

	zone, err := h.zoneService.GetByID(c.Request.Context(), id)
	if errors.Is(err, repository.ErrZoneNotFound) {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "ZONE_NOT_FOUND", "Zone not found")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "ZONE_GET_FAILED", err.Error())
		return
	}

	response.OK(c, zone)
}

func (h *ZoneHandler) GetByFloorID(c *gin.Context) {
	floorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_FLOOR_ID", "Invalid floor id")
		return
	}

	zones, err := h.zoneService.GetAll(c.Request.Context(), &floorID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "ZONE_LIST_FAILED", err.Error())
		return
	}

	response.OK(c, zones)
}

func (h *ZoneHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ZONE_ID", "Invalid zone id")
		return
	}

	var req dto.UpdateZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ZONE_REQUEST", err.Error())
		return
	}

	userCtx := security.FromGin(c)

	var userID *uuid.UUID
	if userCtx != nil {
		userID = userCtx.UserID
	}

	zone, err := h.zoneService.Update(c.Request.Context(), id, req, userID)
	if errors.Is(err, repository.ErrZoneNotFound) {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "ZONE_NOT_FOUND", "Zone not found")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "ZONE_UPDATE_FAILED", err.Error())
		return
	}

	response.OK(c, zone)
}

func (h *ZoneHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_ZONE_ID", "Invalid zone id")
		return
	}

	userCtx := security.FromGin(c)

	var userID *uuid.UUID
	if userCtx != nil {
		userID = userCtx.UserID
	}

	err = h.zoneService.Delete(c.Request.Context(), id, userID)
	if errors.Is(err, repository.ErrZoneNotFound) {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "ZONE_NOT_FOUND", "Zone not found")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "ZONE_DELETE_FAILED", err.Error())
		return
	}

	response.OK(c, gin.H{
		"deleted": true,
	})
}
