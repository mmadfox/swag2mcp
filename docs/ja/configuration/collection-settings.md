# Collection 設定

Collection 設定は、単一の OpenAPI/Swagger/Postman spec ファイルを定義し、その特定のファイルの spec 設定を上書きします。各 collection は spec に属し、1 つの API 仕様ドキュメントを表します。

## Collection セクション

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        llm_instruction: "Use for current and forecast weather data"
        disable: false
        base_url: https://forecast-api.open-meteo.com
        base_mock_url: localhost:8081
        http_client:
          timeout: 5s
```

## パラメーター

### llm_title

- **型:** `string`
- **必須:** いいえ
- **説明:** この collection の人間可読名。MCP ツールのレスポンスに表示されます。
- **ルール:** 最大 120 文字。英字、数字、スペース、基本句読点のみ。
- **例:** `Forecast`、`Air Quality`、`Market Data`

### llm_instruction

- **型:** `string`
- **デフォルト:** `""`
- **説明:** この特定の collection に関する LLM への指示。この collection が提供するエンドポイントを説明します。
- **ルール:** 最大 360 文字。英字、数字、スペース、基本句読点のみ。
- **例:** `"Use for current and forecast weather data."`

### title

- **型:** `string`
- **デフォルト:** `""`
- **説明:** spec ファイルからの生のタイトル。実行時に自動的に入力されます。通常、YAML で設定する必要はありません。

### location

- **型:** `string`
- **必須:** はい
- **説明:** OpenAPI 3.x、Swagger 2.0、または Postman collection の spec ファイルへの URL またはローカルファイルパス。
- **ルール:** 5〜250 文字。
- **例:**
  - URL: `https://raw.githubusercontent.com/org/repo/main/spec.yaml`
  - ローカル: `./specs/my-api.json`
  - ローカル（絶対パス）: `/home/user/.swag2mcp/specs/my-api.yaml`

### disable

- **型:** `bool`
- **デフォルト:** `false`
- **説明:** `true` の場合、この collection は MCP ツールから除外されます。読み込まれず、インデックス化もされません。
- **使用するタイミング:** 設定から削除せずに collection を一時的に無効化。spec ファイルが更新中の場合や API バージョンが非推奨の場合に便利です。

### http_client

- **型:** `object`
- **デフォルト:** spec（またはグローバル）から継承
- **説明:** この collection の HTTP クライアント設定を上書き。グローバル `http_client` のすべての設定を上書き可能：`timeout`、`max_response_size`、`user_agent`、`follow_redirects`、`max_redirects`、`random`、`proxy`、`headers`、`cookies`。
- **例:**
  ```yaml
  http_client:
    timeout: 120s
    headers:
      "X-Custom": "value"
    cookies:
      - name: "session"
        value: "abc123"
  ```

### base_url

- **型:** `string`
- **デフォルト:** `""`（spec から継承）
- **説明:** この collection の spec レベルの `base_url` を上書き。同じ spec 内の異なる collection が異なるベース URL を使用する場合に使用します。
- **例:** spec が `base_url: https://api.open-meteo.com` を持ち、ある collection が `https://air-quality-api.open-meteo.com` を使用する場合、collection レベルで `base_url` を設定します。

### base_mock_url

- **型:** `string`
- **デフォルト:** `""`
- **説明:** `host:port` 形式のモックサーバーアドレス。グローバル設定で `mock_enabled: true` の場合に必須。
- **ルール:** ホストは `localhost`、`127.0.0.1`、または `0.0.0.0` である必要があります。ポートは有効なポート番号である必要があります。
- **例:** `localhost:8081`、`127.0.0.1:9000`
- **使用するタイミング:** `mock_enabled: true` で、この collection を偽のレスポンスでテストしたい場合。

## 1 つの Spec からの複数 Collection

spec は複数の collection を持つことができます — 例えば、API が異なるサービス用に別々の spec ファイルを持つ場合：

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

## Collection の無効化

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: true
```

## HTTP Client の上書き

すべての `http_client` 設定は collection レベルで上書き可能です。Collection の値は、この collection に限り、spec およびグローバルの値より優先されます。

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          timeout: 120s
          headers:
            "X-Custom": "value"
          cookies:
            - name: "session"
              value: "abc123"
```
