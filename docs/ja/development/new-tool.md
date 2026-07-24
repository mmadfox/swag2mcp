# 新しい MCP ツールの追加

## 手順

1. **ツール名定数を追加** `internal/service/service.go` に
2. **リクエスト/レスポンス型を作成** `internal/service/types.go` に
3. **サービスを実装** `internal/service/` に（新規ファイルまたは既存ファイルに追加）
4. **マークダウン定義を作成** `internal/service/definitions/` に — これは `MakeToolDefinitions` が読み取るものです
5. **`Svc` インターフェースにメソッドを追加** `internal/server/mcp/handler.go` に
6. **ハンドラーを追加** `handler.go` に
7. **ツールを登録** `mcp.go` の `registerTools` で
8. **モックを生成**：`go generate ./...`
9. **テストを書く**

## 1. ツール名定数

`internal/service/service.go` に定数を追加：

```go
const MyNewTool = "my_new_tool"
```

## 2. リクエスト/レスポンス型

`internal/service/types.go` で定義：

```go
type MyNewToolRequest struct {
    Param1 string `json:"param1" validate:"required" jsonschema:"required,Description of param1"`
}

type MyNewToolResponse struct {
    Result string `json:"result"`
}
```

## 3. サービスの実装

`internal/service/my_new_tool.go` を作成するか、既存のサービスファイルに追加します。標準のサービスパターンに従います：検証 → 検索 → 実行 → 返却：

```go
func (s *Service) MyNewTool(ctx context.Context, req MyNewToolRequest) (MyNewToolResponse, error) {
    if err := s.validateRequest(req); err != nil {
        return MyNewToolResponse{}, NewLLMError(validationFailedErrCode, err.Error())
    }
    // ビジネスロジック
    return MyNewToolResponse{Result: "ok"}, nil
}
```

## 4. マークダウン定義

`internal/service/definitions/my_new_tool.md` を作成します。このファイルは `MakeToolDefinitions()` によって読み取られ、バイナリに埋め込まれます。フロントマターの `name:` フィールドは定数と一致する必要があります：

```markdown
---
name: my_new_tool
---

# my_new_tool

ツールの説明。

## パラメータ

| パラメータ | 型 | 説明 |
|-----------|------|-------------|
| `param1` | string | 説明 |
```

`tools.go` の `MakeToolDefinitions()` 関数は、埋め込まれた `definitions/` ディレクトリからすべての `.md` ファイルを読み取り、`name` フィールドの YAML フロントマターを解析し、本文をツールの説明として使用します。`instruction.md` ファイルは特別に扱われ、LLM のシステム指示になります。

## 5. Svc インターフェース

`handler.go` の合成 `Svc` インターフェースにメソッドを追加：

```go
type Svc interface {
    // ... 既存のメソッド
    MyNewTool(ctx context.Context, req service.MyNewToolRequest) (service.MyNewToolResponse, error)
}
```

## 6. ハンドラー

`handler.go` の `handler` にハンドラーメソッドを追加します。ハンドラーはサービスに委譲し、結果を `StructuredContent` でラップします：

```go
func (h *handler) handleMyNewTool(
    ctx context.Context,
    _ *sdkmcp.CallToolRequest,
    req service.MyNewToolRequest,
) (*sdkmcp.CallToolResult, any, error) {
    resp, err := h.service.MyNewTool(ctx, req)
    if err != nil {
        return nil, nil, err
    }
    return &sdkmcp.CallToolResult{
        StructuredContent: resp,
    }, nil, nil
}
```

## 7. 登録

`mcp.go` の `registerTools` 関数でツールを登録します。`toolRegistrations` マップにエントリを追加：

```go
service.MyNewTool: {
    addTool[service.MyNewToolRequest](mcpServer, h.handleMyNewTool),
    true, // ツールが可変の場合は false（invoke や auth など）
},
```

`registerTools` 関数のシグネチャは次のとおりです：

```go
func registerTools(mcpServer *sdkmcp.Server, tools []service.Tool, h handler) {
```

`MakeToolDefinitions()` が返すツール定義を反復処理し、それぞれを型付きハンドラーに登録します。`toolRegistrations` マップはツール名定数をハンドラーに接続します。
