package dto

type CreateBuildingRequest struct {
	Name         string   `json:"name" binding:"required"`
	Description  *string  `json:"description"`
	Address      *string  `json:"address"`
	Color        *string  `json:"color"`
	Opacity      *float64 `json:"opacity"`
	Icon         *string  `json:"icon"`
	LabelVisible *bool    `json:"label_visible"`
}

type UpdateBuildingRequest struct {
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	Address      *string  `json:"address"`
	Color        *string  `json:"color"`
	Opacity      *float64 `json:"opacity"`
	Icon         *string  `json:"icon"`
	LabelVisible *bool    `json:"label_visible"`
}
