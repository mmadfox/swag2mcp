# なし

## 目的

認証不要。API はトークンやキーなしでアクセス可能。

## 使用するタイミング

- 公開 API（Open-Meteo、icanhazdadjoke、PokéAPI）
- テストおよびデモ環境
- API が認証を必要としない場合

## 設定

`type: none` を設定するか、単に `auth` セクションを省略します：

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: none
```

## パラメーター

なし。

## 注意点

- 設定から `auth` セクションが完全に存在しない場合、`type: none` と同等です
- リクエストに認証ヘッダーは追加されません
