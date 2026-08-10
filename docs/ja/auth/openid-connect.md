# OpenID Connect (OIDC Discovery)

## 目的

OpenID Connect Discovery に基づく認証。swag2mcp は issuer から OIDC ディスカバリ文書を取得し、トークンエンドポイントを発見して、**client credentials** グラントで Bearer トークンを取得します。トークンは有効期限までキャッシュされます。

## 使用する場合

- OIDC ディスカバリ文書（`.well-known/openid-configuration`）を公開する API
- `client_id` + `client_secret` を持つサーバー間統合
- API が OIDC を使用し、トークンの自動取得が必要な場合

## 設定

```yaml
specs:
  - domain: my-api
    llm_title: My API
    base_url: https://api.example.com
    collections:
      - llm_title: Main
        location: https://example.com/spec.yaml
    auth:
      type: openid-connect
      config:
        issuer: "https://auth.example.com"
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        scopes:
          - openid
          - profile
```

## パラメータ

| パラメータ | 必須 | 説明 |
|-----------|----------|-------------|
| `issuer` | はい | OIDC issuer の URL。swag2mcp は `/.well-known/openid-configuration` を追加してトークンエンドポイントを発見します |
| `client_id` | はい | クライアント識別子 |
| `client_secret` | はい | クライアントシークレット |
| `scopes` | いいえ | 権限のリスト（任意） |

## 注意事項

- 現在のトークンが期限切れになると、swag2mcp は自動的に新しいトークンを要求します
- トークンは有効期限（`expires_in`）までキャッシュされます
- サーバーが `expires_in` を提供しない場合、トークンは 1 時間有効です
- `issuer`、`client_id`、`client_secret` は環境変数に `$(VAR)` 構文を使用できます
- トークンは **client credentials** グラントで取得されます（サーバー間、ユーザーなし）
- ディスカバリ文書には `token_endpoint` フィールドが必要です
