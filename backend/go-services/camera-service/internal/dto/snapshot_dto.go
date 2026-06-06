package dto

type SnapshotResponse struct {
	CameraID    string `json:"camera_id"`
	ObjectPath  string `json:"object_path"`
	ContentType string `json:"content_type"`
	Message     string `json:"message"`
}
