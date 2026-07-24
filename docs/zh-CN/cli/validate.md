# validate

## 用途

检查配置文件和所有引用的规范文件是否存在错误。这是一个**只读**诊断命令 — 它从不修改任何内容。

## 何时使用

- 手动编辑 `swag2mcp.yaml` 后
- 在运行 `mcp` 或 `update` 之前及早发现问题
- 排查为什么 spec 没有加载时
- 在 CI/CD 管道中验证配置更改

## 语法

```bash
swag2mcp validate [path] [flags]
```

## 参数

| 参数 | 位置 | 必需 | 描述 |
|------|------|------|------|
| `path` | 1 | 否 | 工作区目录。如果省略，通过路径解析规则解析。 |

## 标志

| 标志 | 简写 | 类型 | 默认值 | 描述 |
|------|------|------|--------|------|
| `--tags` | `-t` | `string` | `""` | 仅验证具有匹配标签的 spec（逗号分隔） |

## 工作原理

```bash
swag2mcp validate
swag2mcp validate ./my-workspace
swag2mcp validate --tags=public
```

## 检查的内容

| 检查项 | 描述 |
|--------|------|
| YAML 语法 | 配置文件必须是有效的 YAML |
| 配置结构 | 所有必需字段存在，类型正确 |
| 域唯一性 | 没有重复的域 |
| 域格式 | 仅限小写字母、数字、连字符 |
| 规范文件存在性 | `location` 文件或 URL 必须可达 |
| 规范格式 | 文件必须是有效的 OpenAPI 3.x、Swagger 2.0 或 Postman collection |
| 认证设置 | 认证类型和配置对所选方法有效 |
| HTTP 客户端 | HTTP 客户端设置有效 |

## 不检查的内容

| 不检查 | 原因 |
|--------|------|
| 认证端点 | `validate` 检查认证配置语法，但不测试登录/令牌交换 |
| API 端点可用性 | 只检查规范文件 URL，不检查 `base_url` |
| `base_url` 正确性 | 验证格式，但不发送测试请求 |
| 模拟服务器配置 | `base_mock_url` 不验证连接性 |

## 示例输出

```
✅ Configuration is valid.
✓ Spec petstore: OK
✓ Spec meteo: OK
✗ Spec old-api: file not found
```

## 命令后验证

如果验证通过，配置即可用于 `mcp`、`update` 或 `run`。

## 细节

- **无自动初始化：** 与 `add`、`ls` 或 `run` 不同，如果配置缺失，`validate` **不会**自动初始化。它返回错误：`"configuration not found at <path>"`。
- **网络访问：** 验证期间会获取远程规范 URL。如果规范托管在慢速服务器上，命令可能需要更长时间。
- **标签过滤：** 设置 `--tags` 时，只验证匹配指定标签的 spec。其他 spec 被跳过。
