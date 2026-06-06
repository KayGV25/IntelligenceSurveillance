package dto

type StreamInfoResponse struct {
	CameraID       string  `json:"camera_id"`
	ConnectionType string  `json:"connection_type"`
	StreamMode     string  `json:"stream_mode"`
	StreamURL      *string `json:"stream_url"`
	SnapshotURL    *string `json:"snapshot_url"`
	RequiresProxy  bool    `json:"requires_proxy"`
	Message        string  `json:"message"`
}
