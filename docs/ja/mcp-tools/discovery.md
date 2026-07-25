# 発見ツール

発見ツールを使用すると、LLM がスペック階層をナビゲートできます。すべてのスペックを検索し、スペックにドリルダウンしてコレクションを表示し、コレクション内のタグを探索できます。まず `spec_list` で利用可能な API を確認し、ID を使用してさらに深く掘り下げます。

---

## spec_list

### 目的

ワークスペースに登録されているすべての API スペックを一覧表示します。これはセッションの開始点であり、LLM は最初にこれを呼び出して利用可能な API を発見します。

### 使用するタイミング

- セッションの開始時に設定されている API を確認する
- スペックの追加または削除後にリストを更新する
- 他のツールで使用するスペック ID が必要な場合

### 動作方法

すべてのスペックを一意の ID とドメイン名とともに返します。パラメータは不要です。

### パラメータ

なし。

### レスポンス

```json
{
  "specs": [
    {
      "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
      "domain": "meteo"
    },
    {
      "id": "b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7",
      "domain": "dadjoke"
    }
  ]
}
```

| フィールド | 型 | 説明 |
|-------|------|-------------|
| `id` | string | 32 文字の MD5 ハッシュ、スペックの一意識別子 |
| `domain` | string | スペックのドメイン名（例："meteo"、"dadjoke"） |

### 補足

- `id` と `domain` のみを返します。完全な詳細（コレクション、タグ）は `spec_by_id` を使用してください
- すべての ID は 32 文字の MD5 16 進文字列です（`^[0-9a-f]{32}$`）
- スペックが設定されていない場合は空の配列を返します

---

## spec_by_id

### 目的

特定のスペックの詳細情報（ドメイン、すべてのコレクション、その統計情報（タグ数、メソッド数））を取得します。

### 使用するタイミング

- `spec_list` の後にスペック内のコレクションを確認する
- さらにナビゲートするためにコレクション ID が必要な場合

### 動作方法

スペック ID を受け取り、スペックのメタデータとすべてのコレクションをカウント付きで返します。

### パラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|----------|-------------|
| `id` | string | はい | スペックの 32 文字 MD5 ハッシュ |

### レスポンス

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collections": [
    {
      "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "title": "Weather Forecast",
      "llmTitle": "Forecast API",
      "countTags": 3,
      "countMethods": 12
    }
  ]
}
```

| フィールド | 型 | 説明 |
|-------|------|-------------|
| `spec.id` | string | スペック識別子 |
| `spec.domain` | string | スペックドメイン名 |
| `collections[].id` | string | コレクション識別子 |
| `collections[].title` | string | 人間可読なタイトル |
| `collections[].llmTitle` | string | LLM 向けタイトル（オプション） |
| `collections[].countTags` | int | コレクション内のタグ数 |
| `collections[].countMethods` | int | コレクション内の HTTP メソッド数 |

### 補足

- スペック ID が存在しない場合は `not_found` エラーを返します
- `id` は有効な 32 文字の MD5 16 進文字列である必要があります

---

## collection_by_spec

### 目的

特定のスペック内のすべてのコレクションを一覧表示します。`spec_by_id` と似ていますが、余分なスペックメタデータなしでコレクションリストのみを返します。

### 使用するタイミング

- すでにスペック ID を持っており、コレクションリストだけが必要な場合
- `spec_by_id` の軽量な代替として

### パラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|----------|-------------|
| `specId` | string | はい | スペックの 32 文字 MD5 ハッシュ |

### レスポンス

```json
{
  "spec": {
    "id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
    "domain": "meteo"
  },
  "collections": [
    {
      "id": "c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6",
      "title": "Weather Forecast",
      "llmTitle": "Forecast API",
      "countTags": 3,
      "countMethods": 12
    }
  ]
}
```

### 補足

- スペックが存在しない場合は `not_found` を返します
- `spec_by_id` と同じデータですが、余分なスペックラッパーはありません

---

## collection_by_id

### 目的

特定のコレクションの詳細情報（メタデータ、親スペック、コレクション内のすべてのタグ）を取得します。

### 使用するタイミング

- `collection_by_spec` の後にコレクション内のタグを確認する
- `tag_by_id` や `endpoint_by_tag` で使用するタグ ID が必要な場合

### パラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|----------|-------------|
| `id` | string | はい | コレクションの 32 文字 MD5 ハッシュ |

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
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    },
    {
      "id": "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7",
      "title": "current",
      "countMethods": 7
    }
  ]
}
```

| フィールド | 型 | 説明 |
|-------|------|-------------|
| `spec` | object | 親スペック（id、domain） |
| `collection` | object | コレクションメタデータ（id、title、countMethods） |
| `tags[]` | array | id、title、countMethods を持つタグのリスト |

### 補足

- コレクション ID が存在しない場合は `not_found` を返します
- タグは ID 付きで返されます。実際のエンドポイントを表示するには `endpoint_by_tag(tagId)` を使用してください

---

## tag_by_spec

### 目的

すべてのコレクションにわたって、スペック全体のすべてのタグを一覧表示します。利用可能なすべてのタグの俯瞰ビューを取得するのに便利です。

### 使用するタイミング

- 各コレクションにドリルダウンせずにスペック内のすべてのタグを表示したい場合
- 必要なタグがどのコレクションに含まれているかわからない場合

### パラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|----------|-------------|
| `specId` | string | はい | スペックの 32 文字 MD5 ハッシュ |

### レスポンス

```json
{
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    },
    {
      "id": "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7",
      "title": "current",
      "countMethods": 7
    }
  ]
}
```

### 補足

- スペックが存在しない場合は `not_found` を返します
- タグはスペック内のすべてのコレクションから集約されます

---

## tag_by_collection

### 目的

特定のコレクション内のすべてのタグを一覧表示します。`tag_by_spec` とは異なり、親スペックとコレクションのメタデータも返します。

### 使用するタイミング

- `collection_by_id` の後にタグリストを確認する
- 完全なコンテキスト（スペック + コレクション + タグ）が必要な場合

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
  "tags": [
    {
      "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
      "title": "forecast",
      "countMethods": 5
    }
  ]
}
```

### 補足

- コレクションが存在しない場合は `not_found` を返します
- `tag_by_spec` と同じタグデータですが、1 つのコレクションにスコープされています

---

## tag_by_id

### 目的

単一のタグに関する情報（ID、タイトル、含まれるメソッド数）を取得します。これはタグ自体に関する情報です。実際のエンドポイントを表示するには `endpoint_by_tag` を使用してください。

### 使用するタイミング

- タグ ID を持っており、その名前とサイズを確認したい場合
- `endpoint_by_tag` を呼び出す前に、期待されるエンドポイント数を把握するため

### パラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|----------|-------------|
| `id` | string | はい | タグの 32 文字 MD5 ハッシュ |

### レスポンス

```json
{
  "tag": {
    "id": "d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
    "title": "forecast",
    "countMethods": 5
  }
}
```

| フィールド | 型 | 説明 |
|-------|------|-------------|
| `tag.id` | string | タグ識別子 |
| `tag.title` | string | 人間可読なタグ名 |
| `tag.countMethods` | int | このタグ内の HTTP メソッド数 |

### 補足

- タグが存在しない場合は `not_found` を返します
- このツールはタグのメタデータのみを返します。実際のエンドポイントリストを取得するには `endpoint_by_tag` を使用してください
