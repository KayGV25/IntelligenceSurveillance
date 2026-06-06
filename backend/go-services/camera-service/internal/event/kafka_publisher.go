package event

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

const CameraEventsTopic = "camera.events"

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers string) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(strings.Split(brokers, ",")...),
			Topic:        CameraEventsTopic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
			Async:        false,
		},
	}
}

func (p *KafkaPublisher) PublishCameraEvent(ctx context.Context, event CameraEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	key := event.EventID.String()

	if event.CameraID != nil {
		key = event.CameraID.String()
	} else if event.DiscoveredDeviceID != nil {
		key = event.DiscoveredDeviceID.String()
	} else if event.IPAddress != "" {
		key = event.IPAddress
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: payload,
		Time:  event.Timestamp,
	})
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}
