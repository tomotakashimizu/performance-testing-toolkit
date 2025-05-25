# Performance Testing Toolkit

API パフォーマンステストのためのツールキット

## 概要

このプロジェクトは、任意の API に対してリクエスト毎秒（RPS）を固定した負荷試験を実施し、パフォーマンス特性（レイテンシ、エラーレートなど）を計測するためのツールキットです。

## 目的

- 任意の API が指定された RPS で安定して稼働することを確認
- 特定条件下でのパフォーマンス特性（レイテンシ、エラーレートなど）を計測
- 負荷試験結果の可視化と分析
- 汎用的で再利用可能な負荷試験環境の提供

## ディレクトリ構造

```
.
├── README.md                    # このファイル
├── test-targets/               # 負荷試験対象のサンプルAPI
│   ├── README.md
│   ├── docker-compose.yml     # サンプルAPI起動用
│   └── sample-api/            # Go製サンプルAPI
│       ├── Dockerfile
│       ├── README.md
│       ├── go.mod
│       ├── go.sum
│       └── main.go
└── tools/                     # 負荷試験ツール
    └── vegeta/                # Vegetaベースの負荷試験ツール
        ├── README.md
        ├── config/            # 設定ファイル
        │   ├── README.md
        │   ├── local_test_config.txt
        │   ├── sample_config.txt
        │   ├── body1.json
        │   └── body2.json
        ├── load_test_results/ # 試験結果出力先
        └── scripts/
            └── load_test.sh   # メイン実行スクリプト
```

## サポートツール

- [Vegeta](https://github.com/tsenart/vegeta) - HTTP 負荷試験ツール

## クイックスタート

### 1. 前提条件

#### Vegeta のインストール

**macOS:**

```bash
brew install vegeta
```

**Linux:**

```bash
wget https://github.com/tsenart/vegeta/releases/download/v12.12.0/vegeta_12.12.0_linux_amd64.tar.gz
tar -xzf vegeta_12.12.0_linux_amd64.tar.gz
sudo mv vegeta /usr/local/bin/
```

#### Docker のインストール

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) をインストール

### 2. サンプル API の起動

```bash
# test-targetsディレクトリに移動
cd test-targets

# サンプルAPIを起動
docker-compose up -d

# 動作確認
curl http://localhost:8080/health
```

### 3. 負荷試験の実行

```bash
# 負荷試験ツールディレクトリに移動
cd tools/vegeta/scripts

# デフォルト設定で実行（50 RPS、1分間）
./load_test.sh

# カスタム設定で実行
./load_test.sh -r 100 -d 2m -n "my_test"
```

### 4. 結果の確認

```bash
# 結果ファイルを確認
ls -la ../load_test_results/

# テキストレポートを表示
cat ../load_test_results/vegeta_*_report.txt

# HTMLプロットをブラウザで開く
open ../load_test_results/vegeta_*_plot.html
```

## 使用例

### 基本的な負荷試験

```bash
# ローカルサンプルAPIに対する基本的な負荷試験
cd tools/vegeta/scripts
./load_test.sh
```

### 高負荷試験

```bash
# 500 RPS、5分間の高負荷試験
./load_test.sh -r 500 -d 5m -t 10s -n "high_load_test"
```

### 外部 API 試験

```bash
# 1. 設定ファイルを作成
cp ../config/sample_config.txt ../config/my_api_config.txt
# 2. 設定ファイルを編集（URLやヘッダーを変更）
# 3. 負荷試験を実行
./load_test.sh -c ../config/my_api_config.txt -r 10 -d 30s
```

### 段階的負荷増加試験

```bash
# 段階的に負荷を増加させる試験
./load_test.sh -r 50 -d 2m -n "ramp_up_50"
./load_test.sh -r 100 -d 2m -n "ramp_up_100"
./load_test.sh -r 200 -d 2m -n "ramp_up_200"
```

## 詳細ドキュメント

- [Vegeta 負荷試験ツール](./tools/vegeta/README.md) - 詳細な使用方法
- [設定ファイル](./tools/vegeta/config/README.md) - 設定ファイルの作成方法
- [サンプル API](./test-targets/README.md) - テスト対象 API の詳細
- [サンプル API 仕様](./test-targets/sample-api/README.md) - エンドポイント一覧

## 主な機能

### 負荷試験機能

- **固定 RPS 負荷試験**: 指定された RPS で一定時間負荷をかける
- **複数エンドポイント対応**: 複数の API エンドポイントを同時にテスト
- **認証対応**: Bearer Token、API Key 等の認証ヘッダーに対応
- **カスタムリクエスト**: POST/PUT 等のリクエストボディを含むテスト

### レポート機能

- **テキストレポート**: コンソールで確認しやすい形式
- **JSON レポート**: プログラムで処理しやすい形式
- **ヒストグラム**: レイテンシ分布の詳細分析
- **HTML プロット**: 時系列でのレイテンシ変化を可視化

### 設定機能

- **JSON 設定ファイル**: 柔軟なターゲット設定
- **コマンドライン引数**: 実行時パラメータの調整
- **結果ファイル管理**: タイムスタンプ付きファイル名で結果を整理

## トラブルシューティング

### よくある問題

1. **Vegeta がインストールされていない**

   ```bash
   brew install vegeta  # macOS
   ```

2. **サンプル API が起動していない**

   ```bash
   cd test-targets
   docker-compose up -d
   ```

3. **ポートが使用中**

   ```bash
   # docker-compose.ymlでポート番号を変更
   ports:
     - "8081:8080"  # ホスト側を8081に変更
   ```

4. **権限エラー**
   ```bash
   chmod +x tools/vegeta/scripts/load_test.sh
   ```

### デバッグ方法

```bash
# 設定ファイルの確認
python -m json.tool tools/vegeta/config/local_test_config.json

# 手動でのAPI確認
curl -v http://localhost:8080/health

# サンプルAPIのログ確認
cd test-targets
docker-compose logs sample-api
```

## 注意事項

1. **対象 API の利用規約を確認**してから負荷試験を実施してください
2. **本番環境**での負荷試験は事前に関係者と調整してください
3. **適切な RPS**を設定し、対象システムに過度な負荷をかけないよう注意してください
4. **認証情報**を設定ファイルに含める場合は、ファイルの管理に注意してください
