# レスポンスサイズ管理

## 概要

API レスポンスは非常に大きくなる可能性があります — LLM のコンテキストウィンドウに収まらないこともあります。swag2mcp は、大きすぎるレスポンスをディスクに保存し、探索するためのツールを提供することで、自動的にレスポンスサイズを管理します。

## 仕組み

1. **`invoke` を呼び出す** — swag2mcp が API リクエストを実行
2. **レスポンスが小さい場合**（制限内）— インラインで LLM に返される
3. **レスポンスが大きすぎる場合**（制限超過）— `{workspace}/responses/` に JSON ファイルとして保存。LLM は完全なレスポンスの代わりにファイル参照を受け取る

### 例：小さいレスポンス（インライン）

```json
{
  "statusCode": 200,
  "body": {
    "id": 1,
    "name": "Rex",
    "status": "available"
  }
}
```

### 例：大きいレスポンス（ファイル参照）

```json
{
  "statusCode": 200,
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

## 設定

```yaml
http_client:
  max_response_size: 1048576  # 1 MB（バイト単位）
```

### max_response_size

- **型:** `int`（バイト）
- **デフォルト:** `1048576`（1 MB）
- **範囲:** 256 〜 10,485,760 バイト（10 MB）
- **効果:** このサイズを超えるレスポンスはインラインで返されず、ディスクに保存されます
- **増やすタイミング:** 大規模なデータセットを返す API（レポート、ログ、分析）
- **減らすタイミング:** LLM のコンテキストウィンドウが限られている場合、またはファイルベースのアクセスを希望する場合

## 大きなレスポンスの操作

`invoke` が `fileRef` を返した場合、次の 3 つのツールを使用してデータを探索します：

### 1. response_outline — 構造を理解

レスポンスの構造サマリーを取得：キー、型、配列の長さ、ナビゲーションヒント。

```json
→ response_outline(path: "/path/to/file.json")
← {
    "type": "object",
    "size": 1572864,
    "keys": ["data", "meta"],
    "itemCount": 500,
    "compressionHints": [
      "response_compress(path, 'first_of_array', 'data')",
      "response_compress(path, 'sample_array', 'data', arrayHead=3, arrayTail=2)"
    ]
  }
```

### 2. response_compress — 小さなバージョンを取得

データを圧縮してインラインに収めます。複数の圧縮モードから適切なトレードオフを選択できます。

| モード | 説明 | 最適な用途 |
|-------|------|-----------|
| `first_of_array` | 配列の最初の要素のみ保持 | すべての要素が同じ構造の場合 |
| `sample_array` | 配列の先頭（3）と末尾（2）を保持 | 値の範囲を確認したい場合 |
| `truncate_strings` | すべての文字列を N 文字に短縮 | 文字列が非常に長い場合 |
| `keys_only` | 値を型名に置換 | 構造のみが必要な場合 |
| `select_keys` | 指定されたキーのみ保持 | 特定のフィールドのみが必要な場合 |

```json
→ response_compress(path: "/path/to/file.json", mode: "first_of_array", jsonPath: "data")
← {
    "body": [{ "id": 1, "name": "Rex" }],
    "hint": "Compressed array from 500 to 1 item using first_of_array mode"
  }
```

### 3. response_slice — 特定の断片を抽出

JSON パスまたは行範囲で特定の要素または値を取得します。

```json
→ response_slice(path: "/path/to/file.json", jsonPath: "data.0")
← {
    "slice": {
      "value": { "id": 1, "name": "Rex" },
      "jsonPath": "data.0",
      "nextPath": "data.1",
      "prevPath": null
    }
  }
```

## 完全なワークフロー

```
1. invoke(endpoint) → fileRef（レスポンス 1.5 MB）
2. response_outline(path) → 構造: { data: Array(500) }
3. response_compress(path, mode: "first_of_array", jsonPath: "data") → 最初の項目
4. response_slice(path, jsonPath: "data.0") → 最初の項目の詳細
5. response_slice(path, jsonPath: "data.1") → 2 番目の項目
```

## 自動クリーンアップ

MCP サーバーが起動するとき（`swag2mcp mcp`）、48 時間以上経過したレスポンスファイルが自動的に削除されます。手動でクリーンアップすることもできます：

```bash
swag2mcp clean
```

## 重要な注意点

- **制限はバイト単位** — `1048576` = 1 MB、`2097152` = 2 MB など
- **ファイル参照には open コマンドが含まれます** — macOS では `open`、Linux では `xdg-open`
- **レスポンスファイルはランダムなサフィックスで命名されます** — 同時呼び出し間で競合しません
- **responses ディレクトリは自動的に作成されます** — 手動設定は不要です
