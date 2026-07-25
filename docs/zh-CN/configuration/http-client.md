# HTTP 客户端

swag2mcp 使用可配置的 HTTP 客户端进行所有 API 调用。这些设置在全局定义，可以在 spec 和 collection 级别被覆盖。

## 配置

```yaml
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
      domain: ".example.com"
      path: "/"
```

## 超时

控制 swag2mcp 在放弃之前等待 API 响应的时间。

- **类型：** 持续时间（Go 格式：`30s`、`60s`、`2m`）
- **默认值：** `30s`
- **范围：** 1 秒到 5 分钟
- **效果：** 如果 API 在此时间内未响应，请求失败并返回超时错误。
- **何时增加：** 慢速 API、大负载、不可靠的网络。
- **何时减少：** 内部 API、健康检查、快速失败场景。

```yaml
http_client:
  timeout: 60s
```

## 最大响应大小

限制响应在 swag2mcp 将其保存到磁盘而不是内联返回给 LLM 之前的大小。

- **类型：** `int`（字节）
- **默认值：** `1048576`（1 MB）
- **范围：** 256 到 10,485,760 字节（10 MB）
- **效果：** 当响应超过此限制时，保存到 `{workspace}/responses/` 作为 JSON 文件。LLM 收到文件引用，可以使用 `response_outline`、`response_compress` 和 `response_slice` 工具进行探索。
- **何时增加：** 返回大数据集的 API（报告、日志、分析）。
- **何时减少：** LLM 上下文窗口有限，或你更倾向于对所有响应使用基于文件的访问。

```yaml
http_client:
  max_response_size: 4194304  # 4 MB
```

## User-Agent

随每个请求发送的 `User-Agent` 头。某些 API 需要特定的用户代理或阻止已知的机器人用户代理。

- **类型：** `string`
- **默认值：** `"swag2mcp-global/1.0"`
- **效果：** 向 API 服务器标识你的应用程序。
- **何时更改：** API 需要特定的用户代理，或你想为分析标识你的应用程序。

```yaml
http_client:
  user_agent: "MyApp/1.0"
```

## 跟随重定向

控制 swag2mcp 是否自动跟随 HTTP 重定向（3xx 状态码）。

- **类型：** `bool`
- **默认值：** `true`
- **效果：** 当为 `true` 时，swag2mcp 跟随重定向最多 `max_redirects` 次。当为 `false` 时，重定向响应原样返回。
- **何时禁用：** 循环重定向的 API、需要手动检查重定向目标的安全敏感端点。

```yaml
http_client:
  follow_redirects: false
```

## 最大重定向次数

限制 swag2mcp 在停止前跟随的重定向次数。

- **类型：** `int`
- **默认值：** `10`
- **范围：** 0 到 50
- **效果：** 如果 API 重定向次数超过此限制，请求失败。
- **何时更改：** 具有长重定向链的 API，或减少以在重定向循环中更快失败。

```yaml
http_client:
  max_redirects: 5
```

## 随机化器

为每个请求添加类似浏览器的随机头，以避免指纹识别和阻止。

- **类型：** `bool`
- **默认值：** `false`
- **效果：** 当为 `true` 时，swag2mcp 为每个请求生成随机头：`User-Agent`（来自真实浏览器字符串池）、`Accept`、`Accept-Language`、`Accept-Encoding`、`Cache-Control`。这会覆盖 `user_agent` 设置。
- **何时启用：** 基于 User-Agent 或头模式阻止请求的 API、抓取场景。

```yaml
http_client:
  random: true
```

## 代理

代理服务器充当 swag2mcp 和目标 API 之间的中介。所有 HTTP 流量都通过它路由。

**你可能需要代理的情况：**
- **公司网络** — 所有出站流量必须通过公司代理
- **地理限制** — 某些 API 有区域限制，正确区域的代理可以绕过
- **静态 IP** — 需要 IP 白名单的 API
- **匿名性** — 对目标 API 隐藏源 IP

### 代理 URL

- **类型：** `string`
- **默认值：** `""`（无代理）
- **支持的协议：** `http`、`https`、`socks5`、`socks5h`
- **支持 `$(VAR)`：** ✅ 运行时解析

| 协议 | 描述 | 用例 |
|------|------|------|
| `http` | HTTP 流量的 HTTP 代理 | 公司代理、基本代理 |
| `https` | HTTPS 代理（CONNECT 隧道） | 安全的公司代理 |
| `socks5` | SOCKS5 代理（本地 DNS 解析） | 通用，任何协议 |
| `socks5h` | SOCKS5 代理（代理端 DNS 解析） | 代理有更好的 DNS 解析时 |

### 代理认证

如果代理需要认证，提供 `username` 和 `password`：

- **支持 `$(VAR)`：** ✅ 所有三个字段（`url`、`username`、`password`）在运行时解析

```yaml
http_client:
  proxy:
    url: "http://proxy.example.com:8080"
    username: "proxyuser"
    password: "$(PROXY_PASSWORD)"
```

### 代理绕过

不应通过代理的域名列表。适用于内部服务、localhost 或只能直接访问的 API。

```yaml
http_client:
  proxy:
    url: "http://proxy.example.com:8080"
    bypass:
      - "localhost"
      - "127.0.0.1"
      - "*.internal.company.com"
      - "api.local"
```

绕过支持通配符模式（`*.example.com` 匹配任何子域）。

## 头

添加到每个请求的自定义 HTTP 头。头在级联级别之间合并：

```
全局头 → Spec 头（合并）→ Collection 头（合并）
```

Collection 头覆盖 spec 头，后者覆盖相同键的全局头。

```yaml
http_client:
  headers:
    "Accept": "application/json"
    "Accept-Language": "en-US"
```

头值支持 `$(ENV_VAR)` 解析。

## Cookie

随每个请求发送的 Cookie。Cookie 在级联级别之间合并（较低级别覆盖全局中相同 cookie 名称的值）。

```yaml
http_client:
  cookies:
    - name: "session"
      value: "abc123"
      domain: ".example.com"
      path: "/"
      secure: false
      http_only: false
```

### Cookie 字段

| 字段 | 必需 | 描述 |
|------|------|------|
| `name` | 是 | Cookie 名称 |
| `value` | 是 | Cookie 值（支持 `$(ENV_VAR)` 解析） |
| `domain` | 否 | 域范围（例如 `.example.com`） |
| `path` | 否 | 路径范围（例如 `/`） |
| `secure` | 否 | 仅通过 HTTPS 发送 |
| `http_only` | 否 | 不可通过 JavaScript 访问 |

## Spec 级别的自定义头

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    http_client:
      headers:
        "Accept": "application/json"
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Spec 级别的 Cookie

```yaml
specs:
  - domain: example
    llm_title: Example API
    base_url: https://api.example.com
    http_client:
      cookies:
        - name: "session"
          value: "abc123"
        - name: "csrf"
          value: "$(CSRF_TOKEN)"
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## 级联

HTTP 客户端设置从全局级联到 spec 再到 collection。所有设置可以在每个级别被覆盖：

```
全局 (http_client)
    ↓ 覆盖（所有设置）
Spec (specs[].http_client)
    ↓ 覆盖（所有设置）
Collection (specs[].collections[].http_client)
```

**所有 HTTP 客户端设置**（超时、代理、用户代理、重定向、响应大小、随机化器、头、cookie）可以在 spec 和 collection 级别被覆盖。

详情请参见[配置级联](./cascade)。
