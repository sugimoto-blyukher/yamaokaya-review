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

	DB.AutoMigrate(&models.Review{}, &models.Shop{}, &models.User{})
}

func InitDB() {
	//DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Review{})
	println("Created Review table")
	
	DB.AutoMigrate(&models.Shop{})
	println("Created Shop table")
}