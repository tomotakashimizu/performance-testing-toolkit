# Test Targets

負荷試験のターゲットとなる API サーバーを提供するディレクトリです。

## 概要

このディレクトリには、負荷試験の対象となるサンプル API が含まれています。
実際の負荷試験を行う前に、これらのサンプル API を使用して負荷試験の手順や設定を検証できます。

## 含まれるサービス

### sample-api

Go 言語で作成されたシンプルな REST API サーバーです。
様々なパフォーマンス特性を持つエンドポイントを提供し、負荷試験のターゲットとして使用できます。

詳細は [sample-api/README.md](./sample-api/README.md) を参照してください。

## 起動方法

### Docker Compose を使用（推奨）

```bash
# test-targetsディレクトリに移動
cd test-targets

# サービスを起動
docker-compose up -d

# ログを確認
docker-compose logs -f

# サービスを停止
docker-compose down
```

### 個別のサービスのみ起動

```bash
# sample-apiのみ起動
docker-compose up -d sample-api
```

## 動作確認

サービスが正常に起動したら、以下のコマンドで動作確認できます：

```bash
# ヘルスチェック
curl http://localhost:8080/health

# 基本エンドポイント
curl http://localhost:8080/

# アイテム一覧取得
curl http://localhost:8080/api/v1/items
```

## 負荷試験での使用

これらのサービスが起動した状態で、`tools/vegeta` ディレクトリの負荷試験スクリプトを実行できます。

```bash
# 負荷試験の実行例
cd ../tools/vegeta/scripts
./load_test.sh
```

## トラブルシューティング

### ポートが既に使用されている場合

```bash
# 使用中のポートを確認
lsof -i :8080

# docker-compose.ymlでポート番号を変更
# ports:
#   - "8081:8080"  # ホスト側のポートを8081に変更
```

### コンテナのログを確認

```bash
docker-compose logs sample-api
```

### コンテナの状態を確認

```bash
docker-compose ps
```
