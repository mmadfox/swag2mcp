# 快速开始

2 分钟内让 swag2mcp 运行起来。

## 1. 初始化

### 主目录（推荐）

一次性设置，适用于整个系统。配置存储在你的主文件夹中。

::: code-group

```bash [macOS / Linux]
swag2mcp init
# 创建 ~/.swag2mcp/swag2mcp.yaml
```

```powershell [Windows]
swag2mcp.exe init
# 创建 %USERPROFILE%\.swag2mcp\swag2mcp.yaml
```

:::

### 项目目录

用于项目内的隔离工作区。

::: code-group

```bash [macOS / Linux]
mkdir -p ./swag2mcp && swag2mcp init ./swag2mcp
```

```powershell [Windows]
mkdir ./swag2mcp; swag2mcp.exe init ./swag2mcp
```

:::

### 从 ZIP

如果你有现成的工作区（例如来自同事）：

```bash
swag2mcp import --from-zip workspace.zip
```

## 2. 安装智能体技能（推荐）

安装 swag2mcp 技能，教会你的 AI 智能体所有命令、标志、配置格式和真实示例。

询问你的智能体：

```bash
"Add the swag2mcp-cli skill from https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md"
"Add the swag2mcp-format skill from https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md"
```

> 某些 IDE 在添加技能后需要重启。

## 3. LLM 客户端 / IDE 配置

配置你的 IDE 以连接到 swag2mcp。IDE 会在需要时自动启动 MCP 服务器。

::: code-group

```json [OpenCode]
{
  "mcp": {
    "swag2mcp": {
      "type": "local",
      "command": ["swag2mcp", "mcp"],
      "enabled": true
    }
  }
}
```

```json [Claude Desktop]
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

```json [Crush]
{
  "mcp": {
    "swag2mcp": {
      "type": "stdio",
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

:::

对于其他 IDE（Cursor、VS Code、JetBrains），请参见[集成指南](../integration/opencode.md)。

> 如果你在自定义路径（例如 `./swag2mcp`）初始化了工作区，请在命令中使用完整路径：
> `"command": ["swag2mcp", "mcp", "/absolute/path/to/swag2mcp"]`

> **任何配置更改后，重启 MCP 服务器** 以使更改生效。

## 4. 启动 MCP 服务器

### stdio（默认）— 用于本地 IDE

无需配置。你的 IDE 通过上述配置自动启动 swag2mcp。

```bash
swag2mcp mcp
```

### SSE / Streamable HTTP — 用于远程访问

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

或在 `swag2mcp.yaml` 中配置：

```yaml
mcp:
  transport: sse
  addr: ":8080"
  path: "/mcp"
```

所有标志请参见[MCP 服务器参考](../configuration/mcp-server.md)。

### 按标签过滤 spec

```bash
swag2mcp mcp --tags weather,public
```

只有具有匹配标签的 spec 才会对 LLM 可用。

### 验证是否正常工作

连接后，询问你的 LLM 智能体：

```bash
"What MCP tools do you support?"
```

如果智能体列出了 swag2mcp 工具（`spec_list`、`search`、`invoke` 等）— 一切正常。

### 可以尝试的示例查询

| 询问你的智能体 | 发生了什么 |
|-------|-------------|
| "纽约的天气怎么样？" | `invoke` — 调用 Open-Meteo 预报 API |
| "当前 BTC 价格是多少？" | `invoke` — 调用 Binance 行情 API |
| "给我讲个冷笑话" | `invoke` — 调用 icanhazdadjoke API |
| "给我看看皮卡丘" | `invoke` — 按名称调用 PokéAPI |
| "谁是 Rick Sanchez？" | `invoke` — 调用 Rick and Morty 角色 API |
| "北京的空气质量如何？" | `invoke` — 调用 Open-Meteo 空气质量 API |
| "葡萄牙附近的海浪有多高？" | `invoke` — 调用 Open-Meteo 海洋 API |
| "搜索关于狗的笑话" | `invoke` — 调用 dadjoke 搜索端点 |
| "列出所有宝可梦" | `invoke` — 调用 PokéAPI 列表端点 |
| "珠穆朗玛峰的海拔是多少？" | `invoke` — 调用 Open-Meteo 海拔 API |

## 5. 下一步是什么？

- [概念](../concepts/overview.md) — 了解架构
- [配置](../configuration/config-file.md) — 自定义设置
- [CLI 命令](../cli/overview.md) — 完整命令参考
