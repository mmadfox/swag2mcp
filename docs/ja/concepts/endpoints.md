# Endpoints

エンドポイントは、呼び出し可能な特定の HTTP メソッド + パスです（例：`GET /api/users/{id}`）。エンドポイントは、LLM が発見、調査、呼び出しする実際の API 操作です。

## 構造

各エンドポイントには以下が含まれます：

- **HTTP メソッド**: GET、POST、PUT、PATCH、DELETE、HEAD、OPTIONS
- **パス**: `/api/v1/users/{id}`
- **サマリー**: エンドポイントの機能の短い説明 — LLM が一目で目的を理解するのに非常に便利
- **説明**: エンドポイントの動作、パラメーター、ユースケースの詳細な説明
- **パラメーター**: パス、クエリ、ヘッダー、Cookie
- **リクエストボディ**: POST/PUT/PATCH 用
- **レスポンス**: ステータスコードとレスポンススキーマ

`summary` と `description` フィールドは OpenAPI/Swagger/Postman ファイルから取得されます。これらは LLM がエンドポイントの機能を理解する主要な手段です。適切に書かれたサマリーは、エンドポイントの発見をより効果的にします。

## Endpoint 用 MCP ツール

| ツール | 説明 |
|-------|------|
| `endpoint_by_spec` | spec 内のすべてのエンドポイント |
| `endpoint_by_collection` | collection 内のエンドポイント |
| `endpoint_by_tag` | タグ内のエンドポイント |
| `endpoint_by_id` | クイックエンドポイントサマリー |
| `inspect` | 完全なエンドポイント詳細（スキーマ、パラメーター） |
| `invoke` | エンドポイントの呼び出し |
| `search` | テキストでエンドポイントを検索 |

## 非推奨のエンドポイント

spec で `deprecated` とマークされたエンドポイントは、調査時に通知とともに表示されます。

## 設定

エンドポイントは swag2mcp の観点からは**読み取り専用**です。エンドポイント用の YAML 設定はありません — `swag2mcp.yaml` でエンドポイントを追加、削除、名前変更、変更することはできません。

エンドポイントを変更するには（新しいものの追加、サマリーの更新、パラメーターの変更、非推奨のマーク）、元の OpenAPI/Swagger/Postman ファイルを編集し、`swag2mcp update` を実行して再解析と再インデックス化を行います。

## 例

```
Query: "Show details for GET /pet/{petId}"
→ inspect(endpointId: "abc123...")
→ Result:
  GET /pet/{petId}
  Summary: Find pet by ID
  Description: Returns a single pet by its ID
  Parameters:
    - petId (path, integer, required)
  Responses:
    - 200: Pet object
    - 400: Error
    - 404: Not found
```
