# import

## 用途

将规范文件导入工作区的 `specs/` 目录以供本地使用，或从 ZIP 备份恢复完整工作区。三种模式涵盖不同场景：添加单个规范、从现有配置批量导入，或恢复完整工作区。

## 何时使用

- 你有规范 URL 或文件，想将其本地保存到工作区
- 你想下载配置中的所有 collection 规范文件，使工作区自包含
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
| `name` | 3 | 视情况 | 保存的文件名（例如 `example-api.yaml`）。如果省略，从 URL 自动生成。 |

## 标志

| 标志 | 简写 | 类型 | 默认值 | 描述 |
|------|------|------|--------|------|
| `--spec` | `-s` | `string` | `""` | 从配置下载 collection 规范文件。无值时导入所有 specs，或指定域名如 `--spec meteo,github` |
| `--force` | `-f` | `bool` | `false` | 覆盖现有规范文件而不报错 |
| `--from-zip` | | `string` | `""` | 从 swag2mcp 备份 ZIP 恢复工作区 |

## 工作原理

### 模式 1 — 从 URL 或文件单个导入

下载规范文件并保存到 `specs/`：

```bash
swag2mcp import https://example.com/spec.yaml example-api.yaml
swag2mcp import /path/to/workspace https://example.com/spec.yaml example-api.yaml
swag2mcp import ./local-spec.yaml example-api.yaml
```

如果省略 `name`，则从 URL 的文件名自动生成：
```bash
swag2mcp import https://example.com/specs/petstore.yaml
# → 保存为 petstore.yaml
```

使用 `--force` 覆盖现有文件：
```bash
swag2mcp import --force https://example.com/spec.yaml example-api.yaml
```

导入后，输出显示工作区路径、保存的文件以及要添加到 `swag2mcp.yaml` 的 YAML 模板：

```
✅ Imported to /path/to/workspace
   specs/example-api.yaml

   Add to swag2mcp.yaml:
     specs:
       - domain: <your-domain>
         collections:
           - location: specs/example-api.yaml
```

### 模式 2 — 从现有配置批量导入（`--spec`）

从配置的 `location` URL 下载指定域的所有 collection 规范文件，保存到 `specs/`，并更新配置以指向本地副本：

```bash
swag2mcp import --spec                # 所有 specs
swag2mcp import --spec meteo           # 特定 spec
swag2mcp import --spec meteo,github    # 多个 specs
swag2mcp import /path/to/workspace --spec meteo
```

如果指定的域在配置中不存在，命令将返回错误：
```
Error: import_no_match
  Spec "nonexistent" not found in config.
```

这使工作区自包含 — 导入后不再需要远程规范 URL。

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
- **HTTP 客户端：** 导入期间应用配置中的全局 HTTP 客户端设置（超时、代理、头等）。
