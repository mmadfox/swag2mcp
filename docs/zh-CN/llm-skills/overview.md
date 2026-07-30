# Skills for LLM

技能是 Markdown 文件，它们教会您的 LLM 代理如何更有效地使用 swag2mcp。它们作为代理系统提示的一部分加载，为 LLM 提供格式化响应和理解 CLI 命令的精确指令。

## 可用技能

| 技能 | 描述 | 下载 |
|------|------|------|
| **swag2mcp-format** | 将所有 MCP 工具响应格式化为紧凑、易读的 Markdown 表格 | [SKILL.md](https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md) |
| **swag2mcp-cli** | 完整的 CLI 参考 — LLM 知道每个命令、标志和配置选项 | [SKILL.md](https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md) |

## 为什么技能很重要

没有格式化技能，LLM 会自行决定如何显示工具结果 — 通常冗长且不一致。格式化技能确保每个响应遵循相同的简洁模式：列表用紧凑表格，详情用内联标题，模式用紧凑表示。

CLI 技能让 LLM 能够准确回答关于 swag2mcp 命令的任何"如何..."问题，而无需猜测。

## 通过 LLM 代理安装

将以下请求复制到您的 AI 驱动 IDE（OpenCode、Cursor、Claude Desktop、VS Code 等）：

```
将 swag2mcp 技能添加到我的项目：

1. 创建目录 .agents/skills/swag2mcp-format/ 并从 https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md 添加技能
2. 创建目录 .agents/skills/swag2mcp-cli/ 并从 https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md 添加技能
```

代理将下载两个技能文件并将其放置在正确的目录中。

## 手动安装

如果您的 LLM 客户端不支持基于代理的设置，请手动下载文件：

```bash
mkdir -p .agents/skills/swag2mcp-format
mkdir -p .agents/skills/swag2mcp-cli

curl -o .agents/skills/swag2mcp-format/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
curl -o .agents/skills/swag2mcp-cli/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## 配置 LLM 客户端

每个 LLM 客户端和 IDE 都有自己的技能安装方式。以下示例适用于 **OpenCode** — 请查看您的客户端文档以了解正确的方法。

```json
{
  "skills": [
    {
      "name": "swag2mcp-format",
      "sourceURL": "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md"
    },
    {
      "name": "swag2mcp-cli",
      "sourceURL": "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md"
    }
  ]
}
```

## 需要重启

**添加技能后，请重启您的 LLM 客户端或 IDE。** 某些工具仅在启动时加载技能。如果技能似乎没有生效，请尝试：

- **OpenCode**：重启应用程序或重新运行 opencode 命令
- **Cursor**：关闭并重新打开窗口（`Cmd+Shift+W` / `Ctrl+Shift+W`）
- **Claude Desktop**：退出并重新启动应用程序
- **VS Code**：重新加载窗口（`Ctrl+Shift+P` → "Developer: Reload Window"）
