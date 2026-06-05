package domain

import (
	"time"

	"github.com/google/uuid"
)

type DiscoveryMethod string

const (
	DiscoveryMethodCIDRScan DiscoveryMethod = "CIDR_SCAN"
)

type DiscoveredDeviceStatus string

const (
	DiscoveredDeviceStatusDiscovered DiscoveredDeviceStatus = "DISCOVERED"
	DiscoveredDeviceStatusConnecting DiscoveredDeviceStatus = "CONNECTING"
	DiscoveredDeviceStatusConnected  DiscoveredDeviceStatus = "CONNECTED"
	DiscoveredDeviceStatusFailed     DiscoveredDeviceStatus = "FAILED"
	DiscoveredDeviceStatusIgnored    DiscoveredDeviceStatus = "IGNORED"
)

type DiscoveredDevice struct {
	ID              uuid.UUID              `json:"id"`
	IPAddress       string                 `json:"ip_address"`
	MACAddress      *string                `json:"mac_address"`
	Hostname        *string                `json:"hostname"`
	Manufacturer    *string                `json:"manufacturer"`
	Model           *string                `json:"model"`
	FirmwareVersion *string                `json:"firmware_version"`
	ONVIFSupported  bool                   `json:"onvif_supported"`
	RTSPSupported   bool                   `json:"rtsp_supported"`
	HTTPSupported   bool                   `json:"http_supported"`
	HTTPPort        *int                   `json:"http_port"`
	RTSPPort        *int                   `json:"rtsp_port"`
	ONVIFPort       *int                   `json:"onvif_port"`
	DiscoveryMethod DiscoveryMethod        `json:"discovery_method"`
	Status          DiscoveredDeviceStatus `json:"status"`
	DiscoveredAt    time.Time              `json:"discovered_at"`
	LastSeenAt      *time.Time             `json:"last_seen_at"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}
