package event

import (
	"context"
)

type Publisher interface {
	PublishCameraEvent(ctx context.Context, event CameraEvent) error
	Close() error
}
