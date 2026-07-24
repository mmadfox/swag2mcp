# 模拟服务器

## 概述

模拟服务器基于你的 OpenAPI 模式生成模拟 API 响应。它让你无需进行真实 HTTP 调用即可测试 API 集成。这对于开发、测试 LLM 智能体和演示非常有用。

模拟服务器是一个**独立的二进制文件** — `swag2mcp-mock`。它不包含在主 `swag2mcp` 二进制文件中，必须单独安装。

## 安装

```bash
# 选项 1：从 GitHub Releases 下载
# 查找 swag2mcp-mock_&lt;version&gt;_&lt;os&gt;_&lt;arch&gt;.tar.gz

# 选项 2：使用 Go 安装
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

## 配置

在配置中启用模拟服务器：

```yaml
mock_enabled: true
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
        base_mock_url: "127.0.0.1:9090"
```

## 参数

### mock_enabled

- **类型：** `bool`
- **默认值：** `false`
- **效果：** 当为 `true` 时，每个活动的 collection 必须设置 `base_mock_url`。模拟服务器为每个 collection 启动 HTTP 服务器。

### mock_auth

模拟认证服务器的端口。这些模拟 OAuth2、Digest 和 HMAC 认证端点，让你无需真实凭据即可测试经过认证的 API。

| 字段 | 默认值 | 描述 |
|------|--------|------|
| `oauth2_port` | `9090` | 模拟 OAuth2 令牌服务器的端口 |
| `digest_port` | `9091` | 模拟 Digest 认证服务器的端口 |
| `hmac_port` | `9092` | 模拟 HMAC 认证服务器的端口 |

### base_mock_url（每个 collection）

- **类型：** `string`
- **必需：** 是（当 `mock_enabled: true` 时）
- **格式：** `host:port`（例如 `localhost:8080`、`127.0.0.1:9000`）
- **效果：** 每个 collection 在此地址上获得自己的 HTTP 服务器。服务器使用随机生成的数据响应规范中定义的所有端点。

## 启动模拟服务器

```bash
# 使用默认配置启动
swag2mcp-mock mockserver

# 使用 TLS 启动
swag2mcp-mock mockserver --tls

# 使用自定义 TLS 证书启动
swag2mcp-mock mockserver --tls --tls-cert cert.pem --tls-key key.pem
```

### TLS 标志

| 标志 | 描述 |
|------|------|
| `--tls` | 使用自签名证书启用 TLS |
| `--tls-cert` | TLS 证书文件路径 |
| `--tls-key` | TLS 密钥文件路径 |

如果设置了 `--tls` 但没有设置 `--tls-cert` 和 `--tls-key`，则会自动为 `localhost` 生成自签名证书。

## 模拟服务器的功能

当你启动模拟服务器时，它会：

1. **解析所有规范文件** — 读取每个 collection 的 OpenAPI/Swagger 规范
2. **注册处理程序** — 为规范中定义的每个路径和方法创建 HTTP 处理程序
3. **生成模拟数据** — 使用与响应模式匹配的随机生成数据响应（正确的类型、格式和结构）
4. **启动认证服务器** — 模拟 OAuth2、Digest 和 HMAC 认证端点用于测试

### 测试模拟

```bash
# 在一个终端中：
swag2mcp-mock mockserver

# 在另一个终端中：
curl http://localhost:8080/pets
# → [{"id":1,"name":"Pet_name","status":"available"}]
```

## 模拟数据的生成方式

模拟服务器基于 OpenAPI 模式生成逼真的模拟数据：

- **字符串** — 随机单词、句子或特定格式的值（电子邮件、URL、UUID、日期、电话等）
- **数字** — 指定范围内的随机整数和浮点数
- **布尔值** — 随机 true/false
- **数组** — 1 到 3 个随机项
- **对象** — 所有属性填充随机值
- **枚举** — 从枚举列表中随机取值
- **可空字段** — 有时返回 `null`（约 10% 概率）

## 使用场景

- **开发** — 无需真实 API 访问即可测试集成
- **测试 LLM 智能体** — 验证 LLM 能否发现、检查和调用端点
- **演示** — 无需配置真实 API 即可展示 swag2mcp 的工作
- **负载测试** — 在不访问真实 API 的情况下测试 MCP 服务器负载

## 重要说明

- **独立的二进制文件** — `swag2mcp-mock` 不包含在主 `swag2mcp` 二进制文件中。请单独安装。
- **每个 collection 有自己的端口** — 为每个 collection 配置 `base_mock_url`
- **认证模拟服务器是全局的** — 无论你有多少个 collection，OAuth2、Digest 和 HMAC 服务器都在配置的端口上运行
- **规范解析失败不会致命** — 如果 collection 的规范无法解析，会跳过并显示警告
- **自签名 TLS** — 使用 `--tls` 而不提供证书时，仅为 localhost 生成自签名证书
