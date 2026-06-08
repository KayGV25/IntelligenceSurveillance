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

type BuildingHandler struct {
	buildingService *service.BuildingService
}

func NewBuildingHandler(buildingService *service.BuildingService) *BuildingHandler {
	return &BuildingHandler{buildingService: buildingService}
}

func (h *BuildingHandler) Create(c *gin.Context) {
	var req dto.CreateBuildingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BUILDING_REQUEST", err.Error())
		return
	}

	userCtx := security.FromGin(c)

	var userID *uuid.UUID
	if userCtx != nil {
		userID = userCtx.UserID
	}

	building, err := h.buildingService.Create(c.Request.Context(), req, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "BUILDING_CREATE_FAILED", err.Error())
		return
	}

	response.Created(c, building)
}

func (h *BuildingHandler) GetAll(c *gin.Context) {
	buildings, err := h.buildingService.GetAll(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "BUILDING_LIST_FAILED", err.Error())
		return
	}

	response.OK(c, buildings)
}

func (h *BuildingHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BUILDING_ID", "Invalid building id")
		return
	}

	building, err := h.buildingService.GetByID(c.Request.Context(), id)
	if errors.Is(err, repository.ErrBuildingNotFound) {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "BUILDING_NOT_FOUND", "Building not found")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "BUILDING_GET_FAILED", err.Error())
		return
	}

	response.OK(c, building)
}

func (h *BuildingHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BUILDING_ID", "Invalid building id")
		return
	}

	var req dto.UpdateBuildingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BUILDING_REQUEST", err.Error())
		return
	}

	userCtx := security.FromGin(c)

	var userID *uuid.UUID
	if userCtx != nil {
		userID = userCtx.UserID
	}

	building, err := h.buildingService.Update(c.Request.Context(), id, req, userID)
	if errors.Is(err, repository.ErrBuildingNotFound) {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "BUILDING_NOT_FOUND", "Building not found")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "BUILDING_UPDATE_FAILED", err.Error())
		return
	}

	response.OK(c, building)
}

func (h *BuildingHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_BUILDING_ID", "Invalid building id")
		return
	}

	userCtx := security.FromGin(c)

	var userID *uuid.UUID
	if userCtx != nil {
		userID = userCtx.UserID
	}

	err = h.buildingService.Delete(c.Request.Context(), id, userID)
	if errors.Is(err, repository.ErrBuildingNotFound) {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "BUILDING_NOT_FOUND", "Building not found")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "BUILDING_DELETE_FAILED", err.Error())
		return
	}

	response.OK(c, gin.H{
		"deleted": true,
	})
}
