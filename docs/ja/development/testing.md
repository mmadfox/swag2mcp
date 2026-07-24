# テスト

## コマンド

```bash
# ユニットテスト
go test ./...

# 特定のパッケージ
go test ./internal/service/...

# 統合テスト
make integration-tests

# カバレッジ
make cover

# 全テスト
make testall
```

## テスト構造

```
tests/
├── main_test.go              # エントリポイント
├── suite_test.go             # スイートセットアップ
├── suite_auth_test.go        # 認証テスト
├── suite_config_test.go      # 設定テスト
├── suite_mcp_tools_test.go   # MCP ツールテスト
├── suite_search_test.go      # 検索テスト
├── suite_ratelimit_test.go   # レート制限テスト
├── suite_response_test.go    # レスポンステスト
├── suite_export_test.go      # エクスポートテスト
├── suite_import_test.go      # インポートテスト
├── suite_parsing_test.go     # 解析テスト
├── suite_transport_test.go   # トランスポートテスト
├── suite_mock_test.go        # モックサーバーテスト
├── suite_workspace_test.go   # ワークスペーステスト
├── suite_errors_test.go      # エラーテスト
└── suite_version_test.go     # バージョンテスト
```

## カバレッジ

目標：コアパッケージで 80% 以上：

- `auth`
- `cache`
- `config`
- `env`
- `httpclient`
- `id`
- `index`
- `server/mcp`
- `service`
- `spec`
- `workspace`

## モック

MCP サーバーテストには `go.uber.org/mock` を使用：

```bash
go generate ./...
```

`handler.go` から `internal/server/mcp/mock_svc_test.go` を生成します。

## テーブル駆動テスト

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "hello", "HELLO", false},
        {"empty input", "", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := DoSomething(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.Equal(t, tt.want, got)
        })
    }
}
```
