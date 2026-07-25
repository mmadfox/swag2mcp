# Bearer 認証

## 目的

Bearer トークン認証 — 最新の REST API で最も一般的な方法。トークンは `Authorization: Bearer &lt;token&gt;` ヘッダーで送信されます。

## 使用するタイミング

- 最新の REST API
- JWT（JSON Web Token）
- OAuth2 アクセストークン（トークンが既に取得されている場合）
- Bearer トークンを受け付ける任意の API

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
      type: bearer
      config:
        token: "eyJhbGciOiJIUzI1NiIs..."
```

## パラメーター

| パラメーター | 必須 | 説明 |
|-----------|------|------|
| `token` | はい | Bearer トークン（JWT、OAuth2 トークンなど） |

## 注意点

- トークンは静的です — 期限切れになった場合、設定で手動で更新する必要があります
- 自動トークン更新には `oauth2-cc` または `oauth2-pwd` を使用してください
- トークンは環境変数に保存：`token: "$(API_TOKEN)"`
