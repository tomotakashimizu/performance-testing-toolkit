# Vegeta Configuration Files

Vegeta の負荷試験で使用するターゲット設定ファイルです。

## ファイル一覧

### local_test_config.txt

ローカル環境の `sample-api` を対象とした負荷試験設定です。
`test-targets/sample-api` が起動している状態で使用します。

### sample_config.txt

外部 API を対象とした負荷試験設定の例です。
実際の負荷試験を行う際は、このファイルをコピーして対象 API に合わせて編集してください。

### body1.json, body2.json

リクエストボディ用の JSON ファイルです。
POST/PUT リクエストで使用する JSON データを格納しています。

## 設定ファイルの形式

設定ファイルは Vegeta の HTTP 形式（txt 形式）で記述します。

```
HTTPメソッド URL
ヘッダー名: ヘッダー値
@ボディファイル名

HTTPメソッド URL
ヘッダー名: ヘッダー値
```

### フィールドの説明

- **HTTP メソッド URL** (必須): HTTP メソッドとリクエスト先 URL（1 行で記述）
- **ヘッダー** (オプション): HTTP ヘッダー（各行に 1 つずつ記述）
- **@ファイル名** (オプション): リクエストボディを含む JSON ファイルへの参照

### リクエストボディの指定

POST や PUT リクエストで JSON ボディを送信する場合、別途 JSON ファイルを作成し、`@ファイル名`で参照します。

```bash
# JSONファイルを作成
echo '{"name": "Test Item", "description": "This is a test item"}' > test_data.json
```

設定ファイルでは以下のように参照：

```
POST http://localhost:8080/api/items
Content-Type: application/json
@test_data.json
```

## 使用例

### 基本的な GET リクエスト

```
GET http://localhost:8080/api/v1/items
```

### 認証ヘッダー付き GET リクエスト

```
GET https://api.example.com/v1/users
Authorization: Bearer your-token-here
```

### JSON ボディ付き POST リクエスト

まず、JSON ファイルを作成：

```bash
echo '{"name": "Test Item", "description": "This is a test item"}' > item_data.json
```

設定ファイルに記述：

```
POST http://localhost:8080/api/v1/items
Content-Type: application/json
@item_data.json
```

## カスタム設定ファイルの作成

1. `sample_config.txt` をコピーして新しいファイルを作成
2. 対象 API のエンドポイントに合わせて URL とメソッドを変更
3. 必要に応じて認証ヘッダーを追加
4. リクエストボディが必要な場合は、JSON ファイルを作成して `@ファイル名` で参照
5. `load_test.sh` スクリプトで新しい設定ファイルを指定

```bash
# カスタム設定ファイルを使用
./load_test.sh -c custom_config.txt
```

## 注意事項

- URL は実際にアクセス可能なエンドポイントを指定してください
- 認証が必要な API の場合は、有効なトークンを設定してください
- リクエストボディが必要な場合は、JSON ファイルを作成して参照してください
- 負荷試験を行う前に、対象 API の利用規約を確認してください
