package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	Name   string `json:"name" gorm:"not null"`
    Score  int    `json:"score" gorm:"not null;check:score >= 1 AND score <= 5"`
    Body   string `json:"body" gorm:"not null"`
    ShopID uint   `json:"shopID" gorm:"not null;index"`
    Shop   Shop   `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    UserID uint   `json:"userID" gorm:"not null;index"`
    User   User   `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

}

