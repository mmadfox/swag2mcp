# CLI 命令

## 概述

`swag2mcp` CLI 是所有操作的单一入口点 — 从初始化工作区和管理 API 规范到启动用于 LLM 集成的 MCP 服务器。它提供 **13 个命令**，涵盖使用 OpenAPI/Swagger/Postman 规范的完整生命周期。

### CLI 解决的问题

- **工作区生命周期** — 创建（`init`）、检查（`info`、`ls`）、清理（`clean`）、更新（`update`）和删除（`delete`）工作区及其内容
- **Spec 和 collection 管理** — 添加（`add`）、列出（`ls`）和删除（`delete`）API 规范及其 collection
- **运行模式** — 启动 MCP 服务器以进行 LLM 工具访问（`mcp`）或启动交互式 TUI 浏览器（`run`）
- **诊断** — 验证配置（`validate`）、显示版本（`version`）、显示运行时信息（`info`）
- **备份和恢复** — 通过 ZIP 进行完整工作区往返（`export`、`import`）

### 关键细节

- **路径解析** — 接受 `[path]` 的命令期望的是**工作区目录**（不是文件路径）。解析顺序：显式 `[path]` → 当前目录（`./`）→ `~/.swag2mcp/`。CLI 自动追加 `swag2mcp.yaml`。作为服务运行或在 IDE 配置中时，始终传递显式路径，以避免加载错误的工作区。
- **Spec vs Collection** — **spec** 代表逻辑 API 服务（例如"Open-Meteo API"），而 **collection** 是一个 OpenAPI/Swagger/Postman 文件。一个 spec 可以有多个 collection。
- **`--version`** 既支持作为标志（`swag2mcp --version`），也支持作为子命令（`swag2mcp version`）。
- **`add spec` / `add collection`** 通过 `--yaml`（内联字符串或 `-` 表示 stdin）接受 YAML 输入。从文件或 heredoc 管道输入可以避免特殊字符的 shell 引号问题。
- **`delete`** 需要 TTY（交互式终端）。没有 `--force` 或 `--yes` 标志 — 它始终提示选择和确认。
- **`mcp`** 是 LLM 集成的主要命令。它支持三种传输方式：`stdio`（默认）、`sse` 和 `streamable-http`。`--disable-llm-auth` 标志（默认：`true`）从 MCP 工具列表中移除 `auth` 工具，防止 LLM 看到或请求令牌。认证仍然有效 — 令牌通过标准配置机制获取，而不是通过 LLM。此模式推荐用于**生产环境**（LLM 永远无法访问凭据）。对于**调试**或使用短期令牌时，设置 `--disable-llm-auth=false` 让 LLM 通过 `auth` 工具请求新令牌。
- **`validate`** 检查 YAML 语法、配置结构、规范文件存在性、URL 可达性、规范格式（OpenAPI/Swagger/Postman）、认证设置和 HTTP 客户端正确性。它**不**测试认证端点或 API 端点可用性。
- **`export` / `import`** 提供完整的工作区往返 — 配置文件、规范文件、缓存和认证脚本都包含在 ZIP 归档中。
- **`clean`** 删除 `cache/` 和 `responses/` 目录，但保留 `specs/` 和 `auth_scripts/`。旧响应（超过 48 小时）也会在 `mcp` 启动时自动清理。

## 命令

| 命令 | 描述 |
|------|------|
| [`init`](/cli/init) | 使用默认配置初始化工作区目录 |
| [`add`](/cli/add) | 向配置添加 spec 或 collection |
| [`delete`](/cli/delete) | 交互式删除 spec 或 collection |
| [`ls`](/cli/ls) | 列出所有 spec 及其 collection |
| [`run`](/cli/run) | 启动交互式 TUI API 浏览器 |
| [`validate`](/cli/validate) | 验证配置和规范文件 |
| [`clean`](/cli/clean) | 清除缓存的规范和调用响应 |
| [`update`](/cli/update) | 重新验证、重新缓存和重新索引所有 spec |
| [`mcp`](/cli/mcp) | 启动 MCP 服务器以进行 LLM 工具访问 |
| [`version`](/cli/version) | 打印 swag2mcp 版本 |
| [`info`](/cli/info) | 显示详细的配置和运行时信息 |
| [`import`](/cli/import) | 导入规范文件或从 ZIP 恢复工作区 |
| [`export`](/cli/export) | 将工作区导出为可移植的 ZIP 备份 |

## 全局标志

| 标志 | 描述 |
|------|------|
| `--version` | 显示版本（与 `version` 子命令相同） |
| `--help` | 显示任何命令的帮助信息 |
