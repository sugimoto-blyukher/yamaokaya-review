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
			"error": "指定された店舗はありません",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "店舗を保存しました", "data": shop,
	})
}

func FindShops(c *gin.Context) {
	var shops []models.Shop
	config.DB.Find(&shops)

	c.JSON(http.StatusOK, gin.H{"count": len(shops), "data": shops})
}
