# 配置文件

swag2mcp 使用 YAML 配置文件。由 `swag2mcp init` 创建。

## 位置

- **Linux/macOS**：`~/.swag2mcp/swag2mcp.yaml`
- **Windows**：`%USERPROFILE%\.swag2mcp\swag2mcp.yaml`

## 基本结构

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

## 完整示例

```yaml
# ── 全局 HTTP 客户端 ──────────────────────────────────
http_client:
  timeout: 30s
  max_response_size: 1048576
  user_agent: "swag2mcp-global/1.0"
  follow_redirects: true
  max_redirects: 10
  random: false
  proxy:
    url: ""
    username: ""
    password: ""
    bypass: []
  headers:
    "Accept": "application/json"
  cookies:
    - name: "session"
      value: "abc123"

# ── MCP 服务器 ──────────────────────────────────────────
mcp:
  transport: stdio
  addr: ":8080"
  path: "/mcp"
  auth:
    token: ""

# ── 模拟服务器 ─────────────────────────────────────────
mock_enabled: false
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

# ── 速率限制器 ────────────────────────────────────────
disable_ratelimiter: false
rate_limit_interval: 10s

# ── Specs ───────────────────────────────────────────────
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    llm_instruction: "Use this API for weather forecasts and climate data"
    base_url: https://api.open-meteo.com
    disable: false
    tags: ["weather", "climate"]
    http_client:
      timeout: 10s
    auth:
      type: bearer
      config:
        token: "$(TOKEN)"
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: false
        http_client:
          timeout: 5s

  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## 环境变量

使用 `$(VAR_NAME)` 语法引用环境变量。swag2mcp 在启动时解析它们。

```yaml
specs:
  - domain: meteo
    auth:
      type: bearer
      config:
        token: "$(API_TOKEN)"

mcp:
  auth:
    token: "$(MCP_TOKEN)"
```

`$(VAR)` 在以下位置被解析：
- 认证配置字段：`token`、`username`、`password`、`client_id`、`client_secret`、`api_key`、`secret_key`、`domain`
- MCP 服务器认证令牌：`mcp.auth.token`
- HTTP 客户端头和 cookie 值

`$(VAR)` **不会**在基础 URL 或 collection 位置中被解析。

## 验证

```bash
# 验证默认工作区（~/.swag2mcp）
swag2mcp validate

# 验证自定义项目工作区
swag2mcp validate ./my-project
```

如果工作区不在主目录中（例如在项目仓库内），运行 `validate`、`update`、`mcp` 或任何其他命令时始终指定路径。否则 swag2mcp 将使用默认的 `~/.swag2mcp` 工作区。
