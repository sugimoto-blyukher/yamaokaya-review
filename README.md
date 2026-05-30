# yamaokaya-review

個人サイト向けの「山岡家レビューAPI」です。

自分が訪れた山岡家の店舗情報やレビューを記録し、個人サイト上で表示するためのバックエンドAPIとして作成しています。

## 概要

このAPIでは、山岡家の店舗情報、レビュー、ユーザー情報を管理できます。

想定している用途は、個人サイトやポートフォリオサイトからAPIを呼び出し、以下のような情報を表示することです。

* 訪問した山岡家の店舗一覧
* 店舗ごとのレビュー
* 点数・感想・店舗情報
* 今後の訪問記録

現時点では学習・開発途中のAPIであり、GoによるWeb API開発、Ginによるルーティング、GORMによるDB操作の練習も兼ねています。

## 技術スタック

* Go
* Gin
* GORM
* SQLite

## 主な機能

### レビュー機能

* レビューの作成
* レビュー一覧の取得
* ID指定によるレビュー取得

### 店舗機能

* 店舗情報の作成
* 店舗名、住所、緯度、経度の保存

### ユーザー機能

* ユーザーの作成
* ユーザー一覧の取得
* ID指定によるユーザー取得

## データ構造

### Review

```json
{
  "id": 1,
  "name": "函館鍛治店レビュー",
  "score": 4,
  "body": "深夜に食べる醤油ラーメンがうまい",
  "shopID": 1
}
```

### Shop

```json
{
  "name": "ラーメン山岡家 函館鍛治店",
  "address": "北海道函館市鍛治...",
  "lat": 41.0000,
  "lng": 140.0000
}
```

### User

```json
{
  "name": "sugimoto",
  "email": "example@example.com"
}
```

## API エンドポイント

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

---

### レビュー作成

```http
POST /reviews
```

リクエスト例:

```json
{
  "name": "醤油ラーメンレビュー",
  "score": 4,
  "body": "濃い味で深夜に食べるとうまい",
  "shopID": 1
}
```

---

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
      "id": 1,
      "name": "醤油ラーメンレビュー",
      "score": 4,
      "body": "濃い味で深夜に食べるとうまい",
      "shopID": 1
    }
  ]
}
```

---

### レビュー個別取得

```http
GET /reviews/:id
```

例:

```http
GET /reviews/1
```

---

### 店舗作成

```http
POST /shops
```

リクエスト例:

```json
{
  "name": "ラーメン山岡家 函館鍛治店",
  "address": "北海道函館市鍛治...",
  "lat": 41.0000,
  "lng": 140.0000
}
```

---

### ユーザー作成

```http
POST /users
```

リクエスト例:

```json
{
  "name": "sugimoto",
  "email": "example@example.com"
}
```

---

### ユーザー一覧取得

```http
GET /users
```

---

### ユーザー個別取得

```http
GET /users/:id
```

例:

```http
GET /users/1
```

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

デフォルトでは `:3000` で起動します。

```text
http://localhost:3000
```

## DBについて

このAPIではSQLiteを使用しています。

起動時に `test.db` が作成され、GORMの `AutoMigrate` によって必要なテーブルが作成されます。

現時点で作成される主なテーブルは以下です。

* reviews
* shops
* users

## curlでの動作確認

### 店舗を作成

```sh
curl -X POST http://localhost:3000/shops \
  -H "Content-Type: application/json" \
  -d '{
    "name": "ラーメン山岡家 函館鍛治店",
    "address": "北海道函館市鍛治...",
    "lat": 41.0000,
    "lng": 140.0000
  }'
```

### レビューを作成

```sh
curl -X POST http://localhost:3000/reviews \
  -H "Content-Type: application/json" \
  -d '{
    "name": "醤油ラーメンレビュー",
    "score": 4,
    "body": "濃い味でうまい",
    "shopID": 1
  }'
```

### レビュー一覧を取得

```sh
curl http://localhost:3000/reviews
```

### ユーザーを作成

```sh
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "sugimoto",
    "email": "example@example.com"
  }'
```

## 今後の実装予定

* 店舗一覧取得APIの追加
* 店舗ごとのレビュー取得
* レビューの更新・削除
* 評価点のバリデーション
* 店舗画像やラーメン画像の保存
* 個人サイトのフロントエンドとの連携
* Docker対応
* テスト追加
* API仕様書の整備
* 認証機能の追加

## 開発目的

このリポジトリは、単なるCRUD APIの練習ではなく、個人サイト上で自分の山岡家レビューを表示するためのバックエンドとして作成しています。

また、Go、Gin、GORM、SQLiteを使ったWeb API開発の基本を学び、将来的には個人サイトと連携できる形まで発展させることを目的としています。
