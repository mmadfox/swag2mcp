# init

## 用途

`init` 命令创建一个**工作区** — 包含 `swag2mcp.yaml` 配置文件以及缓存、规范、响应和认证脚本子目录的目录。这是设置 swag2mcp 时要运行的第一个命令。

## 何时使用

- 你第一次设置 swag2mcp
- 你想在特定目录中创建新工作区
- 你需要重新初始化损坏或丢失的工作区

## 语法

```bash
swag2mcp init [path] [flags]
```

## 参数

| 参数 | 位置 | 必需 | 描述 |
|------|------|------|------|
| `path` | 1 | 否 | 工作区目录。如果省略，默认为 `~/.swag2mcp`。 |

## 标志

| 标志 | 简写 | 类型 | 默认值 | 描述 |
|------|------|------|--------|------|
| `--interactive` | `-i` | `bool` | `false` | 运行交互式 TUI 向导 |
| `--force` | `-f` | `bool` | `false` | 覆盖非空目录中的现有配置 |

## 工作原理

### 非交互模式（默认）

创建没有 spec 的最小 `swag2mcp.yaml`。之后你手动编辑文件。

```bash
swag2mcp init
# 创建 ~/.swag2mcp/swag2mcp.yaml

swag2mcp init ./my-project
# 创建 ./my-project/swag2mcp.yaml

swag2mcp init /absolute/path
# 创建 /absolute/path/swag2mcp.yaml
```

### 交互模式（`-i`）

启动 18 步 TUI 向导，引导你完成：

1. 选择工作区目录
2. 添加带域、标题、基础 URL 的 spec
3. 配置带位置 URL 的 collection
4. 设置认证（全部 9 种方法）
5. 配置 HTTP 客户端设置（超时、代理、头等）

```bash
swag2mcp init -i
```

### 强制模式（`--force`）

默认情况下，`init` 拒绝在非空目录中运行。使用 `--force` 覆盖：

```bash
swag2mcp init -f
swag2mcp init ./existing-dir -f
```

## 创建的内容

```
~/.swag2mcp/
├── swag2mcp.yaml       # 配置文件
├── cache/               # 下载的远程规范文件
├── specs/               # 本地规范文件
├── responses/           # 保存的 API 调用响应
└── auth_scripts/        # 认证脚本（用于 ScriptAuth 类型）
```

## 命令后验证

```bash
ls ~/.swag2mcp/swag2mcp.yaml
# 如果文件存在，init 成功
```

## 细节

- **路径解析：** `[path]` 是**工作区目录**，不是文件路径。CLI 自动追加 `swag2mcp.yaml`。解析顺序：显式 `[path]` → 当前目录（`./`）→ `~/.swag2mcp/`。
- **非空目录检查：** 没有 `--force` 时，如果目标目录存在且不为空，`init` 返回错误。这防止意外覆盖。
- **认证脚本存根：** 如果任何 spec 使用 `ScriptAuth`，`init` 在 `auth_scripts/` 中创建存根脚本文件（Unix 上为 `.sh`，Windows 上为 `.bat`）。
- **输出：** 成功时，打印配置路径和提示：`"Next step: edit swag2mcp.yaml or run 'swag2mcp ls' to list configured specs"`。
