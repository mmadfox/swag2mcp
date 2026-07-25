# HMAC 認証

## 目的

HMAC-SHA256 リクエスト署名 — 暗号通貨取引所（Binance、Bybit など）で使用される認証方法。各リクエストは秘密鍵で署名されます。

## 使用するタイミング

- Binance API および Binance 互換の取引所
- 暗号通貨取引プラットフォーム
- リクエスト署名が必要な API

## 設定

```yaml
specs:
  - domain: binance
    llm_title: Binance Market Data
    base_url: https://api.binance.com
    collections:
      - llm_title: Market Data
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml
    auth:
      type: hmac
      config:
        api_key: "$(BINANCE_API_KEY)"
        secret_key: "$(BINANCE_SECRET_KEY)"
```

## パラメーター

| パラメーター | 必須 | 説明 |
|-----------|------|------|
| `api_key` | はい | 公開 API キー |
| `secret_key` | はい | 署名用の秘密鍵 |

## 注意点

- swag2mcp はすべてのリクエストに自動的にタイムスタンプ（Unix ミリ秒）を追加します
- 署名はすべてのリクエストパラメーターから計算されます
- キーは環境変数に保存：`api_key: "$(BINANCE_API_KEY)"`
- この方法は Binance API および類似の取引所と互換性があります
