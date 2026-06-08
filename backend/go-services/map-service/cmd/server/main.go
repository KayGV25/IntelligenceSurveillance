package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/config"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/database"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/repository"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/router"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/map-service/internal/service"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	db, err := database.NewPostgresPool(cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer db.Close()

	buildingRepo := repository.NewBuildingRepository(db)
	buildingService := service.NewBuildingService(buildingRepo)

	floorRepo := repository.NewFloorRepository(db)
	floorService := service.NewFloorService(floorRepo)

	zoneRepo := repository.NewZoneRepository(db)
	zoneService := service.NewZoneService(zoneRepo)

	r := router.NewRouter(buildingService, floorService, zoneService)

	addr := fmt.Sprintf(":%s", cfg.AppPort)

	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
