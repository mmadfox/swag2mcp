# OAuth2 Password Grant

## 目的

OAuth2 Resource Owner Password Grant — ユーザーのユーザー名とパスワードを使用した認証。ユーザーがアプリを信頼して認証情報を預けるファーストパーティアプリケーションに適しています。

## 使用するタイミング

- ファーストパーティアプリケーション（モバイル、Web）
- Keycloak および類似の Identity Provider との統合
- API が OAuth2 Password Grant をサポートしている場合

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
      type: oauth2-pwd
      config:
        client_id: "$(CLIENT_ID)"
        client_secret: "$(CLIENT_SECRET)"
        username: "$(USERNAME)"
        password: "$(PASSWORD)"
        token_url: "https://auth.example.com/oauth/token"
        scopes:
          - openid
          - profile
```

## パラメーター

| パラメーター | 必須 | 説明 |
|-----------|------|------|
| `client_id` | はい | クライアント識別子 |
| `username` | はい | ユーザー名 |
| `password` | はい | パスワード |
| `token_url` | はい | トークンエンドポイント URL |
| `client_secret` | いいえ | クライアントシークレット（オプション、パブリッククライアント用） |
| `scopes` | いいえ | 権限のリスト（オプション） |

## 注意点

- `client_secret` はオプション — **パブリッククライアント** がサポートされています（例：Keycloak）
- swag2mcp はトークンが期限切れになると自動的に更新します
- トークンは期限までキャッシュされます
- すべてのパラメーターは環境変数に保存できます
