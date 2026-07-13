package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

r.GET("/ping", func(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H {
		"message" : "pong",
	})
})

func TestPingRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := setupRouter()

	w := httptest.NewRecoder
}