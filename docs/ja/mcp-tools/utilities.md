# ユーティリティツール

ユーティリティツールは、認証トークンの取得、ランタイム情報の取得、インラインに収まらない大規模な API レスポンスの操作など、補助的な機能を提供します。

---

## auth

### 目的

特定のスペックの認証トークン、ヘッダー、またはクエリパラメータを取得します。これにより、LLM は swag2mcp の外部で使用できる認証情報（例：curl コマンドの生成）にアクセスできます。

### 使用するタイミング

- ユーザーが明示的に生のトークンや認証情報を要求した場合のみ
- 認証が必要な curl コマンドやコードスニペットを生成する場合
- ユーザーが設定されている認証方法を確認したい場合

### 使用すべきでないタイミング

- `inspect` や `invoke` の前に `auth` を呼び出さ**ないでください**。`invoke` は自動的に認証を取得して適用します
- 認証が設定されているかどうかを確認するためだけに `auth` を呼び出さ**ないでください**。代わりに `info` を使用してください

### 動作方法

スペックの認証設定を検索し、認証フロー（トークン交換、スクリプト実行など）を実行して現在の認証情報を取得します。

### パラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|----------|-------------|
| `specId` | string | はい | スペックの 32 文字 MD5 ハッシュ |

### レスポンス

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "headers": {
    "Authorization": "Bearer eyJhbGciOiJIUzI1NiIs...",
    "X-API-Key": "my-api-key"
  },
  "queryParams": {
    "api_key": "my-api-key"
  }
}
```

| フィールド | 型 | 説明 |
|-------|------|-------------|
| `token` | string | 生のトークン値（ベアラートークン、API キーなど） |
| `headers` | object | リクエストに含める HTTP ヘッダー |
| `queryParams` | object | リクエストに含めるクエリパラメータ |

### 補足

- **デフォルトで本番環境では無効：** `--disable-llm-auth` フラグ（デフォルト：`true`）は、MCP ツールリストから `auth` ツールを完全に削除します。LLM はトークンを表示したり要求したりできません。デバッグや短期間のトークンには `--disable-llm-auth=false` を設定して有効にします。
- **`invoke` は認証を自動処理：** `invoke` の前に `auth` を呼び出す必要はありません。invoke サービスは自動的に正しい認証を取得して適用します。
- **9 つの認証方法をサポート：** `none`、`basic`、`bearer`、`digest`、`hmac`、`oauth2-cc`（クライアントクレデンシャル）、`oauth2-pwd`（パスワード）、`api-key`、`script`。
- 認証方法が失敗した場合（例：OAuth2 トークンエンドポイントに到達できない、スクリプト実行の失敗）、`auth_error` を返します。

---

## info

### 目的

swag2mcp ランタイムの包括的な概要（バージョン、ワークスペースパス、アクティブなスペック、HTTP クライアント設定、MCP トランスポート設定、認証方法、モックモードステータス）を返します。

### 使用するタイミング

- ユーザーがシステム設定について質問した場合
- ランタイム設定（タイムアウト、レスポンスサイズ制限、トランスポート）を確認する必要がある場合
- どの認証方法が利用可能かを知る必要がある場合
- 設定の問題をトラブルシューティングする場合

### 動作方法

ランタイム状態の事前計算されたスナップショットを返します。パラメータは不要です。

### パラメータ

なし。

### レスポンス

```json
{
  "version": "v1.2.0",
  "workspace": "~/.swag2mcp",
  "uptime": "2h 15m",
  "specs": {
    "total": 4,
    "active": 3,
    "disabled": 1,
    "collections": 6,
    "endpoints": 42
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 KB",
    "follow_redirects": true,
    "max_redirects": 10,
    "randomize": false,
    "proxy": null,
    "headers": {},
    "cookies": []
  },
  "mcp": {
    "transport": "stdio",
    "addr": ":8080",
    "path": "/mcp",
    "auth_enabled": false
  },
  "auth": {
    "methods": ["bearer", "api-key"]
  },
  "mock": {
    "enabled": false
  }
}
```

| フィールド | 型 | 説明 |
|-------|------|-------------|
| `version` | string | swag2mcp バージョン |
| `workspace` | string | ワークスペースディレクトリパス |
| `uptime` | string | サーバー稼働時間（人間可読） |
| `specs` | object | スペック概要：total、active、disabled、collections、endpoints |
| `http_client` | object | HTTP クライアント設定 |
| `http_client.max_response_size` | string | 人間可読形式の最大レスポンスサイズ（例："2 KB"） |
| `mcp` | object | MCP サーバー設定 |
| `auth` | object | 利用可能な認証方法 |
| `mock` | object | モックサーバーステータス |

### 補足

- `max_response_size` は人間可読形式で表示されます（例：`"1 KB"`、`"2 MB"`）
- `uptime` はサーバー起動時間から計算されます
- データはブートストラップ時に取得されたスナップショットであり、MCP サーバー起動時の状態を反映します

---

## response_outline

### 目的

`invoke` によってディスクに保存された大規模な JSON レスポンスファイルの高レベルの構造概要を取得します。実際の値を返さずに、データの形状（キー、型、配列長、ナビゲーションヒント）を返します。

### 使用するタイミング

- `invoke` が `fileRef` を返した直後（レスポンスがインラインに大きすぎる場合）
- これは大規模レスポンスワークフローの**必須の最初のステップ**です

### 動作方法

保存されたレスポンスファイルを読み取り、その構造を分析します：トップレベルの型、キー、配列長、ネスト深度、圧縮ヒント。

### パラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|----------|-------------|
| `path` | string | はい | `fileRef.path` からの絶対パス |
| `maxDepth` | int | いいえ | 最大再帰深度（デフォルト：3） |
| `maxArrayItems` | int | いいえ | 検査する配列アイテム数（デフォルト：5） |

### レスポンス

```json
{
  "outline": {
    "type": "object",
    "size": 1572864,
    "lineCount": 12500,
    "depth": 3,
    "structure": {
      "type": "object",
      "keys": ["data", "meta", "error"],
      "data": {
        "type": "array",
        "length": 500,
        "items": {
          "type": "object",
          "keys": ["id", "name", "status", "createdAt"]
        }
      }
    },
    "schemaHint": "object with 3 keys: data (array[500]), meta (object), error (null)",
    "keys": ["data", "meta", "error"],
    "itemCount": 500,
    "itemType": "object",
    "compressionHints": [
      "response_compress(path, 'first_of_array', 'data')",
      "response_compress(path, 'sample_array', 'data', arrayHead=3, arrayTail=2)",
      "response_compress(path, 'keys_only', 'data')",
      "response_compress(path, 'select_keys', 'data', selectKeys=[id, name])"
    ],
    "navigationHints": {
      "paths": ["data", "meta", "error"],
      "arrays": [
        {"path": "data", "length": 500}
      ]
    }
  }
}
```

| フィールド | 型 | 説明 |
|-------|------|-------------|
| `type` | string | トップレベルの型："object" または "array" |
| `size` | int | ファイルサイズ（バイト） |
| `lineCount` | int | ファイルの行数 |
| `depth` | int | 検査された最大ネスト深度 |
| `structure` | object | キー、型、配列長を持つ再帰的構造 |
| `schemaHint` | string | トップレベルの形状の 1 行サマリー |
| `keys` | array | トップレベルのキー（オブジェクトの場合） |
| `itemCount` | int | 配列の長さ（配列の場合） |
| `compressionHints` | array | パラメータ付きの推奨 `response_compress` 呼び出し |
| `navigationHints` | object | トップレベルのパスと長さ付きの配列 |

### 補足

- パスが無効またはレスポンスディレクトリ内にない場合は `validation_failed` を返します
- ファイルが存在しない場合は `not_found` を返します
- ファイルが有効な JSON でない場合は `validation_failed` を返します
- `compressionHints` フィールドは、`response_compress` 呼び出しのすぐに使用できる提案を提供します

---

## response_compress

### 目的

保存されたレスポンスファイル内の JSON 値を削減してレスポンスサイズ制限内に収め、LLM にインラインで返せるようにします。複数の圧縮モードから、サイズと情報の適切なトレードオフを選択できます。

### 使用するタイミング

- `response_outline` の後に構造を理解するため
- 大規模レスポンスからインラインでデータを取得する必要がある場合
- `response_slice` が狭すぎて、より広いビューが必要な場合

### 動作方法

保存されたレスポンスファイルを読み取り、指定された JSON パスに移動し、圧縮モードを適用して、圧縮結果を返します。結果がまだサイズ制限を超えている場合は、新しいファイルに保存されます。

### パラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|----------|-------------|
| `path` | string | はい | `fileRef.path` からの絶対パス |
| `jsonPath` | string | いいえ | 圧縮する値へのパス（例：`data` または `data.0`） |
| `mode` | string | はい | 圧縮モード（以下の表を参照） |
| `arrayHead` | int | いいえ | `sample_array` モードで保持する先頭アイテム数（デフォルト：3） |
| `arrayTail` | int | いいえ | `sample_array` モードで保持する末尾アイテム数（デフォルト：2） |
| `stringLen` | int | いいえ | `truncate_strings` モードの最大文字列長（デフォルト：80） |
| `selectKeys` | array | いいえ | `select_keys` モードで保持するキー |

### 圧縮モード

| モード | 説明 | 最適な用途 |
|------|-------------|----------|
| `first_of_array` | 配列の最初の要素のみを保持 | すべての要素が同じ構造を持つ場合 |
| `sample_array` | 配列の先頭と末尾を保持 | 値の範囲を確認する必要がある場合 |
| `truncate_strings` | すべての文字列を `stringLen` 文字に短縮 | 文字列が非常に長いが構造が重要な場合 |
| `keys_only` | オブジェクトの値を型名に置換 | 構造のみが必要な場合 |
| `select_keys` | すべてのオブジェクトで指定されたキーのみを保持 | 多くのオブジェクトから特定のフィールドが必要な場合 |

### レスポンス

```json
{
  "body": [
    { "id": 1, "name": "Rex", "status": "available" },
    { "id": 2, "name": "Max", "status": "pending" }
  ],
  "hint": "Compressed array from 500 to 2 items using first_of_array mode"
}
```

| フィールド | 型 | 説明 |
|-------|------|-------------|
| `body` | any | 圧縮された JSON 値（サイズ制限内の場合に存在） |
| `fileRef` | object | ファイル参照（まだ大きすぎる場合に存在） |
| `hint` | string | 何が圧縮されたかの説明 |

### 補足

- 圧縮結果がまだ `max_response_size` を超える場合は、新しいファイルに保存され、`FileReference` が返されます
- デフォルト値：`arrayHead=3`、`arrayTail=2`、`stringLen=80`
- 無効なパス、無効な JSONPath、または JSON 以外のファイルの場合は `validation_failed` を返します
- ファイルが存在しないか JSONPath が一致しない場合は `not_found` を返します

---

## response_slice

### 目的

保存された JSON レスポンスファイルの特定の断片を、論理的な JSON パスまたは行範囲で抽出します。`response_compress` とは異なり、未加工の変更されていないデータを提供します。

### 使用するタイミング

- 大規模レスポンスから特定の要素や値を取得する必要がある場合
- `response_compress` で十分な詳細が得られない場合
- レスポンスを段階的にナビゲートしたい場合

### 動作方法

保存されたレスポンスファイルを読み取り、JSON パス（例：`data.3.name`）または行範囲（例：`120-240`）で断片を抽出します。配列やオブジェクトをステップ実行するためのナビゲーションヒントを返します。

### パラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|----------|-------------|
| `path` | string | はい | `fileRef.path` からの絶対パス |
| `jsonPath` | string | いいえ | 値への論理パス（例：`data.3.name`） |
| `line` | int | いいえ | 断片を中心とする 1 ベースの行番号 |
| `range` | string | いいえ | `start-end` 形式の行範囲（例：`120-240`） |
| `around` | int | いいえ | `line` の前後に含める行数（デフォルト：20） |

### レスポンス

```json
{
  "slice": {
    "lines": [120, 130],
    "fragment": "{\n  \"id\": 1,\n  \"name\": \"Rex\"\n}",
    "value": {
      "id": 1,
      "name": "Rex"
    },
    "jsonPath": "data.0",
    "context": "object",
    "isComplete": true,
    "nextLine": 131,
    "prevLine": 119,
    "nextPath": "data.1",
    "prevPath": null
  }
}
```

| フィールド | 型 | 説明 |
|-------|------|-------------|
| `lines` | array | 1 ベースの行範囲 [start, end] |
| `fragment` | string | 生の JSON テキスト（十分に小さい場合） |
| `value` | any | 抽出された JSON 値 |
| `jsonPath` | string | 使用された JSON パス |
| `context` | string | "object"、"array"、または "value" |
| `isComplete` | bool | 値が有効な JSON 断片の場合は true |
| `nextLine` | int | 行ベースのナビゲーションのための次の推奨行 |
| `prevLine` | int | 推奨される前の行 |
| `nextPath` | string | 配列ナビゲーションのための次の推奨 JSON パス |
| `prevPath` | string | 推奨される前の JSON パス |

### 補足

- **行番号よりも `jsonPath` を優先**してください。JSON パスは安定していて説明的ですが、行番号はファイルが再生成されると変わります
- 抽出された断片が `max_response_size` を超える場合は、新しいファイルに保存され、`FileReference` が返されます
- デフォルトの `around` は 20 行です
- レスポンスには、配列をステップ実行するための `nextPath`/`prevPath` と、行ベースのナビゲーションのための `nextLine`/`prevLine` が含まれます
- 無効なパス、無効な JSONPath、無効な行/範囲、または JSON 以外のファイルの場合は `validation_failed` を返します
- ファイルが存在しないか JSONPath が一致しない場合は `not_found` を返します
