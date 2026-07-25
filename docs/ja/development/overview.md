# 開発概要

## このプロジェクトについて

swag2mcp は、OpenAPI/Swagger/Postman 仕様と LLM エージェントを Model Context Protocol（MCP）を介して橋渡しする Go プロジェクトです。Go 1.23+ で構築され、80 以上のリンターによって強制される厳格なコーディング規約に従っています。

このセクションは、コードベースを理解し、貢献し、新しい認証方法、MCP ツール、または統合機能で swag2mcp を拡張したい**エンジニア**向けに書かれています。

## 開発スキル

プロジェクトには、プロジェクトの規約とパターンをエンコードした 2 つの開発スキルが付属しています。これらを使用することも無視することもできます。これらはツールであり、ルールではありません。

### godeveloper

[godeveloper スキル](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/godeveloper/SKILL.md) は、プロジェクトのすべてのコード規約を定義しています：

- **命名** — パッケージ、ファイル、型、インターフェース、レシーバー、定数
- **フォーマット** — gofmt/gofumpt/goimports/gci、120 行制限、インポート順序
- **エラーハンドリング** — 8 つのエラーコードを持つ `LLMError`、センチネルエラー、エラーラッピング
- **インターフェース** — 小さなインターフェース、合成、コンシューマー側の定義
- **並行性** — ミューテックスの粒度、ゴルーチンのライフタイム、コンテキストの受け渡し
- **テスト** — テーブル駆動テスト、`newTestService()`/`seedTestData()` ヘルパー、モック生成
- **プロジェクトパターン** — サービス層、リクエスト/レスポンス構造体、関数型オプション、MCP ハンドラーパターン

### swag2mcp-cli

[swag2mcp-cli スキル](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md) は、すべての CLI コマンドを構文、フラグ、引数、例とともに文書化しています。CLI コマンドの作業やドキュメントの作成時に役立ちます。

## 主要なアーキテクチャ上の決定

### サービス層パターン

すべての機能は同じ 3 ステップのパターンに従います：

1. **検証**：`s.validateRequest(req)` でリクエストを検証（`go-playground/validator` を使用）
2. **検索**：インメモリインデックスからエンティティを検索（`not_found` コードの `LLMError` を返す）
3. **実行**：ビジネスロジックを実行し、型付きレスポンスまたは `LLMError` を返す

```go
func (s *Service) Search(ctx context.Context, req SearchRequest) (SearchResponse, error) {
    if err := s.validateRequest(req); err != nil {
        return SearchResponse{}, NewLLMError(validationFailedErrCode, err.Error())
    }
    results, err := s.index.Search(req.Query, req.Limit)
    if err != nil {
        return SearchResponse{}, NewLLMError(invokeErrorCode, err.Error())
    }
    return SearchResponse{Results: results}, nil
}
```

### リクエスト/レスポンス構造体

各メソッドには専用の `{Method}Request` および `{Method}Response` 構造体があります。リクエスト構造体は検証用の `validate` タグとドキュメント用の `jsonschema` タグを使用します：

```go
type SearchRequest struct {
    Query string `json:"query" validate:"required,min=1" jsonschema:"description=Search query supporting field filters"`
    Limit int    `json:"limit" validate:"required,min=1,max=50" jsonschema:"description=Maximum results"`
}

type SearchResponse struct {
    Results []EndpointSearchItem `json:"results"`
}
```

### 関数型オプション

設定は関数型オプションパターンを使用します：

```go
type Option func(*Service)

func New(opts ...Option) (*Service, error)

func WithDisableLLMAuth(disable bool) Option {
    return func(s *Service) {
        s.disableLLMAuth.Store(disable)
    }
}
```

### MCP ハンドラーパターン

MCP サーバーは合成インターフェースパターンを使用します。`internal/server/mcp/handler.go` の `Svc` インターフェースは、より小さなインターフェース（`CatalogReader`、`EndpointExplorer`、`EndpointExecutor`、`SystemInfo`、`ResponseManager`）から構成されています。各ハンドラーメソッドはサービス層に委譲します：

```go
type handler struct {
    service Svc
}

func (h *handler) handleSearch(ctx context.Context, _ *sdkmcp.CallToolRequest, req service.SearchRequest) (*sdkmcp.CallToolResult, any, error) {
    resp, err := h.service.Search(ctx, req)
    if err != nil {
        return nil, nil, err
    }
    return &sdkmcp.CallToolResult{StructuredContent: resp}, nil, nil
}
```

### LLMError

LLM に返されるすべてのエラーは、8 つのコードのいずれかを持つ `LLMError` 型を使用します：

| コード | 発生時 |
|------|------|
| `validation_failed` | 無効な入力（間違った ID 形式、必須フィールドの欠落） |
| `not_found` | インデックスにエンティティが見つからない |
| `rate_limit` | エンドポイントごとの 10 秒クールダウンを超過 |
| `invoke_error` | HTTP リクエスト/レスポンスの失敗 |
| `config_error` | 設定の読み込みまたは検証の失敗 |
| `workspace_error` | ワークスペースディレクトリまたはファイル操作の失敗 |
| `parse_error` | スペックファイルの解析の失敗 |
| `auth_error` | 認証トークン取得の失敗 |

メッセージは、何が問題だったか**および**次に何をすべきかを、LLM コンシューマーに適した平易な言葉で説明する必要があります。

### ID 生成

すべての ID は決定論的な MD5 ハッシュです：

```go
id.Domain("meteo")                          // 32 文字の 16 進数
id.Collection("meteo", "Forecast")          // 32 文字の 16 進数
id.Tag("meteo", "Forecast", "pets")         // 32 文字の 16 進数
id.Method("meteo", "Forecast", "pets", "GET", "/v2/pet/{petId}")
```

### 設定カスケード

設定は **グローバル → スペック → コレクション** の 3 つのレベルでカスケードされます。各レベルが前のレベルを上書きします。すべての `http_client` 設定はすべてのレベルで上書き可能です。ヘッダーとクッキーはマージされ、単純な値は置き換えられます。

## クイックリファレンス

| 領域 | 規約 |
|------|------------|
| **Go バージョン** | 1.23+ |
| **フォーマッター** | gofmt、gofumpt、goimports、gci |
| **行長** | 120 文字 |
| **リンター** | `.golangci.yml` で 80 以上 |
| **エラー型** | 8 つのコードを持つ `LLMError` |
| **モックフレームワーク** | `go.uber.org/mock` |
| **テストヘルパー** | `newTestService()`、`seedTestData()` |
| **設定形式** | カスケード付き YAML |
| **認証ディスパッチ** | `UnmarshalYAML` が `type` フィールドを読み取る |
| **ID 生成** | MD5 ベース（`id.Domain()`、`id.Collection()` など） |
| **レート制限** | `invoke` でエンドポイントごとに 10 秒 |
| **レスポンスサイズ** | デフォルト 1 MB、超過時はファイルに保存 |
| **カバレッジ目標** | コアパッケージで 80% 以上 |
| **ビルド** | `make build` |
| **リンター** | `make lint` |
| **テスト** | `go test ./...` |
| **生成** | `go generate ./...` |
