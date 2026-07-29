# swag2mcp-format

**swag2mcp-format** 技能教会您的 LLM 以干净、紧凑、易读的 Markdown 格式显示 swag2mcp MCP 工具响应。没有这个技能，LLM 会自行决定如何格式化响应 — 通常冗长且不一致。

## 覆盖范围

所有 swag2mcp MCP 工具：

- `spec_list`, `spec_by_id` — 规范概览和详情
- `collection_by_spec`, `collection_by_id` — 集合与标签
- `tag_by_spec`, `tag_by_collection`, `tag_by_id` — 标签列表
- `endpoint_by_spec`, `endpoint_by_collection`, `endpoint_by_tag`, `endpoint_by_id` — 端点列表
- `search` — 搜索结果
- `inspect` — 完整操作详情与紧凑模式
- `invoke` — API 调用结果
- `auth` — 身份验证信息
- `info` — 运行时信息

## 通过 LLM 代理安装

将以下请求复制到您的 AI 驱动 IDE：

```
创建目录 .agents/skills/swag2mcp-format/ 并从
https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md 添加技能
```

## 手动安装

```bash
mkdir -p .agents/skills/swag2mcp-format
curl -o .agents/skills/swag2mcp-format/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
```

## 需要重启

添加技能后，请重启您的 LLM 客户端或 IDE（参见[概述](overview.md#需要重启)）。
