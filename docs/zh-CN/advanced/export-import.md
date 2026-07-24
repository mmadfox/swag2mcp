# 导出和导入

## 概述

swag2mcp 支持通过 ZIP 归档进行完整的工作区往返。你可以将整个工作区（配置、规范文件、认证脚本）导出到 ZIP 文件，并在另一台机器上恢复。

## 导出

创建工作区的可移植 ZIP 备份。

```bash
# 导出到默认文件（swag2mcp-backup-<timestamp>.zip）
swag2mcp export

# 使用自定义路径导出
swag2mcp export --output ~/backups/swag2mcp-backup.zip

# 仅导出特定 spec
swag2mcp export --spec meteo
swag2mcp export --spec meteo,store
```

### 导出包含的内容

| 项目 | 描述 |
|------|------|
| `swag2mcp.yaml` | 配置文件 |
| `specs/` | 所有规范文件（OpenAPI/Swagger/Postman） |
| `auth_scripts/` | 认证脚本 |
| `swag2mcp.meta` | 元数据（版本信息，用于兼容性） |

缓存和响应**不会被导出** — 它们是临时数据，恢复时已过时。

### 默认文件名

如果你不指定输出路径，文件将保存为当前目录下的 `swag2mcp-backup-<YYYY-MM-DD-HHMMSS>.zip`（UTC 时间戳）。

## 导入

从 ZIP 备份恢复工作区或导入规范文件。

### 从 ZIP 恢复

```bash
# 恢复完整工作区
swag2mcp import --from-zip /path/to/backup.zip

# 覆盖恢复
swag2mcp import --from-zip /path/to/backup.zip -f
```

ZIP 必须由 `swag2mcp export` 创建 — 任意 ZIP 文件将无法工作。

### 导入单个规范文件

下载规范文件并添加到工作区：

```bash
swag2mcp import https://example.com/spec.yaml myspec
swag2mcp import /path/to/workspace https://example.com/spec.yaml myspec
```

### 从现有配置批量导入

下载指定 spec（域）的所有 collection 规范文件：

```bash
swag2mcp import --spec meteo
swag2mcp import /path/to/workspace --spec meteo,store
```

这会下载每个 collection 的规范文件，保存到 `specs/`，并更新配置以指向本地副本。

## 使用场景

### 备份

```bash
swag2mcp export --output swag2mcp-$(date +%Y-%m-%d).zip
```

### 迁移到另一台机器

```bash
# 在旧机器上
swag2mcp export --output swag2mcp.zip

# 将 ZIP 复制到新机器，然后：
swag2mcp import --from-zip swag2mcp.zip
```

### 共享配置

```bash
swag2mcp init
swag2mcp export --output template.zip
# 与同事共享 template.zip
```

## 导出后验证

始终验证 ZIP 文件是否已创建：

```bash
ls -la swag2mcp-backup-*.zip
```

## 重要说明

- **输出必须是 `.zip` 结尾的文件路径** — 不要传递目录
- **缓存和响应被排除** — 只保留配置、规范和认证脚本
- **ZIP 是自包含的** — 可以在任何安装了 swag2mcp 的机器上恢复
- **Spec 过滤器** — 使用 `--spec` 仅导出或导入特定 spec
