package event

import (
	"time"

	"github.com/google/uuid"
)

const (
	CameraCreatedEvent = "camera.created"
	CameraUpdatedEvent = "camera.updated"
	CameraDeletedEvent = "camera.deleted"
	CameraOnlineEvent  = "camera.online"
	CameraOfflineEvent = "camera.offline"
)

type CameraEvent struct {
	EventID   uuid.UUID  `json:"event_id"`
	EventType string     `json:"event_type"`
	CameraID  uuid.UUID  `json:"camera_id"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	Timestamp time.Time  `json:"timestamp"`
}
