# API キー

## 目的

API キーによる認証。キーは HTTP ヘッダーまたは URL クエリパラメーターとして送信できます。

## 使用するタイミング

- API キーを使用するサービス
- 気象サービス、地理データ、翻訳 API
- API がヘッダー（`X-API-Key`）またはクエリパラメーター（`?api_key=...`）でキーを期待する場合

## 設定

### ヘッダー内のキー

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: api-key
      config:
        key: "X-API-Key"
        in: header
        value: "$(API_KEY)"
```

### クエリパラメーター内のキー

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: api-key
      config:
        key: "api_key"
        in: query
        value: "$(API_KEY)"
```

## パラメーター

| パラメーター | 必須 | 説明 |
|-----------|------|------|
| `key` | はい | ヘッダーまたはクエリパラメーターの名前 |
| `in` | はい | キーの配置場所：`header` または `query` |
| `value` | はい | キーの値 |

## 注意点

- `header` モードでは、キーは HTTP ヘッダーとして追加されます
- `query` モードでは、キーは URL パラメーターとして追加されます
- 値は環境変数に保存：`value: "$(MY_API_KEY)"`
