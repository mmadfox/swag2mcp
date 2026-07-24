# Spec 設定

Spec 設定は API サービスを定義し、その特定の API のグローバル設定を上書きします。各 spec は 1 つの論理 API（例：「Open-Meteo Weather APIs」）を表し、複数の collection（spec ファイル）を含むことができます。

## Spec セクション

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    llm_instruction: "Use this API for weather forecasts and climate data"
    base_url: https://api.open-meteo.com
    disable: false
    tags: ["weather", "climate"]
    http_client:
      timeout: 10s
      max_response_size: 1024
    auth:
      type: bearer
      config:
        token: "$(TOKEN)"
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

## パラメーター

### domain

- **型:** `string`
- **必須:** はい
- **説明:** この API spec の一意識別子。内部で spec を参照するために使用されます。
- **ルール:** 1〜60 文字。小文字（`a-z`）、数字（`0-9`）、ハイフン（`-`）、アンダースコア（`_`）のみ。
- **例:** `meteo`、`binance`、`my-api`

### llm_title

- **型:** `string`
- **必須:** はい
- **説明:** LLM がこの API を参照するための人間可読名。MCP ツールのレスポンスに表示されます。
- **ルール:** 5〜120 文字。英字、数字、スペース、基本句読点のみ。
- **例:** `Open-Meteo Weather APIs`、`Binance Market Data`

### llm_instruction

- **型:** `string`
- **デフォルト:** `""`
- **説明:** この API の使用方法に関する LLM への指示。API の機能と使用タイミングを説明します。
- **ルール:** 最大 500 文字。英字、数字、スペース、基本句読点のみ。
- **例:** `"Use this API for weather forecasts, current conditions, and climate data."`

### base_url

- **型:** `string`
- **必須:** はい
- **説明:** この spec のすべての API リクエストのベース URL。OpenAPI spec のエンドポイントパスがこの URL に追加されます。
- **例:** `https://api.open-meteo.com`、`https://api.binance.com`
- **注:** 異なる collection が異なるベース URL を使用する場合、collection レベルで上書きできます。

### disable

- **型:** `bool`
- **デフォルト:** `false`
- **説明:** `true` の場合、この spec は MCP ツールから除外されます。読み込まれず、インデックス化もされず、LLM が利用できなくなります。
- **使用するタイミング:** 設定から削除せずに API を一時的に無効化。ダウンしている、非推奨の、またはメンテナンス中の API に便利です。

### tags

- **型:** `[]string`（文字列の配列）
- **デフォルト:** `[]`
- **説明:** spec をフィルタリングするためのタグ。CLI コマンド（`ls`、`validate`、`mcp`、`update`）の `--tags` フラグとともに使用します。
- **例:** `["public", "weather"]`、`["internal", "production"]`
- **効果:** `swag2mcp mcp --tags=public` を実行すると、`public` タグを持つ spec のみが読み込まれます。

### http_client

- **型:** `object`
- **デフォルト:** グローバルから継承
- **説明:** この spec のグローバル HTTP クライアント設定を上書き。グローバル `http_client` のすべての設定を上書き可能：`timeout`、`max_response_size`、`user_agent`、`follow_redirects`、`max_redirects`、`random`、`proxy`、`headers`、`cookies`。
- **例:**
  ```yaml
  http_client:
    timeout: 60s
    max_response_size: 4194304
    headers:
      "X-DC": "us-east-1"
  ```

### auth

- **型:** `object`
- **デフォルト:** `none`（認証なし）
- **説明:** この spec の認証設定。全 9 方式とそのパラメーターについては [Authentication](/auth/overview) セクションを参照してください。
- **例:**
  ```yaml
  auth:
    type: bearer
    config:
      token: "$(API_TOKEN)"
  ```

### collections

- **型:** `[]object`（collection の配列）
- **必須:** はい（最低 1 つ）
- **説明:** この spec に属する OpenAPI/Swagger/Postman spec ファイルのリスト。各 collection は 1 つの spec ファイルです。
- **ルール:** spec あたり 1〜30 の collection。
- **参照:** すべての collection パラメーターについては [Collection Settings](./collection-settings) を参照してください。

## Spec の無効化

無効化された spec は読み込まれず、インデックス化もされません。LLM はそれらを表示したり使用したりできません。

```yaml
specs:
  - domain: old-api
    llm_title: Old API
    base_url: https://old-api.example.com
    disable: true
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## HTTP Client の上書き

グローバルレベルのすべての `http_client` 設定は spec レベルで上書き可能です。spec の値は、この spec に限り、グローバルの値より優先されます。

```yaml
specs:
  - domain: slow-api
    llm_title: Slow API
    base_url: https://slow-api.example.com
    http_client:
      timeout: 120s
      max_response_size: 8388608
      headers:
        "X-DC": "us-east-1"
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Proxy の上書き

この spec がグローバルとは異なるプロキシを必要とする場合、spec レベルで設定します：

```yaml
specs:
  - domain: proxied-api
    llm_title: Proxied API
    base_url: https://api.example.com
    http_client:
      proxy:
        url: http://proxy.company.com:8080
        username: $(PROXY_USER)
        password: $(PROXY_PASS)
        bypass:
          - "*.local"
          - "10.0.0.0/8"
    collections:
      - llm_title: Main
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo.json
```
