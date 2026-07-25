# 故障排除

## 安装问题

### swag2mcp: command not found

二进制文件不在你的 PATH 中。

```bash
# 检查是否安装了 Go
go version

# 查找 Go 安装二进制文件的位置
go env GOPATH
# 通常是 ~/go 或 ~/go/bin

# 添加到 PATH（添加到 ~/.zshrc 或 ~/.bashrc）
export PATH=$PATH:$(go env GOPATH)/bin

# 或使用完整路径
~/go/bin/swag2mcp --version
```

如果你从 GitHub Releases 下载了二进制文件，请确保它在 PATH 中的目录中：

```bash
# 移动到 /usr/local/bin（macOS/Linux）
sudo mv swag2mcp /usr/local/bin/
```

### permission denied

二进制文件没有执行权限。

```bash
# 对于 go install（修复所有权）
sudo chown -R $(whoami) $(go env GOPATH)

# 对于下载的二进制文件
chmod +x /path/to/swag2mcp
```

### Go 版本太旧

swag2mcp 需要 Go 1.23+。

```bash
go version
# 如果版本 < 1.23，更新 Go：
# https://go.dev/dl/
```

### 找不到模拟服务器

模拟服务器是一个独立的二进制文件。需要显式安装：

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

## 配置问题

### 找不到配置文件

swag2mcp 找不到 `swag2mcp.yaml`。

```bash
# 创建新配置
swag2mcp init

# 或显式指定路径
swag2mcp mcp /path/to/workspace
swag2mcp ls /path/to/workspace
```

**常见原因：** 你从随机目录运行了 `swag2mcp mcp`，它查找的是 `~/.swag2mcp/` 而不是项目的工作区。始终显式传递路径。

### 加载了错误的工作区

swag2mcp 加载了与预期不同的工作区。

**解析顺序：** 显式 `[path]` → 当前目录（`./`）→ `~/.swag2mcp/`。如果你在没有路径的情况下从没有 `swag2mcp.yaml` 的目录运行 `swag2mcp mcp`，它会回退到 `~/.swag2mcp/`。

**修复：** 始终传递工作区路径：`swag2mcp mcp /path/to/your/workspace`

### YAML 解析错误

配置文件包含无效的 YAML 语法。

```bash
# 验证配置
swag2mcp validate

# 常见错误：
# - 使用制表符代替空格（YAML 需要空格）
# - 嵌套字段缺少缩进
# - 包含特殊字符的字符串未加引号（: # & {）
```

**提示：** 使用 YAML 检查器或支持 YAML 的编辑器来捕获语法错误。

### 验证失败："no specifications defined"

配置文件存在但没有 spec。

```bash
# 添加 spec
swag2mcp add spec

# 或编辑 swag2mcp.yaml 并添加至少一个 spec
```

### 验证失败："duplicate domain"

两个 spec 具有相同的 `domain` 值。域必须唯一。

```bash
# 列出当前 spec
swag2mcp ls

# 检查 swag2mcp.yaml 中的重复域
```

### 验证失败："invalid spec location"

`location` URL 或文件路径不可访问或不是有效的规范文件。

```bash
# 检查 URL 是否可达
curl -I https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml

# 检查本地文件是否存在
ls -la ./specs/my-api.yaml

# 验证文件是有效的 OpenAPI/Swagger/Postman
# （不仅仅是任何 JSON 或 HTML 页面）
```

**常见原因：** `location` 字段指向的是 API 端点本身（例如 `https://api.example.com/v1/users`），而不是规范文件 URL。location 必须指向 OpenAPI/Swagger/Postman 文件。

## MCP 服务器问题

### 端口已被占用

另一个进程正在使用该端口。

```bash
# 查找进程
lsof -i :8080

# 终止它
kill <PID>

# 或使用不同的端口
swag2mcp mcp --transport sse --http-addr :9090
```

### 连接被拒绝

MCP 服务器未运行或无法访问。

```bash
# 确保服务器正在运行
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080

# 在另一个终端中，检查健康端点
curl http://127.0.0.1:8080/health

# 如果使用自定义路径
curl http://127.0.0.1:8080/custom-path/health
```

### MCP 工具未在 LLM 客户端中显示

LLM 客户端看不到任何工具。

```bash
# 检查 spec 是否已加载
swag2mcp ls

# 检查 spec 是否未禁用
swag2mcp validate

# 检查服务器日志
swag2mcp mcp --logfile /tmp/swag2mcp.log
cat /tmp/swag2mcp.log

# 验证 IDE 配置中的工作区路径是否正确
# （必须是绝对路径）
```

**常见原因：**
- IDE 配置中的工作区路径错误
- 所有 spec 都设置了 `disable: true`
- 通过 `--tags` 过滤掉了 spec
- 指定路径下不存在配置文件

### MCP 握手失败（HTTP 传输）

对于 SSE 和 Streamable HTTP 传输，MCP 协议需要在工具调用工作之前进行初始化。

```
步骤 1：POST /mcp → {"method":"initialize", ...}
步骤 2：POST /mcp → {"method":"notifications/initialized"}
步骤 3：POST /mcp → {"method":"tools/list", ...}  ← 现在可以工作
```

确保你的 LLM 客户端在调用工具之前完成握手。

### 健康检查返回 404

健康端点路径可能与 MCP 路径不同。

```bash
# 默认健康端点
curl http://127.0.0.1:8080/health

# 如果你更改了 MCP 路径，健康检查仍在 /health
# （不受 --http-path 影响）
```

### Auth 工具不可用

`auth` MCP 工具没有显示。

`auth` 工具**默认是禁用的**（`--disable-llm-auth=true`）。这是生产环境安全的有意设计。

```bash
# 启用 auth 工具
swag2mcp mcp --disable-llm-auth=false
```

## 认证问题

### 401 Unauthorized

API 因缺少或无效的凭据而拒绝了请求。

```bash
# 检查是否配置了认证
swag2mcp info

# 验证配置
swag2mcp validate

# 检查环境变量是否已设置
echo $MY_TOKEN

# 验证令牌是否未过期（bearer 令牌是静态的）
```

**常见原因：**
- 令牌缺失或为空
- 环境变量未设置
- 令牌已过期（bearer 令牌不会自动刷新）
- 配置了错误的认证类型

### 403 Forbidden

API 因权限不足而拒绝了请求。

- 令牌可能没有所需的范围
- API 密钥可能没有此资源的访问权限
- 查看 API 文档了解所需权限

### OAuth2 令牌端点无法访问

swag2mcp 无法访问 OAuth2 令牌 URL。

```bash
# 检查配置中的 token_url
# 验证 URL 是否正确且可达
curl -X POST https://auth.example.com/oauth/token \
  -d "grant_type=client_credentials" \
  -d "client_id=test" \
  -d "client_secret=test"

# 检查网络连接
# 如果在公司代理后面，检查代理设置
```

### Digest 认证失败

swag2mcp 无法完成 Digest 认证握手。

- 服务器必须在 401 响应中返回 `WWW-Authenticate: Digest ...` 头
- 挑战被缓存 5 分钟 — 如果服务器更改了 nonce，请等待缓存过期
- 检查用户名和密码是否正确

### HMAC 签名不匹配

API 拒绝了 HMAC 签名的请求。

- 验证 `api_key` 和 `secret_key` 是否正确
- 检查 API 是否使用 Binance 风格的 HMAC-SHA256 签名
- 某些交易所使用不同的签名方法 — HMAC 认证专门用于兼容 Binance 的 API

### Script 认证失败

外部认证脚本失败。

```bash
# 检查脚本是否存在
ls -la ~/.swag2mcp/auth_scripts/my-domain.sh

# 手动运行脚本进行测试
sh ~/.swag2mcp/auth_scripts/my-domain.sh

# 检查脚本输出格式（必须是 JSON：{"token": "...", "expires_in": 3600}）
# 检查脚本是否在 30 秒内完成
# 检查脚本是否有执行权限
chmod +x ~/.swag2mcp/auth_scripts/my-domain.sh
```

## 搜索问题

### 没有搜索结果

搜索没有返回任何端点。

```bash
# 检查 spec 是否已加载
swag2mcp ls

# 检查 spec 是否未禁用
swag2mcp validate

# 尝试更简单的查询
# 尝试按方法搜索：method:GET
# 尝试按标签搜索：tag:pets

# 索引在每次 MCP 服务器启动时重建
# 如果你刚刚添加了 spec，请重启服务器
```

### 搜索结果不相关

查询太宽泛或模糊。

- 使用字段过滤器缩小范围：`method:GET +tag:pets`
- 使用精确短语：`"find pet by status"`
- 使用 `limit` 参数获取更集中的结果

## API 调用问题

### invoke 返回错误

API 调用失败。

```bash
# 检查错误消息 — 它包含 HTTP 状态码
# 4xx 错误：检查参数、认证或权限
# 5xx 错误：API 服务器有问题

# 在调用之前始终检查端点
inspect(endpointId: "...")

# 检查是否提供了所有必需参数
# 检查参数类型（字符串、数字、布尔值）
```

### 速率限制错误

LLM 调用同一端点太快。

每个端点有 10 秒的冷却时间。等待后再调用，或禁用速率限制器：

```yaml
disable_ratelimiter: true
```

### 响应太大（返回了 fileRef）

响应超过了 `max_response_size`。

这是正常情况。使用响应工具探索数据：

```
1. response_outline(path) → 了解结构
2. response_compress(path, mode: "first_of_array") → 获取样本
3. response_slice(path, jsonPath: "data.0") → 获取特定数据
```

或增加限制：

```yaml
http_client:
  max_response_size: 4194304  # 4 MB
```

### API 响应慢

API 响应时间太长。

```yaml
http_client:
  timeout: 120s  # 从默认的 30s 增加
```

## 工作区问题

### swag2mcp init 失败："directory is not empty"

目标目录已有文件。

```bash
# 使用 --force 覆盖
swag2mcp init --force

# 或使用不同的目录
swag2mcp init ./new-workspace
```

### swag2mcp update 失败

一个或多个规范文件无法下载。

```bash
# 检查错误消息中哪个 URL 失败
# 验证 URL 是否可访问
curl -I <failed-url>

# 检查网络连接
# 检查代理设置
```

### Export 没有创建 ZIP

`[output]` 参数必须是 `.zip` 结尾的文件路径，而不是目录。

```bash
# 正确
swag2mcp export /path/to/workspace /path/to/backup.zip

# 错误（不会创建 ZIP）
swag2mcp export /path/to/workspace /some/directory
```

### Import 失败："not a valid swag2mcp backup"

ZIP 文件不是由 `swag2mcp export` 创建的。

只有 `swag2mcp export` 创建的 ZIP 归档才能被导入。该归档具有特定的内部结构（`swag2mcp.yaml`、`specs/`、`auth_scripts/`）。

## TUI 问题

### TUI 显示不正确

终端太小或不支持所需功能。

- 最小终端尺寸：80×24 字符
- TUI 使用 Bubbletea，适用于大多数现代终端
- 尝试调整终端窗口大小
- 尝试不同的终端模拟器

### TUI 显示 "no specs found"

工作区没有配置的 spec。

```bash
# 检查 spec
swag2mcp ls

# 添加 spec
swag2mcp add spec
```

## 模拟服务器问题

### 模拟服务器无法启动

```bash
# 检查配置中 mock_enabled: true
# 检查每个 collection 是否设置了 base_mock_url
# 检查端口是否未被占用
lsof -i :9090

# 检查模拟服务器日志
swag2mcp-mock mockserver
```

### 模拟服务器返回空响应

规范文件可能没有定义响应模式。

- 模拟服务器从响应模式生成数据
- 如果找不到模式，返回 `{}`
- 检查你的 OpenAPI 规范是否定义了包含 `schema` 的 `responses`

## 网络问题

### 代理连接失败

swag2mcp 无法通过配置的代理连接。

```bash
# 检查代理 URL 格式（必须包含协议：http://、https://、socks5://）
# 检查代理凭据
# 检查绕过列表 — 目标可能在绕过列表中
# 使用 curl 测试代理
curl -x http://proxy.company.com:8080 https://api.example.com
```

### TLS/SSL 错误

证书验证失败。

- 如果对 MCP 服务器使用自签名证书，客户端必须信任它
- 对于使用 `--tls` 的模拟服务器，会自动生成自签名证书
- 对于 API 调用，swag2mcp 使用系统的证书存储

## 其他问题

### 磁盘使用率高

缓存和响应目录会随时间增长。

```bash
# 清理所有内容
swag2mcp clean

# 旧响应（超过 48 小时）在 MCP 服务器启动时自动清理
# 缓存文件在 1-48 小时内随机过期
```

### go install 后 "command not found"

`go install` 目录不在你的 PATH 中。

```bash
# 查找 Go 安装二进制文件的位置
go env GOPATH
# 添加到 PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### LLM 没有正确使用工具

LLM 可能需要更好的指令或格式化技能。

- 在 spec 配置中使用 `llm_instruction` 来描述 API 的功能
- 考虑使用 [swag2mcp-format 技能](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md) 实现一致的输出格式
- LLM 响应的质量取决于模型及其接收的指令

### 如何报告 bug？

在 [GitHub](https://github.com/mmadfox/swag2mcp/issues) 上提交 issue，包含以下信息：
- swag2mcp 版本（`swag2mcp --version`）
- 你的操作系统和架构
- 你运行的确切命令
- 完整的错误消息
- 你的配置文件（移除密钥）
