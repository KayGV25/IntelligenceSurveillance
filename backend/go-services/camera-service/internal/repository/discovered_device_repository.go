package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/domain"
)

type DiscoveredDeviceRepository struct {
	db *pgxpool.Pool
}

func NewDiscoveredDeviceRepository(db *pgxpool.Pool) *DiscoveredDeviceRepository {
	return &DiscoveredDeviceRepository{db: db}
}

func (r *DiscoveredDeviceRepository) Upsert(
	ctx context.Context,
	device domain.DiscoveredDevice,
) (*domain.DiscoveredDevice, error) {
	if device.ID == uuid.Nil {
		device.ID = uuid.New()
	}

	query := `
		INSERT INTO camera.discovered_devices (
			id,
			ip_address,
			mac_address,
			hostname,
			manufacturer,
			model,
			firmware_version,
			onvif_supported,
			rtsp_supported,
			http_supported,
			http_port,
			rtsp_port,
			onvif_port,
			discovery_method,
			status,
			device_type,
			confidence,
			detection_reason,
			discovered_at,
			last_seen_at
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15,
			$16, $17, $18,
			now(), now()
		)
		ON CONFLICT (ip_address)
		DO UPDATE SET
			mac_address = EXCLUDED.mac_address,
			hostname = EXCLUDED.hostname,
			manufacturer = EXCLUDED.manufacturer,
			model = EXCLUDED.model,
			firmware_version = EXCLUDED.firmware_version,
			onvif_supported = EXCLUDED.onvif_supported,
			rtsp_supported = EXCLUDED.rtsp_supported,
			http_supported = EXCLUDED.http_supported,
			http_port = EXCLUDED.http_port,
			rtsp_port = EXCLUDED.rtsp_port,
			onvif_port = EXCLUDED.onvif_port,
			discovery_method = EXCLUDED.discovery_method,
			status = EXCLUDED.status,
			device_type = EXCLUDED.device_type,
			confidence = EXCLUDED.confidence,
			detection_reason = EXCLUDED.detection_reason,	
			last_seen_at = now(),
			updated_at = now()
		RETURNING
			id,
			ip_address,
			mac_address,
			hostname,
			manufacturer,
			model,
			firmware_version,
			onvif_supported,
			rtsp_supported,
			http_supported,
			http_port,
			rtsp_port,
			onvif_port,
			discovery_method,
			status,
			device_type,
			confidence,
			detection_reason,
			discovered_at,
			last_seen_at,
			created_at,
			updated_at
	`

	saved := &domain.DiscoveredDevice{}

	err := r.db.QueryRow(
		ctx,
		query,
		device.ID,
		device.IPAddress,
		device.MACAddress,
		device.Hostname,
		device.Manufacturer,
		device.Model,
		device.FirmwareVersion,
		device.ONVIFSupported,
		device.RTSPSupported,
		device.HTTPSupported,
		device.HTTPPort,
		device.RTSPPort,
		device.ONVIFPort,
		device.DiscoveryMethod,
		device.Status,
		device.DeviceType,
		device.Confidence,
		device.DetectionReason,
	).Scan(
		&saved.ID,
		&saved.IPAddress,
		&saved.MACAddress,
		&saved.Hostname,
		&saved.Manufacturer,
		&saved.Model,
		&saved.FirmwareVersion,
		&saved.ONVIFSupported,
		&saved.RTSPSupported,
		&saved.HTTPSupported,
		&saved.HTTPPort,
		&saved.RTSPPort,
		&saved.ONVIFPort,
		&saved.DiscoveryMethod,
		&saved.Status,
		&saved.DeviceType,
		&saved.Confidence,
		&saved.DetectionReason,
		&saved.DiscoveredAt,
		&saved.LastSeenAt,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (r *DiscoveredDeviceRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.DiscoveredDevice, error) {
	query := `
		SELECT
			id,
			ip_address,
			mac_address,
			hostname,
			manufacturer,
			model,
			firmware_version,
			onvif_supported,
			rtsp_supported,
			http_supported,
			http_port,
			rtsp_port,
			onvif_port,
			discovery_method,
			status,
			device_type,
			confidence,
			detection_reason,
			discovered_at,
			last_seen_at,
			created_at,
			updated_at
		FROM camera.discovered_devices
		WHERE id = $1
	`

	device := &domain.DiscoveredDevice{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&device.ID,
		&device.IPAddress,
		&device.MACAddress,
		&device.Hostname,
		&device.Manufacturer,
		&device.Model,
		&device.FirmwareVersion,
		&device.ONVIFSupported,
		&device.RTSPSupported,
		&device.HTTPSupported,
		&device.HTTPPort,
		&device.RTSPPort,
		&device.ONVIFPort,
		&device.DiscoveryMethod,
		&device.Status,
		&device.DeviceType,
		&device.Confidence,
		&device.DetectionReason,
		&device.DiscoveredAt,
		&device.LastSeenAt,
		&device.CreatedAt,
		&device.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return device, nil
}
