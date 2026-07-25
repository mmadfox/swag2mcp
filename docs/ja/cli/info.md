# info

## 目的

swag2mcp ランタイムの包括的なサマリーを **JSON** で表示します。バージョン、ワークスペースパス、spec サマリー、HTTP クライアント設定、MCP トランスポート設定、認証方式、モックモードステータスが含まれます。

## 使用するタイミング

- ワークスペースの機械可読な概要が必要な場合
- デバッグのためにランタイム設定を確認する必要がある場合
- アクティブな spec とエンドポイントの数を確認したい場合
- HTTP クライアントまたは MCP トランスポート設定を確認する必要がある場合

## 構文

```bash
swag2mcp info [path]
```

## 引数

| 引数 | 位置 | 必須 | 説明 |
|------|------|------|------|
| `path` | 1 | いいえ | ワークスペースディレクトリ。省略時はパス解決ルールに従います。 |

## フラグ

なし。

## 仕組み

```bash
swag2mcp info
swag2mcp info ./my-workspace
```

## 出力

出力は以下の構造を持つ JSON オブジェクトです：

```json
{
  "version": "v1.2.0",
  "workspace": "~/.swag2mcp",
  "uptime": "2h 15m",
  "specs": {
    "total": 4,
    "active": 3,
    "disabled": 1,
    "collections": 6,
    "endpoints": 42
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "proxy": "none",
    "follow_redirects": true,
    "max_redirects": 10,
    "randomize": false
  },
  "mcp": {
    "transport": "stdio",
    "addr": ":8080",
    "path": "/mcp"
  },
  "auth_methods": ["bearer", "api-key"],
  "mock_enabled": false
}
```

## コマンド実行後の確認

MCP サーバーを起動する前に、`info` を使用してワークスペースが正しく読み込まれ、すべての spec がアクティブであることを確認します。

## ニュアンス

- **自動初期化:** 設定ファイルが存在しない場合、`info` は自動的に init ウィザードを実行します。
- **JSON のみ:** 出力は常に JSON です。人間が読める出力には `ls` を使用してください。
- **`max_response_size`:** 人間が読める形式で表示されます（例：`"1 KB"`、`"2 MB"`）。
- **全文インデックスなし:** `info` は設定と spec のメタデータのみが必要なため、全文インデックス作成を無効にします。
