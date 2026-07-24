# add

## 目的

新しい **spec**（API サービス）または **collection**（OpenAPI/Swagger/Postman ファイル）を既存の設定に追加します。新しい API でワークスペースを拡張するための主要な方法です。

## 使用するタイミング

- LLM エージェントに接続する新しい API がある場合
- OpenAPI spec の URL を見つけて追加したい場合
- 既存の spec に追加の spec ファイル（collection）を追加したい場合
- 対話型ウィザードの代わりに YAML を直接記述したい場合

## 構文

```bash
swag2mcp add spec [path] [flags]
swag2mcp add collection [path] [flags]
```

## 引数

| 引数 | 位置 | 必須 | 説明 |
|------|------|------|------|
| `path` | 1 | いいえ | ワークスペースディレクトリ。省略時はパス解決ルールに従います。 |

## フラグ

### `add spec`

| フラグ | 省略形 | 型 | デフォルト | 説明 |
|-------|--------|-----|-----------|------|
| `--yaml` | `-y` | `string` | `""` | インライン YAML 入力、または `-` で標準入力 |
| `--example` | `-e` | `bool` | `false` | YAML テンプレートを表示して終了 |

### `add collection`

| フラグ | 省略形 | 型 | デフォルト | 説明 |
|-------|--------|-----|-----------|------|
| `--yaml` | `-y` | `string` | `""` | インライン YAML 入力、または `-` で標準入力 |
| `--example` | `-e` | `bool` | `false` | YAML テンプレートを表示して終了 |

## 仕組み

### 対話モード（デフォルト）

spec または collection のフィールドをステップごとに入力できる TUI ウィザードを起動します。

```bash
swag2mcp add spec
swag2mcp add collection
```

### YAML インラインモード

YAML を直接文字列として渡します。**シェルの引用符に注意** — `:`、`#`、`&`、`{` などの特殊文字はコマンドを壊す可能性があります。

```bash
swag2mcp add spec --yaml 'domain: meteo
llm_title: Open-Meteo API
base_url: https://meteo.swagger.io/v2
collections:
  - llm_title: Main
    location: https://example.com/spec.json'
```

### 標準入力からの YAML（複雑な YAML に推奨）

ファイルからパイプするか、ヒアドキュメントを使用してシェルの引用符問題を完全に回避します：

```bash
# ファイルからパイプ
cat spec.yaml | swag2mcp add spec --yaml -

# ヒアドキュメント
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My API
llm_instruction: "Use this API for X & Y # important"
base_url: https://api.example.com/v1
collections:
  - llm_title: Main
    location: https://raw.githubusercontent.com/org/repo/main/spec.yaml
EOF
```

### YAML テンプレート

期待される YAML 構造を表示して終了します：

```bash
swag2mcp add spec --example
swag2mcp add collection --example
```

## YAML 形式

### Spec

```yaml
domain: meteo
llm_title: Open-Meteo API
llm_instruction: Use this API to manage pets.
base_url: https://meteo.swagger.io/v2
tags: [public, demo]
auth:
  type: bearer
  config:
    token: $(TOKEN)
collections:
  - llm_title: Open-Meteo Swagger
    location: https://example.com/spec.json
```

### Collection

```yaml
spec_domain: meteo
llm_title: Orders Collection
location: https://example.com/orders.json
```

## コマンド実行後の確認

```bash
swag2mcp ls [path]
# 新しい spec または collection がリストに表示されるはずです
```

## ニュアンス

- **自動初期化:** 設定ファイルが存在しない場合、`add` は自動的に init ウィザードを実行します。`init` を別途実行する必要はありません。
- **シェルの引用符:** インライン YAML（`--yaml '...'`）は特殊文字で問題が発生しやすいです。単純な値を超える場合は、ヒアドキュメントまたはパイプで `--yaml -` を使用することを推奨します。
- **`--example` は即座に終了** し、既存の設定を確認したり何も変更したりしません。
- **`add spec` と `add collection` の違い:** 新しい API サービス（新しいドメイン）には `add spec` を使用します。既存の spec に別の spec ファイルを追加するには `add collection` を使用します。
