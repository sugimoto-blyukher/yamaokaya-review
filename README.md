# yamaokaya-review

個人サイト向けの「山岡家レビューAPI」です。

訪問した山岡家の店舗、レビュー、ユーザーを記録するためのバックエンドAPIです。GoによるWeb API開発、Ginによるルーティング、GORMによるデータベース操作の学習も目的としています。

## 技術スタック

- Go 1.25
- Gin
- GORM
- SQLite

## 現在の機能

- レビューの作成
- レビュー一覧の取得
- IDを指定したレビューの取得
- 店舗の作成
- ユーザーの作成
- ユーザー一覧の取得
- IDを指定したユーザーの取得

## データモデル

各モデルは `gorm.Model` を埋め込んでいるため、`ID`、`CreatedAt`、`UpdatedAt`、`DeletedAt` を持ちます。

### Review

| フィールド | 型 | 説明 |
| --- | --- | --- |
| `name` | string | レビュー名 |
| `score` | int | 評価点。DB上で1〜5に制限 |
| `body` | string | レビュー本文 |
| `shopID` | uint | レビュー対象店舗のID |
| `userID` | uint | 投稿ユーザーのID |

`Review` は `Shop` と `User` に関連付けられています。店舗またはユーザーを削除した場合、関連するレビューも削除する設定です。

### Shop

| フィールド | 型 | 説明 |
| --- | --- | --- |
| `name` | string | 店舗名 |
| `address` | string | 住所 |
| `lat` | float64 | 緯度 |
| `lng` | float64 | 経度 |

### User

| フィールド | 型 | 説明 |
| --- | --- | --- |
| `name` | string | ユーザー名 |
| `email` | string | メールアドレス。DB上で一意 |

## APIエンドポイント

ベースURLは次のとおりです。

```text
http://localhost:8080
```

### ヘルスチェック

```http
GET /
```

レスポンス例:

```json
{
  "message": "Hello World"
}
```

### 店舗作成

レビューを作成する前に、対象となる店舗を作成します。

```http
POST /shops
Content-Type: application/json
```

リクエスト例:

```json
{
  "name": "ラーメン山岡家 函館鍛治店",
  "address": "北海道函館市鍛治...",
  "lat": 41.0,
  "lng": 140.0
}
```

成功時はステータス `200 OK` と、保存された店舗を返します。

### ユーザー作成

レビューを作成する前に、投稿するユーザーを作成します。

```http
POST /users
Content-Type: application/json
```

リクエスト例:

```json
{
  "name": "sugimoto",
  "email": "example@example.com"
}
```

成功時はステータス `200 OK` と、保存されたユーザーを返します。同じメールアドレスは複数登録できません。

### ユーザー一覧取得

```http
GET /users
```

レスポンスは件数を表す `count` と、ユーザーの配列である `data` を含みます。

### ユーザー個別取得

```http
GET /users/:id
```

例:

```http
GET /users/1
```

対象が存在する場合は `200 OK`、存在しない場合は `404 Not Found` を返します。

### レビュー作成

`shopID` と `userID` には、先に作成した店舗とユーザーのIDを指定します。

```http
POST /reviews
Content-Type: application/json
```

リクエスト例:

```json
{
  "name": "醤油ラーメンレビュー",
  "score": 4,
  "body": "濃い味で深夜に食べるとうまい",
  "shopID": 1,
  "userID": 1
}
```

成功時はステータス `200 OK` と、保存されたレビューを返します。`score` のDB制約に違反した場合など、保存に失敗すると `500 Internal Server Error` を返します。

### レビュー一覧取得

```http
GET /reviews
```

レスポンス例:

```json
{
  "count": 1,
  "data": [
    {
      "ID": 1,
      "CreatedAt": "2026-01-01T00:00:00Z",
      "UpdatedAt": "2026-01-01T00:00:00Z",
      "DeletedAt": null,
      "name": "醤油ラーメンレビュー",
      "score": 4,
      "body": "濃い味で深夜に食べるとうまい",
      "shopID": 1,
      "userID": 1
    }
  ]
}
```

### レビュー個別取得

```http
GET /reviews/:id
```

例:

```http
GET /reviews/1
```

レスポンス例:

```json
{
  "message": "レビューを取得しました",
  "data": {
    "id": 1,
    "name": "醤油ラーメンレビュー",
    "score": 4,
    "body": "濃い味で深夜に食べるとうまい",
    "shopID": 1
  }
}
```

対象が存在しない場合は `404 Not Found` を返します。

## セットアップ

### 1. リポジトリをクローン

```sh
git clone https://github.com/sugimoto-blyukher/yamaokaya-review.git
cd yamaokaya-review
```

### 2. 依存関係を取得

```sh
go mod tidy
```

### 3. サーバーを起動

```sh
go run main.go
```

サーバーは `http://localhost:8080` で起動します。起動時にSQLiteの `test.db` が作成され、GORMの `AutoMigrate` によって次のテーブルが作成されます。

- `reviews`
- `shops`
- `users`

## curlでの動作確認

店舗、ユーザー、レビューの順に作成します。

```sh
curl -X POST http://localhost:8080/shops \
  -H "Content-Type: application/json" \
  -d '{
    "name": "ラーメン山岡家 函館鍛治店",
    "address": "北海道函館市鍛治...",
    "lat": 41.0,
    "lng": 140.0
  }'
```

```sh
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "sugimoto",
    "email": "example@example.com"
  }'
```

```sh
curl -X POST http://localhost:8080/reviews \
  -H "Content-Type: application/json" \
  -d '{
    "name": "醤油ラーメンレビュー",
    "score": 4,
    "body": "濃い味でうまい",
    "shopID": 1,
    "userID": 1
  }'
```

```sh
curl http://localhost:8080/reviews
```

## テスト

コントローラーのテストでは、SQLiteのインメモリDBを使用しています。

```sh
go test ./...
```

現在、レビュー、店舗、ユーザーの作成・一覧取得・個別取得に対応するハンドラーをテストしています。

## 今後の実装予定

- `GET /shops` のルーティング追加
- 店舗ごとのレビュー取得
- レビューの更新・削除
- API入力時の評価点バリデーション
- エラーレスポンスとHTTPステータスの整理
- 店舗画像やラーメン画像の保存
- 個人サイトのフロントエンドとの連携
- Docker対応
- テストの拡充
- API仕様書の整備
- 認証機能の追加
