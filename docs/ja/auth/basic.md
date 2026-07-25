# Basic 認証

## 目的

HTTP Basic 認証 — ユーザー名とパスワードで認証する最もシンプルな方法。

## 使用するタイミング

- Basic 認証のみをサポートするレガシー API
- 複雑なトークン不要のシンプルな認証
- 内部サービス

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
      type: basic
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

- パスワードは `Authorization: Basic ...` ヘッダーで Base64 エンコードされて送信されます — これは**暗号化ではありません**。常に HTTPS を使用してください。
- パスワードは環境変数に保存：`password: "$(MY_PASSWORD)"`
