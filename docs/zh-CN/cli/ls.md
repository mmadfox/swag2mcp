# ls

## 用途

以人类可读格式列出所有配置的 **spec** 及其 **collection**。这是检查工作区中可用 API 的主要方式。

## 何时使用

- 你想查看配置了哪些 API
- 你需要查找 spec 或 collection ID
- 你想检查每个 collection 有多少端点
- 你想按标签过滤 spec

## 语法

```bash
swag2mcp ls [path] [flags]
```

## 参数

| 参数 | 位置 | 必需 | 描述 |
|------|------|------|------|
| `path` | 1 | 否 | 工作区目录。如果省略，通过路径解析规则解析。 |

## 标志

| 标志 | 简写 | 类型 | 默认值 | 描述 |
|------|------|------|--------|------|
| `--tags` | `-t` | `string` | `""` | 按标签过滤 spec（逗号分隔） |

## 工作原理

### 列出所有 spec

显示每个 spec 及其域、collection 和端点计数：

```bash
swag2mcp ls
```

示例输出：

```
Specifications:
  dadjoke (https://icanhazdadjoke.com)
    jokes (3 endpoints)
  meteo (https://meteo.swagger.io/v2)
    forecast (5 endpoints)
    current (8 endpoints)
  binance (https://api.binance.com)
    market-data (12 endpoints)
```

### 按标签过滤

仅显示具有指定标签的 spec：

```bash
swag2mcp ls --tags=public
swag2mcp ls --tags=public,internal
```

## 命令后验证

在 `add`、`delete`、`update` 或 `import` 之后使用 `ls` 确认工作区状态符合预期。

## 细节

- **自动初始化：** 如果不存在配置文件，`ls` 会自动先运行初始化向导。
- **标签过滤：** 标签以逗号分隔。显示匹配**任何**指定标签的 spec（OR 逻辑）。
- **输出格式：** 输出是纯文本，不是 JSON。对于机器可读的输出，使用 `info`。
