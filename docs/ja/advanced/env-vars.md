# 環境変数

## 概要

swag2mcp は設定ファイル内で `$(VAR_NAME)` 構文を使用した環境変数の置換をサポートしています。これにより、機密データ（トークン、パスワード、キー）を YAML ファイルから分離できます。

## 仕組み

swag2mcp が起動すると、設定内の `$(VAR_NAME)` パターンをスキャンし、対応する環境変数の値に置換します。

```yaml
specs:
  - domain: meteo
    auth:
      type: bearer
      config:
        token: "$(API_TOKEN)"
```

環境変数 `API_TOKEN` が設定されている場合、それが代入されます。設定されていない場合、値は空になります。

## `$(VAR)` が解決される場所

| フィールド | 例 |
|-----------|------|
| Auth `token` (bearer) | `token: "$(API_TOKEN)"` |
| Auth `username` / `password` (basic, digest) | `password: "$(API_PASSWORD)"` |
| Auth `client_id` / `client_secret` (oauth2-cc, oauth2-pwd) | `client_secret: "$(OAUTH_SECRET)"` |
| Auth `api_key` / `secret_key` (hmac) | `api_key: "$(BINANCE_API_KEY)"` |
| Auth `domain` (script) | `domain: "$(AUTH_DOMAIN)"` |
| MCP サーバートークン | `token: "$(MCP_TOKEN)"` |
| HTTP クライアントヘッダー | `"X-API-Key": "$(API_KEY)"` |
| HTTP クライアント Cookie 値 | `value: "$(SESSION_TOKEN)"` |

## `$(VAR)` が解決されない場所

- ベース URL（`base_url`）
- Collection の場所（`location`）
- Spec のドメイン名（`domain`）

## 例

```bash
export API_TOKEN="eyJhbGciOiJIUzI1NiIs..."
export MCP_TOKEN="my-secret-token"

swag2mcp mcp
```

## セキュリティのベストプラクティス

- **決して** YAML ファイルに直接シークレットを保存しないでください
- 環境変数または外部シークレットマネージャーを使用してください
- ハードコードされたシークレットが含まれている場合は YAML ファイルを `.gitignore` に追加してください
- シェルプロファイル、IDE 設定、またはデプロイメントパイプラインで環境変数を設定してください

## 構文の詳細

- `$(VAR_NAME)` — 標準構文
- `$( VAR_NAME )` — 括弧内の空白は許可され、トリミングされます
- `$()` — 空の変数名は元の文字列をそのまま返します
- ネストされた `$(...)` パターンは解決されません
