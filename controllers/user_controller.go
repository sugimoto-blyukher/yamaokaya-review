package controllers

import (
	"errors"
	"go-prac/config"
	"go-prac/dtos"
	"go-prac/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		"data": dtos.UserListItem{
			ID:   user.ID,
			Name: user.Name,
		},
	})
}

func FindAllUsers(c *gin.Context) {
	users := make([]dtos.UserListItem, 0)

	if err := config.DB.Model(&models.User{}).Select("id", "name").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ユーザーの取得に失敗しました",
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"message": "全てのユーザーを取得しました",
		"count":   len(users),
		"data":    users,
	})
}

func FindUser(c *gin.Context) {
	var user dtos.UserListItem

	id, err := strconv.ParseUint(c.Param("id"), 10, 63)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "不正な値です",
		})
		return
	}

	if err := config.DB.Model(&models.User{}).Select("id", "name").First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "ユーザーが見つかりません",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ユーザーの取得に失敗しました",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ユーザーを取得しました",
		"data":    user,
	})
}
