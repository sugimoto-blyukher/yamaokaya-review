package dtos

type CreateShopRequest struct {
	Name    string  `json:"name" binding:"required"`
	Address string  `json:"address" binding:"required"`
	Lat     float64 `json:"lat" binding:"required,latitude"`
	Lng     float64 `json:"lng" binding:"required,longitude"`
}
