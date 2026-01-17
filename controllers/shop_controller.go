package controllers

import (
	"net/http"
	"go-prac/config" 
	"go-prac/models" 

	"github.com/gin-gonic/gin"
)

func CreateShop(c *gin.Context) {
	var newShop models.Shop 

	if err := c.ShouldBindJSON(&newShop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := config.DB.Create(&newShop)
	if result.Error !=nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return	
	}

	c.JSON(http.StatusOK, gin.H{"message": "店舗を保存しました", "data": newShop})
}

func FindShops(c *gin.Context) {
	var shops []models.Shop
	config.DB.Find(&shops)

	c.JSON(http.StatusOK, gin.H{"count":len(shops), "data": shops})
}