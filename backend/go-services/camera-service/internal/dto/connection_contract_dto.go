package dto

type ConnectDiscoveredDeviceRequest struct {
	Name            string  `json:"name" binding:"required"`
	Description     *string `json:"description"`
	Username        *string `json:"username"`
	Password        *string `json:"password"`
	PreferredStream string  `json:"preferred_stream"`
}

type ConnectDiscoveredDeviceResponse struct {
	Camera   interface{} `json:"camera"`
	Contract interface{} `json:"contract"`
}
