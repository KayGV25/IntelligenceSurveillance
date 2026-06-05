package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/config"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/database"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/event"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/repository"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/router"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/service"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	db, err := database.NewPostgresPool(cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer db.Close()

	cameraRepo := repository.NewCameraRepository(db)
	deviceRepo := repository.NewDiscoveredDeviceRepository(db)
	contractRepo := repository.NewConnectionContractRepository(db)

	eventPublisher := event.NewKafkaPublisher(cfg.RedpandaBrokers)
	defer eventPublisher.Close()

	cameraService := service.NewCameraService(
		cameraRepo,
		eventPublisher,
		deviceRepo,
		contractRepo,
	)
	discoveryService := service.NewDiscoveryService(deviceRepo)
	r := router.NewRouter(cameraService, discoveryService)

	addr := fmt.Sprintf(":%s", cfg.AppPort)

	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
