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

type FloorHandler struct {
	floorService *service.FloorService
}

func NewFloorHandler(floorService *service.FloorService) *FloorHandler {
	return &FloorHandler{floorService: floorService}
}

func (h *FloorHandler) Create(c *gin.Context) {
	var req dto.CreateFloorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_FLOOR_REQUEST", err.Error())
		return
	}

	userCtx := security.FromGin(c)

	var userID *uuid.UUID
	if userCtx != nil {
		userID = userCtx.UserID
	}

	floor, err := h.floorService.Create(c.Request.Context(), req, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "FLOOR_CREATE_FAILED", err.Error())
		return
	}

	response.Created(c, floor)
}

func (h *FloorHandler) GetAll(c *gin.Context) {
	floors, err := h.floorService.GetAll(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "FLOOR_LIST_FAILED", err.Error())
		return
	}

	response.OK(c, floors)
}

func (h *FloorHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_FLOOR_ID", "Invalid floor id")
		return
	}

	floor, err := h.floorService.GetByID(c.Request.Context(), id)
	if errors.Is(err, repository.ErrFloorNotFound) {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "FLOOR_NOT_FOUND", "Floor not found")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "FLOOR_GET_FAILED", err.Error())
		return
	}

	response.OK(c, floor)
}

func (h *FloorHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_FLOOR_ID", "Invalid floor id")
		return
	}

	var req dto.UpdateFloorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_FLOOR_REQUEST", err.Error())
		return
	}

	userCtx := security.FromGin(c)

	var userID *uuid.UUID
	if userCtx != nil {
		userID = userCtx.UserID
	}

	floor, err := h.floorService.Update(c.Request.Context(), id, req, userID)
	if errors.Is(err, repository.ErrFloorNotFound) {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "FLOOR_NOT_FOUND", "Floor not found")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "FLOOR_UPDATE_FAILED", err.Error())
		return
	}

	response.OK(c, floor)
}

func (h *FloorHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_FLOOR_ID", "Invalid floor id")
		return
	}

	userCtx := security.FromGin(c)

	var userID *uuid.UUID
	if userCtx != nil {
		userID = userCtx.UserID
	}

	err = h.floorService.Delete(c.Request.Context(), id, userID)
	if errors.Is(err, repository.ErrFloorNotFound) {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "FLOOR_NOT_FOUND", "Floor not found")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "FLOOR_DELETE_FAILED", err.Error())
		return
	}

	response.OK(c, gin.H{
		"deleted": true,
	})
}
