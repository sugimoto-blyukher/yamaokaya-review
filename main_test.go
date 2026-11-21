package main

import (
	"bytes"
	"encoding/json"
	"go-prac/config"
	"go-prac/controllers"
	"go-prac/models"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func openTestDatabase(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := config.OpenDatabase(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("テストDBの初期化に失敗しました: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("テストDBの接続プールを取得できません: %v", err)
	}
	t.Cleanup(func() {
		if err := sqlDB.Close(); err != nil {
			t.Errorf("テストDBを閉じられません: %v", err)
		}
	})
	return db
}

func setupTestRouter(t *testing.T) *gin.Engine {
	t.Helper()

	previousDB := config.DB
	config.DB = openTestDatabase(t)
	t.Cleanup(func() { config.DB = previousDB })

	router := gin.New()
	router.POST("/reviews", controllers.CreateReview)
	router.GET("/reviews", controllers.FindAllReviews)
	router.GET("/reviews/:id", controllers.FindReview)
	router.POST("/shops", controllers.CreateShop)
	router.GET("/shops", controllers.FindShops)
	router.POST("/users", controllers.CreateUser)
	router.GET("/users", controllers.FindAllUsers)
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

func TestSQLiteForeignKeys(t *testing.T) {
	db := openTestDatabase(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("テストDBの接続プールを取得できません: %v", err)
	}

	ctx := t.Context()
	for i := 1; i <= 2; i++ {
		conn, err := sqlDB.Conn(ctx)
		if err != nil {
			t.Fatalf("接続%dの取得に失敗しました: %v", i, err)
		}
		// 1本目も保持し、2本目が別のSQLite接続になるようにする。
		defer conn.Close()

		var enabled int
		if err := conn.QueryRowContext(ctx, "PRAGMA foreign_keys").Scan(&enabled); err != nil {
			t.Fatalf("接続%dの外部キー設定を読み取れません: %v", i, err)
		}
		if enabled != 1 {
			t.Errorf("接続%dの PRAGMA foreign_keys = %d、期待値は 1", i, enabled)
		}
	}
}

func TestReviewControllers(t *testing.T) {
	router := setupTestRouter(t)

	status, empty := request(t, router, http.MethodGet, "/reviews", "")
	if status != http.StatusOK || empty["count"] != float64(0) || !reflect.DeepEqual(empty["data"], []any{}) {
		t.Fatalf("空のレビュー一覧: status=%d response=%v", status, empty)
	}

	status, shop := request(t, router, http.MethodPost, "/shops", `{"name":"山岡家","address":"札幌市","lat":43.0,"lng":141.0}`)
	if status != http.StatusOK {
		t.Fatalf("レビュー用店舗の作成に失敗しました: status=%d response=%v", status, shop)
	}
	status, user := request(t, router, http.MethodPost, "/users", `{"name":"太郎","email":"taro@example.com"}`)
	if status != http.StatusCreated {
		t.Fatalf("レビュー用ユーザーの作成に失敗しました: status=%d response=%v", status, user)
	}

	for _, tt := range []struct {
		name string
		body string
	}{
		{"missing shop", `{"name":"太郎","score":5,"body":"おいしい","shopID":999999,"userID":1}`},
		{"missing user", `{"name":"太郎","score":5,"body":"おいしい","shopID":1,"userID":999999}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			status, response := request(t, router, http.MethodPost, "/reviews", tt.body)
			if status != http.StatusInternalServerError {
				t.Errorf("存在しない参照先を指定した投稿: status=%d response=%v", status, response)
			}

			var count int64
			if err := config.DB.Model(&models.Review{}).Count(&count).Error; err != nil {
				t.Fatalf("レビュー件数の取得に失敗しました: %v", err)
			}
			if count != 0 {
				t.Errorf("存在しない参照先を指定したレビューが保存されています: count=%d", count)
			}
		})
	}

	status, created := request(t, router, http.MethodPost, "/reviews", `{"name":"太郎","score":5,"body":"おいしい","shopID":1,"userID":1}`)
	if status != http.StatusCreated {
		t.Fatalf("CreateReview: status=%d response=%v", status, created)
	}
	createdData, ok := created["data"].(map[string]any)
	if !ok || createdData["shopID"] != float64(1) || createdData["userID"] != float64(1) {
		t.Fatalf("CreateReviewの関連IDが不正です: response=%v", created)
	}

	wantReview := map[string]any{
		"name":   "太郎",
		"score":  float64(5),
		"shopID": float64(1),
		"userID": float64(1),
	}
	status, found := request(t, router, http.MethodGet, "/reviews/1", "")
	if status != http.StatusOK || found["message"] != "レビューを取得しました" || !reflect.DeepEqual(found["data"], wantReview) {
		t.Fatalf("FindReview: status=%d response=%v", status, found)
	}

	status, missing := request(t, router, http.MethodGet, "/reviews/999", "")
	if status != http.StatusNotFound || missing["error"] != "レビューが見つかりません" {
		t.Fatalf("存在しないレビューの取得: status=%d response=%v", status, missing)
	}

	status, list := request(t, router, http.MethodGet, "/reviews", "")
	if status != http.StatusOK || list["message"] != "全てのレビューを取得しました" ||
		list["count"] != float64(1) || !reflect.DeepEqual(list["data"], []any{wantReview}) {
		t.Fatalf("FindAllReviews: status=%d response=%v", status, list)
	}
}

func TestFindAllReviewsDatabaseError(t *testing.T) {
	router := setupTestRouter(t)
	if err := config.DB.Migrator().DropTable(&models.Review{}); err != nil {
		t.Fatalf("テストDBのレビューテーブルを削除できません: %v", err)
	}

	status, response := request(t, router, http.MethodGet, "/reviews", "")
	if status != http.StatusInternalServerError || response["error"] != "レビューの取得に失敗しました" {
		t.Fatalf("レビュー一覧の取得失敗時: status=%d response=%v", status, response)
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
	if status != http.StatusCreated {
		t.Fatalf("CreateUser: status=%d response=%v", status, created)
	}
	createdData, ok := created["data"].(map[string]any)
	if !ok || createdData["name"] != "花子" || createdData["email"] != "hanako@example.com" {
		t.Fatalf("CreateUserのレスポンスが不正です: response=%v", created)
	}

	status, found := request(t, router, http.MethodGet, "/users/1", "")
	data, ok := found["data"].(map[string]any)
	if status != http.StatusOK || !ok || len(data) != 2 || data["id"] != float64(1) || data["name"] != "花子" {
		t.Fatalf("FindUser: status=%d response=%v", status, found)
	}

	status, missing := request(t, router, http.MethodGet, "/users/999", "")
	if status != http.StatusNotFound {
		t.Fatalf("存在しないユーザーの取得: status=%d response=%v", status, missing)
	}

	status, list := request(t, router, http.MethodGet, "/users", "")
	users, ok := list["data"].([]any)
	if status != http.StatusOK || list["count"] != float64(1) || !ok || len(users) != 1 {
		t.Fatalf("FindAllUsers: status=%d response=%v", status, list)
	}
	user, ok := users[0].(map[string]any)
	if !ok || len(user) != 2 || user["id"] != float64(1) || user["name"] != "花子" {
		t.Fatalf("ユーザー一覧の項目が不正です: response=%v", list)
	}
}

func TestCreateUserValidation(t *testing.T) {
	router := setupTestRouter(t)

	status, response := request(t, router, http.MethodPost, "/users", `{"name":"花子"}`)
	if status != http.StatusBadRequest {
		t.Fatalf("必須項目不足時のステータスが不正です: status=%d response=%v", status, response)
	}

	var count int64
	if err := config.DB.Model(&models.User{}).Count(&count).Error; err != nil {
		t.Fatalf("ユーザー件数の取得に失敗しました: %v", err)
	}
	if count != 0 {
		t.Fatalf("バリデーションエラー時にユーザーが保存されています: count=%d", count)
	}
}

func TestCreateReviewValidation(t *testing.T) {
	router := setupTestRouter(t)

	status, response := request(t, router, http.MethodPost, "/reviews", `{"name":"太郎","score":6,"body":"おいしい","shopID":1,"userID":1}`)
	if status != http.StatusBadRequest {
		t.Fatalf("スコア範囲外のステータスが不正です: status=%d response=%v", status, response)
	}
}
