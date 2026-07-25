# 実行ツール

実行ツールは swag2mcp の中核です。**search** は ID がない場合にエンドポイントを検索し、**inspect** は完全な OpenAPI 契約を明らかにし、**invoke** は実際の API 呼び出しを行います。常に search → inspect → invoke の順序で使用してください。

---

## search

### 目的

エンドポイント ID がない場合にエンドポイントを検索する唯一のツールです。bluge 検索エンジンを使用して、すべてのスペックのすべてのエンドポイントを全文検索します。

### 使用するタイミング

- エンドポイント ID がわからない場合
- キーワード、メソッド、タグ、パスでエンドポイントを検索したい場合
- 特定の機能にどのエンドポイントが存在するかを発見する必要がある場合

### 動作方法

すべてのスペックにわたって全文インデックスを検索します。フィールドフィルター、ブール演算子、あいまい検索、ワイルドカードなどを使用した構造化クエリをサポートします。

### パラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|----------|-------------|
| `query` | string | はい | 検索クエリ（構造化構文をサポート） |
| `limit` | int | はい | 返す最大結果数（1〜50） |

### クエリ構文

| 例 | 説明 |
|---------|-------------|
| `pet` | すべてのフィールドに対するシンプルなテキスト検索 |
| `method:GET` | HTTP メソッドでフィルタリング |
| `tag:pet` | タグ名でフィルタリング |
| `path:"/api/v1/users"` | 完全一致パス検索 |
| `+method:POST +tag:pet` | 両方の条件に一致する必要あり |
| `-method:DELETE` | DELETE メソッドを除外 |
| `create~` | あいまい検索（タイポ許容） |
| `path:/api/v1/*` | ワイルドカードパス検索 |
| `/pattern/` | 正規表現検索 |
| `term^3` | 用語の関連性をブースト |

**検索可能なフィールド：** `method`（キーワード）、`tag`（キーワード）、`path`（テキスト）、`summary`（テキスト）、`_all`（デフォルトのテキストフィールド）。

**サポート外：** 括弧によるグループ化、明示的な `AND`/`OR` 演算子、フィールドグループ化。

### レスポンス

```json
{
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "tagName": "forecast",
      "collectionId": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "collectionTitle": "Weather Forecast",
      "specId": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
      "specDomain": "meteo",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Get weather forecast for a location"
    }
  ]
}
```

各結果には完全な祖先（スペック → コレクション → タグ）が含まれるため、LLM は関連するエンドポイントにナビゲートできます。

### 補足

- `limit` は 1 〜 50 の間である必要があります（それ以外の場合は `validation_failed` を返します）
- `query` は必須です（空の場合は `validation_failed` を返します）
- 結果は関連性順（最も一致するものが最初）で返されます
- フィールドフィルター（`method:GET`、`tag:pet`）を使用して結果を絞り込みます
- 完全一致パス検索には引用符を使用します：`path:"/v1/forecast"`

---

## inspect

### 目的

エンドポイントの完全な OpenAPI 操作オブジェクト（すべてのパラメータ、リクエストボディスキーマ、レスポンススキーマ、ベース URL、完全な URL）を取得します。これは `invoke` の**前に**呼び出してエンドポイントの契約を理解するためのツールです。

### 使用するタイミング

- 常に `invoke` の前 — 正しい呼び出しを行うには完全な契約が必要です
- API の技術的詳細をユーザーに説明する必要がある場合
- 必須パラメータ、リクエストボディ構造、またはレスポンス形式を知る必要がある場合

### 動作方法

インデックス内のエンドポイントを検索し、解決されたすべてのスキーマを含む完全な OpenAPI 操作オブジェクトを返します。

### パラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|----------|-------------|
| `endpointId` | string | はい | エンドポイントの 32 文字 MD5 ハッシュ |

### レスポンス

```json
{
  "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
  "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
  "collectionId": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
  "specId": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
  "specDomain": "meteo",
  "method": "POST",
  "path": "/pet",
  "baseUrl": "https://meteo.swagger.io/v2",
  "fullUrl": "https://meteo.swagger.io/v2/pet",
  "operation": {
    "id": "addPet",
    "tags": ["pet"],
    "summary": "Add a new pet",
    "description": "Add a new pet to the store",
    "deprecated": false,
    "parameters": [
      {
        "name": "petId",
        "in": "path",
        "description": "ID of the pet",
        "required": true,
        "schema": {
          "type": "integer",
          "format": "int64"
        }
      }
    ],
    "requestBody": {
      "description": "Pet object to add",
      "required": true,
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "name": { "type": "string" },
              "status": { "type": "string", "enum": ["available", "pending", "sold"] }
            },
            "required": ["name"]
          }
        }
      }
    },
    "responses": {
      "200": {
        "description": "Successful operation",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/Pet"
            }
          }
        }
      },
      "405": {
        "description": "Invalid input"
      }
    }
  }
}
```

| フィールド | 型 | 説明 |
|-------|------|-------------|
| `baseUrl` | string | API のベース URL（設定から） |
| `fullUrl` | string | エンドポイントの完全な URL（ベース + パス） |
| `operation.parameters[]` | array | 名前、位置（path/query/header/cookie）、説明、必須フラグ、スキーマを持つパラメータ |
| `operation.requestBody` | object | コンテンツタイプとスキーマを持つリクエストボディ |
| `operation.responses` | map | 説明とスキーマを持つレスポンスコード |
| `operation.deprecated` | bool | エンドポイントが非推奨かどうか |

### 補足

- エンドポイントが存在しない場合は `not_found` を返します
- これは完全な OpenAPI 操作を返す**唯一の**ツールです。`endpoint_by_id` は概要のみを返します
- 必須パラメータとボディ構造を理解するために、`invoke` の前に常に `inspect` を呼び出してください
- `operation` オブジェクトには、完全なスキーマ定義に解決された `$ref` 参照が含まれます

---

## invoke

### 目的

エンドポイントに対して実際の API 呼び出しを実行します。これは実際の HTTP リクエストを行う唯一のツールです。認証は自動的に適用されるため、事前に `auth` を呼び出す必要はありません。

### 使用するタイミング

- `inspect` を呼び出してエンドポイントの契約を理解した後でのみ
- 破壊的操作（POST、PUT、PATCH、DELETE）の場合は明示的なユーザー確認がある場合のみ
- ユーザーが API の呼び出しを要求し、必要なパラメータがすべて揃っている場合

### 動作方法

1. インデックス内のエンドポイントを検索
2. パスパラメータを URL に代入
3. クエリパラメータを追加
4. ヘッダーとクッキーを追加
5. リクエストボディを JSON としてシリアライズ
6. 認証（トークン、ヘッダー、クエリパラメータ）を自動的に取得して適用
7. HTTP リクエストを実行
8. レスポンスを返すか、大きすぎる場合はファイルに保存

### パラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|----------|-------------|
| `endpointId` | string | はい | エンドポイントの 32 文字 MD5 ハッシュ |
| `parameters` | object | いいえ | パス、クエリ、ヘッダーパラメータのキーと値のペア |
| `requestBody` | object | いいえ | POST/PUT/PATCH リクエストのリクエストボディ |
| `headers` | object | いいえ | 送信する追加の HTTP ヘッダー |
| `cookies` | object | いいえ | 送信する追加の HTTP クッキー |

### レスポンス（インライン）

```json
{
  "statusCode": 200,
  "headers": {
    "content-type": "application/json"
  },
  "body": {
    "id": 1,
    "name": "Rex",
    "status": "available"
  }
}
```

### レスポンス（ファイル参照 — ボディがサイズ制限を超えた場合）

```json
{
  "statusCode": 200,
  "headers": {
    "content-type": "application/json"
  },
  "fileRef": {
    "path": "/Users/user/.swag2mcp/responses/response_a1b2c3d4.json",
    "size": 1572864,
    "sizeHint": "1.5 MB",
    "maxSizeHint": "2 KB",
    "message": "Response exceeds the 2 KB limit and has been saved to disk.",
    "openCmd": "open /Users/user/.swag2mcp/responses/response_a1b2c3d4.json"
  }
}
```

| フィールド | 型 | 説明 |
|-------|------|-------------|
| `statusCode` | int | HTTP レスポンスステータスコード |
| `headers` | object | HTTP レスポンスヘッダー |
| `body` | any | レスポンスボディ（サイズ制限内の場合に存在） |
| `fileRef` | object | ファイル参照（ボディがサイズ制限を超えた場合に存在） |

### 大規模レスポンスの操作

`invoke` が `fileRef` を返した場合、レスポンスツールを使用してデータを探索します：

1. **`response_outline(path)`** — 構造概要（キー、型、配列長）を取得
2. **`response_compress(path, mode)`** — データをインラインに収まるよう圧縮
3. **`response_slice(path, jsonPath)`** — 特定の断片を抽出

### 補足

- **認証は自動：** `invoke` ツールはスペックの認証設定から自動的に認証を取得して適用します。事前に `auth` を呼び出す必要は**ありません**。
- **レート制限：** 各エンドポイントには 10 秒のクールダウンがあります。同じエンドポイントへの 10 秒以内の 2 回目の呼び出しは静かにブロックされます（`rate_limit` エラーを返します）。
- **レスポンスサイズ制限：** デフォルトは 2 KB（`max_response_size` で設定可能）。レスポンスがこの制限を超えると、`{workspace}/responses/` に保存され、インラインの `body` の代わりに `FileReference` が返されます。
- **パラメータ処理：** パスパラメータは URL に代入されます。クエリパラメータは追加されます。リクエストからのパラメータは操作スペックのデフォルトを上書きします。
- **リクエストボディ：** POST/PUT/PATCH の場合、ボディは JSON としてシリアライズされます。`Content-Type` は自動的に `application/json` に設定されます。
- **エラーハンドリング：** HTTP エラー（2xx 以外）は、ステータスコードとレスポンスボディを含む `invoke_error` として返されます。
- **破壊的操作：** 明示的なユーザー確認なしに POST/PUT/PATCH/DELETE を呼び出さないでください。
