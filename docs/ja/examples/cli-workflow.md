# CLI ワークフロー

このページでは、初期化から日常的な操作まで、ターミナルから swag2mcp を使用する実際の例を示します。

## クイックスタート

```bash
# 1. ワークスペースを初期化
mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp

# 2. スペックを一覧表示
swag2mcp ls
```

## YAML でスペックを追加

### シンプルなスペック（公開 API）

```bash
swag2mcp add spec --yaml - <<EOF
domain: meteo
llm_title: Open-Meteo Weather API
base_url: https://api.open-meteo.com
collections:
  - llm_title: Weather Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
EOF
```

### 認証付きスペック（環境変数からのベアラートークン）

```bash
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My Protected API
base_url: https://api.example.com/v1
auth:
  type: bearer
  config:
    token: \$(MY_TOKEN)
collections:
  - llm_title: Users
    location: https://raw.githubusercontent.com/my-org/my-api/main/users.yaml
EOF
```

### 複数コレクションのスペック

```bash
swag2mcp add spec --yaml - <<EOF
domain: meteo
llm_title: Open-Meteo APIs
base_url: https://api.open-meteo.com
collections:
  - llm_title: Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
  - llm_title: Air Quality
    base_url: https://air-quality-api.open-meteo.com
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
EOF
```

## 既存のスペックにコレクションを追加

```bash
swag2mcp add collection --yaml - <<EOF
spec_domain: meteo
llm_title: Marine Weather
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
EOF
```

## スペックの一覧表示

```bash
$ swag2mcp ls
Specifications:
  dadjoke (https://icanhazdadjoke.com)
    jokes (3 endpoints)
  meteo (https://api.open-meteo.com)
    forecast (5 endpoints)
    air-quality (8 endpoints)
    marine (4 endpoints)
```

### タグでフィルタリング

```bash
swag2mcp ls --tags=public
```

## ランタイム情報の表示

```bash
$ swag2mcp info
{
  "version": "v1.2.0",
  "workspace": "/home/user/.swag2mcp",
  "specs": {
    "total": 2,
    "active": 2,
    "disabled": 0,
    "collections": 4,
    "endpoints": 20
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "follow_redirects": true
  },
  "mcp": {
    "transport": "stdio"
  },
  "auth": {
    "methods": ["bearer"]
  }
}
```

## 設定の検証

```bash
$ swag2mcp validate
✅ Configuration is valid.
✓ Spec dadjoke: OK
✓ Spec meteo: OK
```

## MCP サーバーの起動

### stdio（IDE 統合用）

```bash
swag2mcp mcp
```

### HTTP（リモートアクセス用）

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

### タグフィルター付き

```bash
swag2mcp mcp --tags=public
```

## スペックの更新

キャッシュされたすべてのスペックファイルを更新：

```bash
swag2mcp update
```

## キャッシュのクリーン

```bash
swag2mcp clean
```

## エクスポートとインポート

### ワークスペースのバックアップ

```bash
swag2mcp export --output ~/backups/swag2mcp-2026-07-24.zip
```

### 別のマシンで復元

```bash
# 新しいマシンで
swag2mcp import --from-zip swag2mcp-2026-07-24.zip
```

## 対話型 TUI エクスプローラー

```bash
swag2mcp run
```

API の検索、ブラウズ、呼び出しのための全画面ターミナル UI が開きます。

## モックサーバー

```bash
# モックバイナリをインストール
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest

# モックサーバーを起動
swag2mcp-mock mockserver
```
