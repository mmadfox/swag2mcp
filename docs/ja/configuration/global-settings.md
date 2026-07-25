# グローバル設定

グローバル設定は `swag2mcp.yaml` のトップレベルの設定ブロックです。spec または collection レベルで上書きされない限り、すべての spec に適用されます。

## 構造

```yaml
http_client:
  # すべての API 呼び出しの HTTP クライアント設定

mcp:
  # MCP サーバー設定

mock_enabled: false
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

disable_ratelimiter: false
rate_limit_interval: 10s
```

## HTTP クライアント

swag2mcp が API に HTTP リクエストを行う方法を制御します：タイムアウト、レスポンスサイズ制限、プロキシ、ヘッダー、Cookie、リダイレクト、ユーザーエージェント。これらの設定は spec と collection にカスケードされます。

すべてのパラメーターと例については [HTTP Client](./http-client) を参照してください。

## MCP サーバー

MCP サーバーが LLM エージェントと通信する方法を制御します：トランスポートタイプ（stdio、SSE、Streamable HTTP）、アドレス、パス、オプションの Bearer トークン認証。

すべてのパラメーター、トランスポート、起動フラグについては [MCP Server](./mcp-server) を参照してください。

## モックサーバー

モックサーバーは OpenAPI スキーマに基づいて偽の API レスポンスを生成します。実際の API にアクセスせずにテストするのに便利です。

```yaml
mock_enabled: true
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092
```

### mock_enabled

- **型:** `bool`
- **デフォルト:** `false`
- **効果:** `true` の場合、swag2mcp は `base_mock_url` が設定されているすべての spec のモックサーバーを起動します。各 collection に `base_mock_url` が設定されている必要があります。
- **有効にするタイミング:** 実際の HTTP 呼び出しを行わずに API 統合をテストしたい場合。モックサーバーは OpenAPI スキーマに基づいて偽のデータを返します。

### mock_auth

モック認証サーバーのポート設定。これらはモックサーバーで認証方式（OAuth2、Digest、HMAC）をテストする際に使用されます。

| フィールド | 型 | デフォルト | 説明 |
|-----------|------|---------|------|
| `oauth2_port` | int | `9090` | モック OAuth2 トークンサーバーのポート（1024〜65535） |
| `digest_port` | int | `9091` | モック Digest 認証サーバーのポート（1024〜65535） |
| `hmac_port` | int | `9092` | モック HMAC 認証サーバーのポート（1024〜65535） |

## レートリミッター

レートリミッターは LLM が同じ API エンドポイントを頻繁に呼び出すのを防ぎます。デフォルトでは、各エンドポイントは 10 秒に 1 回呼び出せます。

```yaml
disable_ratelimiter: false
rate_limit_interval: 10s
```

### disable_ratelimiter

- **型:** `bool`
- **デフォルト:** `false`
- **効果:** `true` の場合、エンドポイントごとのレートリミッターが完全に無効になります。LLM は待機なしで同じエンドポイントを繰り返し呼び出せます。
- **有効にするタイミング:** テスト、デバッグ、または同じエンドポイントを短時間に複数回呼び出す必要がある場合。
- **無効のままにするタイミング（推奨）:** 本番環境。レートリミッターは偶発的な悪用を防ぎ、API のレート制限を尊重します。

### rate_limit_interval

- **型:** 期間（Go 形式：`10s`、`30s`、`1m`）
- **デフォルト:** `10s`
- **効果:** LLM が同じエンドポイントへの呼び出しの間に待機する必要がある時間を設定します。
- **変更するタイミング:** 厳格なレート制限がある API では増やします。負荷を制御できる内部 API では減らします。
- **範囲:** 任意の有効な期間（例：`5s`、`30s`、`1m`、`2m`）。

## カスケード

グローバル設定は spec および collection レベルで上書きできます。すべての `http_client` 設定（タイムアウト、プロキシ、ユーザーエージェント、リダイレクト、レスポンスサイズ、ランダマイザー、ヘッダー、Cookie）は spec と collection の両方のレベルで上書きできます。

```
Global (http_client, mock_enabled, disable_ratelimiter, rate_limit_interval)
    ↓ 上書き（http_client のみ）
Spec (specs[].http_client)
    ↓ 上書き（http_client のみ）
Collection (specs[].collections[].http_client)
```

詳細は [Configuration Cascade](./cascade) を参照してください。
