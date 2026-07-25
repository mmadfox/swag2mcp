# MCP サーバー

MCP サーバーは LLM エージェントの主要な対話ポイントです。設定されたすべての API を LLM が呼び出せる MCP ツールとして公開します。

## 設定

```yaml
mcp:
  transport: stdio
  addr: ":8080"
  path: "/mcp"
  auth:
    token: ""
```

## トランスポート

3 つのトランスポートタイプが利用可能です：

| トランスポート | 説明 | 使用するタイミング |
|-----------|------|----------------|
| `stdio` | 標準入出力 | ローカル LLM クライアント（VS Code、Cursor、Claude Desktop） |
| `sse` | Server-Sent Events | リモートクライアント、HTTP ベースの通信 |
| `streamable-http` | HTTP ストリーミング | Web クライアント、最新の MCP クライアント |

### stdio（デフォルト）

LLM クライアントは swag2mcp を子プロセスとして実行します。通信は標準入出力を介して行われます。ネットワークポートは不要です。

```yaml
mcp:
  transport: stdio
```

```bash
swag2mcp mcp
```

### SSE

HTTP ベースの通信のための Server-Sent Events トランスポート。MCP サーバーは HTTP ポートでリッスンし、LLM クライアントはリモートから接続します。

```yaml
mcp:
  transport: sse
  addr: "127.0.0.1:8080"
  path: "/mcp"
```

```bash
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080
```

### Streamable HTTP

ストリーミングレスポンスをサポートする最新の HTTP トランスポート。SSE と似ていますが、異なるプロトコルを使用します。

```yaml
mcp:
  transport: streamable-http
  addr: "127.0.0.1:8080"
  path: "/mcp"
```

```bash
swag2mcp mcp --transport streamable-http --http-addr 0.0.0.0:8080
```

## パラメーター

### transport

- **型:** `string`
- **デフォルト:** `"stdio"`
- **オプション:** `stdio`、`sse`、`streamable-http`
- **効果:** MCP サーバーが LLM クライアントと通信する方法を決定します。

### addr

- **型:** `string`
- **デフォルト:** `":8080"`
- **説明:** SSE および Streamable HTTP トランスポートのリッスンアドレス。形式：`host:port`。
- **例:** `":8080"`、`"127.0.0.1:8080"`、`"0.0.0.0:9000"`

### path

- **型:** `string`
- **デフォルト:** `"/mcp"`
- **説明:** MCP エンドポイントの URL パス。LLM クライアントは `http://&lt;addr&gt;&lt;path&gt;` にリクエストを送信します。
- **例:** `"/mcp"`、`"/api/mcp"`、`"/v1/mcp"`

### auth.token

- **型:** `string`
- **デフォルト:** `""`（認証なし）
- **説明:** HTTP トランスポート認証用の Bearer トークン。設定すると、LLM クライアントはすべてのリクエストに `Authorization: Bearer &lt;token&gt;` を含める必要があります。
- **注:** `$(ENV_VAR)` 解決をサポートします。

## HTTP 認証

MCP HTTP エンドポイントを Bearer トークンで保護します：

```yaml
mcp:
  auth:
    token: "my-secret-token"
```

または CLI フラグ経由：

```bash
swag2mcp mcp --auth-token "my-secret-token"
```

## ヘルスチェック

MCP サーバーは MCP 初期化なしで動作するヘルスチェックエンドポイントを提供します：

```bash
curl http://127.0.0.1:8080/health
# {"status":"ok","version":"v1.2.0"}
```

## 起動フラグ

CLI フラグは YAML 設定を上書きします。フラグが設定されていない場合、YAML の `mcp` セクションの値がフォールバックとして使用されます。

| フラグ | 型 | デフォルト | 説明 |
|-------|------|---------|------|
| `--transport` | string | `"stdio"` | トランスポートタイプ：`stdio`、`sse`、`streamable-http` |
| `--http-addr` | string | `":8080"` | HTTP サーバーアドレス（SSE および Streamable HTTP 用） |
| `--http-path` | string | `"/mcp"` | MCP ハンドラーの URL パス |
| `--auth-token` | string | `""` | HTTP トランスポート認証用の Bearer トークン |
| `--logfile` | string | `""` | ログファイルパス（未設定の場合は stderr に出力） |
| `--disable-llm-auth` | bool | `true` | MCP ツールリストから `auth` ツールを削除 |
| `--dump-dir` | string | `""` | デバッグ用に HTTP リクエストをダンプするディレクトリ |
| `--tags` | string | `""` | タグで spec をフィルタリング（カンマ区切り） |
