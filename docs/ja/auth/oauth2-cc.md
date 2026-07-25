# OAuth2 Client Credentials

## 目的

OAuth2 Client Credentials Grant — サーバー間通信のための認証。アプリケーションは client_id と client_secret を使用してトークンを取得します。ユーザーの関与は不要です。

## 使用するタイミング

- マイクロサービスおよびサーバー間統合
- マシン間通信
- API が OAuth2 を使用し、client_id + client_secret がある場合

## 設定

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: oauth2-cc
      config:
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        token_url: "https://auth.example.com/oauth/token"
        scopes:
          - read
          - write
```

## パラメーター

| パラメーター | 必須 | 説明 |
|-----------|------|------|
| `client_id` | はい | クライアント識別子 |
| `client_secret` | はい | クライアントシークレット |
| `token_url` | はい | トークンエンドポイント URL |
| `scopes` | いいえ | 権限のリスト（オプション） |

## 注意点

- swag2mcp は現在のトークンが期限切れになると自動的に新しいトークンを要求します
- トークンは有効期限（`expires_in`）までキャッシュされます
- サーバーが `expires_in` を提供しない場合、トークンは 1 時間有効と見なされます
- すべてのパラメーターは環境変数に保存できます
