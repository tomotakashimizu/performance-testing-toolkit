# Vegeta Load Testing Tool

[Vegeta](https://github.com/tsenart/vegeta) を使用した API 負荷試験ツールです。

## 概要

このツールは、任意の API に対してリクエスト毎秒（RPS）を固定した負荷試験を実施し、パフォーマンス特性（レイテンシ、エラーレートなど）を計測するために作成されました。

## 目的

- 任意の API が指定された RPS で安定して稼働することを確認
- 特定条件下でのパフォーマンス特性（レイテンシ、エラーレートなど）を計測
- 負荷試験結果の可視化と分析

## ディレクトリ構造

```
tools/vegeta/
├── README.md                    # このファイル
├── config/                      # 設定ファイル
│   ├── README.md               # 設定ファイルの説明
│   ├── local_test_config.txt   # ローカルsample-api用設定
│   ├── sample_config.txt       # 外部API用設定例
│   ├── body1.json              # POSTリクエスト用JSONボディ
│   └── body2.json              # PUTリクエスト用JSONボディ
├── load_test_results/          # 負荷試験結果
│   ├── vegeta_TIMESTAMP_targets.json
│   ├── vegeta_TIMESTAMP_results.bin
│   ├── vegeta_TIMESTAMP_report.txt
│   ├── vegeta_TIMESTAMP_report.json
│   ├── vegeta_TIMESTAMP_histogram.txt
│   └── vegeta_TIMESTAMP_plot.html
└── scripts/
    └── load_test.sh            # メイン実行スクリプト
```

## 前提条件

### Vegeta のインストール

#### macOS (Homebrew)

```bash
brew install vegeta
```

#### Linux

```bash
# GitHub Releasesからダウンロード
wget https://github.com/tsenart/vegeta/releases/download/v12.12.0/vegeta_12.12.0_linux_amd64.tar.gz
tar -xzf vegeta_12.12.0_linux_amd64.tar.gz
sudo mv vegeta /usr/local/bin/
```

#### Go 環境がある場合

```bash
go install github.com/tsenart/vegeta@latest
```

### 動作確認

```bash
vegeta --version
```

## クイックスタート

### 1. ローカルサンプル API の起動

```bash
# プロジェクトルートから
cd test-targets
docker-compose up -d

# APIの動作確認
curl http://localhost:8080/health
```

### 2. 負荷試験の実行

```bash
# tools/vegeta/scriptsディレクトリに移動
cd tools/vegeta/scripts

# デフォルト設定で実行（50 RPS、1分間）
./load_test.sh

# カスタム設定で実行
./load_test.sh -r 100 -d 2m -n "my_test"
```

### 3. 結果の確認

```bash
# 最新の結果ディレクトリを確認
ls -la ../load_test_results/

# テキストレポートを表示
cat ../load_test_results/vegeta_*_report.txt

# HTMLプロットをブラウザで開く
open ../load_test_results/vegeta_*_plot.html
```

## 使用方法

### 基本的な使用方法

```bash
./load_test.sh [オプション]
```

### オプション

| オプション       | 説明               | デフォルト値                      | 例                               |
| ---------------- | ------------------ | --------------------------------- | -------------------------------- |
| `-c, --config`   | 設定ファイルのパス | `../config/local_test_config.txt` | `-c ../config/sample_config.txt` |
| `-r, --rate`     | リクエスト毎秒     | `50`                              | `-r 100`, `-r 500/s`             |
| `-d, --duration` | 試験時間           | `1m`                              | `-d 30s`, `-d 5m`, `-d 1h`       |
| `-t, --timeout`  | タイムアウト時間   | `30s`                             | `-t 10s`, `-t 1m`                |
| `-o, --output`   | 出力ディレクトリ   | `../load_test_results`            | `-o /tmp/results`                |
| `-n, --name`     | 試験名             | なし                              | `-n "high_load_test"`            |
| `-h, --help`     | ヘルプ表示         | -                                 | `-h`                             |

### 使用例

#### 基本的な負荷試験

```bash
# デフォルト設定（50 RPS、1分間）
./load_test.sh
```

#### 高負荷試験

```bash
# 500 RPS、5分間、タイムアウト10秒
./load_test.sh -r 500 -d 5m -t 10s -n "high_load_test"
```

#### 外部 API 試験

```bash
# カスタム設定ファイルを使用
./load_test.sh -c ../config/production_config.txt -r 10 -d 30s
```

#### 長時間試験

```bash
# 100 RPS、1時間
./load_test.sh -r 100 -d 1h -n "endurance_test"
```

## 設定ファイル

### 設定ファイルの作成

1. `config/sample_config.txt` をコピー
2. 対象 API に合わせてエンドポイントを編集
3. 必要に応じて認証ヘッダーを追加
4. リクエストボディが必要な場合は、JSON ファイルを作成して参照

詳細は [config/README.md](./config/README.md) を参照してください。

### 設定例

```
GET https://api.example.com/v1/users

POST https://api.example.com/v1/users
Content-Type: application/json
Authorization: Bearer YOUR_TOKEN
@user_data.json

GET https://api.example.com/v1/users/123

DELETE https://api.example.com/v1/users/123
Authorization: Bearer YOUR_TOKEN
```

## 出力ファイル

負荷試験実行後、以下のファイルが生成されます：

| ファイル                         | 説明                               |
| -------------------------------- | ---------------------------------- |
| `vegeta_TIMESTAMP_targets.json`  | 使用したターゲット設定             |
| `vegeta_TIMESTAMP_results.bin`   | Vegeta の生データ（バイナリ形式）  |
| `vegeta_TIMESTAMP_report.txt`    | テキスト形式のサマリーレポート     |
| `vegeta_TIMESTAMP_report.json`   | JSON 形式のサマリーレポート        |
| `vegeta_TIMESTAMP_histogram.txt` | レイテンシのヒストグラム           |
| `vegeta_TIMESTAMP_plot.html`     | レイテンシの時系列プロット（HTML） |

### レポートの見方

#### テキストレポート例

```
Requests      [total, rate, throughput]         3000, 50.02, 50.01
Duration      [total, attack, wait]             59.98s, 59.98s, 4.515ms
Latencies     [min, mean, 50, 90, 95, 99, max]  1.265ms, 5.015ms, 4.515ms, 8.125ms, 12.345ms, 25.678ms, 89.123ms
Bytes In      [total, mean]                     720000, 240.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:3000
Error Set:
```

#### 重要な指標

- **Throughput**: 実際の処理レート
- **Latencies**: レスポンス時間の分布
- **Success**: 成功率
- **Status Codes**: HTTP ステータスコードの分布

## トラブルシューティング

### よくある問題

#### 1. Vegeta がインストールされていない

```bash
# エラーメッセージ
エラー: Vegetaがインストールされていません

# 解決方法
brew install vegeta  # macOS
# または GitHub Releases からダウンロード
```

#### 2. 設定ファイルが見つからない

```bash
# エラーメッセージ
エラー: 設定ファイルが見つかりません

# 解決方法
ls ../config/  # 利用可能な設定ファイルを確認
./load_test.sh -c ../config/local_test_config.json
```

#### 3. 対象 API に接続できない

```bash
# 対象APIの動作確認
curl http://localhost:8080/health

# sample-apiが起動していない場合
cd ../../../test-targets
docker-compose up -d
```

#### 4. 権限エラー

```bash
# スクリプトに実行権限を付与
chmod +x load_test.sh
```

### デバッグ方法

#### 1. 設定ファイルの確認

```bash
# 設定ファイルの内容を確認
cat ../config/local_test_config.json

# JSON形式が正しいかチェック
python -m json.tool ../config/local_test_config.json
```

#### 2. 手動で Vegeta コマンドを実行

```bash
# 基本的なテスト
echo "GET http://localhost:8080/" | vegeta attack -rate=1 -duration=5s | vegeta report
```

#### 3. ログの確認

```bash
# sample-apiのログを確認
cd ../../../test-targets
docker-compose logs sample-api
```

## 高度な使用方法

### 分散負荷試験

複数のマシンから同時に負荷をかける場合：

```bash
# マシン1: 25 RPS
./load_test.sh -r 25 -d 5m -n "distributed_1"

# マシン2: 25 RPS
./load_test.sh -r 25 -d 5m -n "distributed_2"

# 合計: 50 RPS
```

### 段階的負荷増加

```bash
# 段階1: 50 RPS
./load_test.sh -r 50 -d 2m -n "ramp_up_50"

# 段階2: 100 RPS
./load_test.sh -r 100 -d 2m -n "ramp_up_100"

# 段階3: 200 RPS
./load_test.sh -r 200 -d 2m -n "ramp_up_200"
```

### 結果の比較分析

```bash
# 複数の結果を比較
vegeta plot \
  ../load_test_results/vegeta_ramp_up_50_*_results.bin \
  ../load_test_results/vegeta_ramp_up_100_*_results.bin \
  > comparison_plot.html
```

## 注意事項

1. **対象 API の利用規約を確認**してから負荷試験を実施してください
2. **本番環境**での負荷試験は事前に関係者と調整してください
3. **適切な RPS**を設定し、対象システムに過度な負荷をかけないよう注意してください
4. **認証情報**を設定ファイルに含める場合は、ファイルの管理に注意してください

## 参考資料

- [Vegeta 公式ドキュメント](https://github.com/tsenart/vegeta)
- [HTTP 負荷試験のベストプラクティス](https://github.com/tsenart/vegeta#usage)
- [設定ファイルの詳細](./config/README.md)
