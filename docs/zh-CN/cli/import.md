# import

## 用途

将规范文件导入工作区，或从 ZIP 备份恢复完整工作区。三种模式涵盖不同场景：添加单个规范、从现有配置批量导入，或恢复完整工作区。

## 何时使用

- 你有规范 URL 或文件，想将其添加到工作区
- 你想下载配置中引用的所有规范文件
- 你需要从 `export` 创建的 ZIP 备份恢复工作区
- 你正在将 swag2mcp 迁移到另一台机器

## 语法

```bash
swag2mcp import [path] [source] [name] [flags]
```

## 参数

| 参数 | 位置 | 必需 | 描述 |
|------|------|------|------|
| `path` | 1 | 否 | 工作区目录。如果省略，通过路径解析规则解析。 |
| `source` | 2 | 视情况 | 规范文件的 URL 或本地路径，或 ZIP 归档的路径 |
| `name` | 3 | 视情况 | 新 spec 的域名 |

## 标志

| 标志 | 简写 | 类型 | 默认值 | 描述 |
|------|------|------|--------|------|
| `--spec` | `-s` | `stringSlice` | `nil` | 从指定 spec 导入 collection（逗号分隔） |
| `--from-zip` | | `string` | `""` | 从 swag2mcp 备份 ZIP 恢复工作区 |

## 工作原理

### 模式 1 — 从 URL 或文件单个导入

下载规范文件并添加域名到工作区：

```bash
swag2mcp import https://example.com/spec.yaml myspec
swag2mcp import /path/to/workspace https://example.com/spec.yaml myspec
swag2mcp import ./local-spec.yaml myspec
```

规范文件保存到 `specs/`，配置更新为新的 spec 条目。

### 模式 2 — 从现有配置批量导入

从配置的 URL 下载指定域的所有 collection：

```bash
swag2mcp import --spec meteo
swag2mcp import /path/to/workspace --spec meteo,store
```

每个 collection 的规范文件被下载并保存到 `specs/`。配置更新为指向本地副本。

### 模式 3 — 从 ZIP 备份恢复

从 `swag2mcp export` 创建的 ZIP 归档恢复完整工作区：

```bash
swag2mcp import --from-zip /path/to/backup.zip
swag2mcp import /path/to/workspace /path/to/backup.zip
```

> **ZIP 必须由 `swag2mcp export` 创建。** 任意 ZIP 文件将无法工作 — 归档具有特定的内部结构（`swag2mcp.yaml`、`specs/`、`auth_scripts/`）。

## 命令后验证

```bash
# 单个或批量导入
swag2mcp ls [path]
# 新的 spec 应出现在列表中

# ZIP 恢复
swag2mcp ls [path]
# 备份中的所有 spec 应出现
```

## 细节

- **批量模式需要配置：** 使用 `--spec` 时，配置文件必须存在。如果需要，先运行 `init`。
- **单个导入创建工作区：** 如果工作区不存在，会自动创建。
- **ZIP 检测：** 以 `.zip` 结尾的位置参数被视为 ZIP 源。`--from-zip` 标志优先于位置检测。
- **`--force`：** 可用于 ZIP 恢复以覆盖现有工作区。
- **HTTP 客户端：** 导入期间应用配置中的全局 HTTP 客户端设置（超时、代理、头等）。
