# 概念

## アーキテクチャ

swag2mcp は API 仕様と LLM エージェントの間のブリッジとして機能します：

<img src="/architecture.svg" width="800" alt="swag2mcp アーキテクチャ">

## コアコンセプト

**Spec** — API ドメインまたはサービスを表す論理コンテナ（例：YouTube、Binance、Open-Meteo）。各 spec は一意の `domain`、`base_url`、オプションの `auth` を持ち、1 つ以上の collection を含みます。また、`llm_instruction` を設定できます — swag2mcp システムプロンプトに注入される短いヒントで、LLM にこの spec の目的と使用タイミングを伝えます。詳細：[Specs](./specs)。

**Collection** — 特定の API を記述する単一の OpenAPI/Swagger/Postman ファイル。`location`（URL またはローカルファイルパス）を指します。1 つの spec は複数の collection を持つことができます — 例えば、"meteo" spec には "Forecast"、"Air Quality"、"Marine" の collection があり、それぞれ異なる spec ファイルを指します。詳細：[Collections](./collections)。

**Tag** — collection 内のエンドポイントのカテゴリ。LLM が適切な操作をより正確に見つけるのに役立ちます。詳細：[Tags](./tags)。

**Endpoint** — 特定の HTTP メソッド + パス（例：`GET /api/users`）。LLM は説明でエンドポイントを見つけ、パラメーターとスキーマを調査し、呼び出すことができます。詳細：[Endpoints](./endpoints)。

**Workspace** — swag2mcp が設定、spec キャッシュ、保存されたレスポンス、認証スクリプトを保存するディレクトリ。詳細：[Workspace](./workspace)。

## 仕組み

1. **spec または collection を追加** — YAML 設定（`~/.swag2mcp/swag2mcp.yaml`）で定義します。例：

   ```yaml
   specs:
     - domain: jokes
       llm_title: Dad Joke API
       base_url: https://icanhazdadjoke.com
       collections:
         - llm_title: Jokes
           location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
   ```
2. **swag2mcp が各 collection を解析** — タグとエンドポイントを作成し、検索用にインデックス化します。
3. **LLM が適切なエンドポイントを見つける** — MCP ツール（`search`、`endpoint_by_tag`、`inspect`）を通じて、LLM は説明で一致するエンドポイントを検索し、パラメーターとリクエストスキーマを確認します。
4. **LLM がエンドポイントを呼び出す** — MCP ツール `invoke` を介して、LLM がリクエストを送信します。swag2mcp は呼び出し前にすべての入力パラメーターをエンドポイントの OpenAPI スキーマに対して検証します（パスパラメーター、クエリパラメーター、ヘッダー、リクエストボディ）。スキーマに一致しないものがある場合、LLM は何が問題かを説明する明確なエラーを受け取ります。検証後、swag2mcp は実際の HTTP 呼び出しを実行し、結果を返します。
5. **結果が LLM に返される** — API レスポンスはエージェントに渡されます。大きなレスポンスはワークスペースに保存され、3 つの専用 MCP ツール（`response_outline`（構造の確認）、`response_compress`（代表的なサンプルに縮小）、`response_slice`（特定の断片の抽出））で探索できます。

swag2mcp は LLM と API の世界の間のブリッジです。API 仕様を追加すると、LLM は MCP プロトコルを通じて適切なエンドポイントを見つけ、そのドキュメントを調査し、呼び出します。必要なのは spec を追加して MCP サーバーを起動することだけです。

> **設定はいつでも編集可能です。** YAML 設定ファイル（`~/.swag2mcp/swag2mcp.yaml`）は手動で編集できます — spec の追加、認証の変更、設定の調整。編集後は毎回 MCP サーバー（`swag2mcp mcp`）を再起動して変更を反映してください。

## 階層

```
Spec (domain, e.g. "meteo")
  └── Collection 1 (spec file, e.g. forecast.yml)
        └── Tag 1 (category)
              └── Endpoint (GET /api/forecast)
              └── Endpoint (POST /api/forecast)
        └── Tag 2
              └── Endpoint (GET /api/forecast/{id})
  └── Collection 2 (spec file, e.g. air-quality.yml)
        └── Tag 3
              └── Endpoint (GET /api/air-quality)
```
