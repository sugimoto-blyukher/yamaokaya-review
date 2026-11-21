package controllers

import (
	"go-prac/config"
	"go-prac/dtos"
	"go-prac/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateShop(c *gin.Context) {

	var req dtos.CreateShopRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	shop := models.Shop{
		Name:    req.Name,
		Address: req.Address,
		Lat:     req.Lat,
		Lng:     req.Lng,
	}

	if err := config.DB.Create(&shop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "店舗の登録に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "店舗を保存しました",
		"data": dtos.ShopListItem{
			ID:      shop.ID,
			Name:    shop.Name,
			Address: shop.Address,
		},
	})
}

func FindShops(c *gin.Context) {
	var shops []dtos.ShopListItem
	config.DB.Model(&models.Shop{}).Select("id", "name", "address").Find(&shops)

	c.JSON(http.StatusOK, gin.H{
		"message": "全ての店舗の取得に成功しました",
		"count": len(shops),
		"data": shops,
	})
}
