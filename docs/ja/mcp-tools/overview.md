# MCP ツール

## 概要

swag2mcp は **19 の MCP ツール** を提供し、LLM エージェントが Model Context Protocol を通じて API に完全にアクセスできるようにします。これらのツールは、利用可能な API の発見、スペック階層のナビゲーション、エンドポイントの検索と検査、API 呼び出しの実行、大規模レスポンスの処理まで、完全なワークフローをカバーします。

### ツールが解決する課題

- **発見** — LLM は ID を事前に知らなくてもスペック、コレクション、タグを見つけられます
- **ナビゲーション** — スペック → コレクション → タグ → エンドポイントの構造化された階層を掘り下げます
- **検索** — ID がない場合に全エンドポイントを全文検索
- **検査** — 呼び出し前に完全な OpenAPI 操作オブジェクトを取得
- **実行** — 自動認証付きで実際の API 呼び出しを実行
- **大規模レスポンス処理** — インラインに収まらない大きなレスポンスの概要表示、圧縮、スライス

### 読み取り専用 vs 可変

| タイプ | 数 | ツール |
|------|-------|-------|
| **読み取り専用** | 17 | すべての発見、エンドポイント、検索、検査、情報、レスポンスツール |
| **可変** | 2 | `invoke`（実際の HTTP 呼び出し）、`auth`（トークン取得） |

読み取り専用ツールは MCP プロトコルで `ReadOnlyHint=true` および `IdempotentHint=true` とマークされ、LLM に副作用なく安全に呼び出せることを示します。

### エラーハンドリング

すべてのツールは、機械可読なコードと人間可読なメッセージを含む構造化された `LLMError` オブジェクトとしてエラーを返します：

| エラーコード | 意味 |
|------------|---------|
| `validation_failed` | 無効な入力（不正な ID 形式、必須フィールドの欠落） |
| `not_found` | インデックスまたはワークスペースにエンティティが見つからない |
| `rate_limit` | 同じエンドポイントへの 10 秒以内の 2 回目の `invoke` 呼び出し |
| `invoke_error` | HTTP 呼び出しの失敗、ダウンロードの失敗 |
| `auth_error` | 認証トークン取得の失敗 |
| `config_error` | 設定ファイルの読み込みまたは保存の失敗 |
| `parse_error` | スペックファイルの解析の失敗 |

## カテゴリ

| カテゴリ | ツール | 説明 |
|----------|-------|-------------|
| **発見** | `spec_list`, `spec_by_id`, `collection_by_spec`, `collection_by_id`, `tag_by_spec`, `tag_by_collection`, `tag_by_id` | スペック階層をナビゲート：スペック、コレクション、タグを検索 |
| **エンドポイント** | `endpoint_by_spec`, `endpoint_by_collection`, `endpoint_by_tag`, `endpoint_by_id` | 階層の異なるレベルでエンドポイントを表示 |
| **実行** | `search`, `inspect`, `invoke` | 検索、完全な契約の検査、API の呼び出し |
| **ユーティリティ** | `auth`, `info`, `response_outline`, `response_compress`, `response_slice` | 認証トークン、ランタイム情報、大規模レスポンス処理 |
| **スキル** | [フォーマットガイド](/mcp-tools/skills) | ツールレスポンスの表示方法をカスタマイズ |

## 全リスト

| ツール | 説明 |
|------|-------------|
| `spec_list` | ワークスペース内のすべての API スペックを一覧表示 |
| `spec_by_id` | コレクションを含む詳細なスペック情報を取得 |
| `collection_by_spec` | スペック内のコレクションを一覧表示 |
| `collection_by_id` | タグを含むコレクションの詳細を取得 |
| `tag_by_spec` | スペック全体のすべてのタグを一覧表示 |
| `tag_by_collection` | コレクション内のタグを一覧表示 |
| `tag_by_id` | タグの詳細（ID、タイトル、メソッド数）を取得 |
| `endpoint_by_spec` | スペック内のすべてのエンドポイントを一覧表示 |
| `endpoint_by_collection` | コレクション内のエンドポイントを一覧表示 |
| `endpoint_by_tag` | タグ内のエンドポイントを一覧表示 |
| `endpoint_by_id` | クイックエンドポイント概要（メソッド、パス、概要） |
| `search` | 全エンドポイントの全文検索 |
| `inspect` | 完全な OpenAPI 操作の詳細（パラメータ、スキーマ） |
| `invoke` | 実際の API 呼び出しを実行 |
| `auth` | スペックの認証トークンまたはヘッダーを取得 |
| `info` | ランタイム情報（バージョン、スペック、設定） |
| `response_outline` | 大規模レスポンスファイルの構造概要 |
| `response_compress` | 大規模レスポンスをインラインに収まるよう圧縮 |
| `response_slice` | 大規模レスポンスの断片を抽出 |

## ナビゲーション階層

```
spec_list
  └── spec_by_id(id)
        └── collection_by_spec(specId)
              └── collection_by_id(id)
                    └── tag_by_collection(collectionId)
                          └── tag_by_id(id)
                                └── endpoint_by_tag(tagId)
                                      └── endpoint_by_id(id)
                                            └── inspect(endpointId)
                                                  └── invoke(endpointId)
```

ID がない場合は `search` を使用してクエリでエンドポイントを検索します。`invoke` が `fileRef` を返した場合（レスポンスが大きすぎる場合）、`response_outline` → `response_compress` または `response_slice` を使用してデータを探索します。
