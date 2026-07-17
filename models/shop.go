package models

import "gorm.io/gorm"

type Shop struct {
	gorm.Model
    Name     string   `json:"name" gorm:"not null;index"`
    Address  string   `json:"address" gorm:"not null"`
    Lat      float64  `json:"lat" gorm:"not null"`
    Lng      float64  `json:"lng" gorm:"not null"`
    Reviews  []Review `json:"-"`
}