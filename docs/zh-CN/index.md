# swag2mcp

<div style="background: #dc2626; color: white; padding: 20px 24px; border-radius: 12px; text-align: center; font-size: 1.4em; font-weight: 700; margin: 24px 0;">
  🚧 开发中 — 即将发布！
</div>

通过模型上下文协议（MCP）将 OpenAPI/Swagger/Postman API 规范与 LLM 智能体连接起来。

<a href="https://www.youtube.com/watch?v=1Da4UmE2f9U" target="_blank">
  <img src="https://raw.githubusercontent.com/mmadfox/swag2mcp/main/docs/cover.png" alt="预览">
</a>

## 你的 API 会说 LLM 语言

一行配置即可将任何 OpenAPI/Swagger/Postman 文件转换为 MCP 服务器。LLM 智能体可以发现、检查和调用你的 API — 无需编写集成代码。

<img src="/architecture.svg" width="700" alt="swag2mcp 架构">

## 告别重复的包装代码

每次将新 API 连接到 LLM 时，你都要编写相同的样板代码：规范解析、认证、错误处理、速率限制。swag2mcp 为你完成这一切 — 19 个现成的 MCP 工具。

## 谁需要它

| 角色 | 原因 |
|------|------|
| **AI 智能体开发者** | 2 分钟连接任何 API，而不是 2 天 |
| **MCP 工程师** | 无需处理代码 — 只需指向规范即可 |
| **架构师** | 为公司所有 LLM 提供统一的 API 集成层 |
| **数据分析师** | 通过自然语言访问 API，无需编码 |
| **DevOps / SRE** | 通过 LLM 进行监控和自动化，无需额外服务 |
| **集成工程师** | 9 种开箱即用的认证方法 — 从 Basic 到 OAuth2 到 HMAC |
| **QA 工程师** | 无需真实 API 即可进行隔离测试的模拟服务器 |
| **产品经理** | 无需后端工作即可快速构建 AI 功能原型 |
| **以及其他许多人** | |

---

## 许可证

基于 **GNU Affero General Public License v3.0**（AGPL v3）许可。

完整许可文本请参见 [LICENSE](https://github.com/mmadfox/swag2mcp/blob/main/LICENSE)。

```
SPDX-License-Identifier: AGPL-3.0-only
```
