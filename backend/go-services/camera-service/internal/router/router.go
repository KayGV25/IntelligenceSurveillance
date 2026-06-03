package router

import (
	"github.com/gin-gonic/gin"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/handler"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/middleware"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/service"
)

func NewRouter(cameraService *service.CameraService) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.UserContextMiddleware())

	healthHandler := handler.NewHealthHandler()
	cameraHandler := handler.NewCameraHandler(cameraService)

	r.GET("/health", healthHandler.Check)

	api := r.Group("/api/v1")
	{
		api.POST("/cameras", cameraHandler.Create)
		api.GET("/cameras", cameraHandler.GetAll)
		api.GET("/cameras/:id", cameraHandler.GetByID)
		api.PUT("/cameras/:id", cameraHandler.Update)
		api.DELETE("/cameras/:id", cameraHandler.Delete)
	}

	return r
}
