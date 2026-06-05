package dto

type ValidateStreamResponse struct {
	CameraID       string `json:"camera_id"`
	Status         string `json:"status"`
	ConnectionType string `json:"connection_type"`
	CheckedURL     string `json:"checked_url"`
	Message        string `json:"message"`
}
