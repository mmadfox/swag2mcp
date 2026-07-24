# Digest 認証

## 目的

HTTP Digest Access Authentication — Basic 認証よりも安全な代替手段。パスワードは平文で送信されず、代わりに MD5 ハッシュが使用されます。

## 使用するタイミング

- Digest のみをサポートするレガシー API
- パスワードを平文で送信せずに認証が必要な場合
- 内部エンタープライズシステム

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
      type: digest
      config:
        username: "admin"
        password: "$(PASSWORD)"
```

## パラメーター

| パラメーター | 必須 | 説明 |
|-----------|------|------|
| `username` | はい | ユーザー名 |
| `password` | はい | パスワード |

## 注意点

- swag2mcp は最初に認証なしでリクエストを送信し、サーバーからチャレンジ（HTTP 401）を受け取り、レスポンスを計算し、`Authorization: Digest ...` ヘッダーで再試行します
- チャレンジは 5 分間キャッシュされます — 後続のリクエストは追加のラウンドトリップを必要としません
- パスワードは環境変数に保存：`password: "$(API_PASSWORD)"`
