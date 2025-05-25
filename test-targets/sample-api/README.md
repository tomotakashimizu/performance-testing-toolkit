# Sample API

負荷試験用のサンプル API サーバーです。

## 概要

この API は、負荷試験のターゲットとして使用するために作成されたシンプルな REST API です。
様々なパフォーマンス特性を持つエンドポイントを提供し、負荷試験の検証に使用できます。

## エンドポイント

### 基本エンドポイント

- `GET /` - ホームページ（API の基本情報を返す）
- `GET /health` - ヘルスチェック

### アイテム管理 API

- `GET /api/v1/items` - 全アイテムの取得
- `POST /api/v1/items` - 新しいアイテムの作成
- `GET /api/v1/items/{id}` - 特定のアイテムの取得
- `PUT /api/v1/items/{id}` - アイテムの更新
- `DELETE /api/v1/items/{id}` - アイテムの削除

### 負荷試験用エンドポイント

- `GET /api/v1/slow` - 意図的に遅いエンドポイント（2 秒の固定遅延）
- `GET /api/v1/random-delay` - ランダムな遅延を持つエンドポイント（0-1 秒）
- `GET /api/v1/cpu-intensive` - CPU 集約的な処理を行うエンドポイント

## 使用例

### アイテムの作成

```bash
curl -X POST http://localhost:8080/api/v1/items \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Item", "description": "This is a test item"}'
```

### アイテムの取得

```bash
curl http://localhost:8080/api/v1/items
```

## 起動方法

### Docker を使用

```bash
docker build -t sample-api .
docker run -p 8080:8080 sample-api
```

### ローカルで直接実行

```bash
go mod download
go run main.go
```

## 負荷試験での使用

この API は以下のような負荷試験シナリオで使用できます：

1. **基本的な GET リクエスト**: `/` や `/health` エンドポイント
2. **CRUD 操作**: `/api/v1/items` エンドポイント群
3. **レイテンシ試験**: `/api/v1/slow` や `/api/v1/random-delay`
4. **CPU 負荷試験**: `/api/v1/cpu-intensive`

各エンドポイントは異なるパフォーマンス特性を持つため、様々な負荷試験シナリオに対応できます。
