# Tags

タグは、collection 内の関連するエンドポイントをグループ化するカテゴリです。タグは存在する場合も存在しない場合もあります — すべての collection にタグがあるわけではなく、collection は任意の数のタグを持つことができます。

タグは OpenAPI/Swagger/Postman ファイル自体から取得されます。タグ用の **YAML 設定はありません** — `swag2mcp.yaml` でタグを作成、名前変更、削除することはできません。タグを変更する唯一の方法は、元の spec ファイルを編集することです。

## 階層

```
Spec (domain, e.g. "meteo")
  └── Collection (spec file, e.g. forecast.yml)
        └── Tag "weather"
              └── GET /forecast
              └── GET /forecast/hourly
        └── Tag "alerts"
              └── GET /alerts
```

## タグの作成方法

タグは解析中に spec ドキュメントから抽出されます：

**OpenAPI 3.x / Swagger 2.0** — 各操作の `tags` リストがタグになります：

```yaml
paths:
  /pet:
    get:
      tags: ["pets"]
      summary: "Find pet by ID"
    post:
      tags: ["pets"]
      summary: "Add a new pet"
  /pet/{petId}/uploadImage:
    post:
      tags: ["pet_images"]
      summary: "Uploads an image"
```

**Postman** — 各トップレベルフォルダーがタグになります。ネストされたフォルダーは最後のフォルダー名を使用します。

エンドポイントにタグがない場合、`"default"` タグの下に配置されます。

## 目的

タグは LLM が関連するエンドポイントのグループを見つけるのに役立ちます。LLM は collection 内のすべてのエンドポイントを検索する代わりに、最初に適切なタグを見つけ、その中のエンドポイントのみを一覧表示できます。

## Tag 用 MCP ツール

| ツール | 説明 |
|-------|------|
| `tag_by_spec` | spec 全体のすべてのタグ |
| `tag_by_collection` | 特定の collection 内のタグ |
| `tag_by_id` | タグの詳細（タイトル、メソッド数） |
| `endpoint_by_tag` | タグの下にグループ化されたエンドポイント |

## 例

```
Query: "Show all tags in the pet collection"
→ tag_by_collection(collectionId: "...")
→ Result: pets (5 methods), pet_images (1 method)
```

## 制限事項

- タグは設定の観点からは読み取り専用です。タグを追加、名前変更、削除するには、元の OpenAPI/Swagger/Postman ファイルを編集し、`swag2mcp update` を実行します。
- YAML 設定で collection ごとにタグをフィルタリングしたり無効にしたりすることはできません。
