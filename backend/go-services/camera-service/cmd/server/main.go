package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/config"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/database"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/event"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/repository"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/router"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/service"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/snapshot"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/storage"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/stream"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	db, err := database.NewPostgresPool(cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer db.Close()

	storageClient, err := storage.NewClient(cfg)
	if err != nil {
		log.Fatalf("failed to create minio client: %v", err)
	}

	snapshotService := snapshot.NewService(storageClient, 5*time.Second)

	cameraRepo := repository.NewCameraRepository(db)
	deviceRepo := repository.NewDiscoveredDeviceRepository(db)
	contractRepo := repository.NewConnectionContractRepository(db)

	streamValidator := stream.NewValidator(5 * time.Second)

	eventPublisher := event.NewKafkaPublisher(cfg.RedpandaBrokers)
	defer eventPublisher.Close()

	cameraService := service.NewCameraService(
		cameraRepo,
		eventPublisher,
		deviceRepo,
		contractRepo,
		streamValidator,
		snapshotService,
	)
	discoveryService := service.NewDiscoveryService(deviceRepo, eventPublisher)
	r := router.NewRouter(cameraService, discoveryService)

	addr := fmt.Sprintf(":%s", cfg.AppPort)

	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
