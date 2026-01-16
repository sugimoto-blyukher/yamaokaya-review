package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	Name string `json:"name"`
	Score int `json:"score"`
	Body string `json:"body"`
}

