# 設定ファイル

swag2mcp は YAML 設定ファイルを使用します。`swag2mcp init` で作成されます。

## 場所

- **Linux/macOS**: `~/.swag2mcp/swag2mcp.yaml`
- **Windows**: `%USERPROFILE%\.swag2mcp\swag2mcp.yaml`

## 基本構造

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

## 完全な例

```yaml
# ── グローバル HTTP クライアント ──────────────────────────────
http_client:
  timeout: 30s
  max_response_size: 1048576
  user_agent: "swag2mcp-global/1.0"
  follow_redirects: true
  max_redirects: 10
  random: false
  proxy:
    url: ""
    username: ""
    password: ""
    bypass: []
  headers:
    "Accept": "application/json"
  cookies:
    - name: "session"
      value: "abc123"

# ── MCP サーバー ──────────────────────────────────────────
mcp:
  transport: stdio
  addr: ":8080"
  path: "/mcp"
  auth:
    token: ""

# ── モックサーバー ─────────────────────────────────────────
mock_enabled: false
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

# ── レートリミッター ────────────────────────────────────────
disable_ratelimiter: false
rate_limit_interval: 10s

# ── Specs ───────────────────────────────────────────────
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    llm_instruction: "Use this API for weather forecasts and climate data"
    base_url: https://api.open-meteo.com
    disable: false
    tags: ["weather", "climate"]
    http_client:
      timeout: 10s
    auth:
      type: bearer
      config:
        token: "$(TOKEN)"
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: false
        http_client:
          timeout: 5s

  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## 環境変数

`$(VAR_NAME)` 構文を使用して環境変数を参照します。swag2mcp は起動時に解決します。

```yaml
specs:
  - domain: meteo
    auth:
      type: bearer
      config:
        token: "$(API_TOKEN)"

mcp:
  auth:
    token: "$(MCP_TOKEN)"
```

`$(VAR)` は以下で解決されます：
- Auth 設定フィールド：`token`、`username`、`password`、`client_id`、`client_secret`、`api_key`、`secret_key`、`domain`
- MCP サーバー認証トークン：`mcp.auth.token`
- HTTP クライアントヘッダーと Cookie 値

`$(VAR)` はベース URL や collection location では**解決されません**。

## 検証

```bash
# デフォルトワークスペースを検証（~/.swag2mcp）
swag2mcp validate

# カスタムプロジェクトワークスペースを検証
swag2mcp validate ./my-project
```

ワークスペースがホームディレクトリにない場合（例：プロジェクトリポジトリ内）、`validate`、`update`、`mcp`、またはその他のコマンドを実行するときは常にパスを指定してください。指定しない場合、swag2mcp はデフォルトの `~/.swag2mcp` ワークスペースを使用します。
