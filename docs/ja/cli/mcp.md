# mcp

## 目的

**MCP（Model Context Protocol）サーバー** を起動します — LLM 統合の主要モードです。これを実行すると、AI エージェント（Claude、Cursor、OpenCode など）が 16 の MCP ツールを通じて API にアクセスできるようになります。

## 使用するタイミング

- LLM エージェントを API に接続したい場合
- IDE（VS Code、Cursor、JetBrains）またはデスクトップアプリ（Claude Desktop）を設定している場合
- MCP プロトコルを介して API を公開する必要がある場合
- 統合前に MCP サーバーをテストしている場合

## 構文

```bash
swag2mcp mcp [path] [flags]
```

## 引数

| 引数 | 位置 | 必須 | 説明 |
|------|------|------|------|
| `path` | 1 | いいえ | ワークスペースディレクトリ。省略時はパス解決ルールに従います。 |

## フラグ

| フラグ | 省略形 | 型 | デフォルト | 説明 |
|-------|--------|-----|-----------|------|
| `--transport` | | `string` | `"stdio"` | MCP トランスポート：`stdio`、`sse`、`streamable-http` |
| `--http-addr` | | `string` | `":8080"` | HTTP サーバーアドレス（`sse` および `streamable-http` 用） |
| `--http-path` | | `string` | `"/mcp"` | MCP ハンドラーの HTTP パス |
| `--auth-token` | | `string` | `""` | HTTP トランスポート認証用の Bearer トークン |
| `--logfile` | `-f` | `string` | `""` | ログファイルパス。未設定の場合は stderr に出力。 |
| `--disable-llm-auth` | | `bool` | `true` | MCP ツールリストから `auth` ツールを削除 |
| `--dump-dir` | | `string` | `""` | デバッグ用に HTTP リクエストをダンプするディレクトリ |
| `--tags` | `-t` | `string` | `""` | タグで spec をフィルタリング（カンマ区切り） |

## 仕組み

### stdio トランスポート（デフォルト）

MCP サーバーが LLM クライアント（IDE、Claude Desktop など）によってサブプロセスとして起動される場合に使用されます。サーバーは標準入出力を介して通信します。

```bash
swag2mcp mcp
```

### SSE トランスポート

HTTP ベースの通信のための Server-Sent Events トランスポート。MCP ハンドシェイクシーケンスが必要です。

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

### Streamable HTTP トランスポート

ストリーミングレスポンスをサポートする最新の HTTP トランスポート。

```bash
swag2mcp mcp --transport streamable-http --http-addr 0.0.0.0:8080
```

### 認証付き

HTTP エンドポイントを Bearer トークンで保護します：

```bash
swag2mcp mcp --transport sse --http-addr :8080 --auth-token "my-secret"
```

### タグフィルタリング付き

特定のタグを持つ spec のみを読み込みます：

```bash
swag2mcp mcp --tags=public
```

### auth ツール有効（デバッグモード）

LLM が `auth` ツールを介して新しいトークンを要求できるようにします：

```bash
swag2mcp mcp --disable-llm-auth=false
```

### リクエストダンプディレクトリ付き

デバッグ用にすべての HTTP リクエストを保存します：

```bash
swag2mcp mcp --dump-dir ./dumps
```

## MCP HTTP トランスポート — ハンドシェイクプロトコル

`sse` または `streamable-http` を使用する場合、MCP プロトコルは特定のハンドシェイクを必要とします。初期化前のツール呼び出しは失敗します：

```
Step 1: POST /mcp → {"method":"initialize", ...}
Step 2: POST /mcp → {"method":"notifications/initialized"}
Step 3: POST /mcp → {"method":"tools/list", ...}   ← これで動作
```

### ヘルスチェック

初期化なしで動作します：

```bash
curl http://localhost:8080/health
# → {"status":"ok","version":"v1.2.0"}
```

## IDE 設定例

### VS Code（`.vscode/settings.json` またはグローバル設定）

```json
{
  "mcp": {
    "servers": {
      "swag2mcp": {
        "command": "swag2mcp",
        "args": ["mcp", "/absolute/path/to/.swag2mcp"]
      }
    }
  }
}
```

### Cursor / Windsurf（`~/.cursor/mcp.json` またはプロジェクト `.cursor/mcp.json`）

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absolute/path/to/.swag2mcp"]
    }
  }
}
```

### Claude Desktop（macOS では `~/Library/Application Support/Claude/claude_desktop_config.json`）

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/absolute/path/to/.swag2mcp"]
    }
  }
}
```

### JetBrains IDE（設定 → ツール → MCP）

- 名前：`swag2mcp`
- コマンド：`swag2mcp`
- 引数：`mcp /absolute/path/to/.swag2mcp`

> **IDE 設定では常に絶対パス** をワークスペースディレクトリに使用してください。相対パスは IDE の作業ディレクトリによっては失敗する可能性があります。

## 出力

成功時に、サーバーは以下を出力します：

```
MCP server listening on http://127.0.0.1:8080/mcp
```

## ニュアンス

- **自動初期化なし:** 設定ファイルが存在しない場合、`mcp` はエラーを返します：`"configuration not found at &lt;path&gt;"`。最初に `init` を実行してください。
- **`--disable-llm-auth`（デフォルト：`true`）:** 有効な場合、`auth` ツールは MCP ツールリストから完全に削除されます。LLM はトークンを表示したり要求したりできません。認証は引き続き機能します — トークンは LLM 経由ではなく、標準の設定メカニズムを通じて取得されます。このモードは**本番環境**に推奨されます。**デバッグ**や短命トークンを使用する場合は、`--disable-llm-auth=false` を設定して LLM が `auth` ツールを介して新しいトークンを要求できるようにします。
- **YAML 設定のフォールバック:** CLI フラグが明示的に設定されていない場合、値は `swag2mcp.yaml` の `mcp` セクションから取得されます（存在する場合）。これにより、毎回フラグを渡す代わりに設定ファイルでサーバーを設定できます。
- **レスポンスクリーンアップ:** 起動時に、48 時間以上経過したレスポンスが `responses/` ディレクトリから自動的に削除されます。
- **パス解決の警告:** `[path]` が省略された場合、`mcp` は最初にカレントディレクトリで `swag2mcp.yaml` を検索し、次に `~/.swag2mcp/` にフォールバックします。間違ったディレクトリからコマンドを実行すると、意図したものとは異なるワークスペースが読み込まれる可能性があります。**サービスとして実行する場合や IDE 設定では、常に `[path]` を明示的に指定してください。**
