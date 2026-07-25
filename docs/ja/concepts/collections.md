# Collections

Collection は、特定の API を記述する単一の OpenAPI/Swagger/Postman ファイルです。`location`（URL またはローカルファイルパス）を指し、spec（ドメイン）に属します。

1 つの spec は複数の collection を持つことができます — 例えば、"meteo" spec には "Forecast"、"Air Quality"、"Marine" の collection があり、それぞれ異なる spec ファイルを指します。

## Collection フィールド

| フィールド | YAML キー | 必須 | 説明 |
|-----------|----------|------|------|
| [LLM Title](#llm-instruction) | `llm_title` | ❌ | LLM 用の collection 表示名（最大 120 文字）。未設定時は spec ドキュメントから自動入力 |
| [LLM Instruction](#llm-instruction) | `llm_instruction` | ❌ | LLM 向けの短いヒント（最大 360 文字）。未設定時は spec ドキュメントから自動入力 |
| Title | `title` | ❌ | 元の spec タイトルの上書き（解析されたドキュメントから自動入力） |
| [Location](#location--how-spec-files-are-resolved) | `location` | ✅ | spec ファイルの URL またはパス（5〜250 文字） |
| [Disable](#disable) | `disable` | ❌ | 読み込み時にこの collection をスキップ |
| [HTTP Client](#http-client-override) | `http_client` | ❌ | collection ごとの HTTP 設定（ヘッダー、Cookie） |
| [Base URL](#base-url-override) | `base_url` | ❌ | この collection の spec のベース URL を上書き |
| [Mock Server](#mock-server) | `base_mock_url` | ❌ | `host:port` 形式のモックサーバーアドレス。`mock_enabled: true` 時に必須 |

## Location — Spec ファイルの解決方法

`location` フィールドは swag2mcp に OpenAPI/Swagger/Postman ファイルの場所を指示します。複数のソースタイプをサポートします：

| ソース | 例 | 説明 |
|-------|-----|------|
| **リモート URL** | `https://raw.githubusercontent.com/.../spec.yaml` | ダウンロードしてキャッシュ |
| **ローカルファイル（絶対パス）** | `/home/user/my-api.yaml` | ファイルシステムから読み取り、キャッシュ |
| **ローカルファイル（相対パス）** | `./my-api.yaml` | 絶対パスに解決、キャッシュ |
| **ワークスペースローカルファイル** | `specs/my-api.yaml` | `~/.swag2mcp/specs/` に保存、直接使用（キャッシュなし） |
| **file:// URI** | `file:///home/user/spec.yaml` | ローカルパスに変換、キャッシュ |

swag2mcp は自動的にソースタイプを検出します：

- `https://` または `http://` → リモート URL（キャッシュ）
- `file://` → ローカルファイル（ファイルシステムパスに変換）
- その他 → ローカルファイル（ホームディレクトリの `~` 展開あり）

### リモート URL

リモート URL を使用すると、swag2mcp はファイルをダウンロードしてローカルにキャッシュします。キャッシュは後続の起動時に再利用され、繰り返しのダウンロードを避けます。

### ローカルファイル

ローカルファイルはファイルシステムから直接読み取られます。ファイルがワークスペースの `specs/` ディレクトリ外にある場合、一貫性のためにキャッシュにコピーされます。

### ワークスペースローカルファイル

ワークスペース内の `specs/` ディレクトリ（`~/.swag2mcp/specs/`）は、ローカル spec ファイルの推奨場所です。ここに保存されたファイルはキャッシュなしで直接使用されます。参照するには `specs/` で始まる相対パスを使用します。

> **注:** `specs/` は単なるディレクトリ名（`cache/` や `responses/` と同様）であり、「spec」という概念ではありません。collection が指す実際の OpenAPI/Swagger/Postman ファイルを保存します。

```bash
# spec ファイルをワークスペースにインポート
swag2mcp import https://example.com/api.yaml myspec

# インポート後、location は次のようになります：
# specs/myspec.yaml
```

## キャッシュシステム

swag2mcp はリモート spec ファイルをキャッシュして、起動のたびにダウンロードするのを避けます。

### 仕組み

1. リモート URL の collection が読み込まれると、swag2mcp はキャッシュをチェックします
2. 有効な（期限切れでない）キャッシュエントリが存在する場合、それが直接使用されます
3. 存在しない場合、ファイルがダウンロードされ、解析され、キャッシュに保存されます

### キャッシュ構造

```
~/.swag2mcp/
  cache/
    {sha256_hash}.spec    # キャッシュされた spec ファイルの内容
    {sha256_hash}.meta    # キャッシュメタデータ（JSON）
```

各キャッシュファイルには、以下を含むメタデータファイルがあります：

```json
{
  "source": "https://example.com/api.yaml",
  "source_type": "url",
  "cached_at": "2024-01-01T00:00:00Z",
  "mod_time": "2024-01-01T00:00:00Z",
  "ttl_sec": 3600
}
```

### キャッシュ TTL

各キャッシュファイルには **1 時間から 48 時間** の間でランダムな TTL が設定されます。これにより、すべてのキャッシュファイルが同時に期限切れになるのを防ぎます（群集問題）。

### キャッシュキー

キャッシュキーは、生の location 文字列の SHA-256 ハッシュです（最初の 16 バイト = 32 桁の 16 進数）。

### キャッシュの管理

```bash
# キャッシュとレスポンスをクリアし、すべての spec ファイルを再ダウンロード
swag2mcp update

# キャッシュとレスポンスのみをクリア
swag2mcp clean
```

- `swag2mcp update` — 設定を検証し、`cache/` と `responses/` をクリアし、すべての collection location を再キャッシュ
- `swag2mcp clean` — `cache/` と `responses/` のすべての内容と、孤立した認証スクリプトを削除
- 古いレスポンスは MCP サーバー起動後 48 時間で自動的にクリーンアップ

## 検証

すべての collection は設定が読み込まれるときに検証されます。検証は `swag2mcp mcp` の起動のたびに実行されます。失敗した場合、MCP サーバーは起動しません — 一部の IDE では、サーバーが単に接続せず、LLM は何を修正すべきかを説明する明確なエラーメッセージを受け取ります。

| チェック | ルール |
|---------|-------|
| **Location** | 必須、5〜250 文字 |
| **Location のアクセス可能性** | 到達可能な URL または既存のファイルである必要があります |
| **Location の有効性** | 有効な OpenAPI 3.x、Swagger 2.0、または Postman ファイルである必要があります |
| **LLM Title** | 最大 120 文字、英字/数字/基本句読点 |
| **LLM Instruction** | 最大 360 文字、タイトルと同じ文字セット |
| **Base URL** | 設定されている場合、有効な URL である必要があります |
| **Base Mock URL** | `host:port` または `host:port/path` 形式で、host は `localhost`、`127.0.0.1`、または `0.0.0.0` |
| **Mock 必須** | `mock_enabled: true` の場合、すべての collection に `base_mock_url` が必要 |
| **重複モックポート** | 2 つの collection が同じモックポートを共有してはいけません |

サーバーを起動する前に問題を診断するには、[`validate`](../cli/validate.md) コマンドを使用します：

```bash
# デフォルトワークスペースを検証（~/.swag2mcp）
swag2mcp validate

# カスタムプロジェクトワークスペースを検証
swag2mcp validate ./my-project
```

## Collection の追加

### YAML 設定経由

`~/.swag2mcp/swag2mcp.yaml` を直接編集します：

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

編集後、変更を反映するために MCP サーバー（`swag2mcp mcp`）を再起動します。

### CLI 経由

```bash
# 対話モード
swag2mcp add collection

# YAML を使用した非対話モード
swag2mcp add collection --yaml 'spec_domain: meteo
llm_title: Forecast
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml'

# 標準入力からパイプ
cat collection.yaml | swag2mcp add collection --yaml -

# YAML 例を表示
swag2mcp add collection --example
```

### Import 経由

```bash
# spec ファイルをワークスペースにインポート
swag2mcp import https://example.com/api.yaml
```

## LLM Instruction

Collection は、より具体的なガイダンスのために独自の `llm_instruction`（最大 360 文字）を持つことができます。これは spec レベルの指示とともに swag2mcp システムプロンプトに注入されます。

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        llm_instruction: "Use this collection for current weather and daily forecasts."
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        llm_instruction: "Use this collection for air quality index and pollution data."
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
```

`llm_title` が設定されていない場合、spec ドキュメントの `title` フィールドから自動的に入力されます。`llm_instruction` が設定されていない場合、spec ドキュメントの `description` フィールドから入力されます。

## Disable

`disable: true` を設定して collection をスキップします。読み込まれず、インデックス化されず、LLM が利用できなくなります。

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
        disable: true
```

## Base URL の上書き

各 collection は spec の `base_url` を上書きできます。これは同じ spec 内の異なる collection が異なる API エンドポイントを使用する場合に便利です。

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

## HTTP Client の上書き

Collection は spec およびグローバルレベルから HTTP 設定（ヘッダー、Cookie）を上書きできます。

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          headers:
            X-API-Version: "2"
          cookies:
            - name: session
              value: abc123
```

設定のカスケード：グローバル → spec → collection。詳細は [Configuration Cascade](../configuration/cascade.md) を参照してください。

## モックサーバー

設定レベルで `mock_enabled: true` が設定されている場合、すべての collection に `base_mock_url` が設定されている必要があります。これは、この collection のモックサーバーがどこで実行されているかを swag2mcp に伝えます。

```yaml
mock_enabled: true
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        base_mock_url: localhost:8080
```

詳細は [Mock Server](../advanced/mock-server.md) を参照してください。

## 例

### 最小限の Collection

```yaml
specs:
  - domain: dadjokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

### 全フィールドの Collection

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        llm_instruction: "Use for current weather and daily forecasts."
        title: "Custom Title"
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        disable: false
        base_url: https://forecast-api.open-meteo.com
        base_mock_url: localhost:8080
        http_client:
          headers:
            X-Custom: value
```

### 1 つの Spec に複数の Collection

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

### ワークスペース内のローカルファイル（specs/ ディレクトリ）

```yaml
specs:
  - domain: myapi
    llm_title: My Internal API
    base_url: https://api.mycompany.com
    collections:
      - llm_title: Users
        location: specs/users.openapi.json
      - llm_title: Orders
        location: specs/orders.openapi.json
```

### 無効化された Collection

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
        disable: true
```

## 関連項目

- [Collection Settings (config)](../configuration/collection-settings.md) — 完全な YAML リファレンス
- [Configuration Cascade](../configuration/cascade.md) — 設定の上書き方法
- [Specs](./specs) — collection の論理コンテナ
- [HTTP Client](../configuration/http-client.md) — HTTP クライアント設定
- [Mock Server](../advanced/mock-server.md) — モックサーバー設定
- [CLI: validate](../cli/validate.md) — validate コマンドリファレンス
- [CLI: update](../cli/update.md) — update コマンドリファレンス
- [CLI: clean](../cli/clean.md) — clean コマンドリファレンス
