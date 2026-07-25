# add

## 用途

向现有配置添加新的 **spec**（API 服务）或 **collection**（OpenAPI/Swagger/Postman 文件）。这是扩展工作区以包含新 API 的主要方式。

## 何时使用

- 你有新的 API 要连接到 LLM 智能体
- 你找到了 OpenAPI 规范 URL 并想添加它
- 你想向现有 spec 添加额外的规范文件（collection）
- 你更倾向于直接编写 YAML 而不是使用交互式向导

## 语法

```bash
swag2mcp add spec [path] [flags]
swag2mcp add collection [path] [flags]
```

## 参数

| 参数 | 位置 | 必需 | 描述 |
|------|------|------|------|
| `path` | 1 | 否 | 工作区目录。如果省略，通过路径解析规则解析。 |

## 标志

### `add spec`

| 标志 | 简写 | 类型 | 默认值 | 描述 |
|------|------|------|--------|------|
| `--yaml` | `-y` | `string` | `""` | YAML 输入内联或 `-` 表示 stdin |
| `--example` | `-e` | `bool` | `false` | 打印 YAML 模板并退出 |

### `add collection`

| 标志 | 简写 | 类型 | 默认值 | 描述 |
|------|------|------|--------|------|
| `--yaml` | `-y` | `string` | `""` | YAML 输入内联或 `-` 表示 stdin |
| `--example` | `-e` | `bool` | `false` | 打印 YAML 模板并退出 |

## 工作原理

### 交互模式（默认）

启动 TUI 向导，让你逐步填写 spec 或 collection 字段。

```bash
swag2mcp add spec
swag2mcp add collection
```

### YAML 内联模式

直接将 YAML 作为字符串传递。**注意 shell 引号** — 特殊字符如 `:`、`#`、`&`、`{` 可能会破坏命令。

```bash
swag2mcp add spec --yaml 'domain: meteo
llm_title: Open-Meteo API
base_url: https://meteo.swagger.io/v2
collections:
  - llm_title: Main
    location: https://example.com/spec.json'
```

### 从 stdin 输入 YAML（推荐用于复杂 YAML）

从文件管道输入或使用 heredoc 以完全避免 shell 引号问题：

```bash
# 从文件管道输入
cat spec.yaml | swag2mcp add spec --yaml -

# Heredoc
swag2mcp add spec --yaml - <<EOF
domain: my-api
llm_title: My API
llm_instruction: "Use this API for X & Y # important"
base_url: https://api.example.com/v1
collections:
  - llm_title: Main
    location: https://raw.githubusercontent.com/org/repo/main/spec.yaml
EOF
```

### YAML 模板

打印预期的 YAML 结构并退出：

```bash
swag2mcp add spec --example
swag2mcp add collection --example
```

## YAML 格式

### Spec

```yaml
domain: meteo
llm_title: Open-Meteo API
llm_instruction: Use this API to manage pets.
base_url: https://meteo.swagger.io/v2
tags: [public, demo]
auth:
  type: bearer
  config:
    token: $(TOKEN)
collections:
  - llm_title: Open-Meteo Swagger
    location: https://example.com/spec.json
```

### Collection

```yaml
spec_domain: meteo
llm_title: Orders Collection
location: https://example.com/orders.json
```

## 命令后验证

```bash
swag2mcp ls [path]
# 新的 spec 或 collection 应出现在列表中
```

## 细节

- **自动初始化：** 如果不存在配置文件，`add` 会自动先运行初始化向导。你不需要单独运行 `init`。
- **Shell 引号：** 内联 YAML（`--yaml '...'`）对于特殊字符很脆弱。对于超出简单值的任何内容，优先使用 `--yaml -` 配合 heredoc 或管道。
- **`--example` 立即退出**，不检查现有配置或修改任何内容。
- **`add spec` vs `add collection`：** 对于新的 API 服务（新域）使用 `add spec`。要向现有 spec 添加另一个规范文件，使用 `add collection`。
