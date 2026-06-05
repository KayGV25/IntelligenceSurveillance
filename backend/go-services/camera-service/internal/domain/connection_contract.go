package domain

import (
	"time"

	"github.com/google/uuid"
)

type ConnectionType string

const (
	ConnectionTypeONVIF          ConnectionType = "ONVIF"
	ConnectionTypeRTSPManual     ConnectionType = "RTSP_MANUAL"
	ConnectionTypeHTTPSnapshot   ConnectionType = "HTTP_SNAPSHOT"
	ConnectionTypeVendorSpecific ConnectionType = "VENDOR_SPECIFIC"
	ConnectionTypeHTTPMJPEG      ConnectionType = "HTTP_MJPEG"
)

type CameraConnectionContract struct {
	ID                 uuid.UUID      `json:"id"`
	CameraID           uuid.UUID      `json:"camera_id"`
	DiscoveredDeviceID *uuid.UUID     `json:"discovered_device_id"`
	ConnectionType     ConnectionType `json:"connection_type"`

	IPAddress string  `json:"ip_address"`
	RTSPUrl   *string `json:"rtsp_url"`
	ONVIFUrl  *string `json:"onvif_url"`

	Username          *string `json:"username"`
	EncryptedPassword *string `json:"-"`

	StreamProfileToken *string `json:"stream_profile_token"`
	MainStreamURL      *string `json:"main_stream_url"`
	SubStreamURL       *string `json:"sub_stream_url"`
	SnapshotURL        *string `json:"snapshot_url"`

	CreatedBy *uuid.UUID `json:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
