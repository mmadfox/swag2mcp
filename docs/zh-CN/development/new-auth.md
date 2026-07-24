# 添加新的认证方法

## 步骤

1. **创建认证客户端** 在 `internal/auth/&lt;name&gt;.go`
2. **实现 `Authenticator` 接口**
3. **添加类型常量** 到 `internal/auth/auth.go`
4. **添加 YAML 解码器** 到 `internal/config/auth.go`
5. **注册解码器** 到 `authDecoders` 映射
6. **编写测试**

## 1. 认证客户端

创建 `internal/auth/my_auth.go`：

```go
package auth

import "net/http"

type MyAuthClient struct {
    Token string `yaml:"token" validate:"required"`
}

func (c *MyAuthClient) New() error {
    c.Token = resolveEnv(c.Token)
    return nil
}

func (c *MyAuthClient) Type() Type {
    return MyAuth
}

func (c *MyAuthClient) Apply(req *http.Request, out *Info) error {
    if c.Token == "" {
        return nil
    }
    setAuthHeader(req, out, "X-My-Auth", c.Token)
    return nil
}

func (c *MyAuthClient) Validate() error {
    return authValidator.Struct(c)
}
```

## 2. Authenticator 接口

每个认证客户端必须实现：

```go
type Authenticator interface {
    New() error                    // 初始化，解析环境变量
    Type() Type                    // 返回认证类型标识符
    Apply(req *http.Request, out *Info) error  // 将认证应用于请求
    Validate() error               // 验证必需字段
}
```

## 3. 类型常量

添加到 `internal/auth/auth.go`：

```go
const MyAuth Type = "my-auth"
```

## 4. YAML 解码器

在 `internal/config/auth.go` 中添加解码器函数。解码器接收 `*yaml.Node`，必须将其解码为你的认证客户端结构体：

```go
func decodeMyAuth(node *yaml.Node) (auth.Authenticator, error) {
    var client auth.MyAuthClient
    if err := decodeConfig(node, &client); err != nil {
        return nil, err
    }
    return &client, nil
}
```

`decodeConfig` 辅助函数处理常见模式：检查节点不为空，将 YAML 解码到结构体，失败时返回描述性错误。

## 5. 注册解码器

将你的解码器添加到 `internal/config/auth.go` 中的 `authDecoders` 映射：

```go
var authDecoders = map[string]authDecoder{
    // ... 现有解码器
    auth.MyAuth.String(): decodeMyAuth,
}
```

`Auth` 上的 `UnmarshalYAML` 方法从 YAML 读取 `type` 字段，将下划线规范化为连字符，在 `authDecoders` 中查找解码器，并使用 `config` 节点调用它。这就是 swag2mcp 知道为每个 spec 实例化哪个认证客户端的方式。

## 6. 测试

创建 `internal/auth/my_auth_test.go`，使用表格驱动测试覆盖：

- `New()` 正确解析环境变量
- `Type()` 返回正确的类型
- `Apply()` 设置正确的头/查询参数
- `Apply()` 优雅处理空值
- `Validate()` 对有效配置通过
- `Validate()` 对缺少必需字段失败
