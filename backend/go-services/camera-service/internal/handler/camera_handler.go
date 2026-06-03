package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/dto"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/repository"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/service"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/common/response"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/common/security"
)

type CameraHandler struct {
	cameraService *service.CameraService
}

func NewCameraHandler(cameraService *service.CameraService) *CameraHandler {
	return &CameraHandler{
		cameraService: cameraService,
	}
}

func (h *CameraHandler) Create(c *gin.Context) {
	var req dto.CreateCameraRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_CAMERA_REQUEST", err.Error())
		return
	}

	userCtx := security.FromGin(c)

	camera, err := h.cameraService.Create(c.Request.Context(), req, userCtx.UserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "CAMERA_CREATE_FAILED", err.Error())
		return
	}

	response.Created(c, camera)
}

func (h *CameraHandler) GetAll(c *gin.Context) {
	cameras, err := h.cameraService.GetAll(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "CAMERA_LIST_FAILED", err.Error())
		return
	}

	response.OK(c, cameras)
}

func (h *CameraHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_CAMERA_ID", "Invalid camera id")
		return
	}

	camera, err := h.cameraService.GetByID(c.Request.Context(), id)
	if errors.Is(err, repository.ErrCameraNotFound) {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "CAMERA_NOT_FOUND", "Camera not found")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "CAMERA_GET_FAILED", err.Error())
		return
	}

	response.OK(c, camera)
}

func (h *CameraHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_CAMERA_ID", "Invalid camera id")
		return
	}

	userCtx := security.FromGin(c)

	err = h.cameraService.Delete(c.Request.Context(), id, userCtx.UserID)
	if errors.Is(err, repository.ErrCameraNotFound) {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "CAMERA_NOT_FOUND", "Camera not found")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "CAMERA_DELETE_FAILED", err.Error())
		return
	}

	response.OK(c, gin.H{
		"deleted": true,
	})
}

func (h *CameraHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_CAMERA_ID", "Invalid camera id")
		return
	}

	var req dto.UpdateCameraRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_CAMERA_REQUEST", err.Error())
		return
	}

	userCtx := security.FromGin(c)

	camera, err := h.cameraService.Update(c.Request.Context(), id, req, userCtx.UserID)
	if errors.Is(err, repository.ErrCameraNotFound) {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "CAMERA_NOT_FOUND", "Camera not found")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "CAMERA_UPDATE_FAILED", err.Error())
		return
	}

	response.OK(c, camera)
}
