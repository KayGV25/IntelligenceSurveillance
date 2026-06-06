package event

import (
	"time"

	"github.com/google/uuid"
)

const (
	CameraCreatedEvent    = "camera.created"
	CameraUpdatedEvent    = "camera.updated"
	CameraDeletedEvent    = "camera.deleted"
	CameraOnlineEvent     = "camera.online"
	CameraOfflineEvent    = "camera.offline"
	CameraDiscoveredEvent = "camera.discovered"
	CameraConnectedEvent  = "camera.connected"
)

type CameraEvent struct {
	EventID            uuid.UUID  `json:"event_id"`
	EventType          string     `json:"event_type"`
	CameraID           *uuid.UUID `json:"camera_id,omitempty"`
	DiscoveredDeviceID *uuid.UUID `json:"discovered_device_id,omitempty"`
	IPAddress          string     `json:"ip_address,omitempty"`
	UserID             *uuid.UUID `json:"user_id,omitempty"`
	Timestamp          time.Time  `json:"timestamp"`
}
