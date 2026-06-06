package dto

type CameraHealthResponse struct {
	CameraID       string  `json:"camera_id"`
	Status         string  `json:"status"`
	ConnectionType string  `json:"connection_type"`
	MainStreamURL  *string `json:"main_stream_url"`
	SnapshotURL    *string `json:"snapshot_url"`
	Message        string  `json:"message"`
}
