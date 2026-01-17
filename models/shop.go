package models

import "gorm.io/gorm"

type Shop struct {
	gorm.Model
	Name string `json:"name"`
	Address string `json:"address"`
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}