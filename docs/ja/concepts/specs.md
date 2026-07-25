# Specs

Spec は API ドメインまたはサービスを表す論理コンテナです（例：YouTube、Binance、Open-Meteo）。各 spec は一意の `domain`、`base_url`、オプションの `auth` を持ち、1 つ以上の collection を含みます。

[Collections](./collections) は OpenAPI/Swagger/Postman ファイルを指します — spec 自体はファイルではなく、それらをグループ化するものです。

## Domain — 命名規則

`domain` は spec の一意識別子です。システム全体で主キーとして使用されます。

| ルール | 制約 |
|-------|------|
| 文字 | `a-z`、`0-9`、`_`、`-` のみ |
| 長さ | 1〜60 文字 |
| 一意性 | **重複不可** — 2 つのアクティブな spec が同じ domain を共有できません |

**有効な例:** `meteo`、`binance`、`github-api`、`my_service`、`openai-v1`

**無効な例:** `Meteo`（大文字）、`my api`（スペース）、`my.api`（ドット）、`a-very-long-domain-name-that-exceeds-sixty-characters`（長すぎる）

## Spec フィールド

| フィールド | YAML キー | 必須 | 説明 |
|-----------|----------|------|------|
| [Domain](#domain--naming-rules) | `domain` | ✅ | 一意の API 識別子（1〜60 文字、`a-z0-9_-`） |
| LLM Title | `llm_title` | ✅ | LLM がこの API を参照するための人間可読名（5〜120 文字） |
| [LLM Instruction](#llm-instruction) | `llm_instruction` | ❌ | swag2mcp システムプロンプトに注入される短いヒント（最大 500 文字） |
| Base URL | `base_url` | ✅ | すべての API リクエストのベース URL（有効な URL） |
| [Disable](#disable) | `disable` | ❌ | 読み込みとインデックス化時にこの spec をスキップ |
| [Tags](#tags) | `tags` | ❌ | フィルタリング用のタグ（例：`["public", "demo"]`） |
| [Auth](#auth) | `auth` | ❌ | 認証設定 |
| [HTTP Client](#http-client) | `http_client` | ❌ | spec ごとの HTTP 設定（ヘッダー、Cookie） |
| [Collections](./collections) | `collections` | ✅ | 1〜30 の collection のリスト |

## 検証

swag2mcp が設定を検証するとき、各 spec に対して以下のルールがチェックされます：

| チェック | ルール |
|---------|-------|
| **重複ドメイン** | 2 つのアクティブな spec が同じ `domain` を共有してはいけません |
| **ドメイン形式** | `^[a-z0-9_-]{1,60}$` に一致する必要があります |
| **LLM Title** | 必須、5〜120 文字、英字/数字/スペース/基本句読点 |
| **LLM Instruction** | 最大 500 文字、タイトルと同じ文字セット |
| **Base URL** | 必須、有効な URL である必要があります |
| **Collections** | 必須、1〜30 項目 |
| **Auth** | 認証タイプごとに検証（例：bearer は `token`、basic は `username` + `password` が必要） |
| **Location** | 各 collection の `location` は有効な URL またはファイルパスである必要があります（5〜250 文字） |

検証は `swag2mcp mcp` の起動のたびに実行されます。失敗した場合、MCP サーバーは起動しません — 一部の IDE では、サーバーが単に接続せず、LLM は何を修正すべきかを説明する明確なエラーメッセージを受け取ります。

サーバーを起動する前に問題を診断するには、[`validate`](../cli/validate.md) コマンドを使用します：

```bash
# デフォルトワークスペースを検証（~/.swag2mcp）
swag2mcp validate

# カスタムプロジェクトワークスペースを検証
swag2mcp validate ./my-project
```

## LLM Instruction

各 spec に `llm_instruction` を設定することを推奨します — この API の目的と使用タイミングを LLM に伝える短いヒント（最大 500 文字）です。この指示は swag2mcp システムプロンプトに注入され、追加のコンテキストなしで LLM が spec の目的を理解するのに役立ちます。

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    llm_instruction: "Use this API to get random dad jokes or search for specific jokes by keyword."
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

Collection も独自の `llm_instruction`（最大 360 文字）を持ち、より具体的なガイダンスを提供できます。

## Auth

認証は spec レベルで設定され、そのすべての collection に適用されます。swag2mcp は 9 つの認証方式をサポートしています：

| 方式 | YAML タイプ | 主要フィールド |
|------|-----------|------------|
| [None](../auth/none.md) | `none` | — |
| [Basic](../auth/basic.md) | `basic` | `username`、`password` |
| [Bearer](../auth/bearer.md) | `bearer` | `token` |
| [Digest](../auth/digest.md) | `digest` | `username`、`password` |
| [OAuth2 Client Credentials](../auth/oauth2-cc.md) | `oauth2-cc` | `client_id`、`client_secret`、`token_url` |
| [OAuth2 Password](../auth/oauth2-pwd.md) | `oauth2-pwd` | `username`、`password`、`client_id`、`token_url` |
| [API Key](../auth/api-key.md) | `api-key` | `key`、`value`、`in`（`header` または `query`） |
| [HMAC](../auth/hmac.md) | `hmac` | `api_key`、`secret_key` |
| [Script](../auth/script.md) | `script` | `domain` |

各方式の詳細は [Auth Overview](../auth/overview.md) を参照してください。

## HTTP Client

spec レベルで HTTP 設定を上書きできます。これらはこの spec の collection が行うすべてのリクエストに適用されます。

```yaml
specs:
  - domain: slow-api
    llm_title: Slow API
    base_url: https://slow-api.example.com
    http_client:
      headers:
        X-API-Version: "2"
      cookies:
        - name: session
          value: abc123
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

設定のカスケード：グローバル → spec → collection。詳細は [Configuration Cascade](../configuration/cascade.md) を参照してください。

## Tags

タグを使用すると、カテゴリで spec をフィルタリングできます。`swag2mcp ls` またはブートストラップ時に `--tags` フラグとともに使用します。

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    tags: ["weather", "public"]
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

```bash
# "weather" タグが付いた spec のみを一覧表示
swag2mcp ls --tags weather
```

## Disable

`disable: true` を設定して spec を完全にスキップします。読み込まれず、インデックス化されず、LLM が利用できなくなります。

```yaml
specs:
  - domain: old-api
    llm_title: Old API (Deprecated)
    base_url: https://old-api.example.com
    disable: true
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## 例

### 最小限の Spec

```yaml
specs:
  - domain: dadjokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

### 認証付き Spec

```yaml
specs:
  - domain: binance
    llm_title: Binance Market Data API
    base_url: https://api.binance.com
    auth:
      type: hmac
      config:
        api_key: $(BINANCE_API_KEY)
        secret_key: $(BINANCE_SECRET_KEY)
    collections:
      - llm_title: Market Data
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml
```

### 複数 Collection の Spec

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

### LLM Instruction と Tags 付き Spec

```yaml
specs:
  - domain: rickandmorty
    llm_title: Rick and Morty API
    llm_instruction: "Use this API to get information about characters, episodes, and locations from the Rick and Morty show."
    base_url: https://rickandmortyapi.com/api
    tags: ["entertainment", "public"]
    collections:
      - llm_title: Characters
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/rick-and-morty.json
```

## 関連項目

- [Spec Settings (config)](../configuration/spec-settings.md) — 完全な YAML リファレンス
- [Configuration Cascade](../configuration/cascade.md) — 設定の上書き方法
- [Auth Overview](../auth/overview.md) — 全 9 認証方式
- [HTTP Client](../configuration/http-client.md) — HTTP クライアント設定
- [Collections](./collections) — spec 内の spec ファイル
