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
			const wantError = "指定された店舗またはユーザーが存在しません"
			if status != http.StatusBadRequest || response["error"] != wantError {
				t.Errorf("status=%d response=%v, want status=%d error=%q",
					status, response, http.StatusBadRequest, wantError)
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

	wantReview := map[string]any{
		"id":     float64(1),
		"name":   "太郎",
		"score":  float64(5),
		"body":   "おいしい",
		"shopID": float64(1),
		"userID": float64(1),
	}
	status, created := request(t, router, http.MethodPost, "/reviews", `{"name":"太郎","score":5,"body":"おいしい","shopID":1,"userID":1}`)
	if status != http.StatusCreated || !reflect.DeepEqual(created["data"], wantReview) {
		t.Fatalf("CreateReview: status=%d response=%v", status, created)
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

func TestCreateDatabaseError(t *testing.T) {
	for _, tt := range []struct {
		path      string
		model     any
		body      string
		wantError string
	}{
		{
			path:      "users",
			model:     &models.User{},
			body:      `{"name":"花子","email":"hanako@example.com"}`,
			wantError: "ユーザーの登録に失敗しました",
		},
		{
			path:      "shops",
			model:     &models.Shop{},
			body:      `{"name":"山岡家","address":"札幌市","lat":43.0,"lng":141.0}`,
			wantError: "店舗の登録に失敗しました",
		},
		{
			path:      "reviews",
			model:     &models.Review{},
			body:      `{"name":"太郎","score":5,"body":"おいしい","shopID":1,"userID":1}`,
			wantError: "投稿に失敗しました",
		},
	} {
		t.Run(tt.path, func(t *testing.T) {
			router := setupTestRouter(t)
			if err := config.DB.Migrator().DropTable(tt.model); err != nil {
				t.Fatalf("テストDBの登録先テーブルを削除できません: %v", err)
			}

			status, response := request(t, router, http.MethodPost, "/"+tt.path, tt.body)
			if status != http.StatusInternalServerError || response["error"] != tt.wantError {
				t.Fatalf("status=%d response=%v, want status=%d error=%q",
					status, response, http.StatusInternalServerError, tt.wantError)
			}
		})
	}
}

func TestSearchDatabaseError(t *testing.T) {
	for _, tt := range []struct {
		path      string
		model     any
		wantError string
	}{
		{"reviews", &models.Review{}, "レビューの取得に失敗しました"},
		{"reviews/1", &models.Review{}, "レビューの取得に失敗しました"},
		{"shops", &models.Shop{}, "店舗の取得に失敗しました"},
		{"users", &models.User{}, "ユーザーの取得に失敗しました"},
		{"users/1", &models.User{}, "ユーザーの取得に失敗しました"},
	} {
		t.Run(tt.path, func(t *testing.T) {
			router := setupTestRouter(t)
			if err := config.DB.Migrator().DropTable(tt.model); err != nil {
				t.Fatalf("テストDBの検索対象テーブルを削除できません: %v", err)
			}

			status, response := request(t, router, http.MethodGet, "/"+tt.path, "")
			if status != http.StatusInternalServerError || response["error"] != tt.wantError {
				t.Fatalf("status=%d response=%v, want status=%d error=%q",
					status, response, http.StatusInternalServerError, tt.wantError)
			}
		})
	}
}

func TestSearchInvalidID(t *testing.T) {
	router := setupTestRouter(t)

	for _, resource := range []string{"reviews", "users"} {
		for _, id := range []string{"abc", "0", "-1", "9223372036854775808"} {
			path := resource + "/" + id
			t.Run(path, func(t *testing.T) {
				status, response := request(t, router, http.MethodGet, "/"+path, "")
				if status != http.StatusBadRequest {
					t.Fatalf("status=%d response=%v, want status=%d", status, response, http.StatusBadRequest)
				}
			})
		}
	}
}

func TestShopControllers(t *testing.T) {
	router := setupTestRouter(t)

	status, empty := request(t, router, http.MethodGet, "/shops", "")
	if status != http.StatusOK || empty["count"] != float64(0) || !reflect.DeepEqual(empty["data"], []any{}) {
		t.Fatalf("空の店舗一覧: status=%d response=%v", status, empty)
	}

	wantShop := map[string]any{
		"id":      float64(1),
		"name":    "山岡家",
		"address": "札幌市",
	}
	status, created := request(t, router, http.MethodPost, "/shops", `{"name":"山岡家","address":"札幌市","lat":43.0,"lng":141.0}`)
	if status != http.StatusOK || created["message"] != "店舗を保存しました" || !reflect.DeepEqual(created["data"], wantShop) {
		t.Fatalf("CreateShop: status=%d response=%v", status, created)
	}

	status, list := request(t, router, http.MethodGet, "/shops", "")
	if status != http.StatusOK || list["count"] != float64(1) || !reflect.DeepEqual(list["data"], []any{wantShop}) {
		t.Fatalf("FindShops: status=%d response=%v", status, list)
	}
}

func TestUserControllers(t *testing.T) {
	router := setupTestRouter(t)

	status, empty := request(t, router, http.MethodGet, "/users", "")
	if status != http.StatusOK || empty["count"] != float64(0) || !reflect.DeepEqual(empty["data"], []any{}) {
		t.Fatalf("空のユーザー一覧: status=%d response=%v", status, empty)
	}

	wantUser := map[string]any{
		"id":   float64(1),
		"name": "花子",
	}
	status, created := request(t, router, http.MethodPost, "/users", `{"name":"花子","email":"hanako@example.com"}`)
	if status != http.StatusCreated || !reflect.DeepEqual(created["data"], wantUser) {
		t.Fatalf("CreateUser: status=%d response=%v", status, created)
	}

	status, duplicate := request(t, router, http.MethodPost, "/users", `{"name":"別のユーザー","email":"hanako@example.com"}`)
	const wantDuplicateError = "メールアドレスが重複しています"
	if status != http.StatusConflict || duplicate["error"] != wantDuplicateError {
		t.Fatalf("status=%d response=%v, want status=%d error=%q",
			status, duplicate, http.StatusConflict, wantDuplicateError)
	}

	status, found := request(t, router, http.MethodGet, "/users/1", "")
	if status != http.StatusOK || !reflect.DeepEqual(found["data"], wantUser) {
		t.Fatalf("FindUser: status=%d response=%v", status, found)
	}

	status, missing := request(t, router, http.MethodGet, "/users/999", "")
	if status != http.StatusNotFound {
		t.Fatalf("存在しないユーザーの取得: status=%d response=%v", status, missing)
	}

	status, list := request(t, router, http.MethodGet, "/users", "")
	if status != http.StatusOK || list["count"] != float64(1) || !reflect.DeepEqual(list["data"], []any{wantUser}) {
		t.Fatalf("FindAllUsers: status=%d response=%v", status, list)
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
