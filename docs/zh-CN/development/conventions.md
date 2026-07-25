# 代码规范

## Go

- **Go 1.26+**
- **gofmt** / **gofumpt** / **goimports** / **gci**
- **每行 120 字符**
- **使用守卫子句** 代替嵌套 if
- **命名**：私有使用 `camelCase`，导出使用 `PascalCase`

## 错误

对 LLM 可见的错误使用 `LLMError`：

```go
type LLMError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

错误代码：
- `validation_failed` — 无效参数
- `not_found` — 资源未找到
- `rate_limit` — 超过速率限制
- `invoke_error` — API 调用错误

## 接口

- 小接口（1-3 个方法）
- 接口组合
- 配置使用函数选项模式

## 测试

- 表格驱动测试
- 测试辅助函数（`newTestService()`、`seedTestData()`）
- 通过 `go.uber.org/mock` 生成模拟
- 核心包覆盖率 80%+

## 配置

- YAML 格式
- 级联：全局 → spec → collection
- 通过 `go-playground/validator` 验证
- 通过 `$(VAR)` 使用环境变量
