package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/config"
	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/database"
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
	cameraService := service.NewCameraService(cameraRepo)

	r := router.NewRouter(cameraService)

	addr := fmt.Sprintf(":%s", cfg.AppPort)

	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
