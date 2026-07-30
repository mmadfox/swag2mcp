# swag2mcp-cli

**swag2mcp-cli** 技能为您的 LLM 提供完整的 swag2mcp CLI 参考 — 每个命令、标志、参数和配置选项。有了这个技能，LLM 可以准确回答"如何..."问题，而无需猜测。

## 覆盖范围

所有 13 个 CLI 命令：

| 命令 | 用途 |
|------|------|
| `init` | 初始化工作区和配置 |
| `add` | 添加规范或集合 |
| `delete` | 删除规范或集合 |
| `ls` | 列出已配置的规范 |
| `run` | 启动 API Explorer TUI |
| `validate` | 验证配置文件 |
| `clean` | 清除缓存数据 |
| `update` | 从配置更新缓存 |
| `mcp` | 启动 MCP 服务器 |
| `version` | 显示版本信息 |
| `info` | 显示运行时信息 |
| `import` | 从 ZIP 文件导入工作区 |
| `export` | 将工作区导出为 ZIP 文件 |

以及所有标志、配置文件结构、身份验证方法和高级选项。

## 直接链接

<https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md>

## 通过 LLM 代理安装

将以下请求复制到您的 AI 驱动 IDE：

```
创建目录 .agents/skills/swag2mcp-cli/ 并从
https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md 添加技能
```

## 手动安装

```bash
mkdir -p .agents/skills/swag2mcp-cli
curl -o .agents/skills/swag2mcp-cli/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## 需要重启

添加技能后，请重启您的 LLM 客户端或 IDE（参见[概述](overview.md#需要重启)）。
