# MCP サーバー

MCP サーバーは LLM エージェントの主要な対話ポイントです。設定されたすべての API を LLM が呼び出せる MCP ツールとして公開します。

## 設定

```yaml
mcp:
  transport: stdio
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

### auth.type

- **Type:** `string`
- **Default:** `""` (no JWT auth)
- **Options:** `jwks`, `oidc`, `introspection`
- **Description:** JWT authentication type for HTTP transport. When set, enables dynamic token verification using JWKS, OIDC Discovery, or token introspection.

### auth.jwks_url

- **Type:** `string`
- **Default:** `""`
- **Description:** URL of the JWKS (JSON Web Key Set) endpoint. Required when `auth.type` is `jwks` or resolved via OIDC discovery.

### auth.issuer

- **Type:** `string`
- **Default:** `""`
- **Description:** Expected JWT issuer (`iss` claim). If set, tokens with a different issuer are rejected.

### auth.audience

- **Type:** `string`
- **Default:** `""`
- **Description:** Expected JWT audience (`aud` claim). If set, tokens without this audience are rejected.

### auth.introspection_url

- **Type:** `string`
- **Default:** `""`
- **Description:** Token introspection endpoint URL. Required when `auth.type` is `introspection`.

### auth.client_id

- **Type:** `string`
- **Default:** `""`
- **Description:** Client ID for introspection auth. Required when `auth.type` is `introspection`.

### auth.client_secret

- **Type:** `string`
- **Default:** `""`
- **Description:** Client secret for introspection auth. Supports `$(ENV_VAR)` resolution.

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

### With JWT authentication (JWKS)

Protect the MCP HTTP endpoint with JWT verification via a JWKS endpoint:

```yaml
mcp:
  auth:
    type: jwks
    jwks_url: "https://auth.example.com/.well-known/jwks.json"
    issuer: "https://auth.example.com/"
    audience: "swag2mcp"
```

```bash
swag2mcp mcp --transport sse --http-addr 0.0.0.0:8080 \
  --auth-type jwks \
  --auth-jwks-url "https://auth.example.com/.well-known/jwks.json" \
  --auth-issuer "https://auth.example.com/" \
  --auth-audience "swag2mcp"
```

### With JWT authentication (OIDC Discovery)

```yaml
mcp:
  auth:
    type: oidc
    issuer: "https://auth.example.com/"
    audience: "swag2mcp"
```

### With JWT authentication (Token Introspection)

```yaml
mcp:
  auth:
    type: introspection
    introspection_url: "https://auth.example.com/introspect"
    client_id: "my-client"
    client_secret: "$(MCP_CLIENT_SECRET)"
```

```bash
swag2mcp mcp --transport sse --http-addr 0.0.0.0:8080 \
  --auth-type introspection \
  --auth-introspection-url "https://auth.example.com/introspect" \
  --auth-client-id "my-client" \
  --auth-client-secret "$(MCP_CLIENT_SECRET)"
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
| `--auth-type` | string | `""` | JWT auth type: `jwks`, `oidc`, `introspection` |
| `--auth-jwks-url` | string | `""` | JWKS URL for JWT auth |
| `--auth-issuer` | string | `""` | JWT issuer for token validation |
| `--auth-audience` | string | `""` | JWT audience for token validation |
| `--auth-introspection-url` | string | `""` | Token introspection URL |
| `--auth-client-id` | string | `""` | Client ID for introspection auth |
| `--auth-client-secret` | string | `""` | Client secret for introspection auth |
