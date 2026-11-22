package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-prac/config" // DBを使うために読み込む
	"go-prac/dtos"
	"go-prac/models" // 構造体を使うために読み込む
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// レビュー投稿 (POST)
func CreateReview(c *gin.Context) {
	var req dtos.CreateReviewRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	review := models.Review{
		Name:   req.Name,
		Score:  req.Score,
		Body:   req.Body,
		ShopID: req.ShopID,
		UserID: req.UserID,
	}

	if err := config.DB.Create(&review).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "指定された店舗またはユーザーが存在しません",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "投稿に失敗しました",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": dtos.ReviewListItem{
			ID:     review.ID,
			Name:   review.Name,
			Score:  review.Score,
			Body:   review.Body,
			ShopID: review.ShopID,
			UserID: review.UserID,
		},
	})

	//c.JSON(http.StatusOK, gin.H{"message": "保存しました", "data": newReview})
}

// 全レビュー取得 (GET)
func FindAllReviews(c *gin.Context) {
	reviews := make([]dtos.ReviewListItem, 0)

	if err := config.DB.Model(&models.Review{}).Select("id", "name", "score", "body", "shop_id", "user_id").Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "レビューの取得に失敗しました",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "全てのレビューを取得しました",
		"count":   len(reviews),
		"data":    reviews,
	})
}

func FindReview(c *gin.Context) {
	var review dtos.ReviewListItem

	id, err := strconv.ParseUint(c.Param("id"), 10, 63)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "不正な値です",
		})
		return
	}

	if err := config.DB.Model(&models.Review{}).Select("id", "name", "score", "body", "shop_id", "user_id").First(&review, "id=?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "レビューが見つかりません",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "レビューの取得に失敗しました",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "レビューを取得しました",
		"data":    review,
	})

}
