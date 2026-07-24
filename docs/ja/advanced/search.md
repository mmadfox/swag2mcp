# 全文検索

## 概要

swag2mcp には組み込みの全文検索エンジン（bluge）が含まれており、すべての spec の全エンドポイントをインデックス化します。LLM はエンドポイント ID を知らなくても、メソッド、パス、サマリー、タグでエンドポイントを検索できます。

## インデックス化の仕組み

spec が追加または更新されると、すべてのエンドポイントがインデックス化されます。以下のフィールドが検索可能です：

| フィールド | 説明 | 例 |
|-----------|------|-----|
| `method` | HTTP メソッド | `GET`、`POST`、`PUT` |
| `path` | API エンドポイントパス | `/api/v1/users/{id}` |
| `summary` | OpenAPI サマリー | "Find pet by ID" |
| `tag` | エンドポイントカテゴリ | "pets"、"users" |
| `_all` | 全フィールドの組み合わせ | method + path + tag + summary |

インデックスは MCP サーバー起動時に毎回再構築されます。高速検索のためにメモリ内に保存されます。

## クエリ構文

検索は正確なフィルタリングのための豊富なクエリ構文をサポートしています：

| 例 | 説明 |
|----|------|
| `pet` | 全フィールドのシンプルなテキスト検索 |
| `method:GET` | すべての GET エンドポイントを検索 |
| `tag:pets` | "pets" タグ内のエンドポイントを検索 |
| `path:"/api/v1/users"` | 正確なパス一致 |
| `+method:POST +tag:pet` | 両方の条件に一致する必要あり |
| `-method:DELETE` | DELETE メソッドを除外 |
| `create~` | あいまい検索（タイポ許容） |
| `cr*` | ワイルドカード検索 |
| `"find pet"` | フレーズ検索 |
| `+summary:pet -method:DELETE` | サマリーに "pet" を含み、DELETE を除外 |

### フィールド固有の検索

`field:value` 構文を使用して特定のフィールド内を検索できます：

```
method:GET
tag:pets
path:"/pet/findByStatus"
summary:"find pet by status"
```

### ブール演算子

- `+` — 用語が一致する必要あり（AND）
- `-` — 用語が一致してはいけない（NOT）
- 用語間のスペース — OR（いずれかの用語が一致）

### あいまい検索とワイルドカード

- `term~` — あいまい検索（類似単語に一致、タイポ対応）
- `te*` — ワイルドカード（任意の文字に一致）
- `te?t` — 1 文字のワイルドカード

## 例

```
# すべての GET リクエストを検索
method:GET

# pet タグ内の POST リクエストを検索
+method:POST +tag:pet

# 正確なパスでエンドポイントを検索
path:"/pet/findByStatus"

# 説明で検索
"find pet by status"

# DELETE 以外をすべて検索
+summary:pet -method:DELETE

# "create" のあいまい検索（タイポ対応）
create~
```

## MCP ツール

`search` MCP ツールは検索エンジンを LLM に公開します：

```
→ search(query: "find pet by status", limit: 5)
← GET /pet/findByStatus — Finds Pets by status
   GET /pet/{petId} — Find pet by ID
```

### パラメーター

| パラメーター | 必須 | 説明 |
|------------|------|------|
| `query` | はい | 検索クエリ（構造化構文をサポート） |
| `limit` | はい | 最大結果数（1〜50） |

## 重要な注意点

- **インデックスはメモリ内** — MCP サーバーが起動するたびに再構築されます。永続的なインデックスファイルはありません。
- **すべてのフィールドは小文字化** — 検索は大文字小文字を区別しません
- **制限は最大 50** — 50 を超える結果を要求することはできません
- **無効なクエリ構文** は例を含む役立つエラーメッセージを返します
- **`_all` フィールド** はシンプルなテキスト検索のために method、path、tag、summary を結合します
