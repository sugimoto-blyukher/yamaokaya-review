package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"go-prac/models"
)

var DB *gorm.DB

func Connect() {
	var err error
	//グローバル変数のDBを代入
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("DB接続失敗")
	}

	DB.AutoMigrate(&models.Review{})
}