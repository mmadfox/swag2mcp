# swag2mcp

**swag2mcp** 是一个 CLI 工具和 MCP（模型上下文协议）服务器，用于将 OpenAPI/Swagger/Postman API 规范与 LLM 代理（Opencode、Crush、Copilot、Cursor 等）连接。

它将您的 API 规范索引到全文检索引擎中，通过 14 个 MCP 工具暴露它们，并让 LLM 能够发现、检查和调用真实的 API 端点——无需编写任何集成代码。

---

## 目录

- [快速开始](#快速开始)
- [配置](#配置)
- [CLI 命令](#cli-命令)
- [MCP 服务器](#mcp-服务器)
- [搜索](#搜索)
- [工作目录 (Workspace)](#工作目录-workspace)
- [缓存](#缓存)
- [开发](#开发)

---

## 快速开始

```bash
# 安装
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest

# 初始化工作目录
swag2mcp init

# 启动 MCP 服务器（用于 LLM 代理）
swag2mcp mcp

# 或使用交互式浏览器
swag2mcp run
```

---

## 配置

### YAML 模式

```yaml
http_client:                        # 可选，全局 HTTP 默认设置
  headers:                          # 可选
    X-API-Version: "2"
  cookies: []                       # 可选
  user_agent: ""                    # 可选
  timeout: 0s                       # 可选
  follow_redirects: true            # 可选
  max_redirects: 10                 # 可选
  max_response_size: 1048           # 可选，字节（默认 1KB，最大 1MB）

specs:
  - domain: petstore                    # 必填，1-60 字符，[a-zA-Z0-9_-]
    llm_title: Petstore API             # 必填，5-120 字符
    llm_instruction: |                  # 可选，最多 500 字符
      使用此 API 管理宠物、订单和用户。
    base_url: https://petstore.swagger.io/v2  # 必填，有效 URL
    disable: false                      # 可选
    tags: [public, demo]                # 可选，用于过滤
    http_client:                        # 可选，覆盖全局设置
      headers:
        X-API-Version: "2"
    auth:                               # 可选
      type: bearer                      # 见认证方法
      config:
        token: $(TOKEN_AUTH)
    collections:
      - llm_title: Petstore Swagger     # 可选，最多 120 字符
        llm_instruction: |             # 可选，最多 360 字符
          Petstore 主要端点
        title: ""                      # 可选，从 spec 自动填充
        location: https://petstore.swagger.io/v2/swagger.json  # 必填，5-250 字符
        disable: false                  # 可选
        base_url: ""                    # 可选，覆盖 spec 的 base_url
        http_client: {}                 # 可选，覆盖 spec
```

### 标签 — 按项目过滤规范

标签允许按项目、环境或团队对规范进行分组。启动 MCP 服务器时，使用 `--tags` 仅加载匹配的规范：

```bash
# 仅启动公共规范的服务器
swag2mcp mcp --tags=public

# 启动多个标签的服务器
swag2mcp mcp --tags=public,internal

# 为不同项目运行多个服务器
swag2mcp mcp --tags=project-alpha --logfile=/tmp/swag2mcp-alpha.log
swag2mcp mcp --tags=project-beta  --logfile=/tmp/swag2mcp-beta.log
```

这允许从单个配置文件为不同项目运行独立的 MCP 服务器。

### 认证方法

| 类型 | 字段 | 配置示例 |
|------|------|----------|
| `none` | — | `type: none` |
| `basic` | `username`, `password` | `username: $(USER)`, `password: $(PASS)` |
| `bearer` | `token` | `token: $(TOKEN)` |
| `digest` | `username`, `password` | `username: admin`, `password: secret` |
| `api-key` | `key`, `value`, `in` (header/query) | `key: X-API-Key`, `value: $(KEY)`, `in: header` |
| `oauth2-cc` | `client_id`, `client_secret`, `token_url`, `scopes` | `client_id: $(ID)`, `token_url: https://auth.example.com/token` |
| `oauth2-pwd` | `username`, `password`, `client_id`, `client_secret`, `token_url`, `scopes` | `username: $(USER)`, `token_url: https://auth.example.com/token` |
| `script` | `source` | `source: 路径/to/auth.sh` |

所有字符串字段支持 `$(ENV_VAR)` 语法——在运行时从环境变量解析。

---

## CLI 命令

所有接受 `[path]` 的命令使用相同的路径解析：

```
swag2mcp <command>          → ~/.swag2mcp/
swag2mcp <command> ./       → {cwd}/.swag2mcp/
swag2mcp <command> path/to  → {cwd}/path/to/.swag2mcp/
```

### `init [path]`

初始化工作目录和配置。

| 标志 | 简写 | 默认值 | 描述 |
|------|------|--------|------|
| `--interactive` | `-i` | `false` | 运行交互式向导 |
| `--force` | `-f` | `false` | 覆盖现有配置 |

```bash
swag2mcp init              # 创建 ~/.swag2mcp/swag2mcp.yaml
swag2mcp init ./           # 创建 ./.swag2mcp/swag2mcp.yaml
swag2mcp init -i           # 交互式向导
```

### `add spec [path]` / `add collection [path]`

向配置添加规范或集合。

| 标志 | 简写 | 默认值 | 描述 |
|------|------|--------|------|
| `--yaml` | `-y` | `""` | YAML 输入（使用 `-` 表示 stdin） |
| `--example` | `-e` | `false` | 显示 YAML 示例 |

```bash
swag2mcp add spec
swag2mcp add spec --yaml 'domain: petstore\nllm_title: Petstore API\nbase_url: https://...'
cat spec.yaml | swag2mcp add spec --yaml -
swag2mcp add spec --example
```

### `delete spec [path]` / `delete collection [path]`

从配置中删除规范或集合。交互式选择。

```bash
swag2mcp delete spec
swag2mcp delete collection
```

### `ls [path]`

列出规范和集合。

| 标志 | 简写 | 默认值 | 描述 |
|------|------|--------|------|
| `--tags` | `-t` | `""` | 按标签过滤（逗号分隔） |

```bash
swag2mcp ls
swag2mcp ls --tags=public,internal
```

### `run [path]`

交互式 API 浏览器（TUI）。搜索、浏览、检查和保存端点。

```bash
swag2mcp run
```

### `validate [path]`

验证配置并检查所有集合位置是否可访问。

| 标志 | 简写 | 默认值 | 描述 |
|------|------|--------|------|
| `--tags` | `-t` | `""` | 按标签过滤规范 |

```bash
swag2mcp validate
swag2mcp validate --tags=public
```

### `clean [path]`

删除 `cache/` 和 `responses/` 目录的所有内容。

```bash
swag2mcp clean
```

### `update [path]`

验证配置，清除缓存，重新缓存所有规范文件。

```bash
swag2mcp update
```

### `mcp [path]`

在无头模式下启动 MCP 服务器（stdio 传输）。这是用于 LLM 集成的主要生产命令。

| 标志 | 简写 | 默认值 | 描述 |
|------|------|--------|------|
| `--logfile` | `-f` | `""` | 日志文件路径 |
| `--tags` | `-t` | `""` | 按标签过滤规范 |
| `--disable-llm-auth` | | `true` | `true` — 认证在后台进行（LLM 看不到令牌）。`false` — LLM 可以通过 `auth` 工具请求令牌 |
| `--dump-dir` | | `""` | HTTP 请求转储目录（调试） |

```bash
swag2mcp mcp
swag2mcp mcp --tags=public --logfile=/var/log/swag2mcp.log
swag2mcp mcp --disable-llm-auth=false
swag2mcp mcp --dump-dir=/tmp/dump
```

---

## MCP 服务器

MCP 服务器通过 stdio 传输暴露 14 个工具。LLM 代理（Opencode、Crush、Copilot、Cursor 等）在配置后自动连接。

### 工具层次结构

```
spec_list                       — 列出所有可用规范
  └─ spec_by_id                 — 按 ID 获取规范详情
       └─ collection_by_spec    — 规范中的集合
            └─ tag_by_collection     — 集合中的标签
                 └─ endpoint_by_tag  — 标签下的端点
                      └─ inspect          — 完整的 OpenAPI 操作
                           └─ invoke       — 执行 API 调用

search                          — 跨所有端点的全文搜索
```

### 工具参考

| 工具 | 参数 | 返回 | 描述 |
|------|------|------|------|
| `spec_list` | — | `Spec[]` | 所有可用规范 |
| `spec_by_id` | `id` | Spec + Collections | 规范详情 |
| `collection_by_spec` | `specId` | Collections | 规范中的集合 |
| `collection_by_id` | `id` | Collection + Tags | 集合详情 |
| `tag_by_collection` | `collectionId` | Tags | 集合中的标签 |
| `tag_by_spec` | `specId` | Tags | 规范中的所有标签 |
| `tag_by_id` | `id` | Tag | 单个标签元数据 |
| `endpoint_by_tag` | `tagId` | Endpoints | 标签下的端点 |
| `endpoint_by_collection` | `collectionId` | Endpoints | 集合中的所有端点 |
| `endpoint_by_spec` | `specId` | Endpoints | 规范中的所有端点 |
| `endpoint_by_id` | `id` | Endpoint | 快速端点摘要 |
| `search` | `query`, `limit` | Endpoints | 全文搜索 |
| `inspect` | `endpointId` | Full Operation | 完整的 OpenAPI 操作对象 |
| `invoke` | `endpointId`, `parameters`, `requestBody` | Response | 执行真实的 API 调用 |
| `auth` | `specId` | Token | 获取规范的认证令牌 |

---

## 搜索

### 查询语法

| 功能 | 语法 | 示例 |
|------|------|------|
| 词条 | `词条` | `宠物` |
| 短语 | `"短语"` | `"添加宠物"` |
| 字段：method | `method:词条` | `method:post` |
| 字段：tag | `tag:词条` | `tag:auth` |
| 字段：path | `path:词条` | `path:/users` |
| 字段：summary | `summary:词条` | `summary:login` |
| 必需 (AND) | `+词条` | `+method:post +tag:user` |
| 排除 (NOT) | `-词条` | `-deprecated` |
| 通配符 | `*` | `path:*/v2/*` |
| 模糊 | `词条~` | `watex~` |
| 正则 | `/模式/` | `/user(s\|sessions)/` |
| 加权 | `词条^N` | `tag:pet^5` |
| 匹配所有 | `*` | `*` |

### 示例

```
# 在 auth 标签中查找 POST 端点
+method:post +tag:auth

# 搜索与登录相关的端点
summary:"login"~

# 查找所有用户相关路径，排除已弃用的
path:*/users/* -deprecated

# 复杂查询
+method:get +tag:pet summary:"find by status"
```

### 索引字段

| 字段 | 类型 | 内容 |
|------|------|------|
| `method` | text | HTTP 方法（小写） |
| `tag` | text | 标签名称（小写） |
| `path` | text | API 路径（小写） |
| `summary` | text（已分析） | 端点摘要/描述（小写） |
| `_all` | text（已分析） | method + path + tag + summary |

---

## 工作目录 (Workspace)

### 目录结构

```
~/.swag2mcp/                    # 或 {项目}/.swag2mcp/
├── swag2mcp.yaml               # 配置文件
├── cache/                      # 缓存的远程规范
│   ├── {hash}.spec             # 规范文件内容
│   └── {hash}.meta             # JSON 元数据
├── specs/                      # 本地规范文件（用户管理）
├── responses/                  # 调用响应文件
└── auth_scripts/               # 认证脚本
```

### 路径解析

```
swag2mcp <command>          → ~/.swag2mcp/
swag2mcp <command> ./       → {cwd}/.swag2mcp/
swag2mcp <command> path/to  → {cwd}/path/to/.swag2mcp/
```

### .gitignore

只应忽略临时数据：

```
.swag2mcp/cache/*
.swag2mcp/responses/*
```

配置文件 `.swag2mcp/swag2mcp.yaml` 和 `.swag2mcp/specs/` 中的规范文件**必须提交到仓库**。

### 建议

将所有规范文件保存在 `.swag2mcp/specs/` 中——这是确保它们被直接使用而不被复制到缓存的唯一方法。

---

## 缓存

### 规则

| 来源 | 行为 |
|------|------|
| HTTP/HTTPS URL | 始终缓存。TTL：随机 1-48 小时。 |
| `specs/` 内的本地路径 | 直接使用，不缓存。 |
| `specs/` 外的本地路径 | 首次访问时复制到缓存。 |
| `file://` URL | 视为本地路径。 |

### 缓存键

规范化位置的 SHA-256 哈希（前 16 字节 = 32 个十六进制字符）。

### 缓存命中逻辑

1. 读取 `.meta` 文件——过期或缺失 → 未命中
2. 对于本地源：`ModTime` 已更改 → 未命中
3. `.spec` 文件缺失 → 未命中
4. 否则 → 命中

---

## 开发

```bash
# 构建
go build ./cmd/swag2mcp/

# 测试
go test ./...

# 代码检查
make lint

# 运行
go run ./cmd/swag2mcp/main.go
```
