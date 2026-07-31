# Script Auth

## 目的

通过外部脚本进行身份验证。swag2mcp 支持两种脚本类型：
- **`.sh`** — Linux 和 macOS 上的 shell 脚本
- **`.bat`** — Windows 上的批处理脚本

脚本必须放置在工作区的 `auth_scripts` 目录中，并将 JSON 令牌输出到 stdout。

## 何时使用

- 自定义或非标准认证方案
- 复杂的令牌获取逻辑（多步骤、带额外检查）
- 当标准方法都不适合你的需求时

## 配置

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: script
      config:
        domain: "jokes"
```

## 参数

| 参数 | 必需 | 描述 |
|------|------|------|
| `domain` | 是 | 脚本文件名（不含扩展名）。必须与 spec 的 `domain` 匹配 — 例如，如果 `domain: jokes`，脚本文件必须是 `jokes.sh`（Unix）或 `jokes.bat`（Windows）。 |

## 脚本位置

脚本必须放在工作区的 `auth_scripts` 目录中：

- **Linux / macOS：** `{workspace}/auth_scripts/{domain}.sh`
- **Windows：** `{workspace}/auth_scripts/{domain}.bat`

## 脚本输出格式

脚本必须将 JSON 输出到 stdout，包含令牌及其过期时间：

```bash
#!/bin/bash
# auth_scripts/jokes.sh

TOKEN=$(curl -s -X POST https://auth.example.com/token \
  -d "grant_type=client_credentials" \
  -d "client_id=$CLIENT_ID" \
  -d "client_secret=$CLIENT_SECRET" | jq -r '.access_token')

echo "{\"token\": \"$TOKEN\", \"expires_in\": 3600}"
```

### JSON 字段

| 字段 | 必需 | 描述 |
|------|------|------|
| `token` | 是 | 认证令牌 |
| `expires_in` | 否 | 令牌生命周期（秒，默认：3600） |

## 说明

- 如果缓存的令牌已过期，swag2mcp 在每个请求上运行脚本
- 脚本必须在 30 秒内完成
- 令牌被缓存直到其过期时间
- 脚本文件名 = `{domain}.sh`（Unix）或 `{domain}.bat`（Windows）
- `domain` 不能包含 `/` 或 `\`
