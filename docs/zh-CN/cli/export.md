# export

## 用途

创建工作区的可移植 ZIP 备份。归档包含配置文件、所有规范文件和认证脚本 — 在另一台机器上恢复工作区所需的一切。

## 何时使用

- 你想在更改之前备份工作区
- 你正在将 swag2mcp 迁移到另一台机器
- 你想与同事共享 API 配置
- 你正在准备可重现的环境

## 语法

```bash
swag2mcp export [path] [output] [flags]
```

## 参数

| 参数 | 位置 | 必需 | 描述 |
|------|------|------|------|
| `path` | 1 | 否 | 工作区目录。如果省略，通过路径解析规则解析。 |
| `output` | 2 | 否 | 输出 ZIP 文件的完整路径。如果省略，默认为 `./swag2mcp-backup-&lt;timestamp&gt;.zip`。 |

## 标志

| 标志 | 简写 | 类型 | 默认值 | 描述 |
|------|------|------|--------|------|
| `--spec` | `-s` | `stringSlice` | `nil` | 仅导出指定的 spec（逗号分隔） |

## 工作原理

### 默认导出

在当前目录中创建带时间戳名称的 ZIP：

```bash
swag2mcp export
# 创建 ./swag2mcp-backup-2026-07-22-143022.zip
```

### 自定义输出路径

```bash
swag2mcp export /path/to/workspace /path/to/backup.zip
```

### 导出特定 spec

```bash
swag2mcp export --spec meteo
swag2mcp export --spec meteo,store
```

## ZIP 中的内容

| 条目 | 描述 |
|------|------|
| `swag2mcp.meta` | 关于导出的元数据 |
| `swag2mcp.yaml` | 配置文件 |
| `specs/` | 所有规范文件（OpenAPI/Swagger/Postman） |
| `auth_scripts/` | 认证脚本 |
| `cache/` | 空（缓存不被导出） |
| `responses/` | 空（响应不被导出） |

## 恢复

使用 `import` 从备份恢复：

```bash
swag2mcp import --from-zip /path/to/backup.zip
```

## 命令后验证

始终验证 ZIP 文件是否已创建：

```bash
ls -la swag2mcp-backup-*.zip
# 或对于自定义输出路径：
ls -la /path/to/backup.zip
```

## 细节

- **输出必须是文件路径：** `[output]` 参数必须是 `.zip` 结尾的完整文件路径。**不要**传递目录 — 如果给定目录路径，命令不会创建 ZIP。
- **默认文件名：** `swag2mcp-backup-<YYYY-MM-DD-HHMMSS>.zip`，使用 UTC 时间戳。
- **`--spec` 过滤器：** 设置后，只包含指定的 spec。其他 spec 被排除在归档之外。
- **无需配置：** `export` 即使没有有效的配置文件也能工作。它导出工作区中存在的任何内容。
- **缓存和响应被排除：** 这些是临时数据，恢复时已过时。只保留配置、规范和认证脚本。
