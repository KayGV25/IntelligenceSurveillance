package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/common/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	response.OK(c, gin.H{
		"status":  "UP",
		"service": "map-service",
	})
}
