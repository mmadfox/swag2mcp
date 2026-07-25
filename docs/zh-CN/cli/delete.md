# delete

## 用途

从配置中删除 **spec**（API 服务）或 **collection**（规范文件）。这是 `add` 的逆操作。

## 何时使用

- API 不再需要
- 你想从 spec 中删除特定的规范文件
- 你正在清理工作区

## 语法

```bash
swag2mcp delete spec [path]
swag2mcp delete collection [path]
```

## 参数

| 参数 | 位置 | 必需 | 描述 |
|------|------|------|------|
| `path` | 1 | 否 | 工作区目录。如果省略，通过路径解析规则解析。 |

## 标志

无。两个子命令都是纯交互式的。

## 工作原理

### 删除 spec

提示你从列表中选择一个 spec，然后在删除前要求确认。

```bash
swag2mcp delete spec
```

### 删除 collection

提示你选择一个 spec，然后选择该 spec 中的一个 collection，然后要求确认。

```bash
swag2mcp delete collection
```

## 查找 ID

交互式提示显示人类可读的名称，而不是 ID。如果你需要 ID 作为参考：

```bash
# 列出所有 spec 及其 ID
swag2mcp ls

# 列出特定 spec 的 collection
swag2mcp ls --tags
```

## 命令后验证

```bash
swag2mcp ls [path]
# 删除的 spec 或 collection 应不再出现
```

## 细节

- **需要 TTY：** 两个命令都需要交互式终端。它们在 CI/CD 管道、cron 作业或非交互式脚本中**无法工作**。
- **没有 `--force` 或 `--yes`：** 无法跳过确认提示。这是有意为之，防止意外删除。
- **自动初始化：** 如果不存在配置文件，`delete` 会自动先运行初始化向导。
- **没有 YAML 模式：** 与 `add` 不同，没有 `--yaml` 标志。删除始终是交互式的。
