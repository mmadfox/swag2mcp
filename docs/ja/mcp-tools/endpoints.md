# エンドポイントツール

エンドポイントツールを使用すると、LLM は階層の異なるレベルで API エンドポイントを表示できます。スペック内のすべてのエンドポイント、コレクション内、タグ内、または単一のエンドポイント概要を表示できます。これらを使用して、検査や呼び出しの前に利用可能な操作を発見します。

---

## endpoint_by_spec

### 目的

すべてのコレクションとタグにわたって、スペック全体のすべてのエンドポイントを一覧表示します。最も包括的なビューを返します。スペック内のすべてのエンドポイントを完全なコンテキスト（タグ、コレクション、スペック）とともに表示します。

### 使用するタイミング

- スペック内のすべてのエンドポイントを表示したい場合
- 必要なエンドポイントがどのコレクションやタグに含まれているかわからない場合
- `spec_by_id` の後に完全なエンドポイントリストを取得する場合

### パラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|----------|-------------|
| `specId` | string | はい | スペックの 32 文字 MD5 ハッシュ |

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

| フィールド | 型 | 説明 |
|-------|------|-------------|
| `id` | string | エンドポイント識別子 |
| `tagId` | string | 親タグ識別子 |
| `tagName` | string | 人間可読なタグ名 |
| `collectionId` | string | 親コレクション識別子 |
| `collectionTitle` | string | 人間可読なコレクションタイトル |
| `specId` | string | 親スペック識別子 |
| `specDomain` | string | スペックドメイン名 |
| `method` | string | HTTP メソッド（GET、POST、PUT、DELETE など） |
| `path` | string | API パス（例：/v1/forecast） |
| `summary` | string | エンドポイントの機能の人間可読な概要 |

### 補足

- スペックが存在しない場合は `not_found` を返します
- 各エンドポイントにはコンテキストのための完全な祖先（スペック → コレクション → タグ）が含まれます
- 単一エンドポイントのクイック概要には `endpoint_by_id` を使用してください

---

## endpoint_by_collection

### 目的

タグに関係なく、特定のコレクション内のすべてのエンドポイントを一覧表示します。スペックとコレクションのメタデータとともに、コレクションごとにグループ化されたエンドポイントを返します。

### 使用するタイミング

- `collection_by_id` の後にコレクション内のすべてのエンドポイントを表示する
- コレクションの完全な API サーフェスを探索したい場合

### パラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|----------|-------------|
| `collectionId` | string | はい | コレクションの 32 文字 MD5 ハッシュ |

### レスポンス

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "tagId": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "tagName": "forecast",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Get weather forecast for a location"
    }
  ]
}
```

### 補足

- コレクションが存在しない場合は `not_found` を返します
- コンテキストのためにスペックとコレクションのメタデータを含みます
- コレクション内のすべてのタグからのエンドポイントが一緒に返されます

---

## endpoint_by_tag

### 目的

特定のタグにグループ化されたすべてのエンドポイントを一覧表示します。これは最も焦点を絞ったビューです。1 つのコレクション内の 1 つのタグのエンドポイントを表示します。

### 使用するタイミング

- `tag_by_id` の後にタグ内の実際のエンドポイントを表示する
- タグがわかっていて、その操作を表示したい場合

### パラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|----------|-------------|
| `tagId` | string | はい | タグの 32 文字 MD5 ハッシュ |

### レスポンス

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  },
  "endpoints": [
    {
      "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "method": "GET",
      "path": "/v1/forecast",
      "summary": "Get weather forecast for a location"
    }
  ]
}
```

### 補足

- タグが存在しない場合は `not_found` を返します
- 完全なコンテキスト（スペック、コレクション、タグのメタデータ）を含みます
- エンドポイントは単一のコレクション内の単一のタグにスコープされます

---

## endpoint_by_id

### 目的

単一のエンドポイントのクイック概要（メソッド、パス、概要、非推奨ステータス）を取得します。これは軽量なツールです。完全な OpenAPI 操作オブジェクト（パラメータ、リクエストボディ、レスポンススキーマ）を取得するには `inspect` を使用してください。

### 使用するタイミング

- エンドポイント ID を持っており、その機能を簡単に確認したい場合
- 完全な詳細のために `inspect` を呼び出すかどうかを判断する前
- 呼び出し前にメソッドとパスを確認する必要がある場合

### パラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|----------|-------------|
| `id` | string | はい | エンドポイントの 32 文字 MD5 ハッシュ |

### レスポンス

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collection": {
    "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
    "title": "Weather Forecast",
    "countMethods": 12
  },
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  },
  "endpoint": {
    "id": "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6",
    "method": "GET",
    "path": "/v1/forecast",
    "summary": "Get weather forecast for a location"
  }
}
```

| フィールド | 型 | 説明 |
|-------|------|-------------|
| `endpoint.id` | string | エンドポイント識別子 |
| `endpoint.method` | string | HTTP メソッド |
| `endpoint.path` | string | API パス |
| `endpoint.summary` | string | 人間可読な概要 |

### 補足

- エンドポイントが存在しない場合は `not_found` を返します
- これは**クイック概要**です。パラメータ、リクエストボディ、レスポンススキーマは返しません
- 完全な技術的詳細（`invoke` の前に必要）については `inspect` を使用してください
