# VimMaster 🎮

Vim学習CUIアプリ — vimtutorよりも楽しく、ゲーム感覚でVimの操作を指に覚えさせるターミナル専用学習プラットフォーム。

## 特徴

- **ステージベースの学習**: 基礎移動からジャンプコマンドまで、段階的にVimを習得
- **リアルなVimエミュレーション**: Normal/Insertモードの切り替え、hjkl移動、dw/dd/x等の編集コマンド
- **評価システム**: 最小ストローク数でクリアするとS評価！
- **美しいTUI**: Bubble Tea + Lip GlossによるモダンなターミナルUI

## 技術スタック

- **言語**: Go
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss)

## 必要要件

- Go 1.24+ または Docker

## セットアップ

### ローカル実行

```bash
go mod download
go run main.go
```

### Docker で実行

```bash
# 開発環境（ホットリロード対応）
docker compose run --rm dev

# プロダクションビルド
docker compose run --rm app
```

### ビルド

```bash
go build -o vimmaster ./main.go
./vimmaster
```

## 操作方法

### タイトル画面
- `Enter` - スタート
- `q` - 終了

### ステージ選択
- `j/k` - 上下選択
- `Enter` - ステージ開始
- `q` - 戻る

### ゲーム画面
- Vimキーバインドでテキストを編集
- `?` - ヒント表示/非表示
- `Ctrl+Q` - ステージ選択に戻る

### サポートされるVimコマンド

| カテゴリ | コマンド | 説明 |
|---------|---------|------|
| 移動 | `h`, `j`, `k`, `l` | 左、下、上、右 |
| 移動 | `w`, `e`, `b` | 単語先頭、単語末尾、前の単語 |
| 移動 | `0`, `$` | 行頭、行末 |
| 移動 | `gg`, `G` | ファイル先頭、ファイル末尾 |
| 移動 | `f{char}`, `t{char}` | 文字検索 |
| 編集 | `x` | 文字削除 |
| 編集 | `dw`, `dd`, `d$` | 単語/行/行末まで削除 |
| 編集 | `cw` | 単語変更 |
| モード | `i`, `a`, `A`, `I`, `o` | 挿入モードに入る |
| モード | `Esc` | ノーマルモードに戻る |

## ステージ一覧

1. 基本移動 (hjkl)
2. 行移動練習
3. 行頭と行末 (0, $)
4. テキスト挿入 (i)
5. 文字削除と追加 (x, a)
6. 新しい行を開く (o)
7. 単語の削除 (dw)
8. 単語移動 (w, b, e)
9. 行末まで削除 (d$)
10. 行ジャンプ (gg, G)
11. 文字検索 (f, t)
12. 行の削除 (dd)

## プロジェクト構成

```
.
├── main.go                    # エントリーポイント
├── internal/
│   ├── app/app.go            # メインのBubble Teaモデル
│   ├── buffer/buffer.go      # 仮想テキストバッファ
│   ├── vim/
│   │   ├── mode.go           # Vimモード定義
│   │   └── command.go        # Vimコマンド処理
│   ├── stage/
│   │   ├── stage.go          # ステージ定義
│   │   ├── loader.go         # ステージ読み込み
│   │   ├── checker.go        # クリア判定
│   │   └── data/*.json       # ステージデータ
│   └── ui/
│       ├── header.go         # ヘッダーUI
│       ├── editor.go         # エディターエリア
│       ├── statusline.go     # ステータスライン
│       └── hint.go           # ヒントエリア
├── Dockerfile                 # プロダクション用
├── Dockerfile.dev             # 開発環境用
├── docker-compose.yml         # Docker Compose設定
└── doc.md                     # 要件定義書
```