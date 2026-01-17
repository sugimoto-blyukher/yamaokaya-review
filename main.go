package main

import (
	"go-prac/config"
	"go-prac/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. DB接続開始
	config.Connect()

	// 2. サーバーの準備
	r := gin.Default()

	// 3. ルーティング（中身の処理はcontrollersに丸投げ）
	r.POST("/reviews", controllers.CreateReview)
	r.GET("/reviews", controllers.FindReviews)
	r.POST("/shops", controllers.CreateShop)
	r.POST("/users", controllers.CreateUser)
	// 4. 起動
	r.Run(":3000")
}
