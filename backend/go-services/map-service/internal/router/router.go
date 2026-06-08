package router

import (
	"github.com/gin-gonic/gin"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/handler"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/middleware"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/service"
)

func NewRouter(
	buildingService *service.BuildingService,
	floorService *service.FloorService,
	zoneService *service.ZoneService,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.UserContextMiddleware())

	healthHandler := handler.NewHealthHandler()
	buildingHandler := handler.NewBuildingHandler(buildingService)
	floorHandler := handler.NewFloorHandler(floorService)
	zoneHandler := handler.NewZoneHandler(zoneService)

	r.GET("/health", healthHandler.Check)

	maps := r.Group("/api/v1/maps")
	{
		maps.POST("/buildings", buildingHandler.Create)
		maps.GET("/buildings", buildingHandler.GetAll)
		maps.GET("/buildings/:id", buildingHandler.GetByID)
		maps.PUT("/buildings/:id", buildingHandler.Update)
		maps.DELETE("/buildings/:id", buildingHandler.Delete)

		maps.POST("/floors", floorHandler.Create)
		maps.GET("/floors", floorHandler.GetAll)
		maps.GET("/floors/:id", floorHandler.GetByID)
		maps.PUT("/floors/:id", floorHandler.Update)
		maps.DELETE("/floors/:id", floorHandler.Delete)

		maps.POST("/zones", zoneHandler.Create)
		maps.GET("/zones", zoneHandler.GetAll)
		maps.GET("/zones/:id", zoneHandler.GetByID)
		maps.PUT("/zones/:id", zoneHandler.Update)
		maps.DELETE("/zones/:id", zoneHandler.Delete)

		maps.GET("/floors/:id/zones", zoneHandler.GetByFloorID)
	}

	return r
}
