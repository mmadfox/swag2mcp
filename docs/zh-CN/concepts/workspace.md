# 工作区

工作区是 swag2mcp 存储所有数据的目录 — 配置、缓存的规范、本地规范文件、保存的响应和认证脚本。

## 结构

```
~/.swag2mcp/                          # 工作区根目录（默认）
├── swag2mcp.yaml                     # 配置文件
├── cache/                            # 缓存的远程规范文件
│   ├── a1b2c3d4e5f6...spec          # 缓存的规范内容
│   └── a1b2c3d4e5f6...meta          # 缓存元数据（JSON）
├── specs/                            # 本地规范文件
│   └── my-api.yaml
├── responses/                        # 保存的 API 响应（大响应）
│   ├── meteo-get-forecast-abc123.json
│   └── response-fragment-def456.json
└── auth_scripts/                     # 认证脚本
    ├── meteo.sh                      # Unix shell 脚本
    └── meteo.bat                     # Windows 批处理脚本
```

## 默认路径

- **Linux/macOS**：`~/.swag2mcp/`
- **Windows**：`%USERPROFILE%\.swag2mcp\`

## 自定义路径

```bash
swag2mcp mcp /path/to/workspace
swag2mcp mcp ./my-workspace
```

## 目录

### cache/

存储下载的远程规范文件。每个文件以其 URL 的 SHA-256 哈希作为文件名进行缓存：

- `{hash}.spec` — 缓存的规范文件内容
- `{hash}.meta` — JSON 元数据（来源 URL、缓存时间、TTL）

每个缓存文件有 1 小时到 48 小时之间的随机 TTL。每次启动时自动检查缓存 — 如果存在有效（未过期）的条目，则重用而不下载。

**命令：**
- `swag2mcp update` — 清除缓存并重新下载所有规范
- `swag2mcp clean` — 清除缓存和响应

### specs/

存储 collection 通过 `location: specs/{name}` 引用的本地规范文件。此处的文件直接使用，无需缓存。

此目录由以下方式填充：
- `swag2mcp import <source> <name>` — 下载远程规范并保存到此
- `swag2mcp export` — 将规范复制到导出 ZIP
- 手动放置 — 你可以自己将规范文件复制到此

### responses/

存储超过 `max_response_size` 限制（默认 1 MB）的 API 响应。当 LLM 调用端点且响应太大时，swag2mcp 将其保存到此并返回文件引用。

命名约定：`{domain}-{method}-{path_with_underscores}-{6char_hex}.json`

旧响应在 MCP 服务器启动后 48 小时自动清理。

### auth_scripts/

存储 `script` 认证类型的认证脚本。每个脚本以 spec 的域命名。

#### 命名约定

| 平台 | 文件名 | 示例 |
|------|--------|------|
| Unix（Linux、macOS） | `{domain}.sh` | `meteo.sh` |
| Windows | `{domain}.bat` | `meteo.bat` |

域不能包含 `/` 或 `\` 字符。

#### 脚本工作原理

1. swag2mcp 以 30 秒超时运行脚本
2. 脚本必须将有效的 JSON 输出到 stdout
3. swag2mcp 解析 JSON 并使用令牌进行 API 请求

#### 预期输出格式

```json
{
  "token": "your-token-here",
  "expires_in": 3600
}
```

| 字段 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `token` | string | ✅ | 认证令牌 |
| `access_token` | string | ❌ | `token` 的替代（优先检查） |
| `token_type` | string | ❌ | 令牌类型（例如"Bearer"） |
| `expires_in` | number | ❌ | 令牌生命周期（秒，默认：3600） |

#### 执行

| 平台 | 命令 |
|------|------|
| Unix | `sh {domain}.sh` |
| Windows | `cmd /c {domain}.bat` |

#### 令牌缓存

令牌在内存中缓存直到过期。每次 API 调用时，swag2mcp 首先检查缓存 — 只有在缓存的令牌过期时才执行脚本。

#### 存根创建

当你配置 `auth: { type: script, config: { domain: "myapi" } }` 时，swag2mcp 会自动创建存根脚本：

**Unix（`auth_scripts/myapi.sh`）：**
```bash
#!/bin/sh
echo '{"token": "your-token-here", "expires_in": 3600}'
```

**Windows（`auth_scripts/myapi.bat`）：**
```bat
@echo off
echo {"token": "your-token-here", "expires_in": 3600}
```

将占位符令牌替换为你的实际认证逻辑。

#### 孤立清理

当你删除 spec 时，其认证脚本会成为孤立文件。swag2mcp 在以下情况下自动删除孤立脚本：
- `swag2mcp update`
- `swag2mcp clean`

## 命令

### update

```bash
swag2mcp update [path]
```

验证配置，清除缓存和响应，然后重新下载所有规范文件。还确保认证脚本存在并删除孤立脚本。

在以下情况后使用此命令：
- 添加或删除 collection
- 更改 collection 位置
- 编辑需要重新缓存的规范文件

### clean

```bash
swag2mcp clean [path]
```

删除 `cache/` 和 `responses/` 的所有内容，以及孤立的认证脚本。不会重新缓存规范 — 使用 `update` 进行重新缓存。

### validate

```bash
swag2mcp validate [path]
```

验证配置，包括所有 collection 位置。请参见[CLI: validate](../cli/validate.md)。

## 导出和导入

```bash
# 将工作区导出到 ZIP（默认名称：swag2mcp-backup-{date}.zip）
swag2mcp export

# 导出到特定路径
swag2mcp export /path/to/workspace /path/to/backup.zip

# 仅导出特定 spec
swag2mcp export --spec meteo

# 从备份恢复
swag2mcp import --from-zip /path/to/backup.zip
swag2mcp import /path/to/workspace /path/to/backup.zip
```

导出包括：`swag2mcp.yaml`、`specs/`、`auth_scripts/`。缓存和响应被排除（它们是本地数据）。

## .gitignore

如果你的工作区在 Git 仓库中，将以下条目添加到 `.gitignore`：

```gitignore
# swag2mcp — 仅本地数据
.swag2mcp/cache/
.swag2mcp/responses/
```

`cache/` 和 `responses/` 目录包含不应提交的本地机器特定数据。其他所有内容（`swag2mcp.yaml`、`specs/`、`auth_scripts/`）应放在仓库中，以便配置在团队中共享。
