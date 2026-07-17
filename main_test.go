package main

import (
	"bytes"
	"encoding/json"
	"go-prac/config"
	"go-prac/controllers"
	"go-prac/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter(t *testing.T) *gin.Engine {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("テストDBへの接続に失敗しました: %v", err)
	}
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("外部キー制約の有効化に失敗しました: %v", err)
	}
	if err := db.AutoMigrate(&models.Shop{}, &models.User{}, &models.Review{}); err != nil {
		t.Fatalf("テストDBの作成に失敗しました: %v", err)
	}
	config.DB = db

	router := gin.New()
	router.POST("/reviews", controllers.CreateReview)
	router.GET("/reviews", controllers.FindReviews)
	router.GET("/reviews/:id", controllers.FindReview)
	router.POST("/shops", controllers.CreateShop)
	router.GET("/shops", controllers.FindShops)
	router.POST("/users", controllers.CreateUser)
	router.GET("/users", controllers.FindUsers)
	router.GET("/users/:id", controllers.FindUser)
	return router
}

func request(t *testing.T, router http.Handler, method, path, body string) (int, map[string]any) {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("レスポンスJSONの解析に失敗しました: %v", err)
	}
	return w.Code, response
}

func TestReviewControllers(t *testing.T) {
	router := setupTestRouter(t)

	status, shop := request(t, router, http.MethodPost, "/shops", `{"name":"山岡家","address":"札幌市","lat":43.0,"lng":141.0}`)
	if status != http.StatusOK {
		t.Fatalf("レビュー用店舗の作成に失敗しました: status=%d response=%v", status, shop)
	}
	status, user := request(t, router, http.MethodPost, "/users", `{"name":"太郎","email":"taro@example.com"}`)
	if status != http.StatusOK {
		t.Fatalf("レビュー用ユーザーの作成に失敗しました: status=%d response=%v", status, user)
	}

	status, created := request(t, router, http.MethodPost, "/reviews", `{"name":"太郎","score":5,"body":"おいしい","shopID":1,"userID":1}`)
	if status != http.StatusOK || created["message"] != "保存しました" {
		t.Fatalf("CreateReview: status=%d response=%v", status, created)
	}
	createdData, ok := created["data"].(map[string]any)
	if !ok || createdData["shopID"] != float64(1) || createdData["userID"] != float64(1) {
		t.Fatalf("CreateReviewの関連IDが不正です: response=%v", created)
	}

	status, found := request(t, router, http.MethodGet, "/reviews/1", "")
	data, ok := found["data"].(map[string]any)
	if status != http.StatusOK || !ok || data["name"] != "太郎" || data["score"] != float64(5) {
		t.Fatalf("FindReview: status=%d response=%v", status, found)
	}

	status, list := request(t, router, http.MethodGet, "/reviews", "")
	if status != http.StatusOK || list["count"] != float64(1) {
		t.Fatalf("FindReviews: status=%d response=%v", status, list)
	}
}

func TestShopControllers(t *testing.T) {
	router := setupTestRouter(t)

	status, created := request(t, router, http.MethodPost, "/shops", `{"name":"山岡家","address":"札幌市","lat":43.0,"lng":141.0}`)
	if status != http.StatusOK || created["message"] != "店舗を保存しました" {
		t.Fatalf("CreateShop: status=%d response=%v", status, created)
	}

	status, list := request(t, router, http.MethodGet, "/shops", "")
	if status != http.StatusOK || list["count"] != float64(1) {
		t.Fatalf("FindShops: status=%d response=%v", status, list)
	}
}

func TestUserControllers(t *testing.T) {
	router := setupTestRouter(t)

	status, created := request(t, router, http.MethodPost, "/users", `{"name":"花子","email":"hanako@example.com"}`)
	if status != http.StatusOK || created["message"] != "保存しました" {
		t.Fatalf("CreateUser: status=%d response=%v", status, created)
	}

	status, found := request(t, router, http.MethodGet, "/users/1", "")
	data, ok := found["data"].(map[string]any)
	if status != http.StatusOK || !ok || data["email"] != "hanako@example.com" {
		t.Fatalf("FindUser: status=%d response=%v", status, found)
	}

	status, list := request(t, router, http.MethodGet, "/users", "")
	if status != http.StatusOK || list["count"] != float64(1) {
		t.Fatalf("FindUsers: status=%d response=%v", status, list)
	}
}
