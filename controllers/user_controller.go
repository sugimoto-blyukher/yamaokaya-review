package controllers

import (
	"github.com/gin-gonic/gin"
	"go-prac/config"
	"go-prac/dtos"
	"go-prac/models"
	"net/http"
)

func CreateUser(c *gin.Context) {
	var req dtos.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ユーザーの登録に失敗しました",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": user,
	})
}

func FindUsers(c *gin.Context) {
	var users []models.User
	config.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{"count": len(users), "data": users})

}

func FindUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDが指定されていません"})
		return
	}

	var user models.User
	result := config.DB.First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ユーザーを取得しました",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
