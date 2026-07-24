# 設定のカスケード

swag2mcp は 3 レベルの設定カスケードを使用します。各レベルは前のレベルを上書きします。これにより、グローバルに適切なデフォルトを設定し、特定の spec や collection の設定を微調整できます。

## レベル

```
Global (http_client, mcp, mock_enabled, disable_ratelimiter, rate_limit_interval)
    ↓ 上書き
Spec (specs[].http_client, specs[].auth, specs[].base_url, specs[].disable, specs[].tags)
    ↓ 上書き
Collection (specs[].collections[].http_client, specs[].collections[].base_url, specs[].collections[].disable)
```

## 何が何を上書きするか

| パラメーター | グローバル | Spec | Collection |
|-----------|----------|------|------------|
| `http_client.timeout` | ✅ | ✅ | ✅ |
| `http_client.max_response_size` | ✅ | ✅ | ✅ |
| `http_client.user_agent` | ✅ | ✅ | ✅ |
| `http_client.follow_redirects` | ✅ | ✅ | ✅ |
| `http_client.max_redirects` | ✅ | ✅ | ✅ |
| `http_client.proxy` | ✅ | ✅ | ✅ |
| `http_client.random` | ✅ | ✅ | ✅ |
| `http_client.headers` | ✅ | ✅ | ✅ |
| `http_client.cookies` | ✅ | ✅ | ✅ |
| `base_url` | ❌ | ✅ | ✅ |
| `auth` | ❌ | ✅ | ❌ |
| `disable` | ❌ | ✅ | ✅ |
| `tags` | ❌ | ✅ | ❌ |
| `mock_enabled` | ✅ | ❌ | ❌ |
| `disable_ratelimiter` | ✅ | ❌ | ❌ |
| `rate_limit_interval` | ✅ | ❌ | ❌ |

すべての `http_client` 設定はすべてのレベルで上書き可能です。Collection レベルの設定は spec およびグローバルよりも完全に優先されます。

## カスケードの例

```yaml
http_client:
  timeout: 30s
  max_response_size: 1048576
  headers:
    "User-Agent": "swag2mcp/1.0"

specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    http_client:
      timeout: 60s  # グローバルタイムアウトを上書き
      headers:
        "X-API-Version": "2"  # グローバルヘッダーに追加
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          timeout: 120s  # spec タイムアウトを上書き
          headers:
            "X-Custom": "value"  # spec + グローバルヘッダーに追加
```

## "Forecast" Collection の有効な設定

```
timeout: 120s（collection から、spec の 60s とグローバルの 30s を上書き）
max_response_size: 1048576（グローバルから）
headers:
  - User-Agent: swag2mcp/1.0（グローバルから）
  - X-API-Version: 2（spec から）
  - X-Custom: value（collection から）
```

## マージの仕組み

### HTTP クライアント設定

単純な値（`timeout`、`max_response_size`、`user_agent`、`follow_redirects`、`max_redirects`、`random`）は各レベルで**置き換え**られます。spec が `timeout: 60s` を設定した場合、グローバルの `30s` を完全に置き換えます。

### ヘッダー

ヘッダーはレベル間で**マージ**されます。3 つのレベルのすべてのヘッダーが結合されます。同じヘッダーキーが複数のレベルに現れる場合、最も低いレベルが優先されます。

### Cookie

Cookie はレベル間で**マージ**されます。同じ Cookie 名が複数のレベルに現れる場合、最も低いレベルが優先されます。

### プロキシ

プロキシは各レベルで**置き換え**られます。spec がプロキシを設定した場合、その spec のグローバルプロキシを完全に置き換えます。
