package main

import (
	"go-prac/config"
	"go-prac/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. DB接続開始
	config.Connect()

	config.InitDB()

	// 2. サーバーの準備
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	// 3. ルーティング（中身の処理はcontrollersに丸投げ）
	r.POST("/reviews", controllers.CreateReview)
	r.GET("/reviews", controllers.FindReviews)
	r.GET("/reviews/:id", controllers.FindReview)
	r.POST("/shops", controllers.CreateShop)
	r.POST("/users", controllers.CreateUser)
	r.GET("/users", controllers.FindUsers)
	r.GET("/users/:id", controllers.FindUser)

	// 4. 起動
	r.Run(":8080")
}
