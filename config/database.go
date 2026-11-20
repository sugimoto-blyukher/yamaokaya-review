package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"fmt"
	"go-prac/models"
)

var DB *gorm.DB

func OpenDatabase(path string) (*gorm.DB, error) {
    db, err := gorm.Open(
        sqlite.Open(path),
        &gorm.Config{},
    )
    if err != nil {
        return nil, fmt.Errorf(
            "DB接続に失敗: %w", err,
        )
    }

    err = db.AutoMigrate(
        &models.Shop{},
        &models.User{},
        &models.Review{},
    )

    if err != nil {
        return nil, fmt.Errorf(
            "マイグレーションに失敗: %w", err,
        )
    }

    return db, nil
}