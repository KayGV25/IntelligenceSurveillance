package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/KayGV25/IntelligenceSurveillance/backend/go-services/camera-service/internal/domain"
)

type ConnectionContractRepository struct {
	db *pgxpool.Pool
}

func NewConnectionContractRepository(db *pgxpool.Pool) *ConnectionContractRepository {
	return &ConnectionContractRepository{db: db}
}

func (r *ConnectionContractRepository) Create(
	ctx context.Context,
	contract domain.CameraConnectionContract,
) (*domain.CameraConnectionContract, error) {
	if contract.ID == uuid.Nil {
		contract.ID = uuid.New()
	}

	query := `
		INSERT INTO camera.camera_connection_contracts (
			id,
			camera_id,
			discovered_device_id,
			connection_type,
			ip_address,
			rtsp_url,
			onvif_url,
			username,
			encrypted_password,
			stream_profile_token,
			main_stream_url,
			sub_stream_url,
			snapshot_url,
			created_by,
			updated_by
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15
		)
		RETURNING
			id,
			camera_id,
			discovered_device_id,
			connection_type,
			ip_address,
			rtsp_url,
			onvif_url,
			username,
			encrypted_password,
			stream_profile_token,
			main_stream_url,
			sub_stream_url,
			snapshot_url,
			created_by,
			updated_by,
			created_at,
			updated_at
	`

	saved := &domain.CameraConnectionContract{}

	err := r.db.QueryRow(
		ctx,
		query,
		contract.ID,
		contract.CameraID,
		contract.DiscoveredDeviceID,
		contract.ConnectionType,
		contract.IPAddress,
		contract.RTSPUrl,
		contract.ONVIFUrl,
		contract.Username,
		contract.EncryptedPassword,
		contract.StreamProfileToken,
		contract.MainStreamURL,
		contract.SubStreamURL,
		contract.SnapshotURL,
		contract.CreatedBy,
		contract.UpdatedBy,
	).Scan(
		&saved.ID,
		&saved.CameraID,
		&saved.DiscoveredDeviceID,
		&saved.ConnectionType,
		&saved.IPAddress,
		&saved.RTSPUrl,
		&saved.ONVIFUrl,
		&saved.Username,
		&saved.EncryptedPassword,
		&saved.StreamProfileToken,
		&saved.MainStreamURL,
		&saved.SubStreamURL,
		&saved.SnapshotURL,
		&saved.CreatedBy,
		&saved.UpdatedBy,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (r *ConnectionContractRepository) FindByCameraID(
	ctx context.Context,
	cameraID uuid.UUID,
) (*domain.CameraConnectionContract, error) {
	query := `
		SELECT
			id,
			camera_id,
			discovered_device_id,
			connection_type,
			ip_address,
			rtsp_url,
			onvif_url,
			username,
			encrypted_password,
			stream_profile_token,
			main_stream_url,
			sub_stream_url,
			snapshot_url,
			created_by,
			updated_by,
			created_at,
			updated_at
		FROM camera.camera_connection_contracts
		WHERE camera_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`

	contract := &domain.CameraConnectionContract{}

	err := r.db.QueryRow(ctx, query, cameraID).Scan(
		&contract.ID,
		&contract.CameraID,
		&contract.DiscoveredDeviceID,
		&contract.ConnectionType,
		&contract.IPAddress,
		&contract.RTSPUrl,
		&contract.ONVIFUrl,
		&contract.Username,
		&contract.EncryptedPassword,
		&contract.StreamProfileToken,
		&contract.MainStreamURL,
		&contract.SubStreamURL,
		&contract.SnapshotURL,
		&contract.CreatedBy,
		&contract.UpdatedBy,
		&contract.CreatedAt,
		&contract.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCameraNotFound
	}

	if err != nil {
		return nil, err
	}

	return contract, nil
}

func (r *ConnectionContractRepository) FindByDiscoveredDeviceID(
	ctx context.Context,
	discoveredDeviceID uuid.UUID,
) (*domain.CameraConnectionContract, error) {
	query := `
		SELECT
			id,
			camera_id,
			discovered_device_id,
			connection_type,
			ip_address,
			rtsp_url,
			onvif_url,
			username,
			encrypted_password,
			stream_profile_token,
			main_stream_url,
			sub_stream_url,
			snapshot_url,
			created_by,
			updated_by,
			created_at,
			updated_at
		FROM camera.camera_connection_contracts
		WHERE discovered_device_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`

	contract := &domain.CameraConnectionContract{}

	err := r.db.QueryRow(ctx, query, discoveredDeviceID).Scan(
		&contract.ID,
		&contract.CameraID,
		&contract.DiscoveredDeviceID,
		&contract.ConnectionType,
		&contract.IPAddress,
		&contract.RTSPUrl,
		&contract.ONVIFUrl,
		&contract.Username,
		&contract.EncryptedPassword,
		&contract.StreamProfileToken,
		&contract.MainStreamURL,
		&contract.SubStreamURL,
		&contract.SnapshotURL,
		&contract.CreatedBy,
		&contract.UpdatedBy,
		&contract.CreatedAt,
		&contract.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return contract, nil
}
