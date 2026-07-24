# 技能

## 自定义输出格式

每个 swag2mcp MCP 工具返回结构化的 JSON 数据。这些数据如何**呈现**给用户取决于 LLM 的格式化技能 — 而你可以完全控制它。

### 默认格式技能

swag2mcp 附带一个内置的格式化技能，为每个工具响应定义了紧凑、人类可读的 markdown：

[swag2mcp-format SKILL.md](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md)

此技能涵盖所有 19 个 MCP 工具，具有：
- 列表的紧凑表格（spec、collection、标签、端点）
- 详情视图的内联头
- `inspect` 的紧凑模式表示
- 所有响应的一致样式

### 为什么技能很重要

相同的数据可以根据技能以截然不同的方式呈现：

| 样式 | 示例输出 |
|------|----------|
| **紧凑表格**（默认） | `GET /pet/{petId}` — Find pet by ID |
| **详细** | `Method: GET, Path: /pet/{petId}, Summary: Find pet by ID, Deprecated: false` |
| **极简** | `GET /pet/{petId}` |
| **技术** | `GET /pet/{petId} → 200: Pet object, 404: Not found` |
| **自定义** | 你能描述的任何格式 |

### 创建你自己的技能

你可以通过描述你想要的精确输出格式来编写自己的格式化技能。该技能是一个 markdown 文件，包含每个工具的格式化规则。以下是一些想法：

- **JSON 输出** — 返回原始 JSON 供机器使用
- **CSV 风格** — 用于电子表格导入的表格数据
- **图表友好** — API 结构的 Mermaid 或 ASCII 图表
- **极简** — 只有方法和路径，没有其他内容
- **文档风格** — 完整的描述、示例和注释

### 唯一的限制是模型

格式化输出的质量完全取决于 LLM 遵循你的格式化规则的能力。编写良好的技能，带有清晰的示例，会产生一致、可靠的输出。模糊的技能会产生不一致的结果。

你可以：
- 按原样使用默认技能
- 分叉并根据自己的喜好调整格式
- 从头开始编写自己的技能
- 根据任务在技能之间切换

### 如何使用技能

技能由 LLM 客户端（OpenCode、Cursor、Claude Desktop 等）作为其系统提示或智能体配置的一部分加载。请参阅客户端的文档以了解如何附加技能文件。

对于 OpenCode，技能在 `opencode.json` 中配置：

```json
{
  "skills": [
    {
      "name": "swag2mcp-format",
      "sourceURL": "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md"
    }
  ]
}
```
