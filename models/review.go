package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	ID uint `json:"id"`
	Name string `json:"name"`
	Score int `json:"score"`
	Body string `json:"body"`
	ShopID uint `json:"shopID"`
	//UserID uint `json:"userID"`
}

