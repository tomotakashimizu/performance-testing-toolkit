#!/bin/bash

# Vegeta Load Testing Script
# 使用方法: ./load_test.sh [オプション]

set -e

# デフォルト設定
DEFAULT_CONFIG="../config/local_test_config.txt"
DEFAULT_RATE="50"
DEFAULT_DURATION="1m"
DEFAULT_TIMEOUT="30s"
DEFAULT_OUTPUT_DIR="../load_test_results"

# 変数の初期化
CONFIG_FILE=""
RATE=""
DURATION=""
TIMEOUT=""
OUTPUT_DIR=""
CUSTOM_NAME=""
HELP=false

# ヘルプ表示関数
show_help() {
    cat <<EOF
Vegeta Load Testing Script

使用方法:
    $0 [オプション]

オプション:
    -c, --config FILE       設定ファイルのパス (デフォルト: $DEFAULT_CONFIG)
    -r, --rate RATE         リクエスト毎秒 (例: 50, 100/s) (デフォルト: $DEFAULT_RATE)
    -d, --duration TIME     試験時間 (例: 1m, 30s, 2h) (デフォルト: $DEFAULT_DURATION)
    -t, --timeout TIME      タイムアウト時間 (例: 30s, 1m) (デフォルト: $DEFAULT_TIMEOUT)
    -o, --output DIR        出力ディレクトリ (デフォルト: $DEFAULT_OUTPUT_DIR)
    -n, --name NAME         試験名（ファイル名に使用）
    -h, --help              このヘルプを表示

使用例:
    # デフォルト設定で実行
    $0

    # カスタム設定で実行
    $0 -c ../config/sample_config.txt -r 100 -d 2m

    # 高負荷試験
    $0 -r 500 -d 5m -t 10s -n "high_load_test"

    # 外部API試験
    $0 -c ../config/production_config.txt -r 10 -d 30s

設定ファイル形式:
    HTTP形式（txt形式）でVegetaのターゲットを定義
    詳細は ../config/README.md を参照

出力ファイル:
    - vegeta_TIMESTAMP_targets.json    使用したターゲット設定
    - vegeta_TIMESTAMP_results.bin     Vegetaの生データ
    - vegeta_TIMESTAMP_report.txt      テキスト形式レポート
    - vegeta_TIMESTAMP_report.json     JSON形式レポート
    - vegeta_TIMESTAMP_histogram.txt   ヒストグラムレポート
    - vegeta_TIMESTAMP_plot.html       時系列プロットHTML

EOF
}

# コマンドライン引数の解析
while [[ $# -gt 0 ]]; do
    case $1 in
    -c | --config)
        CONFIG_FILE="$2"
        shift 2
        ;;
    -r | --rate)
        RATE="$2"
        shift 2
        ;;
    -d | --duration)
        DURATION="$2"
        shift 2
        ;;
    -t | --timeout)
        TIMEOUT="$2"
        shift 2
        ;;
    -o | --output)
        OUTPUT_DIR="$2"
        shift 2
        ;;
    -n | --name)
        CUSTOM_NAME="$2"
        shift 2
        ;;
    -h | --help)
        HELP=true
        shift
        ;;
    *)
        echo "不明なオプション: $1"
        echo "ヘルプを表示するには -h または --help を使用してください"
        exit 1
        ;;
    esac
done

# ヘルプ表示
if [ "$HELP" = true ]; then
    show_help
    exit 0
fi

# デフォルト値の設定
CONFIG_FILE=${CONFIG_FILE:-$DEFAULT_CONFIG}
RATE=${RATE:-$DEFAULT_RATE}
DURATION=${DURATION:-$DEFAULT_DURATION}
TIMEOUT=${TIMEOUT:-$DEFAULT_TIMEOUT}
OUTPUT_DIR=${OUTPUT_DIR:-$DEFAULT_OUTPUT_DIR}

# Vegetaがインストールされているかチェック
if ! command -v vegeta &>/dev/null; then
    echo "エラー: Vegetaがインストールされていません"
    echo "インストール方法:"
    echo "  macOS: brew install vegeta"
    echo "  Linux: https://github.com/tsenart/vegeta/releases からダウンロード"
    exit 1
fi

# 設定ファイルの存在チェック
if [ ! -f "$CONFIG_FILE" ]; then
    echo "エラー: 設定ファイルが見つかりません: $CONFIG_FILE"
    echo "利用可能な設定ファイル:"
    find ../config -name "*.txt" -type f 2>/dev/null || echo "  設定ファイルが見つかりません"
    exit 1
fi

# 出力ディレクトリの作成
mkdir -p "$OUTPUT_DIR"

# タイムスタンプの生成
if [ -n "$CUSTOM_NAME" ]; then
    TIMESTAMP="${CUSTOM_NAME}_$(date +%Y%m%d_%H%M%S)"
else
    TIMESTAMP="$(date +%Y%m%d_%H%M%S)"
fi

# ファイル名の定義
TARGETS_FILE="$OUTPUT_DIR/vegeta_${TIMESTAMP}_targets.json"
RESULTS_FILE="$OUTPUT_DIR/vegeta_${TIMESTAMP}_results.bin"
REPORT_TXT_FILE="$OUTPUT_DIR/vegeta_${TIMESTAMP}_report.txt"
REPORT_JSON_FILE="$OUTPUT_DIR/vegeta_${TIMESTAMP}_report.json"
HISTOGRAM_FILE="$OUTPUT_DIR/vegeta_${TIMESTAMP}_histogram.txt"
PLOT_FILE="$OUTPUT_DIR/vegeta_${TIMESTAMP}_plot.html"

# 設定情報の表示
echo "=========================================="
echo "Vegeta Load Testing"
echo "=========================================="
echo "設定ファイル: $CONFIG_FILE"
echo "レート: $RATE/s"
echo "期間: $DURATION"
echo "タイムアウト: $TIMEOUT"
echo "出力ディレクトリ: $OUTPUT_DIR"
echo "タイムスタンプ: $TIMESTAMP"
echo "=========================================="

# 設定ファイルをコピー（記録用）
cp "$CONFIG_FILE" "$TARGETS_FILE"
echo "ターゲット設定を保存: $TARGETS_FILE"

# 負荷試験の実行
echo "負荷試験を開始します..."
echo "実行コマンド: vegeta attack -targets=$CONFIG_FILE -rate=$RATE -duration=$DURATION -timeout=$TIMEOUT"

if vegeta attack \
    -targets="$CONFIG_FILE" \
    -rate="$RATE" \
    -duration="$DURATION" \
    -timeout="$TIMEOUT" \
    -format=http \
    >"$RESULTS_FILE"; then

    echo "負荷試験が完了しました: $RESULTS_FILE"
else
    echo "エラー: 負荷試験が失敗しました"
    exit 1
fi

# レポートの生成
echo "レポートを生成しています..."

# テキスト形式のレポート
if vegeta report <"$RESULTS_FILE" >"$REPORT_TXT_FILE"; then
    echo "テキストレポートを生成: $REPORT_TXT_FILE"
else
    echo "警告: テキストレポートの生成に失敗しました"
fi

# JSON形式のレポート
if vegeta report -type=json <"$RESULTS_FILE" >"$REPORT_JSON_FILE"; then
    echo "JSONレポートを生成: $REPORT_JSON_FILE"
else
    echo "警告: JSONレポートの生成に失敗しました"
fi

# ヒストグラムレポート
if vegeta report -type=hist <"$RESULTS_FILE" >"$HISTOGRAM_FILE"; then
    echo "ヒストグラムレポートを生成: $HISTOGRAM_FILE"
else
    echo "警告: ヒストグラムレポートの生成に失敗しました"
fi

# HTMLプロット
if vegeta plot <"$RESULTS_FILE" >"$PLOT_FILE"; then
    echo "HTMLプロットを生成: $PLOT_FILE"
else
    echo "警告: HTMLプロットの生成に失敗しました"
fi

# 結果の概要表示
echo "=========================================="
echo "負荷試験結果の概要"
echo "=========================================="

if [ -f "$REPORT_TXT_FILE" ]; then
    cat "$REPORT_TXT_FILE"
else
    echo "レポートファイルが見つかりません"
fi

echo "=========================================="
echo "生成されたファイル:"
echo "  ターゲット設定: $TARGETS_FILE"
echo "  生データ: $RESULTS_FILE"
echo "  テキストレポート: $REPORT_TXT_FILE"
echo "  JSONレポート: $REPORT_JSON_FILE"
echo "  ヒストグラム: $HISTOGRAM_FILE"
echo "  HTMLプロット: $PLOT_FILE"
echo "=========================================="

# HTMLプロットをブラウザで開くかの確認
if command -v open &>/dev/null && [ -f "$PLOT_FILE" ]; then
    read -p "HTMLプロットをブラウザで開きますか？ (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        open "$PLOT_FILE"
    fi
fi

echo "負荷試験が完了しました！"
