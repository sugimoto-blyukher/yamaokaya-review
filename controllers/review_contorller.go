package controllers

import (
	"net/http"
	"go-prac/config" // ★DBを使うために読み込む
	"go-prac/models" // ★構造体を使うために読み込む

	"github.com/gin-gonic/gin"
)

// レビュー投稿 (POST)
func CreateReview(c *gin.Context) {
	var newReview models.Review // modelsパッケージのReviewを使う

	if err := c.ShouldBindJSON(&newReview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// configパッケージのDB変数を使う
	result := config.DB.Create(&newReview)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "保存しました", "data": newReview})
}

// 全レビュー取得 (GET)
func FindReviews(c *gin.Context) {
	var reviews []models.Review
	config.DB.Find(&reviews)
	
	c.JSON(http.StatusOK, gin.H{"count": len(reviews), "data": reviews})
}