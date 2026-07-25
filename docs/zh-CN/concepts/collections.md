# Collections

一个 collection 是描述特定 API 的单个 OpenAPI/Swagger/Postman 文件。它指向一个 `location`（URL 或本地文件路径），并属于一个 spec（域）。

一个 spec 可以有多个 collection — 例如，"meteo" spec 可能有"Forecast"、"Air Quality"和"Marine" collection，每个指向不同的规范文件。

## Collection 字段

| 字段 | YAML 键 | 必需 | 描述 |
|------|---------|------|------|
| [LLM 标题](#llm-instruction) | `llm_title` | ❌ | 给 LLM 看的 collection 显示名称（最多 120 字符）。如果未设置，从规范文档自动填充 |
| [LLM 指令](#llm-instruction) | `llm_instruction` | ❌ | 给 LLM 的简短提示（最多 360 字符）。如果未设置，从规范文档自动填充 |
| 标题 | `title` | ❌ | 原始规范标题覆盖（从解析的文档自动填充） |
| [位置](#location--规范文件的解析方式) | `location` | ✅ | 规范文件的 URL 或路径（5–250 字符） |
| [禁用](#disable) | `disable` | ❌ | 加载时跳过此 collection |
| [HTTP 客户端](#http-client-覆盖) | `http_client` | ❌ | 每个 collection 的 HTTP 设置（头、cookie） |
| [基础 URL](#base-url-覆盖) | `base_url` | ❌ | 为此 collection 覆盖 spec 的基础 URL |
| [模拟服务器](#mock-server) | `base_mock_url` | ❌ | 模拟服务器地址，格式为 `host:port`。当 `mock_enabled: true` 时必需 |

## Location — 规范文件的解析方式

`location` 字段告诉 swag2mcp 在哪里找到 OpenAPI/Swagger/Postman 文件。它支持多种来源类型：

| 来源 | 示例 | 描述 |
|------|------|------|
| **远程 URL** | `https://raw.githubusercontent.com/.../spec.yaml` | 下载并缓存 |
| **本地文件（绝对路径）** | `/home/user/my-api.yaml` | 从文件系统读取，缓存 |
| **本地文件（相对路径）** | `./my-api.yaml` | 解析为绝对路径，缓存 |
| **工作区本地文件** | `specs/my-api.yaml` | 存储在 `~/.swag2mcp/specs/`，直接使用（不缓存） |
| **file:// URI** | `file:///home/user/spec.yaml` | 转换为本地路径，缓存 |

swag2mcp 自动检测来源类型：

- `https://` 或 `http://` → 远程 URL（缓存）
- `file://` → 本地文件（转换为文件系统路径）
- 其他所有内容 → 本地文件（支持 `~` 展开到主目录）

### 远程 URL

当你使用远程 URL 时，swag2mcp 下载文件并在本地缓存。缓存会在后续启动时重用，避免重复下载。

### 本地文件

本地文件直接从文件系统读取。如果文件在工作区 `specs/` 目录之外，为了保持一致性，会复制到缓存。

### 工作区本地文件

工作区内的 `specs/` 目录（`~/.swag2mcp/specs/`）是本地规范文件的推荐位置。存储在此处的文件直接使用，无需缓存。使用以 `specs/` 开头的相对路径来引用它们。

> **注意：** `specs/` 只是一个目录名称（如 `cache/` 或 `responses/`），不是"spec"的概念。它存储 collection 指向的实际 OpenAPI/Swagger/Postman 文件。

```bash
# 将规范文件导入工作区
swag2mcp import https://example.com/api.yaml myspec

# 导入后，location 变为：
# specs/myspec.yaml
```

## 缓存系统

swag2mcp 缓存远程规范文件，避免每次启动时下载。

### 工作原理

1. 当加载带有远程 URL 的 collection 时，swag2mcp 检查缓存
2. 如果存在有效（未过期）的缓存条目，直接使用
3. 如果不存在，则下载、解析文件并存储在缓存中

### 缓存结构

```
~/.swag2mcp/
  cache/
    {sha256_hash}.spec    # 缓存的规范文件内容
    {sha256_hash}.meta    # 缓存元数据（JSON）
```

每个缓存文件都有一个包含以下内容的元数据文件：

```json
{
  "source": "https://example.com/api.yaml",
  "source_type": "url",
  "cached_at": "2024-01-01T00:00:00Z",
  "mod_time": "2024-01-01T00:00:00Z",
  "ttl_sec": 3600
}
```

### 缓存 TTL

每个缓存文件获得 **1 小时到 48 小时** 之间的随机 TTL。这防止所有缓存文件同时过期（惊群问题）。

### 缓存键

缓存键是原始 location 字符串的 SHA-256 哈希（前 16 字节 = 32 个十六进制字符）。

### 管理缓存

```bash
# 清除缓存和响应，重新下载所有规范文件
swag2mcp update

# 仅清除缓存和响应
swag2mcp clean
```

- `swag2mcp update` — 验证配置，清除 `cache/` 和 `responses/`，然后重新缓存所有 collection 位置
- `swag2mcp clean` — 删除 `cache/` 和 `responses/` 的所有内容，以及孤立的认证脚本
- 旧响应在 MCP 服务器启动后 48 小时自动清理

## 验证

每个 collection 在加载配置时都会被验证。验证在每次 `swag2mcp mcp` 启动时运行。如果失败，MCP 服务器将不会启动 — 在某些 IDE 中，这意味着服务器根本无法连接，LLM 会收到清晰的错误消息，说明需要修复什么。

| 检查项 | 规则 |
|--------|------|
| **位置** | 必需，5–250 字符 |
| **位置可访问性** | 必须是可达的 URL 或存在的文件 |
| **位置有效性** | 必须是有效的 OpenAPI 3.x、Swagger 2.0 或 Postman 文件 |
| **LLM 标题** | 最多 120 字符，字母/数字/基本标点 |
| **LLM 指令** | 最多 360 字符，与标题相同的字符集 |
| **基础 URL** | 如果设置，必须是有效的 URL |
| **基础模拟 URL** | 必须是 `host:port` 或 `host:port/path`，其中 host 为 `localhost`、`127.0.0.1` 或 `0.0.0.0` |
| **模拟必需** | 如果 `mock_enabled: true`，每个 collection 必须有 `base_mock_url` |
| **重复模拟端口** | 没有两个 collection 可以共享相同的模拟端口 |

要在启动服务器之前诊断问题，使用 [`validate`](../cli/validate.md) 命令：

```bash
# 验证默认工作区（~/.swag2mcp）
swag2mcp validate

# 验证自定义项目工作区
swag2mcp validate ./my-project
```

## 添加 Collection

### 通过 YAML 配置

直接编辑 `~/.swag2mcp/swag2mcp.yaml`：

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

编辑后，重启 MCP 服务器（`swag2mcp mcp`）以使更改生效。

### 通过 CLI

```bash
# 交互模式
swag2mcp add collection

# 非交互模式，使用 YAML
swag2mcp add collection --yaml 'spec_domain: meteo
llm_title: Forecast
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml'

# 从 stdin 管道输入
cat collection.yaml | swag2mcp add collection --yaml -

# 显示 YAML 示例
swag2mcp add collection --example
```

### 通过导入

```bash
# 将规范文件导入工作区
swag2mcp import https://example.com/api.yaml
```

## LLM 指令

Collection 可以有自己 `llm_instruction`（最多 360 字符）以提供更具体的指导。这会与 spec 级别的指令一起注入到 swag2mcp 系统提示中。

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        llm_instruction: "Use this collection for current weather and daily forecasts."
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        llm_instruction: "Use this collection for air quality index and pollution data."
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
```

如果未设置 `llm_title`，它会从规范文档的 `title` 字段自动填充。如果未设置 `llm_instruction`，它会从规范文档的 `description` 字段填充。

## 禁用

设置 `disable: true` 跳过 collection。它不会被加载、索引或提供给 LLM。

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: true
```

## 基础 URL 覆盖

每个 collection 可以覆盖 spec 的 `base_url`。当同一 spec 中的不同 collection 使用不同的 API 端点时，这很有用。

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

## HTTP 客户端覆盖

Collection 可以覆盖 spec 和全局级别的 HTTP 设置（头、cookie）。

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          headers:
            X-API-Version: "2"
          cookies:
            - name: session
              value: abc123
```

设置级联：全局 → spec → collection。详情请参见[配置级联](../configuration/cascade.md)。

## 模拟服务器

当在配置级别设置 `mock_enabled: true` 时，每个 collection 必须设置 `base_mock_url`。这告诉 swag2mcp 模拟服务器正在为此 collection 运行的位置。

```yaml
mock_enabled: true
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        base_mock_url: localhost:8080
```

完整详情请参见[模拟服务器](../advanced/mock-server.md)。

## 示例

### 最小 Collection

```yaml
specs:
  - domain: dadjokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

### 包含所有字段的完整 Collection

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        llm_instruction: "Use for current weather and daily forecasts."
        title: "Custom Title"
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        disable: false
        base_url: https://forecast-api.open-meteo.com
        base_mock_url: localhost:8080
        http_client:
          headers:
            X-Custom: value
```

### 每个 Spec 多个 Collection

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

### 工作区中的本地文件（specs/ 目录）

```yaml
specs:
  - domain: myapi
    llm_title: My Internal API
    base_url: https://api.mycompany.com
    collections:
      - llm_title: Users
        location: specs/users.openapi.json
      - llm_title: Orders
        location: specs/orders.openapi.json
```

### 禁用的 Collection

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: true
```

## 相关

- [Collection 设置（配置）](../configuration/collection-settings.md) — 完整 YAML 参考
- [配置级联](../configuration/cascade.md) — 设置如何相互覆盖
- [Specs](./specs) — collection 的逻辑容器
- [HTTP 客户端](../configuration/http-client.md) — HTTP 客户端配置
- [模拟服务器](../advanced/mock-server.md) — 模拟服务器设置
- [CLI: validate](../cli/validate.md) — validate 命令参考
- [CLI: update](../cli/update.md) — update 命令参考
- [CLI: clean](../cli/clean.md) — clean 命令参考
