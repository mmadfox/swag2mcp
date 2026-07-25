# 测试

## 命令

```bash
# 单元测试
go test ./...

# 特定包
go test ./internal/service/...

# 集成测试
make integration-tests

# 覆盖率
make cover

# 所有测试
make testall
```

## 测试结构

```
tests/
├── main_test.go              # 入口点
├── suite_test.go             # 套件设置
├── suite_auth_test.go        # 认证测试
├── suite_config_test.go      # 配置测试
├── suite_mcp_tools_test.go   # MCP 工具测试
├── suite_search_test.go      # 搜索测试
├── suite_ratelimit_test.go   # 速率限制测试
├── suite_response_test.go    # 响应测试
├── suite_export_test.go      # 导出测试
├── suite_import_test.go      # 导入测试
├── suite_parsing_test.go     # 解析测试
├── suite_transport_test.go   # 传输测试
├── suite_mock_test.go        # 模拟服务器测试
├── suite_workspace_test.go   # 工作区测试
├── suite_errors_test.go      # 错误测试
└── suite_version_test.go     # 版本测试
```

## 覆盖率

目标：核心包 80%+：

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

## 模拟

使用 `go.uber.org/mock` 进行 MCP 服务器测试：

```bash
go generate ./...
```

从 `handler.go` 生成 `internal/server/mcp/mock_svc_test.go`。

## 表格驱动测试

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
